package cmds

import (
	"context"
	"fmt"
	"io"

	client "github.com/aserto-dev/go-aserto/client"
	dsc2 "github.com/aserto-dev/go-directory/aserto/directory/common/v2"
	dse2 "github.com/aserto-dev/go-directory/aserto/directory/exporter/v2"
	dsr2 "github.com/aserto-dev/go-directory/aserto/directory/reader/v2"
	"github.com/rs/zerolog/log"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const pageSize int32 = 100

type ObjRelSub struct {
	Object   string
	Relation string
	Subject  string
}

func (ors ObjRelSub) Key() string {
	return ors.Object + "|" + ors.Relation + "|" + ors.Subject
}

type MigrateCmd struct {
}

func (a *MigrateCmd) Run(c *Common) error {
	ctx := context.Background()

	opts := []client.ConnectionOption{
		client.WithAddr(c.Host),
		client.WithAPIKeyAuth(c.APIKey),
		client.WithTenantID(c.TenantID),
		client.WithInsecure(c.Insecure),
	}

	clnt, err := client.NewConnection(ctx, opts...)
	if err != nil {
		return err
	}

	e2 := dse2.NewExporterClient(clnt.Conn)
	r2 := dsr2.NewReaderClient(clnt.Conn)

	ots := a.getObjectTypes(ctx, r2)

	rts := a.getRelationTypes(ctx, r2)

	ors, rsc := a.getObjectRelationSubject(ctx, e2)

	fmt.Println("rts", len(rts))
	fmt.Println("ots", len(ots))
	fmt.Println("rsc", rsc)
	fmt.Println("ors", len(ors))

	for i, r := range ors {
		fmt.Printf("%d  %s %s %s\n", i, r.Object, r.Relation, r.Subject)
	}

	return nil
}

func (a *MigrateCmd) getObjectTypes(ctx context.Context, r2 dsr2.ReaderClient) []*dsc2.ObjectType {
	token := ""
	results := []*dsc2.ObjectType{}
	for {
		resp, err := r2.GetObjectTypes(ctx, &dsr2.GetObjectTypesRequest{Page: &dsc2.PaginationRequest{Size: pageSize, Token: token}})
		if err != nil {
			log.Error().Err(err).Msg("GetObjectTypes")
			return []*dsc2.ObjectType{}
		}
		results = append(results, resp.Results...)
		if resp.Page.NextToken == "" {
			break
		}
		token = resp.Page.NextToken
	}
	return results
}

func (a *MigrateCmd) getRelationTypes(ctx context.Context, r2 dsr2.ReaderClient) []*dsc2.RelationType {
	token := ""
	results := []*dsc2.RelationType{}
	for {
		resp, err := r2.GetRelationTypes(ctx, &dsr2.GetRelationTypesRequest{Page: &dsc2.PaginationRequest{Size: pageSize, Token: token}})
		if err != nil {
			log.Error().Err(err).Msg("GetRelationTypes")
			return []*dsc2.RelationType{}
		}
		results = append(results, resp.Results...)
		if resp.Page.NextToken == "" {
			break
		}
		token = resp.Page.NextToken
	}
	return results
}

func (a *MigrateCmd) getObjectRelationSubject(ctx context.Context, e2 dse2.ExporterClient) ([]*ObjRelSub, int) {
	stream, err := e2.Export(ctx, &dse2.ExportRequest{
		Options:   uint32(dse2.Option_OPTION_DATA_RELATIONS),
		StartFrom: &timestamppb.Timestamp{},
	})
	if err != nil {
		log.Error().Err(err).Msg("ExportRequest")
		return []*ObjRelSub{}, 0
	}

	rsc := 0
	orsMap := map[string]*ObjRelSub{}

	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Error().Err(err).Msg("Recv")
			return []*ObjRelSub{}, 0
		}

		r, ok := msg.Msg.(*dse2.ExportResponse_Relation)
		if !ok {
			log.Warn().Msg("unknown message type, skipped")
			continue
		}

		rsc++

		i := ObjRelSub{
			Object:   r.Relation.Object.GetType(),
			Relation: r.Relation.Relation,
			Subject:  r.Relation.Subject.GetType(),
		}

		if _, ok := orsMap[i.Key()]; !ok {
			orsMap[i.Key()] = &i
		}
	}

	results := []*ObjRelSub{}
	for _, v := range orsMap {
		results = append(results, v)
	}

	return results, rsc
}
