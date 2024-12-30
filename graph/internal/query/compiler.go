package query

import (
	"github.com/aserto-dev/azm/model"
	"github.com/aserto-dev/go-directory/pkg/derr"
	"github.com/pkg/errors"
)

var (
	CompileError = errors.New("compile error")
)

type compiler struct {
	m     *model.Model
	funcs Functions
}

func Compile(m *model.Model, set *Load) (*Plan, error) {
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

func (c *compiler) compile(set *Load) (Expression, error) {
	obj := c.m.Objects[set.OT]
	if obj == nil {
		return nil, derr.ErrObjectTypeNotFound.Msg(set.OT.String())
	}

	switch {
	case obj.HasRelation(set.RT):
		return c.compileRelation(obj.Relations[set.RT], set)

	case obj.HasPermission(set.RT):
		return c.compilePermission(obj.Permissions[set.RT], set)
	}

	return nil, derr.ErrRelationNotFound.Msg(set.RT.String())
}

func (c *compiler) compileRelation(r *model.Relation, set *Load) (Expression, error) {
	// steps := c.m.StepRelation(r, set.ST)
	// for _, step := range steps {
	// }

	return nil, nil
}

func (c *compiler) compilePermission(p *model.Permission, set *Load) (Expression, error) {
	return nil, nil
}
