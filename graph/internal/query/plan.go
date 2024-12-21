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

type VisitOption bool

const (
	StepInto VisitOption = false
	StepOver VisitOption = true
)

type ExpressionVisitor interface {
	OnSet(*Set) error
	OnCallStart(*Call) (VisitOption, error)
	OnCallEnd(*Call)
	OnCompositeStart(*Composite) error
	OnCompositeEnd(*Composite)
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

func (s *Set) isExpression() {}

// Function call.
type Call struct {
	Signature *Set
	Param     Expression
}

func (c *Call) isExpression() {}

// Composite applies an operator to a set of expressions.
type Composite struct {
	Operator Operator
	Operands []Expression
}

func (c *Composite) isExpression() {}

type Functions map[Set]Expression

type Plan struct {
	Expression Expression
	Functions  Functions
}

func (p *Plan) Visit(visitor ExpressionVisitor) error {
	backlog := ds.NewStack(p.Expression)

	for !backlog.IsEmpty() {
		switch e := backlog.Pop().(type) {
		case *Set:
			if err := visitor.OnSet(e); err != nil {
				return err
			}

		case *Call:
			abort, err := visitor.OnCallStart(e)
			if err != nil {
				return err
			}

			if !abort {
				backlog.Push(unwind{e})
				backlog.Push(p.Functions[*e.Signature])
				backlog.Push(e.Param)
			}

		case *Composite:
			if err := visitor.OnCompositeStart(e); err != nil {
				return err
			}

			backlog.Push(unwind{e})
			for _, op := range slices.Backward(e.Operands) {
				backlog.Push(op)
			}

		case unwind:
			visitUnwind(visitor, e.expr)
		}
	}

	return nil
}

func visitUnwind(visitor ExpressionVisitor, e Expression) {
	switch e := e.(type) {
	case *Call:
		visitor.OnCallEnd(e)
	case *Composite:
		visitor.OnCompositeEnd(e)
	}
}

type unwind struct {
	expr Expression
}

func (u unwind) isExpression() {}

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
