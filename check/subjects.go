package check

import (
	"github.com/aserto-dev/azm/model"
	dsc "github.com/aserto-dev/go-directory/aserto/directory/common/v3"
	dsr "github.com/aserto-dev/go-directory/aserto/directory/reader/v3"
	"github.com/aserto-dev/go-directory/pkg/derr"
	"github.com/samber/lo"
)

type SubjectSearch struct {
	graphSearch
}

func NewSubjectSearch(m *model.Model, req *dsr.GetGraphRequest, reader RelationReader, explain, trace bool) *SubjectSearch {
	return &SubjectSearch{graphSearch{
		m: m,
		params: &relation{
			ot:   model.ObjectName(req.ObjectType),
			oid:  ObjectID(req.ObjectId),
			rel:  model.RelationName(req.Relation),
			st:   model.ObjectName(req.SubjectType),
			sid:  ObjectID(req.SubjectId),
			srel: model.RelationName(req.SubjectRelation),
		},
		getRels: reader,
		memo:    newSearchMemo(trace),
		explain: explain,
	}}
}

func (s *SubjectSearch) Search() ([]*dsc.ObjectIdentifier, error) {
	if err := s.validate(); err != nil {
		return nil, err
	}

	res, err := s.search(s.params)
	if err != nil {
		return nil, err
	}

	return res.Subjects(), nil
}

func (s *SubjectSearch) search(params *relation) (searchResults, error) {
	status := s.memo.MarkVisited(params)
	switch status {
	case searchStatusComplete:
		return s.memo.visited[*params], nil
	case searchStatusPending:
		return nil, derr.ErrCycleDetected
	case searchStatusUnknown:
	}

	o := s.m.Objects[params.ot]

	var (
		results searchResults
		err     error
	)

	if o.HasRelation(params.rel) {
		results, err = s.searchRelation(params)
	} else {
		return nil, derr.ErrNotImplemented.Msg("search on permissions is not supported")
	}

	s.memo.MarkComplete(params, results)

	return results, err
}

func (s *SubjectSearch) searchRelation(params *relation) (searchResults, error) {
	r := s.m.Objects[params.ot].Relations[params.rel]
	steps := s.m.StepRelation(r, params.st)

	results := searchResults{}

	for _, step := range steps {
		var (
			res searchResults
			err error
		)
		switch {
		case step.IsDirect(), step.IsWildcard():
			res, err = s.findNeighbor(step, params)
		case step.IsSubject():
			res, err = s.searchSubjectRelation(step, params)
		}

		if err != nil {
			return results, err
		}

		results = lo.Assign(results, res)
	}

	return results, nil
}

func (s *SubjectSearch) findNeighbor(step *model.RelationRef, params *relation) (searchResults, error) {
	sid := params.sid.String()
	if step.IsWildcard() {
		sid = "*"
	}

	req := &dsc.Relation{
		ObjectType:  params.ot.String(),
		ObjectId:    params.oid.String(),
		Relation:    params.rel.String(),
		SubjectType: step.Object.String(),
		SubjectId:   sid,
	}

	results := searchResults{}

	rels, err := s.getRels(req)
	if err != nil {
		return results, err
	}

	for _, rel := range rels {
		if rel.SubjectId == "*" || params.oid == ObjectID(rel.ObjectId) {
			result := relation{
				ot:  model.ObjectName(rel.ObjectType),
				oid: ObjectID(rel.ObjectId),
				rel: model.RelationName(rel.Relation),
				st:  model.ObjectName(rel.SubjectType),
				sid: ObjectID(rel.SubjectId),
			}

			var path []searchPath
			if s.explain {
				path = append(results[result], searchPath{&result}) //nolint: gocritic
			}

			results[result] = path
		}
	}

	return results, nil
}

func (s *SubjectSearch) searchSubjectRelation(step *model.RelationRef, params *relation) (searchResults, error) {
	results := searchResults{}

	req := &dsc.Relation{
		ObjectType:      params.ot.String(),
		ObjectId:        params.oid.String(),
		Relation:        params.rel.String(),
		SubjectType:     step.Object.String(),
		SubjectRelation: step.Relation.String(),
	}
	rels, err := s.getRels(req)
	if err != nil {
		return results, err
	}

	for _, rel := range rels {
		current := relationFromProto(rel)

		if params.srel == model.RelationName(rel.SubjectRelation) && params.st == model.ObjectName(rel.SubjectType) {
			// We're searching for a subject relation (not a Check call) and we have a match.

			var path []searchPath
			if s.explain {
				path = append(results[*current], searchPath{current}) //nolint: gocritic
			}
			results[*current] = path
		}

		check := &relation{
			ot:   step.Object,
			oid:  ObjectID(rel.SubjectId),
			rel:  step.Relation,
			st:   params.st,
			sid:  params.sid,
			srel: params.srel,
		}

		res, err := s.search(check)
		if err != nil {
			return results, err
		}

		if s.explain {
			res = lo.MapValues(res, func(paths []searchPath, _ relation) []searchPath {
				return lo.Map(paths, func(p searchPath, _ int) searchPath {
					return append(searchPath{current}, p...)
				})
			})
		}

		results = lo.Assign(results, res)
	}

	return results, nil
}
