package query

import (
	"fmt"
	"slices"

	"github.com/aserto-dev/azm/internal/ds"
	"github.com/aserto-dev/azm/model"
)

// Operator represents one of the three set operations:
// union, intersection, or difference.
type Operator int

const (
	// Set union.
	Union Operator = iota
	// Set intersection.
	Intersection
	// Set difference.
	Difference
)

type ScopeModifier int

const (
	Unmodified      ScopeModifier = 0x0
	SubjectWildcard ScopeModifier = 0x1
	ObjectWildcard  ScopeModifier = 0x2
)

func (s ScopeModifier) Has(mod ScopeModifier) bool {
	return s&mod != 0
}

// A relation type is identified by its object type, relation type, subject type, and subject relation.
type RelationType struct {
	OT  model.ObjectName
	RT  model.RelationName
	ST  model.ObjectName
	SRT model.RelationName
}

func (r *RelationType) String() string {
	srt := ""
	if r.SRT != "" {
		srt = "#" + r.SRT.String()
	}

	return fmt.Sprintf("%s#%s@%s%s", r.OT, r.RT, r.ST, srt)
}

// Expression is a node in the query-plan's AST.
type Expression interface {
	isExpression()
}

// A Load expression specifies a relation-set to load.
type Load struct {
	*RelationType
	Modifier ScopeModifier
}

func (s *Load) isExpression() {}

// Pipe expressions perform set expansions.
// The results of From are forwarded to To.
type Pipe struct {
	From Expression
	To   Expression
}

func (c *Pipe) isExpression() {}

// Call expressions execute a function.
// Functions aren't named and are identified by their signature.
type Call struct {
	Signature *RelationType
}

func (c *Call) isExpression() {}

// Composite applies an operator to a set of expressions.
type Composite struct {
	Operator Operator
	Operands []Expression
}

func (c *Composite) isExpression() {}

type Functions map[RelationType]Expression

type Plan struct {
	Expression Expression
	Functions  Functions
}

type ExpressionVisitor interface {
	OnLoad(*Load) error
	OnCallStart(*Call) (StepOption, error)
	OnCallEnd(*Call)
	OnCompositeStart(*Composite) (StepOption, error)
	OnCompositeEnd(*Composite)
	OnPipeStart(*Pipe) (StepOption, error)
	OnPipeEnd(*Pipe)
}

type StepOption bool

const (
	StepInto StepOption = false
	StepOver StepOption = true
)

func (p *Plan) Visit(visitor ExpressionVisitor) error {
	backlog := ds.NewStack(p.Expression)

	for !backlog.IsEmpty() {
		switch expr := backlog.Pop().(type) {
		case *Load:
			if err := visitor.OnLoad(expr); err != nil {
				return err
			}

		case *Pipe:
			step, err := visitor.OnPipeStart(expr)
			if err != nil {
				return err
			}

			if step == StepInto {
				backlog.Push(unwind{expr})
				backlog.Push(expr.To)
				backlog.Push(expr.From)
			}

		case *Call:
			step, err := visitor.OnCallStart(expr)
			if err != nil {
				return err
			}

			if step == StepInto {
				backlog.Push(unwind{expr})
				backlog.Push(p.Functions[*expr.Signature])
			}

		case *Composite:
			step, err := visitor.OnCompositeStart(expr)
			if err != nil {
				return err
			}

			if step == StepInto {
				backlog.Push(unwind{expr})
				for _, op := range slices.Backward(expr.Operands) {
					backlog.Push(op)
				}
			}

		case unwind:
			visitUnwind(visitor, expr.expr)
		}
	}

	return nil
}

func visitUnwind(visitor ExpressionVisitor, e Expression) {
	switch e := e.(type) {
	case *Call:
		visitor.OnCallEnd(e)
	case *Pipe:
		visitor.OnPipeEnd(e)
	case *Composite:
		visitor.OnCompositeEnd(e)
	}
}

type unwind struct {
	expr Expression
}

func (u unwind) isExpression() {}
