package query

import (
	"github.com/aserto-dev/azm/internal/ds"
	"github.com/aserto-dev/azm/model"
)

type RelationType struct {
	OT  model.ObjectName
	RT  model.RelationName
	ST  model.ObjectName
	SRT model.RelationName
}

type Relation struct {
	RelationType
	ObjectID  model.ObjectID
	SubjectID model.ObjectID
}

type Operator int

const (
	Union Operator = iota
	Intersection
	Negation
)

type Term interface {
	isTerm()
}

func (rt *RelationType) isTerm() {}

func (o Operator) isTerm() {}

type Plan interface {
	isPlan()
}

type Single RelationType

func (s Single) isPlan() {}

type Composite struct {
	Operator Operator
	Operands []Plan
}

func (c Composite) isPlan() {}

func BuildQueryPlan(m *model.Model, qry *RelationType) Plan {
	in := ds.NewStack[*RelationType]()
	//nolint:gocritic
	// out := ds.NewStack[Term]()

	in.Push(qry)

	for !in.IsEmpty() {
		cur := in.Pop()

		ot := m.Objects[cur.OT]
		if ot.HasRelation(cur.RT) {
			rt := ot.Relations[cur.RT]
			steps := m.StepRelation(rt, cur.ST)
			if len(steps) == 0 {
				panic("todo")
			}
		} else {
			continue
		}
	}

	return nil
}
