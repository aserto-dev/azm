package cache

import (
	"github.com/aserto-dev/azm/graph"
	"github.com/aserto-dev/azm/mempool"
	dsr "github.com/aserto-dev/go-directory/aserto/directory/reader/v3"
	"github.com/aserto-dev/go-directory/pkg/pb"
	"github.com/aserto-dev/go-directory/pkg/prop"
	"google.golang.org/protobuf/types/known/structpb"
)

func (c *Cache) Check(req *dsr.CheckRequest, relReader graph.RelationReader) (*dsr.CheckResponse, error) {
	relsPool := mempool.NewRelationsPool()
	checker := graph.NewCheck(c.model, req, relReader, relsPool)

	ctx := pb.NewStruct()

	ok, err := checker.Check()
	if err != nil {
		ctx.Fields[prop.Reason] = structpb.NewStringValue(err.Error())
	}

	return &dsr.CheckResponse{Check: ok, Trace: checker.Trace(), Context: ctx}, nil
}

type graphSearch interface {
	Search() (*dsr.GetGraphResponse, error)
}

func (c *Cache) GetGraph(req *dsr.GetGraphRequest, relReader graph.RelationReader) (*dsr.GetGraphResponse, error) {
	var (
		search graphSearch
		err    error
	)

	relsPool := mempool.NewRelationsPool()

	if req.ObjectId == "" {
		search, err = graph.NewObjectSearch(c.model, req, relReader, relsPool)
	} else {
		search, err = graph.NewSubjectSearch(c.model, req, relReader, relsPool)
	}

	if err != nil {
		return nil, err
	}

	return search.Search()
}
