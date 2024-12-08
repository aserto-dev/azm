package graph_test

import (
	"regexp"
	"testing"

	"github.com/aserto-dev/azm/graph"
	"github.com/aserto-dev/azm/model"
	dsc "github.com/aserto-dev/go-directory/aserto/directory/common/v3"
	dsr "github.com/aserto-dev/go-directory/aserto/directory/reader/v3"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

func checkReq(expr string, trace bool) *dsr.CheckRequest {
	return parseRelation(expr).checkReq(trace)
}

func graphReq(expr string) *dsr.GetGraphRequest {
	return parseRelation(expr).graphReq()
}

type relation struct {
	ObjectType      model.ObjectName
	ObjectID        model.ObjectID
	Relation        model.RelationName
	SubjectType     model.ObjectName
	SubjectID       model.ObjectID
	SubjectRelation model.RelationName
}

var rx = regexp.MustCompile(`^(\w+):([\w\?]+)#(\w+)@(\w+):([\w\?\*]+)(#?\w*)$`)

func parseRelation(r string) *relation {
	matches := rx.FindStringSubmatch(r)
	if len(matches) < 7 {
		return nil
	}
	rel := &relation{
		ObjectType:  model.ObjectName(matches[1]),
		ObjectID:    model.ObjectID(lo.Ternary(matches[2] == "?", "", matches[2])),
		Relation:    model.RelationName(matches[3]),
		SubjectType: model.ObjectName(matches[4]),
		SubjectID:   model.ObjectID(lo.Ternary(matches[5] == "?", "", matches[5])),
	}
	if matches[6] != "" {
		rel.SubjectRelation = model.RelationName(matches[6][1:])
	}

	return rel
}

func (r *relation) proto(pool graph.MessagePool[dsc.Relation, *dsc.Relation]) *dsc.Relation {
	rel := pool.Get()
	rel.ObjectType = r.ObjectType.String()
	rel.ObjectId = r.ObjectID.String()
	rel.Relation = r.Relation.String()
	rel.SubjectType = r.SubjectType.String()
	rel.SubjectId = r.SubjectID.String()
	rel.SubjectRelation = r.SubjectRelation.String()
	return rel
}

func (r *relation) checkReq(trace bool) *dsr.CheckRequest {
	return &dsr.CheckRequest{
		ObjectType:  r.ObjectType.String(),
		ObjectId:    r.ObjectID.String(),
		Relation:    r.Relation.String(),
		SubjectType: r.SubjectType.String(),
		SubjectId:   r.SubjectID.String(),
		Trace:       trace,
	}
}

func (r *relation) graphReq() *dsr.GetGraphRequest {
	return &dsr.GetGraphRequest{
		ObjectType:      r.ObjectType.String(),
		ObjectId:        r.ObjectID.String(),
		Relation:        r.Relation.String(),
		SubjectType:     r.SubjectType.String(),
		SubjectId:       r.SubjectID.String(),
		SubjectRelation: r.SubjectRelation.String(),
		Explain:         true,
		Trace:           true,
	}
}

type RelationsReader []*relation

func NewRelationsReader(rels ...string) RelationsReader {
	return lo.Map(rels, func(rel string, _ int) *relation {
		r := parseRelation(rel)
		if r == nil {
			panic("invalid relation: " + rel)
		}
		return r
	})
}

func (r RelationsReader) GetRelations(req *dsc.Relation, pool graph.MessagePool[dsc.Relation, *dsc.Relation], out *graph.Relations) error {
	ot := model.ObjectName(req.ObjectType)
	oid := model.ObjectID(req.ObjectId)
	rn := model.RelationName(req.Relation)
	st := model.ObjectName(req.SubjectType)
	sid := model.ObjectID(req.SubjectId)
	sr := model.RelationName(req.SubjectRelation)

	matches := lo.Filter(r, func(rel *relation, _ int) bool {
		return (ot == "" || rel.ObjectType == ot) &&
			(oid == "" || rel.ObjectID == oid) &&
			(rn == "" || rel.Relation == rn) &&
			(st == "" || rel.SubjectType == st) &&
			(sid == "" || rel.SubjectID == sid) &&
			(sr == "" || rel.SubjectRelation == sr)
	})

	for _, rel := range matches {
		*out = append(*out, rel.proto(pool))
	}

	return nil
}

type parseTest struct {
	expr     string
	expected [6]string
}

func TestParseRelation(t *testing.T) {
	for _, test := range []parseTest{
		{"obj:oid#rel@subj:sid", [6]string{"obj", "oid", "rel", "subj", "sid", ""}},
		{"obj:oid#rel@subj:sid#srel", [6]string{"obj", "oid", "rel", "subj", "sid", "srel"}},
		{"obj:?#rel@subj:sid", [6]string{"obj", "", "rel", "subj", "sid", ""}},
		{"obj:?#rel@subj:sid#srel", [6]string{"obj", "", "rel", "subj", "sid", "srel"}},
		{"obj:oid#rel@subj:?", [6]string{"obj", "oid", "rel", "subj", "", ""}},
		{"obj:oid#rel@subj:?#srel", [6]string{"obj", "oid", "rel", "subj", "", "srel"}},
	} {
		t.Run(test.expr, func(tt *testing.T) {
			r := parseRelation(test.expr)
			assert.NotNil(tt, r)
			assert.Equal(tt, relation{
				ObjectType:      model.ObjectName(test.expected[0]),
				ObjectID:        model.ObjectID(test.expected[1]),
				Relation:        model.RelationName(test.expected[2]),
				SubjectType:     model.ObjectName(test.expected[3]),
				SubjectID:       model.ObjectID(test.expected[4]),
				SubjectRelation: model.RelationName(test.expected[5]),
			}, *r)
		})
	}
}
