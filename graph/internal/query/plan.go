package query

import (
	"github.com/aserto-dev/azm/model"
)

type Operator int

const (
	Union Operator = iota
	Intersection
	Difference
)

type Expression interface {
	isExpression()
}

type Set struct {
	OT  model.ObjectName
	RT  model.RelationName
	ST  model.ObjectName
	SRT model.RelationName
}

func (s Set) isExpression() {}

type ComputedSet struct {
	Set
	Expansion model.RelationName
}

func (cs ComputedSet) isExpression() {}

// Function call.
type Call struct {
	Signature Set
	Param     Expression
}

func (c Call) isExpression() {}

// Composite applies an operator to a set of expressions.
type Composite struct {
	Operator Operator
	Operands []Expression
}

func (c Composite) isExpression() {}

type Plan struct {
	Expression Expression
	Functions  map[Set]Expression
}

// func BuildQueryPlan(m *model.Model, qry *RelationType) Plan {
// 	in := ds.NewStack[*RelationType]()
// 	//nolint:gocritic
// 	// out := ds.NewStack[Term]()
//
// 	in.Push(qry)
//
// 	for !in.IsEmpty() {
// 		cur := in.Pop()
//
// 		ot := m.Objects[cur.OT]
// 		if ot.HasRelation(cur.RT) {
// 			rt := ot.Relations[cur.RT]
// 			steps := m.StepRelation(rt, cur.ST)
// 			if len(steps) == 0 {
// 				panic("todo")
// 			}
// 		} else {
// 			continue
// 		}
// 	}
//
// 	return nil
// }
