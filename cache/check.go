package cache

import (
	"github.com/aserto-dev/azm/graph"
	dsr "github.com/aserto-dev/go-directory/aserto/directory/reader/v3"
)

func (c *Cache) Check(req *dsr.CheckRequest, relReader graph.RelationReader) (*dsr.CheckResponse, error) {
	checker := graph.NewCheck(c.model, req, relReader)

	ok, err := checker.Check()
	if err != nil {
		return nil, err
	}

	return &dsr.CheckResponse{Check: ok}, nil
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
		search, err = graph.NewObjectSearch(c.model, req, relReader)
	} else {
		search, err = graph.NewSubjectSearch(c.model, req, relReader)
	}

	if err != nil {
		return nil, err
	}

	return search.Search()
}
