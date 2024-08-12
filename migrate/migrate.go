// nolint: staticcheck
package migrate

import (
	"context"
	"fmt"
	"io"
	"os"
	"sort"

	"github.com/aserto-dev/azm/model"
	dsc2 "github.com/aserto-dev/go-directory/aserto/directory/common/v2"
	dse2 "github.com/aserto-dev/go-directory/aserto/directory/exporter/v2"
	dsr2 "github.com/aserto-dev/go-directory/aserto/directory/reader/v2"
	"github.com/samber/lo"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const pageSize int32 = 100

type Migrator struct {
	Metadata      *Metadata
	RelationMap   *ObjRelSubContainer
	PermissionMap *ObjPermRelContainer
}

func NewMigrator() *Migrator {
	return &Migrator{
		Metadata: &Metadata{
			ObjectTypes:   []*dsc2.ObjectType{},
			RelationTypes: []*dsc2.RelationType{},
			Permissions:   []*dsc2.Permission{},
		},
		RelationMap:   NewObjRelSubContainer(),
		PermissionMap: NewObjPermRelContainer(),
	}
}

// Load hydrates the migrator from given gRPC connection.
func (m *Migrator) Load(conn grpc.ClientConnInterface) error {
	e2 := dse2.NewExporterClient(conn)
	r2 := dsr2.NewReaderClient(conn)

	m.Metadata.ObjectTypes = m.getObjectTypes(context.Background(), r2)
	m.Metadata.RelationTypes = m.getRelationTypes(context.Background(), r2)
	m.Metadata.Permissions = m.getPermissions(context.Background(), r2)

	m.RelationMap = m.getObjectRelationSubject(context.Background(), e2)

	return nil
}

func (m *Migrator) Process() error {
	m.removeObsoleteObjectTypes()
	m.removeObsoleteRelationTypes()

	m.consolidatePermissions()

	m.fixupObjectTypes()
	m.sortObjectTypes()

	m.invertUnions()
	m.fixupRelationTypes()

	m.trueUpRelations()

	if err := m.normalize(); err != nil {
		return err
	}

	if err := m.validate(); err != nil {
		return err
	}

	return nil
}

func (m *Migrator) normalize() error {
	for i := 0; i < len(m.Metadata.ObjectTypes); i++ {
		if !model.IsValidIdentifier(m.Metadata.ObjectTypes[i].Name) {
			if normalized, err := model.NormalizeIdentifier(m.Metadata.ObjectTypes[i].Name); err == nil {
				m.Metadata.ObjectTypes[i].Name = normalized
			}
		}
	}

	for i := 0; i < len(m.Metadata.RelationTypes); i++ {
		if !model.IsValidIdentifier(m.Metadata.RelationTypes[i].Name) {
			if normalized, err := model.NormalizeIdentifier(m.Metadata.RelationTypes[i].Name); err == nil {
				m.Metadata.RelationTypes[i].Name = normalized
			}
		}

		for j := 0; j < len(m.Metadata.RelationTypes[i].Unions); j++ {
			if !model.IsValidIdentifier(m.Metadata.RelationTypes[i].Unions[j]) {
				if normalized, err := model.NormalizeIdentifier(m.Metadata.RelationTypes[i].Unions[j]); err == nil {
					m.Metadata.RelationTypes[i].Unions[j] = normalized
				}
			}
		}

		for j := 0; j < len(m.Metadata.RelationTypes[i].Permissions); j++ {
			if !model.IsValidIdentifier(m.Metadata.RelationTypes[i].Permissions[j]) {
				if normalized, err := model.NormalizeIdentifier(m.Metadata.RelationTypes[i].Permissions[j]); err == nil {
					m.Metadata.RelationTypes[i].Permissions[j] = normalized
				}
			}
		}
	}

	return nil
}

func (m *Migrator) validate() error {
	for _, ot := range m.Metadata.ObjectTypes {
		if !model.IsValidIdentifier(ot.Name) {
			fmt.Fprintf(os.Stderr, "ot %s is not a valid identifier\n", ot.Name)
		}
	}

	for _, rt := range m.Metadata.RelationTypes {
		if !model.IsValidIdentifier(rt.Name) {
			fmt.Fprintf(os.Stderr, "rt %s is not a valid identifier\n", rt.Name)
		}
		for _, u := range rt.Unions {
			if !model.IsValidIdentifier(u) {
				fmt.Fprintf(os.Stderr, "rt %s union %s is not a valid identifier\n", rt.Name, u)
			}
		}
		for _, p := range rt.Permissions {
			if !model.IsValidIdentifier(p) {
				fmt.Fprintf(os.Stderr, "rt %s permission %s is not a valid identifier\n", rt.Name, p)
			}
		}
	}

	return nil
}

func (m *Migrator) Write(w io.Writer, opts ...WriterOption) error {
	args := &WriterArgs{header: true, filename: "", description: "", timestamp: false}
	for _, opt := range opts {
		opt(args)
	}

	if args.header {
		writeManifestHeader(w)
	}

	writeManifestInfo(w, args)
	writeManifestModel(w, 2, 3)
	writeManifestTypes(w)

	for _, ot := range m.Metadata.ObjectTypes {
		rts := lo.Filter(m.Metadata.RelationTypes, func(item *dsc2.RelationType, index int) bool {
			return item.ObjectType == ot.Name
		})

		this := false
		if len(rts) == 0 && len(m.PermissionMap.GetPerms(ot.Name)) == 0 {
			this = true
		}

		writeTypeInstance(w, 2, ot, this)

		if len(rts) > 0 {
			writeRelations(w, 4)
		}

		for _, rt := range rts {
			writeRelationInstance(w, 6, rt)
		}

		if len(m.PermissionMap.GetPerms(ot.Name)) > 0 {
			writePermissions(w, 4)
		}

		writePermissionInstance(w, 6, m.PermissionMap.GetPerms(ot.Name))

		fmt.Fprintln(w)
	}

	return nil
}

func (m *Migrator) getObjectTypes(ctx context.Context, r2 dsr2.ReaderClient) []*dsc2.ObjectType {
	token := ""
	results := []*dsc2.ObjectType{}
	for {
		resp, err := r2.GetObjectTypes(ctx, &dsr2.GetObjectTypesRequest{Page: &dsc2.PaginationRequest{Size: pageSize, Token: token}})
		if err != nil {
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

func (m *Migrator) getRelationTypes(ctx context.Context, r2 dsr2.ReaderClient) []*dsc2.RelationType {
	token := ""
	results := []*dsc2.RelationType{}
	for {
		resp, err := r2.GetRelationTypes(ctx, &dsr2.GetRelationTypesRequest{Page: &dsc2.PaginationRequest{Size: pageSize, Token: token}})
		if err != nil {
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

func (m *Migrator) getPermissions(ctx context.Context, r2 dsr2.ReaderClient) []*dsc2.Permission {
	token := ""
	results := []*dsc2.Permission{}
	for {
		resp, err := r2.GetPermissions(ctx, &dsr2.GetPermissionsRequest{Page: &dsc2.PaginationRequest{Size: pageSize, Token: token}})
		if err != nil {
			return []*dsc2.Permission{}
		}
		results = append(results, resp.Results...)
		if resp.Page.NextToken == "" {
			break
		}
		token = resp.Page.NextToken
	}
	return results
}

func (m *Migrator) getObjectRelationSubject(ctx context.Context, e2 dse2.ExporterClient) *ObjRelSubContainer {
	result := NewObjRelSubContainer()

	stream, err := e2.Export(ctx, &dse2.ExportRequest{
		Options:   uint32(dse2.Option_OPTION_DATA_RELATIONS),
		StartFrom: &timestamppb.Timestamp{},
	})
	if err != nil {
		return result
	}

	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return NewObjRelSubContainer()
		}

		r, ok := msg.Msg.(*dse2.ExportResponse_Relation)
		if !ok {
			continue
		}

		result.Add(&ObjRelSub{
			Object:   r.Relation.Object.GetType(),
			Relation: r.Relation.Relation,
			Subject:  r.Relation.Subject.GetType(),
		})
	}

	return result
}

func (m *Migrator) removeObsoleteObjectTypes() {
	m.Metadata.ObjectTypes = lo.Filter(m.Metadata.ObjectTypes, func(x *dsc2.ObjectType, index int) bool {
		if len(m.RelationMap.GetRels(x.Name)) > 0 {
			return true
		}
		if len(m.PermissionMap.GetPerms(x.Name)) > 0 {
			return true
		}
		if ot, ok := RefObjectTypes[x.Name]; ok {
			return !IsObsolete(dsc2.Flag(ot.Status))
		}
		return true
	})
}

func (m *Migrator) fixupObjectTypes() {
	for i := 0; i < len(m.Metadata.ObjectTypes); i++ {
		ot, ok := RefObjectTypes[m.Metadata.ObjectTypes[i].Name]
		if !ok {
			rels := m.RelationMap.GetRels(m.Metadata.ObjectTypes[i].Name)
			perms := m.PermissionMap.GetPerms(m.Metadata.ObjectTypes[i].Name)

			ordinal := 3 + len(rels) + len(perms)
			m.Metadata.ObjectTypes[i].Ordinal = int32(ordinal)

			continue
		}

		if ot.DisplayName != "" {
			m.Metadata.ObjectTypes[i].DisplayName = ot.DisplayName
		}

		if ot.Ordinal != 0 {
			m.Metadata.ObjectTypes[i].Ordinal = ot.Ordinal
		}
	}

	for i := 0; i < len(m.Metadata.ObjectTypes); i++ {
		if m.Metadata.ObjectTypes[i].Ordinal == 0 {
			rels := m.RelationMap.GetRels(m.Metadata.ObjectTypes[i].Name)
			perms := m.PermissionMap.GetPerms(m.Metadata.ObjectTypes[i].Name)

			ordinal := len(rels) + len(perms[m.Metadata.ObjectTypes[i].Name])

			m.Metadata.ObjectTypes[i].Ordinal = int32(ordinal)
		}
	}
}

func (m *Migrator) sortObjectTypes() {
	sort.Slice(m.Metadata.ObjectTypes, func(i, j int) bool {
		return m.Metadata.ObjectTypes[i].Ordinal < m.Metadata.ObjectTypes[j].Ordinal
	})
}

func (m *Migrator) removeObsoleteRelationTypes() {
	m.Metadata.RelationTypes = lo.Filter(m.Metadata.RelationTypes, func(x *dsc2.RelationType, index int) bool {
		if rtMap, ok := RefRelationTypes[x.ObjectType]; ok {
			if rt, ok := rtMap[x.Name]; ok {
				return !IsObsolete(dsc2.Flag(rt.Status))
			}
		}
		return true
	})
}

func (m *Migrator) invertUnions() {
	for i := 0; i < len(m.Metadata.ObjectTypes); i++ {
		obj := m.Metadata.ObjectTypes[i]

		rels := lo.Filter(m.Metadata.RelationTypes, func(rel *dsc2.RelationType, _ int) bool {
			return rel.ObjectType == obj.Name
		})

		inverted := lo.Associate(rels, func(rel *dsc2.RelationType) (string, []string) {
			return rel.Name, []string{}
		})

		for _, rel := range rels {
			for _, union := range rel.Unions {
				if _, ok := inverted[union]; ok {
					inverted[union] = append(inverted[union], rel.Name)
				}
			}
		}

		for _, rel := range rels {
			rel.Unions = inverted[rel.Name]
		}
	}
}

func (m *Migrator) fixupRelationTypes() {
	for i := 0; i < len(m.Metadata.RelationTypes); i++ {
		rtMap, ok := RefRelationTypes[m.Metadata.RelationTypes[i].ObjectType]
		if !ok {
			continue
		}
		rt, ok := rtMap[m.Metadata.RelationTypes[i].Name]
		if !ok {
			continue
		}

		if rt.DisplayName != "" {
			m.Metadata.RelationTypes[i].DisplayName = rt.DisplayName
		}

		if len(rt.Unions) > 0 {
			m.Metadata.RelationTypes[i].Unions = lo.Union(m.Metadata.RelationTypes[i].Unions, rt.Unions)
		}
	}
}

func (m *Migrator) trueUpRelations() {
	for i := 0; i < len(m.Metadata.RelationTypes); i++ {
		subs := m.RelationMap.GetSubs(m.Metadata.RelationTypes[i].ObjectType, m.Metadata.RelationTypes[i].Name)
		for _, osr := range subs {
			if !lo.Contains(m.Metadata.RelationTypes[i].Unions, osr.Subject) {
				m.Metadata.RelationTypes[i].Unions = append(m.Metadata.RelationTypes[i].Unions, osr.Subject)
			}
		}
	}
}

func (m *Migrator) consolidatePermissions() {
	m.PermissionMap = NewObjPermRelContainer()

	for i := 0; i < len(m.Metadata.RelationTypes); i++ {
		for _, pn := range m.Metadata.RelationTypes[i].Permissions {
			m.PermissionMap.Add(&ObjPermRel{
				Object:     m.Metadata.RelationTypes[i].ObjectType,
				Permission: pn,
				Relation:   m.Metadata.RelationTypes[i].Name,
			})
		}
	}
}

type WriterOption func(*WriterArgs)

type WriterArgs struct {
	header      bool
	filename    string
	description string
	timestamp   bool
}

// Filename label.
func WithFilename(filename string) WriterOption {
	return func(a *WriterArgs) {
		a.filename = filename
	}
}

// Description label.
func WithDescription(description string) WriterOption {
	return func(a *WriterArgs) {
		a.description = description
	}
}

// Include schema header.
func WithHeader(header bool) WriterOption {
	return func(a *WriterArgs) {
		a.header = header
	}
}

// Include timestamp.
func WithTimestamp(timestamp bool) WriterOption {
	return func(a *WriterArgs) {
		a.timestamp = timestamp
	}
}
