package graph_test

import (
	"testing"

	azmgraph "github.com/aserto-dev/azm/graph"
	"github.com/aserto-dev/azm/model"
	v3 "github.com/aserto-dev/azm/v3"
	dsc "github.com/aserto-dev/go-directory/aserto/directory/common/v3"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

func TestSearchSubjects(t *testing.T) {
	rels := relations()

	m, err := v3.LoadFile("./check_test.yaml")
	assert.NoError(t, err)
	assert.NotNil(t, m)

	for _, test := range searchSubjectsTests {
		t.Run(test.search, func(tt *testing.T) {
			assert := assert.New(tt)

			subjSearch := azmgraph.NewSubjectSearch(m, graphReq(test.search), rels.GetRelations)

			res, err := subjSearch.Search()
			assert.NoError(err)
			tt.Logf("explanation: +%v\n", res.Explanation)
			tt.Logf("trace: +%v\n", res.Trace)

			subjects := lo.Map(res.Results, func(s *dsc.ObjectIdentifier, _ int) object {
				return object{
					Type: model.ObjectName(s.ObjectType),
					ID:   model.ObjectID(s.ObjectId),
				}
			})

			for _, e := range test.expected {
				assert.Contains(subjects, e)
			}

			assert.Equal(len(test.expected), len(subjects), subjects)
		})
	}
}

var searchSubjectsTests = []searchTest{
	// Relations
	{"folder:folder1#owner@user:?", []object{{"user", "f1_owner"}}},
	{"folder:folder1#viewer@user:?", []object{{"user", "f1_viewer"}}},
	{"folder:folder2#owner@user:?", []object{}},
	{"folder:folder2#viewer@user:?", []object{}},
	{"group:f1_viewers#member@user:?", []object{{"user", "f1_viewer"}}},
	{"folder:folder1#viewer@group:?#member", []object{{"group", "f1_viewers"}}},
	{"doc:doc1#parent@folder:?", []object{{"folder", "folder1"}}},
	{"folder:folder2#viewer@group:?#member", []object{}},
	{"doc:doc1#owner@user:?", []object{{"user", "d1_owner"}}},
	{"doc:doc1#viewer@group:?#member", []object{{"group", "d1_viewers"}, {"group", "d1_subviewers"}}},
	{"doc:doc1#viewer@user:?", []object{{"user", "user1"}, {"user", "user2"}, {"user", "user3"}, {"user", "f1_owner"}}},
	{"doc:doc2#viewer@user:?", []object{{"user", "*"}, {"user", "user2"}}},
	{"group:d1_viewers#member@group:?#member", []object{{"group", "d1_subviewers"}}},
	{"group:d1_viewers#member@user:?", []object{{"user", "user2"}, {"user", "user3"}}},
	{"group:root#member@user:?", []object{{"user", "leaf_user"}}},
	{"group:root#member@group:?#member", []object{{"group", "leaf"}, {"group", "branch"}, {"group", "trunk"}}},

	// Permissions
	{"folder:folder1#can_create_file@user:?", []object{{"user", "f1_owner"}}},
	{"folder:folder1#can_read@user:?", []object{{"user", "f1_owner"}, {"user", "f1_viewer"}}},
	{"folder:folder1#can_read@group:?#member", []object{{"group", "f1_viewers"}}},
	{"folder:folder1#can_share@user:?", []object{{"user", "f1_owner"}}},
	{"doc:doc1#can_change_owner@user:?", []object{{"user", "d1_owner"}, {"user", "f1_owner"}}},
	{"doc:doc1#can_write@user:?", []object{{"user", "d1_owner"}, {"user", "f1_owner"}}},
	{"doc:doc1#can_read@user:?", []object{{"user", "d1_owner"}, {"user", "f1_owner"}, {"user", "user1"}, {"user", "user2"}, {"user", "user3"}, {"user", "f1_viewer"}}},
	{"doc:doc1#can_read@group:?#member", []object{{"group", "f1_viewers"}, {"group", "d1_viewers"}, {"group", "d1_subviewers"}}},
	{"doc:doc1#can_share@user:?", []object{{"user", "f1_owner"}}},
	{"doc:doc1#can_invite@user:?", []object{{"user", "f1_viewer"}}},
	{"doc:doc1#can_invite@group:?#member", []object{{"group", "f1_viewers"}}},
	{"doc:doc2#can_change_owner@user:?", []object{{"user", "f1_owner"}}},
	{"doc:doc2#can_write@user:?", []object{{"user", "f1_owner"}}},
	{"doc:doc2#can_read@user:?", []object{{"user", "*"}, {"user", "user2"}, {"user", "f1_owner"}, {"user", "f1_viewer"}}},
	{"doc:doc2#can_share@user:?", []object{{"user", "f1_owner"}}},
	{"doc:doc2#can_invite@user:?", []object{{"user", "f1_owner"}, {"user", "f1_viewer"}}},
	{"doc:doc2#can_read@group:?#member", []object{{"group", "f1_viewers"}}},
	{"doc:doc3#can_change_owner@user:?", []object{{"user", "f1_owner"}}},
	{"doc:doc3#can_write@user:?", []object{{"user", "f1_owner"}}},
	{"doc:doc3#can_read@user:?", []object{{"user", "f1_owner"}, {"user", "f1_viewer"}}},
	{"doc:doc3#can_share@user:?", []object{{"user", "f1_owner"}}},
	{"doc:doc3#can_invite@user:?", []object{{"user", "f1_owner"}, {"user", "f1_viewer"}}},
	{"doc:doc3#can_read@group:?#member", []object{{"group", "f1_viewers"}}},
}
