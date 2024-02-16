package cache

import (
	"github.com/aserto-dev/azm/check"
	dsr "github.com/aserto-dev/go-directory/aserto/directory/reader/v3"
)

func (c *Cache) Check(req *dsr.CheckRequest, relReader check.RelationReader) (*dsr.CheckResponse, error) {
	checker := check.New(c.model, req, relReader)

	ok, err := checker.Check()
	if err != nil {
		return nil, err
	}

	return &dsr.CheckResponse{Check: ok}, nil
}

type graphSearch interface {
	Search() (*dsr.GetGraphResponse, error)
}

func (c *Cache) GetGraph(req *dsr.GetGraphRequest, relReader check.RelationReader) (*dsr.GetGraphResponse, error) {
	var search graphSearch
	if req.ObjectId != "" {
		search = check.NewSubjectSearch(c.model, req, relReader)
	} else {
		search = check.NewObjectSearch(c.model, req, relReader)
	}
	return search.Search()
}
