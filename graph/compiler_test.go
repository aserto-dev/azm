package graph_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aserto-dev/azm/internal/ds"
	"github.com/aserto-dev/azm/internal/query"
	"github.com/aserto-dev/azm/mempool"
	"github.com/aserto-dev/azm/model"
	v3 "github.com/aserto-dev/azm/v3"
	"github.com/stretchr/testify/assert"
	rqur "github.com/stretchr/testify/require"
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

type compileTest struct {
	signature *query.RelationType
	tests     []checkTest
}

var compileTests = []compileTest{
	simpleRelationTests,
	compositeRelationTests,
	simplePermissionTests,
	unionPermissionTests,
	intersectionPermissionTests,
	negationPermissionTests,
}

func TestCompile(t *testing.T) {
	require := rqur.New(t)
	m, err := loadManifest(testManifest)
	require.NoError(err)

	relPool := mempool.NewRelationsPool()
	objSetPool := ds.NewSetPool[model.ObjectID]()

	for _, suite := range compileTests {
		plan, err := query.Compile(m, suite.signature, nil)
		require.NoError(err)
		require.NotNil(plan)

		for _, test := range suite.tests {
			t.Run(fmt.Sprintf("%s---%s", suite.signature, test.check), func(tt *testing.T) {
				interpreter := query.NewInterpreter(plan, execRels.GetRelations, relPool, objSetPool)
				result, err := interpreter.Run(checkReq(test.check, false))

				assert.NoError(tt, err)
				assert.Equal(tt, test.expected, !result.IsEmpty(), test.check)
			})
		}
	}
}

var simpleRelationTests = compileTest{
	signature: rel("doc", "owner", "user"),
	tests: []checkTest{
		{"doc:doc1#owner@user:doc1owner", true},
		{"doc:doc1#owner@user:doc1editor", false},
	},
}

var compositeRelationTests = compileTest{
	signature: rel("doc", "editor", "user"),
	tests: []checkTest{
		{"doc:doc1#editor@user:doc1owner", true},
		{"doc:doc1#editor@user:doc1editor", true},
		{"doc:doc1#editor@user:someuser", false},
		{"doc:doc2#editor@user:someuser", true},
	},
}

var simplePermissionTests = compileTest{
	signature: rel("doc", "can_delete", "user"),
	tests: []checkTest{
		{"doc:doc1#can_delete@user:doc1owner", true},
		{"doc:doc1#can_delete@user:doc1editor", false},
	},
}

var unionPermissionTests = compileTest{
	signature: rel("doc", "can_edit", "user"),
	tests: []checkTest{
		{"doc:doc1#can_delete@user:doc1owner", true},
		{"doc:doc1#editor@user:doc1editor", true},
		{"doc:doc2#editor@user:someuser", true},
		{"doc:doc1#can_edit@user:folder1editor", true},
		{"doc:doc1#can_edit@user:folder1editorsmember", true},
	},
}

var intersectionPermissionTests = compileTest{
	signature: rel("doc", "can_share", "user"),
	tests: []checkTest{
		{"doc:doc1#can_share@user:doc1owner", true},
		{"doc:doc1#can_share@user:doc1editor", false},
	},
}

var negationPermissionTests = compileTest{
	signature: rel("doc", "only_editor", "user"),
	tests: []checkTest{
		{"doc:doc1#only_editor@user:doc1owner", false},
		{"doc:doc1#only_editor@user:doc1editor", true},
	},
}

func loadManifest(manifest string) (*model.Model, error) {
	r := strings.NewReader(manifest)
	return v3.Load(r)
}

func rel(ot, rt, st string, srt ...string) *query.RelationType { //nolint:unparam
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
