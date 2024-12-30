package graph_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/aserto-dev/azm/graph/internal/query"
	"github.com/aserto-dev/azm/mempool"
	"github.com/aserto-dev/azm/model"
)

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
	"doc:doc2#editor@user:*",
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
		Expression: load("doc", "owner", "user"),
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

var Functions = map[query.RelationType]query.Expression{
	*load("group", "member", "user").RelationType: &query.Composite{
		Operator: query.Union,
		Operands: []query.Expression{
			load("group", "member", "user"),
			&query.Pipe{From: load("group", "member", "group", "member"), To: &query.Call{Signature: load("group", "member", "user").RelationType}},
			&query.Pipe{From: load("group", "member", "team", "mate"), To: &query.Call{Signature: load("team", "mate", "user").RelationType}},
		},
	},

	*load("team", "mate", "user").RelationType: &query.Composite{
		Operator: query.Union,
		Operands: []query.Expression{
			load("team", "mate", "user"),
			&query.Pipe{From: load("team", "mate", "team", "mate"), To: &query.Call{Signature: load("team", "mate", "user").RelationType}},
			&query.Pipe{From: load("team", "mate", "group", "user"), To: &query.Call{Signature: load("group", "member", "user").RelationType}},
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
				load("doc", "owner", "user"),
				load("doc", "editor", "user"),
				&query.Pipe{From: load("doc", "editor", "group", "member"), To: &query.Call{Signature: load("group", "member", "user").RelationType}},
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
				load("doc", "owner", "user"),
				load("doc", "editor", "user"),
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
				load("doc", "editor", "user"),
				load("doc", "owner", "user"),
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
				load("doc", "owner", "user"),
				load("doc", "editor", "user"),
				&query.Pipe{
					From: load("doc", "parent", "folder"),
					To: &query.Composite{
						Operator: query.Union,
						Operands: []query.Expression{
							load("folder", "editor", "user"),
							&query.Pipe{
								From: load("folder", "editor", "group", "member"),
								To:   &query.Call{Signature: load("group", "member", "user").RelationType},
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

var wildcardTests = []checkTest{
	{"doc:doc2#editor@user:any", true},
}

func TestExecWildcard(t *testing.T) {
	assert := assert.New(t)
	pool := mempool.NewRelationsPool()

	plan := &query.Plan{
		Expression: &query.Composite{
			Operator: query.Union,
			Operands: []query.Expression{
				load("doc", "editor", "user"),
				wildcard("doc", "editor", "user"),
			},
		},
	}

	for _, test := range wildcardTests {
		t.Run(test.check, func(tt *testing.T) {
			interpreter := query.NewInterpreter(plan, execRels.GetRelations, pool)
			result, err := interpreter.Run(checkReq(test.check, false))
			assert.NoError(err)
			assert.Equal(test.expected, !result.IsEmpty(), test.check)
		})
	}
}

func load(ot, rt, st string, srt ...string) *query.Load {
	var sr model.RelationName
	switch len(srt) {
	case 0:
		break
	case 1:
		sr = rn(srt[0])
	default:
		panic("only one subject relation type allowed")
	}

	return &query.Load{
		RelationType: &query.RelationType{
			OT:  on(ot),
			RT:  rn(rt),
			ST:  on(st),
			SRT: sr,
		},
	}
}

func wildcard(ot, rt, st string) *query.Load {
	return &query.Load{
		RelationType: &query.RelationType{
			OT: on(ot),
			RT: rn(rt),
			ST: on(st),
		},
		Modifier: query.SubjectWildcard,
	}
}
