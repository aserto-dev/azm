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
