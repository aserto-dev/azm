package parser_test

import (
	"testing"

	"github.com/aserto-dev/azm/parser"
	"github.com/aserto-dev/azm/types"
	"github.com/stretchr/testify/assert"
)

type ObjectName = types.ObjectName
type RelationName = types.RelationName
type RelationRef = types.RelationRef
type Permission = types.Permission

func TestRelationParser(t *testing.T) {
	for _, test := range relationTests {
		t.Run(test.input, func(tt *testing.T) {
			rel := parser.ParseRelation(test.input)
			test.validate(rel, assert.New(tt))
		})
	}
}

func TestPermissionParser(t *testing.T) {
	for _, test := range permissionTests {
		t.Run(test.input, func(tt *testing.T) {
			perm := parser.ParsePermission(test.input)
			test.validate(perm, assert.New(tt))
		})
	}

}

type relationTest struct {
	input    string
	validate func([]*RelationRef, *assert.Assertions)
}

type permissionTest struct {
	input    string
	validate func(*Permission, *assert.Assertions)
}

var relationTests = []relationTest{
	{
		"user",
		func(rel []*RelationRef, assert *assert.Assertions) {
			assert.Len(rel, 1)
			term := rel[0]
			assert.True(term.IsDirect())
			assert.Equal(ObjectName("user"), term.Object)
			assert.Empty(term.Relation)
		},
	},
	{
		"name-with-dashes",
		func(rel []*RelationRef, assert *assert.Assertions) {
			assert.Len(rel, 1)
			term := rel[0]
			assert.True(term.IsDirect())
			assert.Equal(ObjectName("name-with-dashes"), term.Object)
			assert.Empty(term.Relation)
		},
	},
	{
		"group#member",
		func(rel []*RelationRef, assert *assert.Assertions) {
			assert.Len(rel, 1)
			term := rel[0]
			assert.True(term.IsSubject())
			assert.Equal(ObjectName("group"), term.Object)
			assert.Equal(RelationName("member"), term.Relation)
		},
	},
	{
		"user:*",
		func(rel []*RelationRef, assert *assert.Assertions) {
			assert.Len(rel, 1)
			term := rel[0]
			assert.True(term.IsWildcard())
			assert.Equal(ObjectName("user"), term.Object)
			assert.Equal(RelationName("*"), term.Relation)
		},
	},
	{
		"user | group",
		func(rel []*RelationRef, assert *assert.Assertions) {
			assert.Len(rel, 2)

			assert.True(rel[0].IsDirect())
			assert.Equal(ObjectName("user"), rel[0].Object)
			assert.Empty(rel[0].Relation)

			assert.True(rel[1].IsDirect())
			assert.Equal(ObjectName("group"), rel[1].Object)
			assert.Empty(rel[1].Relation)
		},
	},
	{
		"user | group | user:* | group#member",
		func(rel []*RelationRef, assert *assert.Assertions) {
			assert.Len(rel, 4)

			assert.True(rel[0].IsDirect())
			assert.Equal(ObjectName("user"), rel[0].Object)
			assert.Empty(rel[0].Relation)

			assert.True(rel[0].IsDirect())
			assert.Equal(ObjectName("group"), rel[1].Object)
			assert.Empty(rel[1].Relation)

			assert.True(rel[2].IsWildcard())
			assert.Equal(ObjectName("user"), rel[2].Object)

			assert.True(rel[3].IsSubject())
			assert.Equal(ObjectName("group"), rel[3].Object)
			assert.Equal(RelationName("member"), rel[3].Relation)
		},
	},
}

var permissionTests = []permissionTest{
	{
		"can_write",
		func(perm *Permission, assert *assert.Assertions) {
			assert.Equal(RelationName("can_write"), perm.Union[0].RelOrPerm)
			assert.Empty(perm.Union[0].Base)
			assert.Empty(perm.Intersection)
			assert.Nil(perm.Exclusion)
		},
	},
	{
		"can_write | parent->can_read",
		func(perm *Permission, assert *assert.Assertions) {
			assert.Len(perm.Union, 2)
			assert.Equal(RelationName("can_write"), perm.Union[0].RelOrPerm)
			assert.Empty(perm.Union[0].Base)
			assert.Equal(RelationName("parent"), perm.Union[1].Base)
			assert.Equal(RelationName("can_read"), perm.Union[1].RelOrPerm)
		},
	},
	{
		"can_write & can_read",
		func(perm *Permission, assert *assert.Assertions) {
			assert.Len(perm.Intersection, 2)
			assert.Equal(RelationName("can_write"), perm.Intersection[0].RelOrPerm)
			assert.Empty(perm.Intersection[0].Base)
			assert.Equal(RelationName("can_read"), perm.Intersection[1].RelOrPerm)
			assert.Empty(perm.Intersection[1].Base)
		},
	},
	{
		"can_write - can_read",
		func(perm *Permission, assert *assert.Assertions) {
			assert.Equal(RelationName("can_write"), perm.Exclusion.Include.RelOrPerm)
			assert.Empty(perm.Exclusion.Include.Base)
			assert.Equal(RelationName("can_read"), perm.Exclusion.Exclude.RelOrPerm)
			assert.Empty(perm.Exclusion.Exclude.Base)
		},
	},
	{
		"parent->can_read - parent->can_write",
		func(perm *Permission, assert *assert.Assertions) {
			assert.Equal(RelationName("parent"), perm.Exclusion.Include.Base)
			assert.Equal(RelationName("can_read"), perm.Exclusion.Include.RelOrPerm)
			assert.Equal(RelationName("parent"), perm.Exclusion.Exclude.Base)
			assert.Equal(RelationName("can_write"), perm.Exclusion.Exclude.RelOrPerm)
		},
	},
}
