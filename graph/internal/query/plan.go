package query

import (
	"slices"

	"github.com/aserto-dev/azm/internal/ds"
	"github.com/aserto-dev/azm/model"
)

type Operator int

const (
	Union Operator = iota
	Intersection
	Difference
)

type ExpressionVisitor interface {
	VisitSet(Set) (bool, error)
	VisitCall(Call) (bool, error)
	VisitComposite(Composite) (bool, error)
}

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

type Functions map[Set]Expression

type Plan struct {
	Expression Expression
	Functions  Functions
}

func (p *Plan) Visit(visitor ExpressionVisitor) {
	backlog := ds.NewStack(p.Expression)

	for !backlog.IsEmpty() {
		switch e := backlog.Pop().(type) {
		case Set:
			visitor.VisitSet(e)
		case Call:
			visitor.VisitCall(e)
			backlog.Push(e.Param)
			backlog.Push(p.Functions[e.Signature])
		case Composite:
			visitor.VisitComposite(e)
			for _, op := range slices.Backward(e.Operands) {
				backlog.Push(op)
			}
		}
	}
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
