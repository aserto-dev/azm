package graph_test

import (
	"os"
	"testing"

	azmgraph "github.com/aserto-dev/azm/graph"
	"github.com/aserto-dev/azm/model"
	v3 "github.com/aserto-dev/azm/v3"
	dsc "github.com/aserto-dev/go-directory/aserto/directory/common/v3"
	dsr "github.com/aserto-dev/go-directory/aserto/directory/reader/v3"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

func TestSearchObjects(t *testing.T) {
	rels := relations()

	r, err := os.Open("./check_test.yaml")
	assert.NoError(t, err)
	assert.NotNil(t, r)

	m, err := v3.Load(r)
	assert.NoError(t, err)
	assert.NotNil(t, m)

	for _, test := range searchObjectsTests {
		t.Run(test.name, func(tt *testing.T) {
			assert := assert.New(tt)

			objSearch := azmgraph.NewObjectSearch(m, test.search, rels.GetRelations)

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
	name     string
	search   *dsr.GetGraphRequest
	expected []object
}

var searchObjectsTests = []searchTest{
	// Relations
	{name: "groups where members of leaf are members", search: graph("group", "", "member", "group", "leaf", "member"),
		expected: []object{{Type: "group", ID: "branch"}, {Type: "group", ID: "trunk"}, {Type: "group", ID: "root"}},
	},
	{name: "groups where members of leaf are guests", search: graph("group", "", "guest", "group", "leaf", "member"),
		expected: []object{},
	},
	{name: "docs where members of leaf are viewers", search: graph("doc", "", "viewer", "group", "leaf", "member"),
		expected: []object{{Type: "doc", ID: "doc_tree"}},
	},
	{name: "groups where members of yang are members", search: graph("group", "", "member", "group", "yang", "member"),
		expected: []object{{Type: "group", ID: "yin"}, {Type: "group", ID: "yang"}},
	},
	{name: "groups where yin_user is a member", search: graph("group", "", "member", "user", "yin_user", ""),
		expected: []object{{Type: "group", ID: "yin"}, {Type: "group", ID: "yang"}},
	},
	{name: "folders owned by f1_owner", search: graph("folder", "", "owner", "user", "f1_owner", ""),
		expected: []object{{Type: "folder", ID: "folder1"}},
	},
	{name: "folders where members of f1_viewers are viewers ", search: graph("folder", "", "viewer", "group", "f1_viewers", "member"),
		expected: []object{{Type: "folder", ID: "folder1"}},
	},
	{name: "docs where members of d1_subviewers are viewers", search: graph("doc", "", "viewer", "group", "d1_subviewers", "member"),
		expected: []object{{Type: "doc", ID: "doc1"}},
	},
	{name: "groups where f1_viewer is a member", search: graph("group", "", "member", "user", "f1_viewer", ""),
		expected: []object{{Type: "group", ID: "f1_viewers"}},
	},
	{name: "folders where folder1 is parent", search: graph("folder", "", "parent", "folder", "folder1", ""),
		expected: []object{{Type: "folder", ID: "folder2"}},
	},
	{name: "docs where folder1 is parent", search: graph("doc", "", "parent", "folder", "folder1", ""),
		expected: []object{{Type: "doc", ID: "doc1"}, {Type: "doc", ID: "doc2"}},
	},
	{name: "docs where d1_owner is owner", search: graph("doc", "", "owner", "user", "d1_owner", ""),
		expected: []object{{Type: "doc", ID: "doc1"}},
	},
	{name: "docs where members of d1_viewers are viewers", search: graph("doc", "", "viewer", "group", "d1_viewers", "member"),
		expected: []object{{Type: "doc", ID: "doc1"}},
	},
	{name: "docs where user1 is viewer", search: graph("doc", "", "viewer", "user", "user1", ""),
		expected: []object{{Type: "doc", ID: "doc1"}, {Type: "doc", ID: "doc2"}},
	},
	{name: "docs where f1_owner is viewer", search: graph("doc", "", "viewer", "user", "f1_owner", ""),
		expected: []object{{Type: "doc", ID: "doc1"}, {Type: "doc", ID: "doc2"}},
	},
	{name: "groups where user2 is a member", search: graph("group", "", "member", "user", "user2", ""),
		expected: []object{{Type: "group", ID: "d1_viewers"}},
	},
	{name: "docs where f1_owner is viewer", search: graph("doc", "", "viewer", "user", "f1_owner", ""),
		expected: []object{{Type: "doc", ID: "doc1"}, {Type: "doc", ID: "doc2"}},
	},
	{name: "docs where every user is viewer", search: graph("doc", "", "viewer", "user", "*", ""),
		expected: []object{{Type: "doc", ID: "doc2"}},
	},
	{name: "docs where user2 is viewer", search: graph("doc", "", "viewer", "user", "user2", ""),
		expected: []object{{Type: "doc", ID: "doc1"}, {Type: "doc", ID: "doc2"}},
	},
	{name: "groups where members of d1_subviewers are members", search: graph("group", "", "member", "group", "d1_subviewers", "member"),
		expected: []object{{Type: "group", ID: "d1_viewers"}},
	},
	{name: "groups where user3 is a member", search: graph("group", "", "member", "user", "user3", ""),
		expected: []object{{Type: "group", ID: "d1_subviewers"}, {Type: "group", ID: "d1_viewers"}},
	},

	// Permissions
	// {name: "folders where f1_owner is_owner", search: graph("folder", "", "is_owner", "user", "f1_owner", ""),
	//     expected: []object{{Type: "folder", ID: "folder1"}},
	// },
	// {name: "folders where f1_owner can_create_file", search: graph("folder", "", "can_create_file", "user", "f1_owner", ""),
	//     expected: []object{{Type: "folder", ID: "folder1"}},
	// },
	// {name: "folders where f1_owner can_read", search: graph("folder", "", "can_read", "user", "f1_owner", ""),
	//     expected: []object{{Type: "folder", ID: "folder1"}},
	// },
	// {name: "folders where f1_owner can_share", search: graph("folder", "", "can_share", "user", "f1_owner", ""),
	//     expected: []object{{Type: "folder", ID: "folder1"}},
	// },
	// {name: "docs where f1_owner can_change_owner", search: graph("doc", "", "can_change_owner", "user", "f1_owner", ""),
	//     expected: []object{{Type: "doc", ID: "doc1"}, {Type: "doc", ID: "doc2"}, {Type: "doc", ID: "doc3"}},
	// },
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
	return RelationsReader{
		{"folder", "folder1", "owner", "user", "f1_owner", ""},           // folder:folder1#owner@user:f1_owner
		{"folder", "folder2", "parent", "folder", "folder1", ""},         // folder:folder2#parent@folder:folder1
		{"folder", "folder1", "viewer", "group", "f1_viewers", "member"}, // folder:folder1#viewer@group:f1_viewers#member
		{"group", "f1_viewers", "member", "user", "f1_viewer", ""},       // group:f1_viewers#member@user:f1_viewer
		{"doc", "doc1", "parent", "folder", "folder1", ""},               // doc:doc1#parent@folder:folder1
		{"doc", "doc1", "owner", "user", "d1_owner", ""},                 // doc:doc1#owner@user:d1_owner
		{"doc", "doc1", "viewer", "group", "d1_viewers", "member"},       // doc:doc1#viewer@group:d1_viewers#member
		{"doc", "doc1", "viewer", "user", "user1", ""},                   // doc:doc1#viewer@user:user1
		{"doc", "doc1", "viewer", "user", "f1_owner", ""},                // doc:doc1#viewer@user:f1_owner
		{"group", "d1_viewers", "member", "user", "user2", ""},           // group:d1_viewers#member@user:user2
		{"doc", "doc2", "parent", "folder", "folder1", ""},               // doc:doc2#parnet@folder:folder1
		{"doc", "doc2", "viewer", "user", "*", ""},                       // doc:doc2#viewer@user:*
		{"doc", "doc2", "viewer", "user", "user2", ""},                   // doc:doc2#viewer@user:user2
		{"doc", "doc3", "parent", "folder", "folder2", ""},               // doc:doc3#parnet@folder:folder2

		{"group", "d1_viewers", "member", "group", "d1_subviewers", "member"}, // group:d1_viewers#member@group:d1_subviewers#member
		{"group", "d1_subviewers", "member", "user", "user3", ""},             // group:d1_subviewers#member@user:user3
		// {"group", "f1_viewers", "member", "group", "f1_subviewers", "member"}, // group:f1_viewers#member@group:f1_subviewers#member
		// {"group", "f1_subviewers", "member", "user", "user4", ""},             // group:d1_subviewers#member@user:user4

		// nested groups
		{"group", "leaf", "member", "user", "leaf_user", ""},
		{"group", "branch", "member", "group", "leaf", "member"},
		{"group", "trunk", "member", "group", "branch", "member"},
		{"group", "root", "member", "group", "trunk", "member"},
		{"doc", "doc_tree", "viewer", "group", "root", "member"},

		// mutually recursive groups with users
		{"group", "yin", "member", "group", "yang", "member"}, // group:yin#member@group:yang#member
		{"group", "yang", "member", "group", "yin", "member"}, // group:yang#member@group:yin#member
		{"group", "yin", "member", "user", "yin_user", ""},    // group:yin#member@user:yin_user
		{"group", "yang", "member", "user", "yang_user", ""},  // group:yang#member@user:yang_user

		// mutually recursive groups with no users
		{"group", "alpha", "member", "group", "omega", "member"}, // group:alpha#member@group:omega#member
		{"group", "omega", "member", "group", "alpha", "member"}, // group:omega#member@group:alpha#member
	}
}

func graph(
	objectType model.ObjectName, objectID string, // nolint: unparam
	relation model.RelationName,
	subjectType model.ObjectName, subjectID string,
	subjectRelation model.RelationName,
) *dsr.GetGraphRequest {
	return &dsr.GetGraphRequest{
		ObjectType:      objectType.String(),
		ObjectId:        objectID,
		Relation:        relation.String(),
		SubjectType:     subjectType.String(),
		SubjectId:       subjectID,
		SubjectRelation: subjectRelation.String(),
		Explain:         true,
		Trace:           true,
	}
}
