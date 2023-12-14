package cache

import (
	"github.com/aserto-dev/azm/walk"
	dsr "github.com/aserto-dev/go-directory/aserto/directory/reader/v3"
)

func (c *Cache) Check(req *dsr.CheckRequest, relReader walk.RelationReader) (*dsr.CheckResponse, error) {
	w := walk.New(c.model, req, relReader)

	ok, err := w.Check()
	if err != nil {
		return nil, err
	}

	return &dsr.CheckResponse{Check: ok}, nil
}
