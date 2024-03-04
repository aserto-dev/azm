package graph_test

import (
	"strings"
	"testing"

	"github.com/aserto-dev/azm/graph"
	"github.com/aserto-dev/azm/model"
	v3 "github.com/aserto-dev/azm/v3"
	dsc "github.com/aserto-dev/go-directory/aserto/directory/common/v3"
	"github.com/samber/lo"
	req "github.com/stretchr/testify/require"
)

func TestReverseSearch(t *testing.T) {
	require := req.New(t)

	m, err := v3.Load(strings.NewReader(folderModel))
	require.NoError(err)
	require.NotNil(m)

	rm, err := v3.Load(strings.NewReader(inverseModel))
	require.NoError(err)
	require.NotNil(rm)

	rels := NewRelationsReader(
		"folder:leaf#parent@folder:branch",
		"folder:branch#parent@folder:root",
		"folder:root#owner@group:admins",
		"doc:leaf_doc#parent@folder:leaf",
		"doc:root_doc#parent@folder:root",
		"group:admins#member@user:admin",
		"group:admins#member@group:superusers#member",
		"group:superusers#member@user:su",
	)

	for _, test := range reverseSearchTests {
		t.Run(test.search, func(tt *testing.T) {
			assert := req.New(tt)

			s := graph.NewSubjectSearch(rm, invertedGraphReq(test.search), reverseLookup(rm, rels.GetRelations))

			res, err := s.Search()
			assert.NoError(err)

			tt.Logf("explanation: +%v\n", res.Explanation)
			tt.Logf("trace: +%v\n", res.Trace)

			objects := lo.Map(res.Results, func(s *dsc.ObjectIdentifier, _ int) object {
				return object{
					Type: model.ObjectName(s.ObjectType),
					ID:   model.ObjectID(s.ObjectId),
				}
			})

			for _, e := range test.expected {
				assert.Contains(objects, e)
			}

			assert.Equal(len(test.expected), len(objects), objects)
		})
	}

}

var reverseSearchTests = []searchTest{
	{"folder:?#is_owner@user:admin", []object{{"folder", "leaf"}, {"folder", "branch"}, {"folder", "root"}}},
	{"doc:?#is_owner@user:admin", []object{{"doc", "leaf_doc"}, {"doc", "root_doc"}}},
	{"folder:?#is_owner@user:su", []object{{"folder", "leaf"}, {"folder", "branch"}, {"folder", "root"}}},
	{"folder:?#is_owner@group:superusers#member", []object{{"folder", "leaf"}, {"folder", "branch"}, {"folder", "root"}}},
}

const folderModel = `
model:
  version: 3

types:
  user: {}

  group:
    relations:
      member: user | group#member

  folder:
    relations:
      parent: folder
      owner: user | group#member
    permissions:
      is_owner: owner | parent->is_owner

  doc:
    relations:
      parent: folder
      owner: user | group#member
    permissions:
      is_owner: owner | parent->is_owner
`

const inverseModel = `
model:
  version: 3

types:
  doc: {}

  folder:
    relations:
      doc_parent: doc | folder#doc_parent
      folder_parent: folder | folder#folder_parent

  group:
    relations:
      group_member: user | group#group_member
      folder_owner: folder
      doc_owner: doc
    permissions:
      folder_is_owner: folder_owner | folder_owner->folder_parent | group_member->folder_is_owner
      doc_is_owner: doc_owner | folder_owner->doc_parent | group_member->doc_is_owner

  user:
    relations:
      group_member: group
      folder_owner: folder | group#folder_owner
      doc_owner: doc
    permissions:
      folder_is_owner: folder_owner | folder_owner->folder_parent | group_member->folder_is_owner
      doc_is_owner: doc_owner | folder_owner->doc_parent | group_member->doc_is_owner
`

func reverseLookup(_ *model.Model, reader graph.RelationReader) graph.RelationReader {
	return func(r *dsc.Relation) ([]*dsc.Relation, error) {
		x := strings.SplitN(r.Relation, "_", 2)

		rr := &dsc.Relation{
			ObjectType:  r.SubjectType,
			ObjectId:    r.SubjectId,
			Relation:    x[1],
			SubjectType: r.ObjectType,
			SubjectId:   r.ObjectId,
		}

		res, err := reader(rr)
		if err != nil {
			return nil, err
		}

		return lo.Map(res, func(r *dsc.Relation, _ int) *dsc.Relation {
			return &dsc.Relation{
				ObjectType:  r.SubjectType,
				ObjectId:    r.SubjectId,
				Relation:    r.Relation,
				SubjectType: r.ObjectType,
				SubjectId:   r.ObjectId,
			}
		}), nil
	}
}
