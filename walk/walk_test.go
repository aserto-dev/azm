package walk_test

import (
	"os"
	"testing"

	"github.com/aserto-dev/azm/model"
	v3 "github.com/aserto-dev/azm/v3"
	"github.com/aserto-dev/azm/walk"
	dsc "github.com/aserto-dev/go-directory/aserto/directory/common/v3"
	dsr "github.com/aserto-dev/go-directory/aserto/directory/reader/v3"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

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

func TestCheck(t *testing.T) {
	rels := RelationsReader{
		{"folder", "folder1", "owner", "user", "f1_owner", ""},           // folder:folder1#owner@user:f1_owner
		{"folder", "folder1", "viewer", "group", "f1_viewers", "member"}, // folder:folder1#viewer@group:f1_viewers#member
		{"group", "f1_viewers", "member", "user", "f1_viewer", ""},       // group:f1_viewers#member@user:f1_viewer
		{"doc", "doc1", "parent", "folder", "folder1", ""},               // doc:doc1#parent@folder:folder1
		{"doc", "doc1", "owner", "user", "d1_owner", ""},                 // doc:doc1#owner@user:d1_owner
		{"doc", "doc1", "viewer", "group", "d1_viewers", "member"},       // doc:doc1#viewer@group:d1_viewers#member
		{"doc", "doc1", "viewer", "user", "user1", ""},                   // doc:doc1#viewer@user:user1
		{"group", "d1_viewers", "member", "user", "user2", ""},           // group:d1_viewers#member@user:user2
		{"doc", "doc2", "viewer", "user", "*", ""},                       // doc:doc2#viewer@user:*

		{"group", "d1_viewers", "member", "group", "d1_subviewers", "member"}, // group:d1_viewers#member@group:d1_subviewers#member
		{"group", "d1_subviewers", "member", "user", "user3", ""},             // group:d1_subviewers#member@user:user3

		// mutually recursive groups with users
		{"group", "yin", "member", "group", "yang", "member"}, // group:yin#member@group:yang#member
		{"group", "yang", "member", "group", "yin", "member"}, // group:yang#member@group:yin#member
		{"group", "yin", "member", "user", "yin_user", ""},    // group:yin#member@user:yin_user
		{"group", "yang", "member", "user", "yang_user", ""},  // group:yang#member@user:yang_user

		// mutually recursive groups with no users
		{"group", "alpha", "member", "group", "omega", "member"}, // group:alpha#member@group:omega#member
		{"group", "omega", "member", "group", "alpha", "member"}, // group:omega#member@group:alpha#member
	}

	tests := []struct {
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

		// // // intersection
		{name: "writer cannot share", check: check("doc", "doc1", "can_share", "user", "d1_owner"), expected: false},
		{name: "parent owner can share", check: check("doc", "doc1", "can_share", "user", "f1_owner"), expected: true},

		// // negation
		{name: "f1_owner can read folder1", check: check("folder", "folder1", "can_read", "user", "f1_owner"), expected: true},
		{name: "f1_owner isn't doc1 viewer", check: check("doc", "doc1", "viewer", "user", "f1_owner"), expected: false},
		{name: "parent owner can invite", check: check("doc", "doc1", "can_invite", "user", "f1_owner"), expected: true},
	}

	r, err := os.Open("./walk_test.yaml")
	assert.NoError(t, err)
	assert.NotNil(t, r)

	m, err := v3.Load(r)
	assert.NoError(t, err)
	assert.NotNil(t, m)

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			assert := assert.New(tt)

			walker := walk.New(m, test.check, rels.GetRelations)

			res, err := walker.Check()
			assert.NoError(err)
			assert.Equal(test.expected, res)
		})
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
	}

}
