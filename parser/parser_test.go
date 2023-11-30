package parser_test

import (
	"testing"

	"github.com/antlr4-go/antlr/v4"
	"github.com/aserto-dev/azm/model"
	"github.com/aserto-dev/azm/parser"
	"github.com/stretchr/testify/assert"
)

func TestRelationParser(t *testing.T) {
	tests := []struct {
		input    string
		validate func([]*model.Relation, *assert.Assertions)
	}{
		{
			"user",
			func(rel []*model.Relation, assert *assert.Assertions) {
				assert.Len(rel, 1)
				term := rel[0]
				assert.Equal(model.ObjectName("user"), term.Direct)
				assert.Nil(term.Subject)
				assert.Empty(term.Wildcard)
				assert.Nil(term.Arrow)
			},
		},
		{
			"group#member",
			func(rel []*model.Relation, assert *assert.Assertions) {
				assert.Len(rel, 1)
				term := rel[0]
				assert.Equal(model.ObjectName("group"), term.Subject.Object)
				assert.Equal(model.RelationName("member"), term.Subject.Relation)
				assert.Empty(term.Direct)
				assert.Empty(term.Wildcard)
				assert.Nil(term.Arrow)
			},
		},
		{
			"user:*",
			func(rel []*model.Relation, assert *assert.Assertions) {
				assert.Len(rel, 1)
				term := rel[0]
				assert.Equal(model.ObjectName("user"), term.Wildcard)
				assert.Nil(term.Subject)
				assert.Empty(term.Direct)
				assert.Nil(term.Arrow)
			},
		},
		{
			"parent->can_read",
			func(rel []*model.Relation, assert *assert.Assertions) {
				assert.Len(rel, 1)
				term := rel[0]
				assert.Equal(model.RelationName("parent"), term.Arrow.Base)
				assert.Equal("can_read", term.Arrow.Relation)
				assert.Nil(term.Subject)
				assert.Empty(term.Direct)
				assert.Empty(term.Wildcard)
			},
		},
		{
			"user | group",
			func(rel []*model.Relation, assert *assert.Assertions) {
				assert.Len(rel, 2)
				assert.Equal(model.ObjectName("user"), rel[0].Direct)
				assert.Equal(model.ObjectName("group"), rel[1].Direct)
			},
		},
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
			"can_write",
			func(perm *model.Permission, assert *assert.Assertions) {
				assert.Equal("can_write", perm.Union[0].Relation)
				assert.Empty(perm.Union[0].Base)
				assert.Empty(perm.Intersection)
				assert.Nil(perm.Exclusion)
			},
		},
		{
			"can_write | parent->can_read",
			func(perm *model.Permission, assert *assert.Assertions) {
				assert.Len(perm.Union, 2)
				assert.Equal("can_write", perm.Union[0].Relation)
				assert.Empty(perm.Union[0].Base)
				assert.Equal(model.RelationName("parent"), perm.Union[1].Base)
				assert.Equal("can_read", perm.Union[1].Relation)
			},
		},
		{
			"can_write & can_read",
			func(perm *model.Permission, assert *assert.Assertions) {
				assert.Len(perm.Intersection, 2)
				assert.Equal("can_write", perm.Intersection[0].Relation)
				assert.Empty(perm.Intersection[0].Base)
				assert.Equal("can_read", perm.Intersection[1].Relation)
				assert.Empty(perm.Intersection[1].Base)
			},
		},
		{
			"can_write - can_read",
			func(perm *model.Permission, assert *assert.Assertions) {
				assert.Equal("can_write", perm.Exclusion.Base.Relation)
				assert.Empty(perm.Exclusion.Base.Base)
				assert.Equal("can_read", perm.Exclusion.Subtract.Relation)
				assert.Empty(perm.Exclusion.Subtract.Base)
			},
		},
		{
			"parent->can_read - parent->can_write",
			func(perm *model.Permission, assert *assert.Assertions) {
				assert.Equal(model.RelationName("parent"), perm.Exclusion.Base.Base)
				assert.Equal("can_read", perm.Exclusion.Base.Relation)
				assert.Equal(model.RelationName("parent"), perm.Exclusion.Subtract.Base)
				assert.Equal("can_write", perm.Exclusion.Subtract.Relation)
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

	var v parser.RelationVisitor
	return v.Visit(rTree).([]*model.Relation)
}

func parsePermission(input string) *model.Permission {
	p := newParser(input)
	pTree := p.Permission()

	var v parser.PermissionVisitor
	return v.Visit(pTree).(*model.Permission)
}

func newParser(input string) *parser.AzmParser {
	lexer := parser.NewAzmLexer(antlr.NewInputStream(input))
	stream := antlr.NewCommonTokenStream(lexer, 0)
	p := parser.NewAzmParser(stream)
	p.AddErrorListener(antlr.NewDiagnosticErrorListener(true))
	return p
}
