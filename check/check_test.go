package check_test

import (
	"os"
	"strings"
	"testing"

	azmcheck "github.com/aserto-dev/azm/check"
	"github.com/aserto-dev/azm/model"
	v3 "github.com/aserto-dev/azm/v3"
	dsc "github.com/aserto-dev/go-directory/aserto/directory/common/v3"
	dsr "github.com/aserto-dev/go-directory/aserto/directory/reader/v3"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

func TestCheck(t *testing.T) {
	rels := relations()

	r, err := os.Open("./check_test.yaml")
	assert.NoError(t, err)
	assert.NotNil(t, r)

	m, err := v3.Load(r)
	assert.NoError(t, err)
	assert.NotNil(t, m)

	for _, test := range checkTests {
		t.Run(test.name, func(tt *testing.T) {
			assert := assert.New(tt)

			checker := azmcheck.New(m, test.check, rels.GetRelations)

			res, err := checker.Check()
			assert.NoError(err)
			tt.Log("trace:\n", strings.Join(checker.Trace(), "\n"))
			assert.Equal(test.expected, res)
		})
	}
}

func TestSearchSubjects(t *testing.T) {
	rels := relations()

	r, err := os.Open("./check_test.yaml")
	assert.NoError(t, err)
	assert.NotNil(t, r)

	m, err := v3.Load(r)
	assert.NoError(t, err)
	assert.NotNil(t, m)

	for _, test := range searchSubjectTests {
		t.Run(test.name, func(tt *testing.T) {
			assert := assert.New(tt)

			checker := azmcheck.NewGraph(m, test.search, rels.GetRelations)

			res, err := checker.Search()
			assert.NoError(err)
			tt.Log("trace:\n", strings.Join(checker.Trace(), "\n"))

			subjects := lo.Uniq(lo.Map(res, func(s azmcheck.CheckParams, _ int) object {
				return object{
					Type: s.ST,
					ID:   s.SID,
				}
			}))

			for _, e := range test.expected {
				assert.Contains(subjects, e)
			}

			assert.Equal(len(test.expected), len(subjects), subjects)

		})
	}
}

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

			checker := azmcheck.NewGraph(m, test.search, rels.GetRelations)

			res, err := checker.Search()
			assert.NoError(err)
			tt.Log("trace:\n", strings.Join(checker.Trace(), "\n"))

			subjects := lo.Uniq(lo.Map(res, func(s azmcheck.CheckParams, _ int) object {
				return object{
					Type: s.OT,
					ID:   s.OID,
				}
			}))

			for _, e := range test.expected {
				assert.Contains(subjects, e)
			}

			assert.Equal(len(test.expected), len(subjects), subjects)

		})
	}
}

var checkTests = []struct {
	name     string
	check    *dsr.CheckRequest
	expected bool
}{
	// Relations
	{name: "no assignment", check: check("doc", "doc1", "owner", "user", "user1"), expected: false},
	{name: "direct assignment", check: check("doc", "doc1", "viewer", "user", "user1"), expected: true},
	{name: "wildcard", check: check("doc", "doc2", "viewer", "user", "user1"), expected: true},
	{name: "wildcard", check: check("doc", "doc2", "viewer", "user", "userX"), expected: true},
	{name: "subject relation", check: check("doc", "doc1", "viewer", "user", "user2"), expected: true},
	{name: "nested groups", check: check("doc", "doc1", "viewer", "user", "user3"), expected: true},
	{name: "container not in set", check: check("doc", "doc1", "viewer", "group", "d1_viewers"), expected: false},

	{name: "recursive groups - yin/yin", check: check("group", "yin", "member", "user", "yin_user"), expected: true},
	{name: "recursive groups - yin/yang", check: check("group", "yin", "member", "user", "yang_user"), expected: true},
	{name: "recursive groups - yang/yin", check: check("group", "yang", "member", "user", "yin_user"), expected: true},
	{name: "recursive groups - yang/yang", check: check("group", "yang", "member", "user", "yang_user"), expected: true},

	{name: "recursive groups - alpha/omega", check: check("group", "alpha", "member", "user", "user1"), expected: false},

	// Permissions
	{name: "owner can change owner", check: check("doc", "doc1", "can_change_owner", "user", "d1_owner"), expected: true},
	{name: "viewer cannot change owner", check: check("doc", "doc1", "can_change_owner", "user", "user1"), expected: false},
	{name: "unrelated cannot change owner", check: check("doc", "doc1", "can_change_owner", "user", "userX"), expected: false},

	{name: "owner can read", check: check("doc", "doc1", "can_read", "user", "d1_owner"), expected: true},
	{name: "parent owner can read", check: check("doc", "doc1", "can_read", "user", "f1_owner"), expected: true},
	{name: "direct viewer can read", check: check("doc", "doc1", "can_read", "user", "user1"), expected: true},
	{name: "parent viewer can read", check: check("doc", "doc1", "can_read", "user", "f1_viewer"), expected: true},
	{name: "unrelated cannot read", check: check("doc", "doc1", "can_read", "user", "userX"), expected: false},

	{name: "owner can write", check: check("doc", "doc1", "can_write", "user", "d1_owner"), expected: true},
	{name: "parent owner can write", check: check("doc", "doc1", "can_write", "user", "f1_owner"), expected: true},
	{name: "viewer cannot write", check: check("doc", "doc1", "can_write", "user", "user2"), expected: false},

	{name: "folder owner", check: check("folder", "folder1", "owner", "user", "f1_owner"), expected: true},
	{name: "folder owner can create file", check: check("folder", "folder1", "can_create_file", "user", "f1_owner"), expected: true},
	{name: "folder owner can share", check: check("folder", "folder1", "can_share", "user", "f1_owner"), expected: true},

	// intersection
	{name: "writer cannot share", check: check("doc", "doc1", "can_share", "user", "d1_owner"), expected: false},
	{name: "parent owner can share", check: check("doc", "doc1", "can_share", "user", "f1_owner"), expected: true},

	// // negation
	{name: "f1_owner can read folder1", check: check("folder", "folder1", "can_read", "user", "f1_owner"), expected: true},
	{name: "f1_owner is doc1 viewer", check: check("doc", "doc1", "viewer", "user", "f1_owner"), expected: true},
	{name: "parent owner and doc viewer cannot invite", check: check("doc", "doc1", "can_invite", "user", "f1_owner"), expected: false},
}

type object struct {
	Type model.ObjectName
	ID   model.ObjectID
}

var searchSubjectTests = []struct {
	name     string
	search   *dsr.GetGraphRequest
	expected []object
}{
	// Relations
	{name: "users that are folder1 owners", search: graph("folder", "folder1", "owner", "user", "", ""),
		expected: []object{{Type: "user", ID: "f1_owner"}},
	},
	{name: "users that are folder1 viewers", search: graph("folder", "folder1", "viewer", "user", "", ""),
		expected: []object{{Type: "user", ID: "f1_viewer"}},
	},
	{name: "users that are folder2 owners", search: graph("folder", "folder2", "owner", "user", "", ""),
		expected: []object{{Type: "user", ID: "f1_owner"}},
	},
	{name: "users that are folder2 viewers", search: graph("folder", "folder2", "viewer", "user", "", ""),
		expected: []object{{Type: "user", ID: "f1_viewer"}},
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
		expected: []object{{Type: "group", ID: "f1_viewers"}},
	},
	{name: "users that are doc1 owners", search: graph("doc", "doc1", "owner", "user", "", ""),
		expected: []object{{Type: "user", ID: "d1_owner"}},
	},
	{name: "groups where members are doc1 viewers", search: graph("doc", "doc1", "viewer", "group", "", "member"),
		expected: []object{{Type: "group", ID: "d1_viewers"}},
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
		expected: []object{{Type: "group", ID: "f1_viewers"}, {Type: "group", ID: "d1_viewers"}},
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
		expected: []object{
			{Type: "user", ID: "*"}, {Type: "user", ID: "user2"},
			{Type: "user", ID: "f1_owner"}, {Type: "user", ID: "f1_viewer"},
		},
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

var searchObjectsTests = []struct {
	name     string
	search   *dsr.GetGraphRequest
	expected []object
}{
	// Relations
	// {name: "folders owned by f1_owner", search: graph("folder", "", "owner", "user", "f1_owner", ""),
	//     expected: []object{{Type: "folder", ID: "folder1"}},
	// },
	// {name: "folders where members of f1_viewers are viewers ", search: graph("folder", "", "viewer", "group", "f1_viewers", "member"),
	//     expected: []object{{Type: "folder", ID: "folder1"}},
	// },
	// {name: "docs where members of d1_subviewers are viewers", search: graph("doc", "", "viewer", "group", "d1_subviewers", "member"),
	//     expected: []object{{Type: "doc", ID: "doc1"}},
	// },
	// {name: "groups where f1_viewer is a member", search: graph("group", "", "member", "user", "f1_viewer", ""),
	//     expected: []object{{Type: "group", ID: "f1_viewers"}},
	// },
	// {name: "folders where folder1 is parent", search: graph("folder", "", "parent", "folder", "folder1", ""),
	//     expected: []object{{Type: "folder", ID: "folder2"}},
	// },
	// {name: "docs where folder1 is parent", search: graph("doc", "", "parent", "folder", "folder1", ""),
	//     expected: []object{{Type: "doc", ID: "doc1"}, {Type: "doc", ID: "doc2"}},
	// },
	// {name: "docs where d1_owner is owner", search: graph("doc", "", "owner", "user", "d1_owner", ""),
	//     expected: []object{{Type: "doc", ID: "doc1"}},
	// },
	// {name: "docs where members of d1_viewers are viewers", search: graph("doc", "", "viewer", "group", "d1_viewers", "member"),
	//     expected: []object{{Type: "doc", ID: "doc1"}},
	// },
	// {name: "docs where user1 is viewer", search: graph("doc", "", "viewer", "user", "user1", ""),
	//     expected: []object{{Type: "doc", ID: "doc1"}, {Type: "doc", ID: "doc2"}},
	// },
	// {name: "docs where f1_owner is viewer", search: graph("doc", "", "viewer", "user", "f1_owner", ""),
	//     expected: []object{{Type: "doc", ID: "doc1"}, {Type: "doc", ID: "doc2"}},
	// },
	// {name: "groups where user2 is a member", search: graph("group", "", "member", "user", "user2", ""),
	//     expected: []object{{Type: "group", ID: "d1_viewers"}},
	// },
	// {name: "docs where f1_owner is viewer", search: graph("doc", "", "viewer", "user", "f1_owner", ""),
	//     expected: []object{{Type: "doc", ID: "doc1"}, {Type: "doc", ID: "doc2"}},
	// },
	// {name: "docs where every user is viewer", search: graph("doc", "", "viewer", "user", "*", ""),
	//     expected: []object{{Type: "doc", ID: "doc2"}},
	// },
	// {name: "docs where user2 is viewer", search: graph("doc", "", "viewer", "user", "user2", ""),
	//     expected: []object{{Type: "doc", ID: "doc1"}, {Type: "doc", ID: "doc2"}},
	// },
	// {name: "groups where members of d1_subviewers are members", search: graph("group", "", "member", "group", "d1_subviewers", "member"),
	//     expected: []object{{Type: "group", ID: "d1_viewers"}},
	// },
	// {name: "groups where user3 is a member", search: graph("group", "", "member", "user", "user3", ""),
	//     expected: []object{{Type: "group", ID: "d1_subviewers"}, {Type: "group", ID: "d1_viewers"}},
	// },

	// // Permissions
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
	{name: "docs where f1_owner can_change_owner", search: graph("doc", "", "can_change_owner", "user", "f1_owner", ""),
		expected: []object{{Type: "doc", ID: "doc1"}, {Type: "doc", ID: "doc2"}, {Type: "doc", ID: "doc3"}},
	},
	{name: "docs where f1_owner can_write", search: graph("doc", "", "can_write", "user", "f1_owner", ""),
		expected: []object{{Type: "doc", ID: "doc1"}, {Type: "doc", ID: "doc2"}},
	},
	{name: "docs where f1_owner can_read", search: graph("doc", "", "can_read", "user", "f1_owner", ""),
		expected: []object{{Type: "doc", ID: "doc1"}, {Type: "doc", ID: "doc2"}},
	},
	{name: "docs where f1_owner can_share", search: graph("doc", "", "can_share", "user", "f1_owner", ""),
		expected: []object{{Type: "doc", ID: "doc1"}, {Type: "doc", ID: "doc2"}},
	},
	{name: "docs where f1_owner can_invite", search: graph("doc", "", "can_invite", "user", "f1_owner", ""),
		expected: []object{},
	},
	{name: "folders where members of f1_viewers is_owner ", search: graph("folder", "", "is_owner", "group", "f1_viewers", "member"),
		expected: []object{},
	},
	{name: "folders where members of f1_viewers can_create_file ", search: graph("folder", "", "can_create_file", "group", "f1_viewers", "member"),
		expected: []object{},
	},
	{name: "folders where members of f1_viewers can_read ", search: graph("folder", "", "can_read", "group", "f1_viewers", "member"),
		expected: []object{{Type: "folder", ID: "folder1"}},
	},
	{name: "folders where f1_viewers can_read ", search: graph("folder", "", "can_read", "group", "f1_viewers", ""),
		expected: []object{},
	},
	{name: "folders where members of f1_viewers can_share ", search: graph("folder", "", "can_share", "group", "f1_viewers", "member"),
		expected: []object{},
	},
	{name: "docs where members of f1_viewers can_change_owner", search: graph("doc", "", "can_change_owner", "group", "f1_viewers", "member"),
		expected: []object{},
	},
	{name: "docs where members of f1_viewers can_write", search: graph("doc", "", "can_write", "group", "f1_viewers", "member"),
		expected: []object{},
	},
	{name: "docs where members of f1_viewers can_read", search: graph("doc", "", "can_read", "group", "f1_viewers", "member"),
		expected: []object{{Type: "doc", ID: "doc1"}, {Type: "doc", ID: "doc2"}},
	},
	{name: "docs where members of f1_viewers can_share", search: graph("doc", "", "can_share", "group", "f1_viewers", "member"),
		expected: []object{},
	},
	{name: "docs where members of f1_viewers can_invite", search: graph("doc", "", "can_invite", "group", "f1_viewers", "member"),
		expected: []object{{Type: "doc", ID: "doc1"}, {Type: "doc", ID: "doc2"}},
	},
	{name: "folders where f1_viewer is_owner ", search: graph("folder", "", "is_owner", "user", "f1_viewer", ""),
		expected: []object{},
	},
	{name: "folders where f1_viewer can_create_file ", search: graph("folder", "", "can_create_file", "user", "f1_viewer", ""),
		expected: []object{},
	},
	{name: "folders where f1_viewer can_read ", search: graph("folder", "", "can_read", "user", "f1_viewer", ""),
		expected: []object{{Type: "folder", ID: "folder1"}},
	},
	{name: "folders where f1_viewer can_share ", search: graph("folder", "", "can_share", "user", "f1_viewer", ""),
		expected: []object{},
	},
	{name: "docs where f1_viewer can_change_owner", search: graph("doc", "", "can_change_owner", "user", "f1_viewer", ""),
		expected: []object{},
	},
	{name: "docs where f1_viewer can_write", search: graph("doc", "", "can_write", "user", "f1_viewer", ""),
		expected: []object{},
	},
	{name: "docs where f1_viewer can_read", search: graph("doc", "", "can_read", "user", "f1_viewer", ""),
		expected: []object{{Type: "doc", ID: "doc1"}, {Type: "doc", ID: "doc2"}},
	},
	{name: "docs where f1_viewer can_share", search: graph("doc", "", "can_share", "user", "f1_viewer", ""),
		expected: []object{},
	},
	{name: "docs where f1_viewer can_invite", search: graph("doc", "", "can_invite", "user", "f1_viewer", ""),
		expected: []object{{Type: "doc", ID: "doc1"}},
	},
	{name: "folders where d1_owner is_owner ", search: graph("folder", "", "is_owner", "user", "d1_owner", ""),
		expected: []object{},
	},
	{name: "folders where d1_owner can_create_file ", search: graph("folder", "", "can_create_file", "user", "d1_owner", ""),
		expected: []object{},
	},
	{name: "folders where d1_owner can_read ", search: graph("folder", "", "can_read", "user", "d1_owner", ""),
		expected: []object{},
	},
	{name: "folders where d1_owner can_share ", search: graph("folder", "", "can_share", "user", "d1_owner", ""),
		expected: []object{},
	},
	{name: "docs where d1_owner can_change_owner", search: graph("doc", "", "can_change_owner", "user", "d1_owner", ""),
		expected: []object{{Type: "doc", ID: "doc1"}},
	},
	{name: "docs where d1_owner can_write", search: graph("doc", "", "can_write", "user", "d1_owner", ""),
		expected: []object{{Type: "doc", ID: "doc1"}},
	},
	{name: "docs where d1_owner can_read", search: graph("doc", "", "can_read", "user", "d1_owner", ""),
		expected: []object{{Type: "doc", ID: "doc1"}, {Type: "doc", ID: "doc2"}},
	},
	{name: "docs where d1_owner can_share", search: graph("doc", "", "can_share", "user", "d1_owner", ""),
		expected: []object{},
	},
	{name: "docs where d1_owner can_invite", search: graph("doc", "", "can_invite", "user", "d1_owner", ""),
		expected: []object{},
	},
	{name: "folders where members of d1_viewers is_owner ", search: graph("folder", "", "is_owner", "group", "d1_viewers", "member"),
		expected: []object{},
	},
	{name: "folders where members of d1_viewers can_create_file ", search: graph("folder", "", "can_create_file", "group", "d1_viewers", "member"),
		expected: []object{},
	},
	{name: "folders where members of d1_viewers can_read ", search: graph("folder", "", "can_read", "group", "d1_viewers", "member"),
		expected: []object{},
	},
	{name: "folders where d1_viewers can_read ", search: graph("folder", "", "can_read", "group", "d1_viewers", ""),
		expected: []object{},
	},
	{name: "folders where members of d1_viewers can_share ", search: graph("folder", "", "can_share", "group", "d1_viewers", "member"),
		expected: []object{},
	},
	{name: "docs where members of d1_viewers can_change_owner", search: graph("doc", "", "can_change_owner", "group", "d1_viewers", "member"),
		expected: []object{},
	},
	{name: "docs where members of d1_viewers can_write", search: graph("doc", "", "can_write", "group", "d1_viewers", "member"),
		expected: []object{},
	},
	{name: "docs where members of d1_viewers can_read", search: graph("doc", "", "can_read", "group", "d1_viewers", "member"),
		expected: []object{{Type: "doc", ID: "doc1"}},
	},
	{name: "docs where members of d1_viewers can_share", search: graph("doc", "", "can_share", "group", "d1_viewers", "member"),
		expected: []object{},
	},
	{name: "docs where members of d1_viewers can_invite", search: graph("doc", "", "can_invite", "group", "d1_viewers", "member"),
		expected: []object{},
	},
	{name: "folders where user1 is_owner ", search: graph("folder", "", "is_owner", "user", "user1", ""),
		expected: []object{},
	},
	{name: "folders where user1 can_create_file ", search: graph("folder", "", "can_create_file", "user", "user1", ""),
		expected: []object{},
	},
	{name: "folders where user1 can_read ", search: graph("folder", "", "can_read", "user", "user1", ""),
		expected: []object{},
	},
	{name: "folders where user1 can_share ", search: graph("folder", "", "can_share", "user", "user1", ""),
		expected: []object{},
	},
	{name: "docs where user1 can_change_owner", search: graph("doc", "", "can_change_owner", "user", "user1", ""),
		expected: []object{},
	},
	{name: "docs where user1 can_write", search: graph("doc", "", "can_write", "user", "user1", ""),
		expected: []object{},
	},
	{name: "docs where user1 can_read", search: graph("doc", "", "can_read", "user", "user1", ""),
		expected: []object{{Type: "doc", ID: "doc1"}, {Type: "doc", ID: "doc2"}},
	},
	{name: "docs where user1 can_share", search: graph("doc", "", "can_share", "user", "user1", ""),
		expected: []object{},
	},
	{name: "docs where user1 can_invite", search: graph("doc", "", "can_invite", "user", "user1", ""),
		expected: []object{},
	},
	{name: "folders where user2 is_owner ", search: graph("folder", "", "is_owner", "user", "user2", ""),
		expected: []object{},
	},
	{name: "folders where user2 can_create_file ", search: graph("folder", "", "can_create_file", "user", "user2", ""),
		expected: []object{},
	},
	{name: "folders where user2 can_read ", search: graph("folder", "", "can_read", "user", "user2", ""),
		expected: []object{},
	},
	{name: "folders where user2 can_share ", search: graph("folder", "", "can_share", "user", "user2", ""),
		expected: []object{},
	},
	{name: "docs where user2 can_change_owner", search: graph("doc", "", "can_change_owner", "user", "user2", ""),
		expected: []object{},
	},
	{name: "docs where user2 can_write", search: graph("doc", "", "can_write", "user", "user2", ""),
		expected: []object{},
	},
	{name: "docs where user2 can_read", search: graph("doc", "", "can_read", "user", "user2", ""),
		expected: []object{{Type: "doc", ID: "doc1"}, {Type: "doc", ID: "doc2"}},
	},
	{name: "docs where user2 can_share", search: graph("doc", "", "can_share", "user", "user2", ""),
		expected: []object{},
	},
	{name: "docs where user2 can_invite", search: graph("doc", "", "can_invite", "user", "user2", ""),
		expected: []object{},
	},
	{name: "folders where user3 is_owner ", search: graph("folder", "", "is_owner", "user", "user3", ""),
		expected: []object{},
	},
	{name: "folders where user3 can_create_file ", search: graph("folder", "", "can_create_file", "user", "user3", ""),
		expected: []object{},
	},
	{name: "folders where user3 can_read ", search: graph("folder", "", "can_read", "user", "user3", ""),
		expected: []object{},
	},
	{name: "folders where user3 can_share ", search: graph("folder", "", "can_share", "user", "user3", ""),
		expected: []object{},
	},
	{name: "docs where user3 can_change_owner", search: graph("doc", "", "can_change_owner", "user", "user3", ""),
		expected: []object{},
	},
	{name: "docs where user3 can_write", search: graph("doc", "", "can_write", "user", "user3", ""),
		expected: []object{},
	},
	{name: "docs where user3 can_read", search: graph("doc", "", "can_read", "user", "user3", ""),
		expected: []object{{Type: "doc", ID: "doc1"}, {Type: "doc", ID: "doc2"}},
	},
	{name: "docs where user3 can_share", search: graph("doc", "", "can_share", "user", "user3", ""),
		expected: []object{},
	},
	{name: "docs where user3 can_invite", search: graph("doc", "", "can_invite", "user", "user3", ""),
		expected: []object{},
	},
}

type relation struct {
	ObjectType      model.ObjectName
	ObjectID        string
	Relation        model.RelationName
	SubjectType     model.ObjectName
	SubjectID       string
	SubjectRelation model.RelationName
}

func (r *relation) AsProto() *dsc.Relation {
	return &dsc.Relation{
		ObjectType:      r.ObjectType.String(),
		ObjectId:        r.ObjectID,
		Relation:        r.Relation.String(),
		SubjectType:     r.SubjectType.String(),
		SubjectId:       r.SubjectID,
		SubjectRelation: r.SubjectRelation.String(),
	}
}

type RelationsReader []*relation

func (r RelationsReader) GetRelations(req *dsc.Relation) ([]*dsc.Relation, error) {
	ot := model.ObjectName(req.ObjectType)
	rn := model.RelationName(req.Relation)
	st := model.ObjectName(req.SubjectType)
	sr := model.RelationName(req.SubjectRelation)

	matches := lo.Filter(r, func(rel *relation, _ int) bool {
		return (ot == "" || rel.ObjectType == ot) &&
			(req.ObjectId == "" || rel.ObjectID == req.ObjectId) &&
			(rn == "" || rel.Relation == rn) &&
			(st == "" || rel.SubjectType == st) &&
			(req.SubjectId == "" || rel.SubjectID == req.SubjectId) &&
			(sr == "" || rel.SubjectRelation == sr)
	})

	return lo.Map(matches, func(r *relation, _ int) *dsc.Relation {
		return r.AsProto()
	}), nil
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

func check(
	objectType model.ObjectName, objectID string,
	relation model.RelationName,
	subjectType model.ObjectName, subjectID string,
) *dsr.CheckRequest {
	return &dsr.CheckRequest{
		ObjectType:  objectType.String(),
		ObjectId:    objectID,
		Relation:    relation.String(),
		SubjectType: subjectType.String(),
		SubjectId:   subjectID,
		Trace:       true,
	}

}

func graph(
	objectType model.ObjectName, objectID string,
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
	}
}
