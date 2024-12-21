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
//   doc:
//     relations:
//       owner: user
//       editor: user | group#member
//     permissions:
//       can_delete: owner
//       can_edit: can_delete | editor
//       can_share: owner & editor
//       only_editor: editor - owner

var execRels = NewRelationsReader(
	"doc:doc1#owner@user:user1",
	"doc:doc1#editor@user:user1",
	"doc:doc1#editor@user:user2",
	"doc:doc1#editor@group:group1#member",
	"group:group1#member@user:user3",
)

var evalTests = []checkTest{
	{"doc:doc1#owner@user:user1", true},
	{"doc:doc1#owner@user:user2", false},
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

			// result, err := query.Exec(checkReq(test.check, false), plan, execRels.GetRelations, pool)
			assert.NoError(err)
			assert.Equal(test.expected, !result.Empty())
		})
	}
}

var unionTests = []checkTest{
	// {"doc:doc1#can_edit@user:user3", true},
	// {"doc:doc1#can_edit@user:user1", true},
	// {"doc:doc1#can_edit@user:user2", true},
	{"doc:doc2#can_edit@user:user1", false},
}

func TestExecUnion(t *testing.T) {
	assert := assert.New(t)
	pool := mempool.NewRelationsPool()

	groupMember := &query.Composite{
		Operator: query.Union,
		Operands: []query.Expression{
			set("group", "member", "user"),
			&query.Call{Signature: set("group", "member", "user"), Param: computed("group", "member", "group", "member")},
			&query.Call{Signature: set("team", "mate", "user"), Param: computed("group", "member", "team", "mate")},
		},
	}

	teamMate := &query.Composite{
		Operator: query.Union,
		Operands: []query.Expression{
			set("team", "mate", "user"),
			&query.Call{Signature: set("team", "mate", "user"), Param: computed("team", "mate", "team", "mate")},
			&query.Call{Signature: set("group", "member", "user"), Param: computed("team", "mate", "group", "user")},
		},
	}

	plan := &query.Plan{
		Expression: &query.Composite{
			Operator: query.Union,
			Operands: []query.Expression{
				set("doc", "owner", "user"),
				set("doc", "editor", "user"),
				&query.Call{Signature: set("group", "member", "user"), Param: computed("doc", "editor", "group", "member")},
			},
		},
		Functions: map[query.Set]query.Expression{
			*set("group", "member", "user"): groupMember,
			*set("team", "mate", "user"):    teamMate,
		},
	}

	for _, test := range unionTests {
		t.Run(test.check, func(tt *testing.T) {
			interpreter := query.NewInterpreter(plan, execRels.GetRelations, pool)
			result, err := interpreter.Run(checkReq(test.check, false))
			assert.NoError(err)
			assert.Equal(test.expected, !result.Empty())
		})
	}
}

var intersectionTests = []checkTest{
	{"doc:doc1#can_share@user:user1", true},
	{"doc:doc1#can_share@user:user2", false},
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
			assert.Equal(test.expected, !result.Empty())
		})
	}
}

var negationTests = []checkTest{
	{"doc:doc1#only_editor@user:user1", false},
	{"doc:doc1#only_editor@user:user2", true},
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
			assert.Equal(test.expected, !result.Empty())
		})
	}
}

func set(ot, rt, st string) *query.Set {
	return computed(ot, rt, st)
}

func computed(ot, rt, st string, srt ...string) *query.Set {
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
