package graph

import (
	"github.com/aserto-dev/azm/graph/internal/query"
	"github.com/aserto-dev/azm/mempool"
	"github.com/aserto-dev/azm/model"
	dsr "github.com/aserto-dev/go-directory/aserto/directory/reader/v3"
)

type PlannedCheck struct {
	m       *model.Model
	getRels RelationReader
	plan    query.Plan
	pool    *mempool.RelationsPool
}

func NewPlannedCheck(m *model.Model, plan query.Plan, reader RelationReader, pool *mempool.RelationsPool) *PlannedCheck {
	return &PlannedCheck{m: m, plan: plan, getRels: reader, pool: pool}
}

func (c *PlannedCheck) Check(req *dsr.CheckRequest) (bool, error) {
	return query.Exec(req, c.m, c.plan, c.getRels, c.pool)
}
