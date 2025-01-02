package cache

import (
	"github.com/aserto-dev/azm/graph"
	"github.com/aserto-dev/azm/internal/query"
	"github.com/aserto-dev/azm/mempool"
	dsr "github.com/aserto-dev/go-directory/aserto/directory/reader/v3"
	"github.com/aserto-dev/go-directory/pkg/pb"
	"github.com/aserto-dev/go-directory/pkg/prop"
	"google.golang.org/protobuf/types/known/structpb"
)

// If true, use a shared memory pool for all requests.
// Othersise, each call gets its own pool.
const sharedPool = true

func (c *Cache) Check(req *dsr.CheckRequest, relReader graph.RelationReader) (*dsr.CheckResponse, error) {
	checker := graph.NewCheck(c.model.Load(), req, relReader, c.relationsPool())

	ctx := pb.NewStruct()

	ok, err := checker.Check()
	if err != nil {
		ctx.Fields[prop.Reason] = structpb.NewStringValue(err.Error())
	}

	return &dsr.CheckResponse{Check: ok, Trace: checker.Trace(), Context: ctx}, nil
}

func (c *Cache) PlannedCheck(req *dsr.CheckRequest, relReader graph.RelationReader) (*dsr.CheckResponse, error) {
	plan, err := query.Compile(
		c.model.Load(),
		query.NewRelationType(req.ObjectType, req.Relation, req.SubjectType),
		c.queryCache,
	)
	if err != nil {
		return nil, err
	}

	ctx := pb.NewStruct()
	interpreter := query.NewInterpreter(plan, relReader, c.relationsPool())

	result, err := interpreter.Run(req)
	if err != nil {
		ctx.Fields[prop.Reason] = structpb.NewStringValue(err.Error())
	}

	return &dsr.CheckResponse{Check: !result.IsEmpty(), Context: ctx}, nil
}

type graphSearch interface {
	Search() (*dsr.GetGraphResponse, error)
}

func (c *Cache) GetGraph(req *dsr.GetGraphRequest, relReader graph.RelationReader) (*dsr.GetGraphResponse, error) {
	var (
		search graphSearch
		err    error
	)

	if req.ObjectId == "" {
		search, err = graph.NewObjectSearch(c.model.Load(), req, relReader, c.relationsPool())
	} else {
		search, err = graph.NewSubjectSearch(c.model.Load(), req, relReader, c.relationsPool())
	}

	if err != nil {
		return nil, err
	}

	return search.Search()
}

func (c *Cache) relationsPool() *mempool.RelationsPool {
	if sharedPool {
		return c.relsPool
	}

	return mempool.NewRelationsPool()
}
