package cache_test

import (
	"os"
	"testing"

	"github.com/aserto-dev/azm/model"
	v3 "github.com/aserto-dev/azm/v3"
	"github.com/aserto-dev/azm/walk"
	dsc "github.com/aserto-dev/go-directory/aserto/directory/common/v3"
	dsr "github.com/aserto-dev/go-directory/aserto/directory/reader/v3"
	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
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
	matches := lo.Filter(r, func(rel *relation, _ int) bool {
		return rel.ObjectType == model.ObjectName(req.ObjectType) &&
			rel.ObjectID == req.ObjectId &&
			rel.Relation == model.RelationName(req.Relation) &&
			rel.SubjectType == model.ObjectName(req.SubjectType) &&
			rel.SubjectRelation == model.RelationName(req.SubjectRelation)
	})

	return lo.Map(matches, func(r *relation, _ int) *dsc.Relation {
		return r.AsProto()
	}), nil
}

func TestCheckRelation(t *testing.T) {
	rels := RelationsReader{
		{"doc", "doc1", "viewer", "group", "doc1_viewers", "member"}, // viewers group
		{"doc", "doc1", "viewer", "user", "user1", ""},               // direct viewer
		{"group", "doc1_viewers", "member", "user", "user2", ""},     // group member
		{"doc", "doc2", "viewer", "user", "*", ""},                   // wildcard

		{"group", "doc1_viewers", "member", "group", "d1_subviewers", "member"},
		{"group", "d1_subviewers", "member", "user", "user3", ""},

		// mutually recursive groups with users
		{"group", "yin", "member", "group", "yang", "member"},
		{"group", "yang", "member", "group", "yin", "member"},
		{"group", "yin", "member", "user", "yin_user", ""},
		{"group", "yang", "member", "user", "yang_user", ""},

		// mutually recursive groups with no users
		{"group", "alpha", "member", "group", "omega", "member"},
		{"group", "omega", "member", "group", "alpha", "member"},
	}

	tests := []struct {
		name     string
		check    *dsr.CheckRequest
		expected bool
	}{
		{name: "no assignment", check: check("doc", "doc1", "owner", "user", "user1"), expected: false},
		{name: "direct assignment", check: check("doc", "doc1", "viewer", "user", "user1"), expected: true},
		{name: "wildcard", check: check("doc", "doc2", "viewer", "user", "user1"), expected: true},
		{name: "wildcard", check: check("doc", "doc2", "viewer", "user", "userX"), expected: true},
		{name: "subject relation", check: check("doc", "doc1", "viewer", "user", "user2"), expected: true},
		{name: "nested groups", check: check("doc", "doc1", "viewer", "user", "user3"), expected: true},
		{name: "container not in set", check: check("doc", "doc1", "viewer", "group", "doc1_viewers"), expected: false},

		{name: "recursive groups - yin/yin", check: check("group", "yin", "member", "user", "yin_user"), expected: true},
		{name: "recursive groups - yin/yang", check: check("group", "yin", "member", "user", "yang_user"), expected: true},
		{name: "recursive groups - yang/yin", check: check("group", "yang", "member", "user", "yin_user"), expected: true},
		{name: "recursive groups - yang/yang", check: check("group", "yang", "member", "user", "yang_user"), expected: true},

		{name: "recursive groups - alpha/omega", check: check("group", "alpha", "member", "user", "user1"), expected: false},
	}

	r, err := os.Open("./path_test.yaml")
	require.NoError(t, err)
	require.NotNil(t, r)

	m, err := v3.Load(r)
	require.NoError(t, err)
	require.NotNil(t, m)

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			assert := require.New(tt)

			walker := walk.New(m, nil, test.check, rels.GetRelations)

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
