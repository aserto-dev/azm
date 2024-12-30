package query

import (
	"github.com/aserto-dev/azm/model"
	"github.com/aserto-dev/go-directory/pkg/derr"
	"github.com/pkg/errors"
	"github.com/samber/lo"
)

var (
	CompileError = errors.New("compile error")
)

// Empty set expressions are used during compilations but are optimized out of the
// final query plan.
type emptySet struct{}

func (emptySet) isExpression() {}

type compiler struct {
	m     *model.Model
	funcs Functions
}

func Compile(m *model.Model, set *RelationType) (*Plan, error) {
	plan := &Plan{}

	c := &compiler{m, Functions{}}

	expr, err := c.compile(set)
	if err != nil {
		return plan, err
	}

	plan.Expression = expr
	plan.Functions = c.funcs

	return plan, nil
}

func (c *compiler) compile(set *RelationType) (Expression, error) {
	if expr, ok := c.funcs[*set]; ok {
		return expr, nil
	}

	obj := c.m.Objects[set.OT]
	if obj == nil {
		return nil, derr.ErrObjectTypeNotFound.Msg(set.OT.String())
	}

	var (
		expr Expression
		err  error
	)
	switch {
	case obj.HasRelation(set.RT):
		expr, err = c.compileRelation(obj.Relations[set.RT], set)

	case obj.HasPermission(set.RT):
		expr, err = c.compilePermission(obj.Permissions[set.RT], set)
	default:
		return nil, derr.ErrRelationNotFound.Msg(set.RT.String())
	}

	if err != nil {
		return nil, err
	}

	c.funcs[*set] = expr
	return expr, nil
}

func (c *compiler) compileRelation(r *model.Relation, set *RelationType) (Expression, error) {
	steps := c.m.StepRelation(r, set.ST)
	ops := make([]Expression, len(steps))

	for i, step := range steps {
		if step.IsSubject() {
			signature := RelationType{
				OT: step.Object,
				RT: step.Relation,
				ST: set.ST,
			}

			if _, err := c.compile(&signature); err != nil {
				return nil, err
			}

			ops[i] = &Pipe{
				From: &Load{
					RelationType: &RelationType{
						OT:  set.OT,
						RT:  set.RT,
						ST:  step.Object,
						SRT: step.Relation,
					},
				},
				To: &Call{
					Signature: &signature,
				},
			}
		} else {
			ops[i] = &Load{
				RelationType: set,
				Modifier:     lo.Ternary(step.IsWildcard(), SubjectWildcard, Unmodified),
			}
		}
	}

	if len(ops) == 1 {
		return ops[0], nil
	}

	return &Composite{Operator: Union, Operands: ops}, nil
}

func (c *compiler) compilePermission(p *model.Permission, set *RelationType) (Expression, error) {
	if !lo.Contains(p.SubjectTypes, set.ST) {
		return emptySet{}, nil
	}

	terms := p.Terms()
	ops := make([]Expression, len(terms))

	// index of the first term that resolves to an empty set or -1 if all terms are non-empty.
	firstEmpty := -1

	for i, term := range terms {
		switch {
		case !lo.Contains(term.SubjectTypes, set.ST):
			ops[i] = emptySet{}

			if firstEmpty == -1 {
				firstEmpty = i
			}

		case term.IsArrow():
			baseRel := c.m.Objects[set.OT].Relations[term.Base]
			baseRelTypes := baseRel.Types()
			paths := make([]Expression, len(baseRelTypes))
			for i, baseType := range baseRelTypes {
				expr, err := c.compile(&RelationType{OT: baseType.Object, RT: term.RelOrPerm, ST: set.ST})
				if err != nil {
					return nil, err
				}

				if isEmptySet(expr) {
					paths[i] = emptySet{}
				} else {
					paths[i] = &Pipe{
						From: &Load{
							RelationType: &RelationType{OT: set.OT, RT: term.Base, ST: baseType.Object, SRT: baseType.Relation},
						},
						To: expr,
					}
				}
			}

		default:
			expr, err := c.compile(&RelationType{OT: set.OT, RT: term.RelOrPerm, ST: set.ST})
			if err != nil {
				return nil, err
			}
			ops[i] = expr
		}
	}

	if firstEmpty != -1 && (p.IsIntersection() || p.IsExclusion() && firstEmpty == 0) {
		// short circuit.
		// An empty set in an intersection or as the first term of an exclusion.
		return emptySet{}, nil
	}

	ops = lo.Filter(ops, func(expr Expression, _ int) bool { return !isEmptySet(expr) })
	if len(ops) == 1 {
		return ops[0], nil
	}

	var op Operator
	switch {
	case p.IsUnion():
		op = Union
	case p.IsIntersection():
		op = Intersection
	case p.IsExclusion():
		op = Difference
	}

	return &Composite{Operator: op, Operands: ops}, nil
}

func isEmptySet(expr Expression) bool {
	_, ok := expr.(emptySet)
	return ok
}
