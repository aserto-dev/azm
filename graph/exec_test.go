package graph_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/aserto-dev/azm/graph/internal/query"
	"github.com/aserto-dev/azm/mempool"
	"github.com/aserto-dev/azm/model"
)

// Manifest
//
// model:
//   version: 3
//
// types:
//   user:
//
//   group:
//     relations:
//       member: user | group#member | team#mate
//
//   team:
//     relations:
//       mate: user | team#mate | group#member
//
//   folder:
//     relations:
//       editor: user | group#member
//     permissions:
//       can_edit: editor
//
//   doc:
//     relations:
//       parent: folder
//       owner: user
//       editor: user | group#member
//     permissions:
//       can_delete: owner
//       can_edit: can_delete | editor | parent->can_edit
//       can_share: owner & editor
//       only_editor: editor - owner

var execRels = NewRelationsReader(
	"doc:doc1#parent@folder:folder1",
	"doc:doc1#owner@user:doc1owner",
	"doc:doc1#editor@user:doc1owner",
	"doc:doc1#editor@user:doc1editor",
	"doc:doc1#editor@group:doc1owners#member",
	"group:doc1owners#member@user:doc1ownersmember",
	"folder:folder1#editor@user:folder1editor",
	"folder:folder1#editor@group:folder1editors#member",
	"group:folder1editors#member@user:folder1editorsmember",
)

var evalTests = []checkTest{
	{"doc:doc1#owner@user:doc1owner", true},
	{"doc:doc1#owner@user:doc1editor", false},
}

type (
	on = model.ObjectName
	rn = model.RelationName
)

func TestExecSet(t *testing.T) {
	assert := assert.New(t)
	pool := mempool.NewRelationsPool()

	plan := &query.Plan{
		Expression: set("doc", "owner", "user"),
	}

	for _, test := range evalTests {
		t.Run(test.check, func(tt *testing.T) {
			interpreter := query.NewInterpreter(plan, execRels.GetRelations, pool)
			result, err := interpreter.Run(checkReq(test.check, false))

			assert.NoError(err)
			assert.Equal(test.expected, !result.IsEmpty(), test.check)
		})
	}
}

var Functions = map[query.Set]query.Expression{
	*set("group", "member", "user"): &query.Composite{
		Operator: query.Union,
		Operands: []query.Expression{
			set("group", "member", "user"),
			&query.Pipe{From: set("group", "member", "group", "member"), To: &query.Call{Signature: set("group", "member", "user")}},
			&query.Pipe{From: set("group", "member", "team", "mate"), To: &query.Call{Signature: set("team", "mate", "user")}},
		},
	},

	*set("team", "mate", "user"): &query.Composite{
		Operator: query.Union,
		Operands: []query.Expression{
			set("team", "mate", "user"),
			&query.Pipe{From: set("team", "mate", "team", "mate"), To: &query.Call{Signature: set("team", "mate", "user")}},
			&query.Pipe{From: set("team", "mate", "group", "user"), To: &query.Call{Signature: set("group", "member", "user")}},
		},
	},
}

var unionTests = []checkTest{
	{"doc:doc1#can_edit@user:doc1ownersmember", true},
	{"doc:doc1#can_edit@user:doc1owner", true},
	{"doc:doc1#can_edit@user:doc1editor", true},
	{"doc:doc2#can_edit@user:doc1owner", false},
}

func TestExecUnion(t *testing.T) {
	assert := assert.New(t)
	pool := mempool.NewRelationsPool()

	plan := &query.Plan{
		Expression: &query.Composite{
			Operator: query.Union,
			Operands: []query.Expression{
				set("doc", "owner", "user"),
				set("doc", "editor", "user"),
				&query.Pipe{From: set("doc", "editor", "group", "member"), To: &query.Call{Signature: set("group", "member", "user")}},
			},
		},
		Functions: Functions,
	}

	for _, test := range unionTests {
		t.Run(test.check, func(tt *testing.T) {
			interpreter := query.NewInterpreter(plan, execRels.GetRelations, pool)
			result, err := interpreter.Run(checkReq(test.check, false))
			assert.NoError(err)
			assert.Equal(test.expected, !result.IsEmpty(), test.check)
		})
	}
}

var intersectionTests = []checkTest{
	{"doc:doc1#can_share@user:doc1owner", true},
	{"doc:doc1#can_share@user:doc1editor", false},
}

func TestExecIntersection(t *testing.T) {
	assert := assert.New(t)
	pool := mempool.NewRelationsPool()

	plan := &query.Plan{
		Expression: &query.Composite{
			Operator: query.Intersection,
			Operands: []query.Expression{
				set("doc", "owner", "user"),
				set("doc", "editor", "user"),
			},
		},
	}

	for _, test := range intersectionTests {
		t.Run(test.check, func(tt *testing.T) {
			interpreter := query.NewInterpreter(plan, execRels.GetRelations, pool)
			result, err := interpreter.Run(checkReq(test.check, false))
			assert.NoError(err)
			assert.Equal(test.expected, !result.IsEmpty())
		})
	}
}

var negationTests = []checkTest{
	{"doc:doc1#only_editor@user:doc1owner", false},
	{"doc:doc1#only_editor@user:doc1editor", true},
}

func TestExecNegation(t *testing.T) {
	assert := assert.New(t)
	pool := mempool.NewRelationsPool()

	plan := &query.Plan{
		Expression: &query.Composite{
			Operator: query.Difference,
			Operands: []query.Expression{
				set("doc", "editor", "user"),
				set("doc", "owner", "user"),
			},
		},
	}

	for _, test := range negationTests {
		t.Run(test.check, func(tt *testing.T) {
			interpreter := query.NewInterpreter(plan, execRels.GetRelations, pool)
			result, err := interpreter.Run(checkReq(test.check, false))
			assert.NoError(err)
			assert.Equal(test.expected, !result.IsEmpty())
		})
	}
}

var arrowTests = []checkTest{
	{"doc:doc1#can_edit@user:folder1editor", true},
	{"doc:doc1#can_edit@user:folder1editorsmember", true},
}

func TestExecArrow(t *testing.T) {
	assert := assert.New(t)
	pool := mempool.NewRelationsPool()

	plan := &query.Plan{
		Expression: &query.Composite{
			Operator: query.Union,
			Operands: []query.Expression{
				set("doc", "owner", "user"),
				set("doc", "editor", "user"),
				&query.Pipe{
					From: set("doc", "parent", "folder"),
					To: &query.Composite{
						Operator: query.Union,
						Operands: []query.Expression{
							set("folder", "editor", "user"),
							&query.Pipe{
								From: set("folder", "editor", "group", "member"),
								To:   &query.Call{Signature: set("group", "member", "user")},
							},
						},
					},
				},
			},
		},
		Functions: Functions,
	}
	for _, test := range arrowTests {
		t.Run(test.check, func(tt *testing.T) {
			interpreter := query.NewInterpreter(plan, execRels.GetRelations, pool)
			result, err := interpreter.Run(checkReq(test.check, false))
			assert.NoError(err)
			assert.Equal(test.expected, !result.IsEmpty(), test.check)
		})
	}
}

func set(ot, rt, st string, srt ...string) *query.Set {
	var sr model.RelationName
	switch len(srt) {
	case 0:
		break
	case 1:
		sr = rn(srt[0])
	default:
		panic("only one subject relation type allowed")
	}

	return &query.Set{
		OT:  on(ot),
		RT:  rn(rt),
		ST:  on(st),
		SRT: sr,
	}
}
