package parser_test

import (
	"testing"

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
			},
		},
		{
			"name-with-dashes",
			func(rel []*model.Relation, assert *assert.Assertions) {
				assert.Len(rel, 1)
				term := rel[0]
				assert.Equal(model.ObjectName("name-with-dashes"), term.Direct)
				assert.Nil(term.Subject)
				assert.Empty(term.Wildcard)
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
			rel := parser.ParseRelation(test.input)
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
				assert.Equal("can_write", perm.Union[0].RelOrPerm)
				assert.Empty(perm.Union[0].Base)
				assert.Empty(perm.Intersection)
				assert.Nil(perm.Exclusion)
			},
		},
		{
			"can_write | parent->can_read",
			func(perm *model.Permission, assert *assert.Assertions) {
				assert.Len(perm.Union, 2)
				assert.Equal("can_write", perm.Union[0].RelOrPerm)
				assert.Empty(perm.Union[0].Base)
				assert.Equal(model.RelationName("parent"), perm.Union[1].Base)
				assert.Equal("can_read", perm.Union[1].RelOrPerm)
			},
		},
		{
			"can_write & can_read",
			func(perm *model.Permission, assert *assert.Assertions) {
				assert.Len(perm.Intersection, 2)
				assert.Equal("can_write", perm.Intersection[0].RelOrPerm)
				assert.Empty(perm.Intersection[0].Base)
				assert.Equal("can_read", perm.Intersection[1].RelOrPerm)
				assert.Empty(perm.Intersection[1].Base)
			},
		},
		{
			"can_write - can_read",
			func(perm *model.Permission, assert *assert.Assertions) {
				assert.Equal("can_write", perm.Exclusion.Include.RelOrPerm)
				assert.Empty(perm.Exclusion.Include.Base)
				assert.Equal("can_read", perm.Exclusion.Exclude.RelOrPerm)
				assert.Empty(perm.Exclusion.Exclude.Base)
			},
		},
		{
			"parent->can_read - parent->can_write",
			func(perm *model.Permission, assert *assert.Assertions) {
				assert.Equal(model.RelationName("parent"), perm.Exclusion.Include.Base)
				assert.Equal("can_read", perm.Exclusion.Include.RelOrPerm)
				assert.Equal(model.RelationName("parent"), perm.Exclusion.Exclude.Base)
				assert.Equal("can_write", perm.Exclusion.Exclude.RelOrPerm)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.input, func(tt *testing.T) {
			perm := parser.ParsePermission(test.input)
			test.validate(perm, assert.New(tt))
		})
	}

}