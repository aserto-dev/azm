package parser_test

import (
	"testing"

	"github.com/antlr4-go/antlr/v4"
	"github.com/aserto-dev/azm/model"
	"github.com/aserto-dev/azm/parser"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

type RelationVisitor struct {
	parser.BaseAzmVisitor
}

func (v *RelationVisitor) Visit(tree antlr.ParseTree) interface{} {
	switch t := tree.(type) {
	case *parser.RelationContext:
		return t.Accept(v)
	default:
		panic("RelationVisitor can only visit relations")
	}
}

func (v *RelationVisitor) VisitRelation(c *parser.RelationContext) interface{} {
	return lo.Map(c.AllRel(), func(rel parser.IRelContext, _ int) *model.Relation {
		return rel.Accept(v).(*model.Relation)
	})
}

func (v *RelationVisitor) VisitSingleRel(c *parser.SingleRelContext) interface{} {
	return &model.Relation{Direct: model.ObjectName(c.Single().ID().GetText())}
}

func (v *RelationVisitor) VisitWildcardRel(c *parser.WildcardRelContext) interface{} {
	return &model.Relation{Wildcard: model.ObjectName(c.Wildcard().ID().GetText())}
}

func (v *RelationVisitor) VisitSubjectRel(c *parser.SubjectRelContext) interface{} {
	return &model.Relation{Subject: &model.SubjectRelation{
		Object:   model.ObjectName(c.Subject().ID(0).GetText()),
		Relation: model.RelationName(c.Subject().ID(1).GetText()),
	}}
}

type PermissionVisitor struct {
	parser.BaseAzmVisitor
}

func (v *PermissionVisitor) Visit(tree antlr.ParseTree) interface{} {
	switch t := tree.(type) {
	case *parser.UnionPermContext, *parser.IntersectionPermContext, *parser.ExclusionPermContext:
		return t.Accept(v)
	default:
		panic("PermissionVisitor can only visit permissions")
	}
}

func (v *PermissionVisitor) VisitUnionPerm(c *parser.UnionPermContext) interface{} {
	return &model.Permission{
		Union: lo.Map(c.Union().AllPerm(), func(perm parser.IPermContext, _ int) string {
			return perm.Accept(v).(string)
		}),
	}
}

func (v *PermissionVisitor) VisitIntersectionPerm(c *parser.IntersectionPermContext) interface{} {
	return &model.Permission{
		Intersection: lo.Map(c.Intersection().AllPerm(), func(perm parser.IPermContext, _ int) string {
			return perm.Accept(v).(string)
		}),
	}
}

func (v *PermissionVisitor) VisitExclusionPerm(c *parser.ExclusionPermContext) interface{} {
	return &model.Permission{
		Exclusion: &model.ExclusionPermission{
			Base:     c.Exclusion().Perm(0).Accept(v).(string),
			Subtract: c.Exclusion().Perm(1).Accept(v).(string),
		},
	}
}

func (v *PermissionVisitor) VisitSinglePerm(c *parser.SinglePermContext) interface{} {
	return c.Single().ID().GetText()
}

func (v *PermissionVisitor) VisitArrowPerm(c *parser.ArrowPermContext) interface{} {
	return c.Arrow().GetText()
}

func TestRelationParser(t *testing.T) {
	tests := []struct {
		input    string
		validate func([]*model.Relation, *assert.Assertions)
	}{
		{
			"user | group | user:* | group#member",
			func(rel []*model.Relation, assert *assert.Assertions) {
				assert.Len(rel, 4)
				assert.Equal(model.ObjectName("user"), rel[0].Direct)
				assert.Equal(model.ObjectName("group"), rel[1].Direct)
				assert.Equal(model.ObjectName("user"), rel[2].Wildcard)
				assert.Equal(model.ObjectName("group"), rel[3].Subject.Object)
				assert.Equal(model.RelationName("member"), rel[3].Subject.Relation)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.input, func(tt *testing.T) {
			rel := parseRelation(test.input)
			test.validate(rel, assert.New(tt))
		})
	}
}

func TestPermissionParser(t *testing.T) {
	tests := []struct {
		input    string
		validate func(*model.Permission, *assert.Assertions)
	}{
		{
			"can_write | parent->can_read",
			func(perm *model.Permission, assert *assert.Assertions) {
				assert.Len(perm.Union, 2)
				assert.Equal("can_write", perm.Union[0])
				assert.Equal("parent->can_read", perm.Union[1])
			},
		},
		{
			"can_write & can_read",
			func(perm *model.Permission, assert *assert.Assertions) {
				assert.Len(perm.Intersection, 2)
				assert.Equal("can_write", perm.Intersection[0])
				assert.Equal("can_read", perm.Intersection[1])
			},
		},
		{
			"can_write - can_read",
			func(perm *model.Permission, assert *assert.Assertions) {
				assert.Equal("can_write", perm.Exclusion.Base)
				assert.Equal("can_read", perm.Exclusion.Subtract)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.input, func(tt *testing.T) {
			perm := parsePermission(test.input)
			test.validate(perm, assert.New(tt))
		})
	}

}

func parseRelation(input string) []*model.Relation {
	p := newParser(input)
	rTree := p.Relation()

	var v RelationVisitor
	return v.Visit(rTree).([]*model.Relation)
}

func parsePermission(input string) *model.Permission {
	p := newParser(input)
	pTree := p.Permission()

	var v PermissionVisitor
	return v.Visit(pTree).(*model.Permission)
}

func newParser(input string) *parser.AzmParser {
	lexer := parser.NewAzmLexer(antlr.NewInputStream(input))
	stream := antlr.NewCommonTokenStream(lexer, 0)
	p := parser.NewAzmParser(stream)
	p.AddErrorListener(antlr.NewDiagnosticErrorListener(true))
	return p
}
