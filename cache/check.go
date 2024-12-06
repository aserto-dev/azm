package cache

import (
	"github.com/aserto-dev/azm/graph"
	dsr "github.com/aserto-dev/go-directory/aserto/directory/reader/v3"
	"github.com/aserto-dev/go-directory/pkg/pb"
	"github.com/aserto-dev/go-directory/pkg/prop"
	"google.golang.org/protobuf/types/known/structpb"
)

func (c *Cache) Check(req *dsr.CheckRequest, relReader graph.RelationReader) (*dsr.CheckResponse, error) {
	c.mtx.RLock()
	defer c.mtx.RUnlock()

	checker := graph.NewCheck(c.model, req, relReader, c.relsPool)

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
	c.mtx.RLock()
	defer c.mtx.RUnlock()

	var (
		search graphSearch
		err    error
	)

	if req.ObjectId == "" {
		search, err = graph.NewObjectSearch(c.model, req, relReader, c.relsPool)
	} else {
		search, err = graph.NewSubjectSearch(c.model, req, relReader, c.relsPool)
	}

	if err != nil {
		return nil, err
	}

	return search.Search()
}
