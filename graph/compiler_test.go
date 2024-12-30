package graph_test

import (
	"strings"
	"testing"

	"github.com/aserto-dev/azm/graph/internal/query"
	"github.com/aserto-dev/azm/mempool"
	"github.com/aserto-dev/azm/model"
	v3 "github.com/aserto-dev/azm/v3"
	"github.com/stretchr/testify/assert"
)

const testManifest = `
model:
  version: 3

types:
  user:

  group:
    relations:
      member: user | group#member | team#mate

  team:
    relations:
      mate: user | team#mate | group#member

  folder:
    relations:
      editor: user | group#member
    permissions:
      can_edit: editor

  doc:
    relations:
      parent: folder
      owner: user
      editor: user | user:* | group#member
    permissions:
      can_delete: owner
      can_edit: can_delete | editor | parent->can_edit
      can_share: owner & editor
      only_editor: editor - owner
`

var relationTests = []checkTest{
	{"doc:doc1#owner@user:doc1owner", true},
	{"doc:doc1#owner@user:doc1editor", false},
}

func TestCompileRelation(t *testing.T) {
	assert := assert.New(t)
	m, err := loadManifest(testManifest)
	assert.NoError(err)

	pool := mempool.NewRelationsPool()

	plan, err := query.Compile(m, rel("doc", "owner", "user"))
	assert.NoError(err)
	assert.NotNil(plan)

	for _, test := range relationTests {
		t.Run(test.check, func(tt *testing.T) {
			interpreter := query.NewInterpreter(plan, execRels.GetRelations, pool)
			result, err := interpreter.Run(checkReq(test.check, false))

			assert.NoError(err)
			assert.Equal(test.expected, !result.IsEmpty(), test.check)
		})
	}
}

func loadManifest(manifest string) (*model.Model, error) {
	r := strings.NewReader(manifest)
	return v3.Load(r)
}

func rel(ot, rt, st string, srt ...string) *query.RelationType {
	sr := ""
	if len(srt) > 0 {
		sr = srt[0]
	}

	return &query.RelationType{
		OT:  on(ot),
		RT:  rn(rt),
		ST:  on(st),
		SRT: rn(sr),
	}
}
