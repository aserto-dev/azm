package graph_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/aserto-dev/azm/graph/internal/query"
	"github.com/aserto-dev/azm/mempool"
	"github.com/aserto-dev/azm/model"
	v3 "github.com/aserto-dev/azm/v3"
)

const execManifest = `
model:
  version: 3

types:
  user:

  group:
    relations:
      member: user | group#member

  doc:
    relations:
      owner: user
      editor: user | group#member
    permissions:
      can_delete: owner
      can_edit: can_delete | editor
      can_share: owner & editor
      only_editor: editor - owner
`

var execRels = NewRelationsReader(
	"doc:doc1#owner@user:user1",
	"doc:doc1#editor@user:user1",
	"doc:doc1#editor@user:user2",
)

var singleTests = []checkTest{
	{"doc:doc1#owner@user:user1", true},
	{"doc:doc1#owner@user:user2", false},
}

func TestExecSingle(t *testing.T) {
	assert := assert.New(t)
	pool := mempool.NewRelationsPool()

	m, err := v3.Load(strings.NewReader(execManifest))
	assert.NoError(err)

	plan := single("doc", "owner", "user")

	for _, test := range singleTests {
		t.Run(test.check, func(tt *testing.T) {
			result, err := query.Exec(checkReq(test.check, false), m, plan, execRels.GetRelations, pool)
			assert.NoError(err)
			assert.Equal(test.expected, result)
		})
	}
}

var unionTests = []checkTest{
	{"doc:doc1#can_edit@user:user1", true},
	{"doc:doc1#can_edit@user:user2", true},
	{"doc:doc2#can_edit@user:user1", false},
}

func TestExecUnion(t *testing.T) {
	assert := assert.New(t)
	pool := mempool.NewRelationsPool()

	m, err := v3.Load(strings.NewReader(execManifest))
	assert.NoError(err)

	plan := query.Composite{
		Operator: query.Union,
		Operands: []query.Plan{
			single("doc", "owner", "user"),
			single("doc", "editor", "user"),
		},
	}

	for _, test := range unionTests {
		t.Run(test.check, func(tt *testing.T) {
			result, err := query.Exec(checkReq(test.check, false), m, plan, execRels.GetRelations, pool)
			assert.NoError(err)
			assert.Equal(test.expected, result)
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

	m, err := v3.Load(strings.NewReader(execManifest))
	assert.NoError(err)

	plan := query.Composite{
		Operator: query.Intersection,
		Operands: []query.Plan{
			single("doc", "owner", "user"),
			single("doc", "editor", "user"),
		},
	}

	for _, test := range intersectionTests {
		t.Run(test.check, func(tt *testing.T) {
			result, err := query.Exec(checkReq(test.check, false), m, plan, execRels.GetRelations, pool)
			assert.NoError(err)
			assert.Equal(test.expected, result)
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

	m, err := v3.Load(strings.NewReader(execManifest))
	assert.NoError(err)

	plan := query.Composite{
		Operator: query.Intersection,
		Operands: []query.Plan{
			single("doc", "editor", "user"),
			query.Composite{
				Operator: query.Negation,
				Operands: []query.Plan{
					single("doc", "owner", "user"),
				},
			},
		},
	}

	for _, test := range negationTests {
		t.Run(test.check, func(tt *testing.T) {
			result, err := query.Exec(checkReq(test.check, false), m, plan, execRels.GetRelations, pool)
			assert.NoError(err)
			assert.Equal(test.expected, result)
		})
	}
}

func single(ot, rt, st string) query.Plan {
	return query.Single{
		OT: model.ObjectName(ot),
		RT: model.RelationName(rt),
		ST: model.ObjectName(st),
	}
}
