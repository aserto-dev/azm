package graph_test

import (
	"os"
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

	r, err := os.Open("./check_test.yaml")
	assert.NoError(t, err)
	assert.NotNil(t, r)

	m, err := v3.Load(r)
	assert.NoError(t, err)
	assert.NotNil(t, m)

	for _, test := range searchSubjectsTests {
		t.Run(test.name, func(tt *testing.T) {
			assert := assert.New(tt)

			subjSearch := azmgraph.NewSubjectSearch(m, test.search, rels.GetRelations)

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
	{name: "users that are folder1 owners", search: graph("folder", "folder1", "owner", "user", "", ""),
		expected: []object{{Type: "user", ID: "f1_owner"}},
	},
	{name: "users that are folder1 viewers", search: graph("folder", "folder1", "viewer", "user", "", ""),
		expected: []object{{Type: "user", ID: "f1_viewer"}},
	},
	{name: "users that are folder2 owners", search: graph("folder", "folder2", "owner", "user", "", ""),
		expected: []object{},
	},
	{name: "users that are folder2 viewers", search: graph("folder", "folder2", "viewer", "user", "", ""),
		expected: []object{},
	},
	{name: "users that are f1_viewers members", search: graph("group", "f1_viewers", "member", "user", "", ""),
		expected: []object{{Type: "user", ID: "f1_viewer"}},
	},
	{name: "groups where members are folder1 viewers", search: graph("folder", "folder1", "viewer", "group", "", "member"),
		expected: []object{{Type: "group", ID: "f1_viewers"}},
	},
	{name: "folders that are doc1 parents", search: graph("doc", "doc1", "parent", "folder", "", ""),
		expected: []object{{Type: "folder", ID: "folder1"}},
	},
	{name: "groups where members are folder2 viewers", search: graph("folder", "folder2", "viewer", "group", "", "member"),
		expected: []object{},
	},
	{name: "users that are doc1 owners", search: graph("doc", "doc1", "owner", "user", "", ""),
		expected: []object{{Type: "user", ID: "d1_owner"}},
	},
	{name: "groups where members are doc1 viewers", search: graph("doc", "doc1", "viewer", "group", "", "member"),
		expected: []object{{Type: "group", ID: "d1_viewers"}, {Type: "group", ID: "d1_subviewers"}},
	},
	{name: "users that are doc1 viewers", search: graph("doc", "doc1", "viewer", "user", "", ""),
		expected: []object{
			{Type: "user", ID: "user1"}, {Type: "user", ID: "user2"}, {Type: "user", ID: "user3"},
			{Type: "user", ID: "f1_owner"},
		},
	},
	{name: "users that are doc2 viewers (wildcard)", search: graph("doc", "doc2", "viewer", "user", "", ""),
		expected: []object{{Type: "user", ID: "*"}, {Type: "user", ID: "user2"}},
	},
	{name: "groups where members are members of d1_viewers", search: graph("group", "d1_viewers", "member", "group", "", "member"),
		expected: []object{{Type: "group", ID: "d1_subviewers"}},
	},
	{name: "users that are d1_viewers members", search: graph("group", "d1_viewers", "member", "user", "", ""),
		expected: []object{{Type: "user", ID: "user2"}, {Type: "user", ID: "user3"}},
	},
	{name: "users that are members of root group", search: graph("group", "root", "member", "user", "", ""),
		expected: []object{{Type: "user", ID: "leaf_user"}},
	},
	{name: "groups where members are are members of root group", search: graph("group", "root", "member", "group", "", "member"),
		expected: []object{{Type: "group", ID: "leaf"}, {Type: "group", ID: "branch"}, {Type: "group", ID: "trunk"}},
	},

	// Permissions
	{name: "users that can_create_file on folder1", search: graph("folder", "folder1", "can_create_file", "user", "", ""),
		expected: []object{{Type: "user", ID: "f1_owner"}},
	},
	{name: "users that can_read folder1", search: graph("folder", "folder1", "can_read", "user", "", ""),
		expected: []object{{Type: "user", ID: "f1_owner"}, {Type: "user", ID: "f1_viewer"}},
	},
	{name: "groups where members can_read folder1", search: graph("folder", "folder1", "can_read", "group", "", "member"),
		expected: []object{{Type: "group", ID: "f1_viewers"}},
	},
	{name: "users that cah_share folder1", search: graph("folder", "folder1", "can_share", "user", "", ""),
		expected: []object{{Type: "user", ID: "f1_owner"}},
	},
	{name: "user that can can_change_owner on doc1", search: graph("doc", "doc1", "can_change_owner", "user", "", ""),
		expected: []object{{Type: "user", ID: "d1_owner"}, {Type: "user", ID: "f1_owner"}},
	},
	{name: "users that can_write doc1", search: graph("doc", "doc1", "can_write", "user", "", ""),
		expected: []object{{Type: "user", ID: "d1_owner"}, {Type: "user", ID: "f1_owner"}},
	},
	{name: "users that can_read doc1", search: graph("doc", "doc1", "can_read", "user", "", ""),
		expected: []object{
			{Type: "user", ID: "d1_owner"}, {Type: "user", ID: "f1_owner"},
			{Type: "user", ID: "user1"}, {Type: "user", ID: "user2"}, {Type: "user", ID: "user3"},
			{Type: "user", ID: "f1_viewer"},
		},
	},
	{name: "groups where members can_read doc1", search: graph("doc", "doc1", "can_read", "group", "", "member"),
		expected: []object{{Type: "group", ID: "f1_viewers"}, {Type: "group", ID: "d1_viewers"}, {Type: "group", ID: "d1_subviewers"}},
	},
	{name: "users that can_share doc1", search: graph("doc", "doc1", "can_share", "user", "", ""),
		expected: []object{{Type: "user", ID: "f1_owner"}},
	},
	{name: "users that can_invite doc1", search: graph("doc", "doc1", "can_invite", "user", "", ""),
		expected: []object{{Type: "user", ID: "f1_viewer"}},
	},
	{name: "groups with members that can_invite doc1", search: graph("doc", "doc1", "can_invite", "group", "", "member"),
		expected: []object{{Type: "group", ID: "f1_viewers"}},
	},
	{name: "users that can_change_owner on doc2", search: graph("doc", "doc2", "can_change_owner", "user", "", ""),
		expected: []object{{Type: "user", ID: "f1_owner"}},
	},
	{name: "users that can_write doc2", search: graph("doc", "doc2", "can_write", "user", "", ""),
		expected: []object{{Type: "user", ID: "f1_owner"}},
	},
	{name: "users that can_read doc2", search: graph("doc", "doc2", "can_read", "user", "", ""),
		expected: []object{
			{Type: "user", ID: "*"}, {Type: "user", ID: "user2"},
			{Type: "user", ID: "f1_owner"}, {Type: "user", ID: "f1_viewer"},
		},
	},
	{name: "users that can_share doc2", search: graph("doc", "doc2", "can_share", "user", "", ""),
		expected: []object{{Type: "user", ID: "f1_owner"}},
	},
	{name: "users that can_invite doc2", search: graph("doc", "doc2", "can_invite", "user", "", ""),
		expected: []object{{Type: "user", ID: "f1_owner"}, {Type: "user", ID: "f1_viewer"}},
	},
	{name: "groups where members can_read doc2", search: graph("doc", "doc2", "can_read", "group", "", "member"),
		expected: []object{{Type: "group", ID: "f1_viewers"}},
	},
	{name: "users that can_change_owner on doc3", search: graph("doc", "doc3", "can_change_owner", "user", "", ""),
		expected: []object{{Type: "user", ID: "f1_owner"}},
	},
	{name: "users that can_write doc3", search: graph("doc", "doc3", "can_write", "user", "", ""),
		expected: []object{{Type: "user", ID: "f1_owner"}},
	},
	{name: "users that can_read doc3", search: graph("doc", "doc3", "can_read", "user", "", ""),
		expected: []object{{Type: "user", ID: "f1_owner"}, {Type: "user", ID: "f1_viewer"}},
	},
	{name: "users that can_share doc3", search: graph("doc", "doc3", "can_share", "user", "", ""),
		expected: []object{{Type: "user", ID: "f1_owner"}},
	},
	{name: "users that can_invite doc3", search: graph("doc", "doc3", "can_invite", "user", "", ""),
		expected: []object{{Type: "user", ID: "f1_owner"}, {Type: "user", ID: "f1_viewer"}},
	},
	{name: "groups where members can_read doc3", search: graph("doc", "doc3", "can_read", "group", "", "member"),
		expected: []object{{Type: "group", ID: "f1_viewers"}},
	},
}
