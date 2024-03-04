package graph_test

import (
	"testing"

	"github.com/aserto-dev/azm/graph"
	"github.com/aserto-dev/azm/model"
	v3 "github.com/aserto-dev/azm/v3"
	dsc "github.com/aserto-dev/go-directory/aserto/directory/common/v3"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

func TestSearchObjects(t *testing.T) {
	rels := relations()

	m, err := v3.LoadFile("./check_test.yaml")
	assert.NoError(t, err)
	assert.NotNil(t, m)

	for _, test := range searchObjectsTests {
		t.Run(test.search, func(tt *testing.T) {
			assert := assert.New(tt)

			objSearch := graph.NewObjectSearch(m, graphReq(test.search), rels.GetRelations)

			res, err := objSearch.Search()
			assert.NoError(err)
			tt.Logf("explanation: +%v\n", res.Explanation.AsMap())
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

type object struct {
	Type model.ObjectName
	ID   model.ObjectID
}

type searchTest struct {
	search   string
	expected []object
}

var searchObjectsTests = []searchTest{
	// Relations
	{"group:?#member@group:leaf#member", []object{{"group", "branch"}, {"group", "trunk"}, {"group", "root"}}},
	{"group:?#guest@group:leaf#member", []object{}},
	{"doc:?#viewer@group:leaf#member", []object{{"doc", "doc_tree"}}},
	{"group:?#member@group:yang#member", []object{{"group", "yin"}, {"group", "yang"}}},
	{"group:?#member@user:yin_user", []object{{"group", "yin"}, {"group", "yang"}}},
	{"folder:?#owner@user:f1_owner", []object{{"folder", "folder1"}}},
	{"folder:?#viewer@group:f1_viewers#member", []object{{"folder", "folder1"}}},
	{"doc:?#viewer@group:d1_subviewers#member", []object{{"doc", "doc1"}}},
	{"group:?#member@user:f1_viewer", []object{{"group", "f1_viewers"}}},
	{"folder:?#parent@folder:folder1", []object{{"folder", "folder2"}}},
	{"doc:?#parent@folder:folder1", []object{{"doc", "doc1"}, {"doc", "doc2"}}},
	{"doc:?#owner@user:d1_owner", []object{{"doc", "doc1"}}},
	{"doc:?#viewer@group:d1_viewers#member", []object{{"doc", "doc1"}}},
	{"doc:?#viewer@user:user1", []object{{"doc", "doc1"}, {"doc", "doc2"}}},
	{"doc:?#viewer@user:f1_owner", []object{{"doc", "doc1"}, {"doc", "doc2"}}},
	{"group:?#member@user:user2", []object{{"group", "d1_viewers"}}},
	{"doc:?#viewer@user:f1_owner", []object{{"doc", "doc1"}, {"doc", "doc2"}}},
	{"doc:?#viewer@user:*", []object{{"doc", "doc2"}}},
	{"doc:?#viewer@user:user2", []object{{"doc", "doc1"}, {"doc", "doc2"}}},
	{"group:?#member@group:d1_subviewers#member", []object{{"group", "d1_viewers"}}},
	{"group:?#member@user:user3", []object{{"group", "d1_subviewers"}, {"group", "d1_viewers"}}},

	// // Permissions
	// {"folder:?#is_owner@user:f1_owner", []object{{"folder", "folder1"}, {"folder", "folder2"}}},
	// {"folder:?#can_create_file@user:f1_owner", []object{{"folder", "folder1"}}},
	// {"folder:?#can_read@user:f1_owner", []object{{"folder", "folder1"}}},
	// {"folder:?#can_share@user:f1_owner", []object{{"folder", "folder1"}}},
	// {"doc:?#can_change_owner@user:f1_owner", []object{{"doc", "doc1"}, {"doc", "doc2"}, {"doc", "doc3"}}},
	// {name: "docs where f1_owner can_write", search: graph("doc", "", "can_write", "user", "f1_owner", ""),
	//     expected: []object{{Type: "doc", ID: "doc1"}, {Type: "doc", ID: "doc2"}},
	// },
	// {name: "docs where f1_owner can_read", search: graph("doc", "", "can_read", "user", "f1_owner", ""),
	//     expected: []object{{Type: "doc", ID: "doc1"}, {Type: "doc", ID: "doc2"}},
	// },
	// {name: "docs where f1_owner can_share", search: graph("doc", "", "can_share", "user", "f1_owner", ""),
	//     expected: []object{{Type: "doc", ID: "doc1"}, {Type: "doc", ID: "doc2"}},
	// },
	// {name: "docs where f1_owner can_invite", search: graph("doc", "", "can_invite", "user", "f1_owner", ""),
	//     expected: []object{},
	// },
	// {name: "folders where members of f1_viewers is_owner ", search: graph("folder", "", "is_owner", "group", "f1_viewers", "member"),
	//     expected: []object{},
	// },
	// {name: "folders where members of f1_viewers can_create_file ", search: graph("folder", "", "can_create_file", "group", "f1_viewers", "member"),
	//     expected: []object{},
	// },
	// {name: "folders where members of f1_viewers can_read ", search: graph("folder", "", "can_read", "group", "f1_viewers", "member"),
	//     expected: []object{{Type: "folder", ID: "folder1"}},
	// },
	// {name: "folders where f1_viewers can_read ", search: graph("folder", "", "can_read", "group", "f1_viewers", ""),
	//     expected: []object{},
	// },
	// {name: "folders where members of f1_viewers can_share ", search: graph("folder", "", "can_share", "group", "f1_viewers", "member"),
	//     expected: []object{},
	// },
	// {name: "docs where members of f1_viewers can_change_owner", search: graph("doc", "", "can_change_owner", "group", "f1_viewers", "member"),
	//     expected: []object{},
	// },
	// {name: "docs where members of f1_viewers can_write", search: graph("doc", "", "can_write", "group", "f1_viewers", "member"),
	//     expected: []object{},
	// },
	// {name: "docs where members of f1_viewers can_read", search: graph("doc", "", "can_read", "group", "f1_viewers", "member"),
	//     expected: []object{{Type: "doc", ID: "doc1"}, {Type: "doc", ID: "doc2"}},
	// },
	// {name: "docs where members of f1_viewers can_share", search: graph("doc", "", "can_share", "group", "f1_viewers", "member"),
	//     expected: []object{},
	// },
	// {name: "docs where members of f1_viewers can_invite", search: graph("doc", "", "can_invite", "group", "f1_viewers", "member"),
	//     expected: []object{{Type: "doc", ID: "doc1"}, {Type: "doc", ID: "doc2"}},
	// },
	// {name: "folders where f1_viewer is_owner ", search: graph("folder", "", "is_owner", "user", "f1_viewer", ""),
	//     expected: []object{},
	// },
	// {name: "folders where f1_viewer can_create_file ", search: graph("folder", "", "can_create_file", "user", "f1_viewer", ""),
	//     expected: []object{},
	// },
	// {name: "folders where f1_viewer can_read ", search: graph("folder", "", "can_read", "user", "f1_viewer", ""),
	//     expected: []object{{Type: "folder", ID: "folder1"}},
	// },
	// {name: "folders where f1_viewer can_share ", search: graph("folder", "", "can_share", "user", "f1_viewer", ""),
	//     expected: []object{},
	// },
	// {name: "docs where f1_viewer can_change_owner", search: graph("doc", "", "can_change_owner", "user", "f1_viewer", ""),
	//     expected: []object{},
	// },
	// {name: "docs where f1_viewer can_write", search: graph("doc", "", "can_write", "user", "f1_viewer", ""),
	//     expected: []object{},
	// },
	// {name: "docs where f1_viewer can_read", search: graph("doc", "", "can_read", "user", "f1_viewer", ""),
	//     expected: []object{{Type: "doc", ID: "doc1"}, {Type: "doc", ID: "doc2"}},
	// },
	// {name: "docs where f1_viewer can_share", search: graph("doc", "", "can_share", "user", "f1_viewer", ""),
	//     expected: []object{},
	// },
	// {name: "docs where f1_viewer can_invite", search: graph("doc", "", "can_invite", "user", "f1_viewer", ""),
	//     expected: []object{{Type: "doc", ID: "doc1"}},
	// },
	// {name: "folders where d1_owner is_owner ", search: graph("folder", "", "is_owner", "user", "d1_owner", ""),
	//     expected: []object{},
	// },
	// {name: "folders where d1_owner can_create_file ", search: graph("folder", "", "can_create_file", "user", "d1_owner", ""),
	//     expected: []object{},
	// },
	// {name: "folders where d1_owner can_read ", search: graph("folder", "", "can_read", "user", "d1_owner", ""),
	//     expected: []object{},
	// },
	// {name: "folders where d1_owner can_share ", search: graph("folder", "", "can_share", "user", "d1_owner", ""),
	//     expected: []object{},
	// },
	// {name: "docs where d1_owner can_change_owner", search: graph("doc", "", "can_change_owner", "user", "d1_owner", ""),
	//     expected: []object{{Type: "doc", ID: "doc1"}},
	// },
	// {name: "docs where d1_owner can_write", search: graph("doc", "", "can_write", "user", "d1_owner", ""),
	//     expected: []object{{Type: "doc", ID: "doc1"}},
	// },
	// {name: "docs where d1_owner can_read", search: graph("doc", "", "can_read", "user", "d1_owner", ""),
	//     expected: []object{{Type: "doc", ID: "doc1"}, {Type: "doc", ID: "doc2"}},
	// },
	// {name: "docs where d1_owner can_share", search: graph("doc", "", "can_share", "user", "d1_owner", ""),
	//     expected: []object{},
	// },
	// {name: "docs where d1_owner can_invite", search: graph("doc", "", "can_invite", "user", "d1_owner", ""),
	//     expected: []object{},
	// },
	// {name: "folders where members of d1_viewers is_owner ", search: graph("folder", "", "is_owner", "group", "d1_viewers", "member"),
	//     expected: []object{},
	// },
	// {name: "folders where members of d1_viewers can_create_file ", search: graph("folder", "", "can_create_file", "group", "d1_viewers", "member"),
	//     expected: []object{},
	// },
	// {name: "folders where members of d1_viewers can_read ", search: graph("folder", "", "can_read", "group", "d1_viewers", "member"),
	//     expected: []object{},
	// },
	// {name: "folders where d1_viewers can_read ", search: graph("folder", "", "can_read", "group", "d1_viewers", ""),
	//     expected: []object{},
	// },
	// {name: "folders where members of d1_viewers can_share ", search: graph("folder", "", "can_share", "group", "d1_viewers", "member"),
	//     expected: []object{},
	// },
	// {name: "docs where members of d1_viewers can_change_owner", search: graph("doc", "", "can_change_owner", "group", "d1_viewers", "member"),
	//     expected: []object{},
	// },
	// {name: "docs where members of d1_viewers can_write", search: graph("doc", "", "can_write", "group", "d1_viewers", "member"),
	//     expected: []object{},
	// },
	// {name: "docs where members of d1_viewers can_read", search: graph("doc", "", "can_read", "group", "d1_viewers", "member"),
	//     expected: []object{{Type: "doc", ID: "doc1"}},
	// },
	// {name: "docs where members of d1_viewers can_share", search: graph("doc", "", "can_share", "group", "d1_viewers", "member"),
	//     expected: []object{},
	// },
	// {name: "docs where members of d1_viewers can_invite", search: graph("doc", "", "can_invite", "group", "d1_viewers", "member"),
	//     expected: []object{},
	// },
	// {name: "folders where user1 is_owner ", search: graph("folder", "", "is_owner", "user", "user1", ""),
	//     expected: []object{},
	// },
	// {name: "folders where user1 can_create_file ", search: graph("folder", "", "can_create_file", "user", "user1", ""),
	//     expected: []object{},
	// },
	// {name: "folders where user1 can_read ", search: graph("folder", "", "can_read", "user", "user1", ""),
	//     expected: []object{},
	// },
	// {name: "folders where user1 can_share ", search: graph("folder", "", "can_share", "user", "user1", ""),
	//     expected: []object{},
	// },
	// {name: "docs where user1 can_change_owner", search: graph("doc", "", "can_change_owner", "user", "user1", ""),
	//     expected: []object{},
	// },
	// {name: "docs where user1 can_write", search: graph("doc", "", "can_write", "user", "user1", ""),
	//     expected: []object{},
	// },
	// {name: "docs where user1 can_read", search: graph("doc", "", "can_read", "user", "user1", ""),
	//     expected: []object{{Type: "doc", ID: "doc1"}, {Type: "doc", ID: "doc2"}},
	// },
	// {name: "docs where user1 can_share", search: graph("doc", "", "can_share", "user", "user1", ""),
	//     expected: []object{},
	// },
	// {name: "docs where user1 can_invite", search: graph("doc", "", "can_invite", "user", "user1", ""),
	//     expected: []object{},
	// },
	// {name: "folders where user2 is_owner ", search: graph("folder", "", "is_owner", "user", "user2", ""),
	//     expected: []object{},
	// },
	// {name: "folders where user2 can_create_file ", search: graph("folder", "", "can_create_file", "user", "user2", ""),
	//     expected: []object{},
	// },
	// {name: "folders where user2 can_read ", search: graph("folder", "", "can_read", "user", "user2", ""),
	//     expected: []object{},
	// },
	// {name: "folders where user2 can_share ", search: graph("folder", "", "can_share", "user", "user2", ""),
	//     expected: []object{},
	// },
	// {name: "docs where user2 can_change_owner", search: graph("doc", "", "can_change_owner", "user", "user2", ""),
	//     expected: []object{},
	// },
	// {name: "docs where user2 can_write", search: graph("doc", "", "can_write", "user", "user2", ""),
	//     expected: []object{},
	// },
	// {name: "docs where user2 can_read", search: graph("doc", "", "can_read", "user", "user2", ""),
	//     expected: []object{{Type: "doc", ID: "doc1"}, {Type: "doc", ID: "doc2"}},
	// },
	// {name: "docs where user2 can_share", search: graph("doc", "", "can_share", "user", "user2", ""),
	//     expected: []object{},
	// },
	// {name: "docs where user2 can_invite", search: graph("doc", "", "can_invite", "user", "user2", ""),
	//     expected: []object{},
	// },
	// {name: "folders where user3 is_owner ", search: graph("folder", "", "is_owner", "user", "user3", ""),
	//     expected: []object{},
	// },
	// {name: "folders where user3 can_create_file ", search: graph("folder", "", "can_create_file", "user", "user3", ""),
	//     expected: []object{},
	// },
	// {name: "folders where user3 can_read ", search: graph("folder", "", "can_read", "user", "user3", ""),
	//     expected: []object{},
	// },
	// {name: "folders where user3 can_share ", search: graph("folder", "", "can_share", "user", "user3", ""),
	//     expected: []object{},
	// },
	// {name: "docs where user3 can_change_owner", search: graph("doc", "", "can_change_owner", "user", "user3", ""),
	//     expected: []object{},
	// },
	// {name: "docs where user3 can_write", search: graph("doc", "", "can_write", "user", "user3", ""),
	//     expected: []object{},
	// },
	// {name: "docs where user3 can_read", search: graph("doc", "", "can_read", "user", "user3", ""),
	//     expected: []object{{Type: "doc", ID: "doc1"}, {Type: "doc", ID: "doc2"}},
	// },
	// {name: "docs where user3 can_share", search: graph("doc", "", "can_share", "user", "user3", ""),
	//     expected: []object{},
	// },
	// {name: "docs where user3 can_invite", search: graph("doc", "", "can_invite", "user", "user3", ""),
	//     expected: []object{},
	// },
}

func relations() RelationsReader {
	return NewRelationsReader(
		"folder:folder1#owner@user:f1_owner",
		"folder:folder2#parent@folder:folder1",
		"folder:folder1#viewer@group:f1_viewers#member",
		"group:f1_viewers#member@user:f1_viewer",
		"doc:doc1#parent@folder:folder1",
		"doc:doc1#owner@user:d1_owner",
		"doc:doc1#viewer@group:d1_viewers#member",
		"doc:doc1#viewer@user:user1",
		"doc:doc1#viewer@user:f1_owner",
		"group:d1_viewers#member@user:user2",
		"doc:doc2#parent@folder:folder1",
		"doc:doc2#viewer@user:*",
		"doc:doc2#viewer@user:user2",
		"doc:doc3#parent@folder:folder2",

		"group:d1_viewers#member@group:d1_subviewers#member",
		"group:d1_subviewers#member@user:user3",
		// "group:f1_viewers#member@group:f1_subviewers#member",
		// "group:d1_subviewers#member@user:user4",

		// nested groups
		"group:leaf#member@user:leaf_user",
		"group:branch#member@group:leaf#member",
		"group:trunk#member@group:branch#member",
		"group:root#member@group:trunk#member",
		"doc:doc_tree#viewer@group:root#member",

		// mutually recursive groups with users
		"group:yin#member@group:yang#member",
		"group:yang#member@group:yin#member",
		"group:yin#member@user:yin_user",
		"group:yang#member@user:yang_user",

		// mutually recursive groups with no users
		"group:alpha#member@group:omega#member",
		"group:omega#member@group:alpha#member",
	)
}
