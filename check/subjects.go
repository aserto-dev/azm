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

func NewSubjectSearch(m *model.Model, req *dsr.GetGraphRequest, reader RelationReader) *SubjectSearch {
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
		memo:    newSearchMemo(req.Trace),
		explain: req.Explain,
	}}
}

func (s *SubjectSearch) Search() (*dsr.GetGraphResponse, error) {
	resp := &dsr.GetGraphResponse{}

	if err := s.validate(); err != nil {
		return resp, err
	}

	res, err := s.search(s.params)
	if err != nil {
		return resp, err
	}

	resp.Results = res.Subjects()

	if s.explain {
		resp.Explanation = res.Explain()
	}

	resp.Trace = s.memo.Trace()

	return resp, nil
}

func (s *SubjectSearch) search(params *relation) (searchResults, error) {
	status := s.memo.MarkVisited(params)
	switch status {
	case searchStatusComplete:
		return s.memo.visited[*params], nil
	case searchStatusPending:
		// We have a cycle.
		return nil, nil
	case searchStatusNew:
	}

	o := s.m.Objects[params.ot]

	var (
		results searchResults
		err     error
	)

	if o.HasRelation(params.rel) {
		results, err = s.searchRelation(params)
	} else {
		results, err = s.searchPermission(params)
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
		if rel.SubjectId != "*" && params.oid != ObjectID(rel.ObjectId) {
			continue
		}

		edge := relation{
			ot:  model.ObjectName(rel.ObjectType),
			oid: ObjectID(rel.ObjectId),
			rel: model.RelationName(rel.Relation),
			st:  model.ObjectName(rel.SubjectType),
			sid: ObjectID(rel.SubjectId),
		}

		subj := edge.subject()

		var path []searchPath
		if s.explain {
			path = append(results[*subj], searchPath{&edge}) //nolint: gocritic
		}

		results[*subj] = path
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

			subj := current.subject()

			var path []searchPath
			if s.explain {
				path = append(results[*subj], searchPath{current}) //nolint: gocritic
			}
			results[*subj] = path
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
			res = lo.MapValues(res, func(paths []searchPath, _ object) []searchPath {
				return lo.Map(paths, func(p searchPath, _ int) searchPath {
					return append(searchPath{current}, p...)
				})
			})
		}

		results = lo.Assign(results, res)
	}

	return results, nil
}

func (s *SubjectSearch) searchPermission(params *relation) (searchResults, error) {
	p := s.m.Objects[params.ot].Permissions[params.rel]

	results := searchResults{}

	subjTypes := []model.ObjectName{}
	if params.srel == "" {
		subjTypes = append(subjTypes, params.st)
	} else {
		subjTypes = s.m.Objects[params.st].Relations[params.srel].SubjectTypes
	}

	if len(lo.Intersect(subjTypes, p.SubjectTypes)) == 0 {
		// The subject type cannot have this permission.
		return results, nil
	}

	terms := p.Terms()
	termChecks := make([][]*relation, 0, len(terms))
	for _, pt := range terms {
		// expand arrow operators.
		expanded, err := s.expandTerm(pt, params)
		if err != nil {
			return results, err
		}
		termChecks = append(termChecks, expanded)
	}

	switch {
	case p.IsUnion():
		return s.union(termChecks)
	case p.IsIntersection():
		return s.intersection(termChecks)
	case p.IsExclusion():
		include, err := s.union(termChecks[:1])
		switch {
		case err != nil:
			return results, err
		case include == nil:
			// We have a cycle.
			return nil, nil
		case len(include) == 0:
			// Short-circuit: The include term is false, so the permission is false.
			return results, nil
		}

		exclude, err := s.union(termChecks[1:])
		if err != nil {
			return results, err
		}

		return lo.OmitByKeys(include, lo.Keys(exclude)), nil
	}

	return results, derr.ErrUnknown.Msg("unknown permission operator")
}

func (s *SubjectSearch) expandTerm(pt *model.PermissionTerm, params *relation) ([]*relation, error) {
	if pt.IsArrow() {
		// Resolve the base of the arrow.
		rels, err := s.getRels(&dsc.Relation{
			ObjectType: params.ot.String(),
			ObjectId:   params.oid.String(),
			Relation:   pt.Base.String(),
		})
		if err != nil {
			return []*relation{}, err
		}

		expanded := lo.Map(rels, func(rel *dsc.Relation, _ int) *relation {
			return &relation{
				ot:   model.ObjectName(rel.SubjectType),
				oid:  ObjectID(rel.SubjectId),
				rel:  pt.RelOrPerm,
				st:   params.st,
				sid:  params.sid,
				srel: params.srel,
			}
		})

		return expanded, nil
	}

	return []*relation{{ot: params.ot, oid: params.oid, rel: pt.RelOrPerm, st: params.st, sid: params.sid, srel: params.srel}}, nil
}

func (s *SubjectSearch) union(checks [][]*relation) (searchResults, error) {
	results := searchResults{}
	status := searchStatusPending

	for _, check := range checks {
		var (
			res searchResults
			err error
		)

		switch len(check) {
		case 0:
			res, err = searchResults{}, nil
		case 1:
			res, err = s.search(check[0])
		default:
			res, err = s.union(lo.Map(check, func(params *relation, _ int) []*relation {
				return []*relation{params}
			}))
		}

		switch {
		case err != nil:
			return res, err
		case res == nil:
			// We have a cycle.
			continue
		}

		results = lo.Assign(results, res)
		status = searchStatusComplete
	}

	// return nil if all checks result in a cycle
	return lo.Ternary(status == searchStatusComplete, results, nil), nil
}

func (s *SubjectSearch) intersection(checks [][]*relation) (searchResults, error) {
	results := []searchResults{}
	status := searchStatusPending

	for _, check := range checks {
		// if the base of an arrow operator resolves to multiple objects (e.g. multiple "parents")
		// then a match on any of them is sufficient.
		result, err := s.union([][]*relation{check})
		switch {
		case err != nil:
			return searchResults{}, err
		case result == nil:
			// We have a cycle.
			continue
		case len(result) == 0:
			return result, nil
		}

		status = searchStatusComplete
		results = append(results, result)
	}

	if status == searchStatusPending {
		// All checks result in a cycle.
		return nil, nil
	}

	intersection := lo.Reduce(results, func(agg searchResults, item searchResults, i int) searchResults {
		if i == 0 {
			return item
		}

		for subj, paths := range agg {
			itemPaths, inBoth := item[subj]
			if inBoth {
				// add the paths from the current item to the intersection.
				agg[subj] = append(paths, itemPaths...)
			} else {
				// the subject is not in the intersection.
				delete(agg, subj)
			}
		}

		return agg
	}, searchResults{})

	return intersection, nil
}
