package graph

import (
	"github.com/aserto-dev/azm/internal/query"
	"github.com/aserto-dev/azm/mempool"
	dsr "github.com/aserto-dev/go-directory/aserto/directory/reader/v3"
)

type PlannedCheck struct {
	interpreter *query.Interpreter
}

func NewPlannedCheck(plan *query.Plan, reader RelationReader, relPool *mempool.RelationsPool, setPool *query.ObjSetPool) *PlannedCheck {
	return &PlannedCheck{interpreter: query.NewInterpreter(plan, reader, relPool, setPool)}
}

func (c *PlannedCheck) Check(req *dsr.CheckRequest) (bool, error) {
	result, err := c.interpreter.Run(req)
	if err != nil {
		return false, err
	}

	return !result.IsEmpty(), nil
}
