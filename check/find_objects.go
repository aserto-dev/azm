package check

import (
	"sort"

	"github.com/aserto-dev/azm/model"
	dsc "github.com/aserto-dev/go-directory/aserto/directory/common/v3"
	dsr "github.com/aserto-dev/go-directory/aserto/directory/reader/v3"
	"github.com/aserto-dev/go-directory/pkg/derr"
	"github.com/samber/lo"
)

type searchPath []*checkParams

type searchResults map[checkParams][]searchPath

type searchStatus int

const (
	searchStatusUnknown searchStatus = iota
	searchStatusPending
	searchStatusComplete
)

type searchMemo struct {
	visited map[checkParams]searchResults
	history []*checkParams
}

func newSearchMemo(trace bool) *searchMemo {
	return &searchMemo{
		visited: map[checkParams]searchResults{},
		history: lo.Ternary(trace, []*checkParams{}, nil),
	}
}

func (m *searchMemo) MarkVisited(params checkParams) searchStatus {
	results, ok := m.visited[params]
	switch {
	case !ok:
		m.visited[params] = nil
		if m.history != nil {
			m.history = append(m.history, &params)
		}
		return searchStatusUnknown
	case results == nil:
		return searchStatusPending
	default:
		return searchStatusComplete
	}
}

func (m *searchMemo) Status(params checkParams) searchStatus {
	results, ok := m.visited[params]
	switch {
	case !ok:
		return searchStatusUnknown
	case results == nil:
		return searchStatusPending
	default:
		return searchStatusComplete
	}
}

func (m *searchMemo) MarkComplete(params checkParams, results searchResults) {
	m.visited[params] = results
}

type ObjectSearch struct {
	m       *model.Model
	params  *checkParams
	getRels RelationReader

	memo    *searchMemo
	explain bool
}

func NewObjectSearch(m *model.Model, req *dsr.GetGraphRequest, reader RelationReader, explain, trace bool) *ObjectSearch {
	return &ObjectSearch{
		m: m,
		params: &checkParams{
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
	}
}

func (f *ObjectSearch) Search() ([]*dsc.ObjectIdentifier, error) {
	o := f.m.Objects[f.params.ot]
	if o == nil {
		return nil, derr.ErrObjectTypeNotFound.Msg(f.params.ot.String())
	}

	if !o.HasRelOrPerm(f.params.rel) {
		return nil, derr.ErrRelationNotFound.Msg(f.params.rel.String())
	}

	if _, ok := f.m.Objects[f.params.st]; !ok {
		return nil, derr.ErrObjectTypeNotFound.Msg(f.params.st.String())
	}

	res, err := f.search(f.params)
	if err != nil {
		return nil, err
	}

	return lo.MapToSlice(res, func(p checkParams, _ []searchPath) *dsc.ObjectIdentifier {
		return &dsc.ObjectIdentifier{
			ObjectType: p.ot.String(),
			ObjectId:   p.oid.String(),
		}
	}), nil
}

func (f *ObjectSearch) Paths() searchResults {
	return f.memo.visited[*f.params]
}

func (f *ObjectSearch) search(params *checkParams) (searchResults, error) {
	status := f.memo.MarkVisited(*params)
	switch status {
	case searchStatusComplete:
		return f.memo.visited[*params], nil
	case searchStatusPending:
		panic("cycle detected")
	}

	o := f.m.Objects[params.ot]

	var (
		results searchResults
		err     error
	)

	if o.HasRelation(params.rel) {
		results, err = f.searchRelation(params)
	} else {
		results, err = f.findPermission(params)
	}

	f.memo.MarkComplete(*params, results)

	return results, err
}

func (f *ObjectSearch) searchRelation(params *checkParams) (searchResults, error) {
	results := searchResults{}

	r := f.m.Objects[params.ot].Relations[params.rel]
	steps := f.stepRelation(r, params.st)

	for _, step := range steps {
		var (
			res searchResults
			err error
		)
		switch {
		case step.IsDirect(), step.IsWildcard():
			res, err = f.findNeighbor(step, params)
		case step.IsSubject():
			res, err = f.searchSubjectSet(step, params)
		}

		if err != nil {
			return results, err
		}

		results = lo.Assign(results, res)
	}

	return results, nil
}

func (f *ObjectSearch) stepRelation(r *model.Relation, subjs ...model.ObjectName) []*model.RelationRef {
	steps := lo.FilterMap(r.Union, func(rr *model.RelationRef, _ int) (*model.RelationRef, bool) {
		if rr.IsDirect() || rr.IsWildcard() {
			// include direct or wildcard with the expected types.
			return rr, len(subjs) == 0 || lo.Contains(subjs, rr.Object)
		}

		// include subject relations that match or can resolve to the expected types.
		include := len(subjs) == 0 ||
			len(lo.Intersect(f.m.Objects[rr.Object].Relations[rr.Relation].SubjectTypes, subjs)) > 0 ||
			lo.Contains(subjs, rr.Object)

		return rr, include
	})

	sort.Slice(steps, func(i, j int) bool {
		switch {
		case steps[i].Assignment() > steps[j].Assignment():
			// Wildcard < Subject < Direct
			return true
		case steps[i].Assignment() == steps[j].Assignment():
			return steps[i].String() < steps[j].String()
		default:
			return false
		}
	})

	return steps
}

func (f *ObjectSearch) findNeighbor(step *model.RelationRef, params *checkParams) (searchResults, error) {
	sid := params.sid.String()
	if step.IsWildcard() {
		sid = "*"
	}

	req := &dsc.Relation{
		ObjectType:  params.ot.String(),
		Relation:    params.rel.String(),
		SubjectType: step.Object.String(),
		SubjectId:   sid,
	}

	results := searchResults{}

	rels, err := f.getRels(req)
	if err != nil {
		return results, err
	}

	for _, rel := range rels {
		if rel.SubjectId == "*" || params.sid == ObjectID(rel.SubjectId) {
			result := checkParams{
				ot:  model.ObjectName(rel.ObjectType),
				oid: ObjectID(rel.ObjectId),
				rel: model.RelationName(rel.Relation),
				st:  model.ObjectName(rel.SubjectType),
				sid: ObjectID(rel.SubjectId),
			}

			var path []searchPath
			if f.explain {
				path = append(results[result], []*checkParams{&result})
			}

			results[result] = path
		}
	}

	return results, nil
}

func (f *ObjectSearch) searchSubjectSet(step *model.RelationRef, params *checkParams) (searchResults, error) {
	expansion := &checkParams{
		ot:   step.Object,
		rel:  step.Relation,
		st:   f.params.st,
		sid:  f.params.sid,
		srel: f.params.srel,
	}

	subjSet := searchResults{}

	switch f.memo.Status(*expansion) {
	case searchStatusUnknown:
		set, err := f.search(expansion)
		if err != nil {
			return nil, err
		}

		if expansion.rel == expansion.srel {
			self := &checkParams{
				ot:   expansion.ot,
				oid:  expansion.sid,
				rel:  expansion.rel,
				st:   expansion.st,
				sid:  expansion.sid,
				srel: expansion.srel,
			}
			set[*self] = nil
		}

		subjSet = set

	case searchStatusPending:
		// This is a recursive structure.
		// Expand the subject set to find all sets that contain it.
		set, err := f.expandSubjectSet(expansion)
		if err != nil {
			return nil, err
		}

		subjSet = set

	case searchStatusComplete:
		subjSet = f.memo.visited[*expansion]
	}

	if *params == *expansion {
		return subjSet, nil
	}

	results := searchResults{}
	for subj, paths := range subjSet {
		search := &dsc.Relation{
			ObjectType:      params.ot.String(),
			Relation:        params.rel.String(),
			SubjectType:     subj.ot.String(),
			SubjectId:       subj.oid.String(),
			SubjectRelation: params.srel.String(),
		}
		objects, err := f.getRels(search)
		if err != nil {
			return nil, err
		}

		for _, rel := range objects {
			p := checkParams{
				ot:  model.ObjectName(rel.ObjectType),
				oid: ObjectID(rel.ObjectId),
				rel: model.RelationName(rel.Relation),
				st:  model.ObjectName(rel.SubjectType),
				sid: ObjectID(rel.SubjectId),
			}

			var resPaths []searchPath
			if f.explain {
				resPaths = append(results[p], paths...)
			}

			results[p] = resPaths
		}
	}

	return results, nil
}

func (f *ObjectSearch) expandSubjectSet(params *checkParams) (searchResults, error) {
	results := searchResults{}
	backlog := map[checkParams][]*checkParams{*params: nil}

	for len(backlog) > 0 {
		// pop item from backlog
		var (
			cur  checkParams
			path []*checkParams
		)

		for k, v := range backlog {
			cur = k
			path = v
			break
		}

		delete(backlog, cur)

		rels, err := f.getRels(cur.AsRelation())
		if err != nil {
			return nil, err
		}

		for _, rel := range rels {
			result := paramsFromRel(rel)
			if _, ok := results[*result]; ok {
				continue
			}

			step := checkParams{
				ot:   params.ot,
				rel:  params.rel,
				st:   result.ot,
				sid:  result.oid,
				srel: result.rel,
			}
			stepPath := append(path, result)
			backlog[step] = stepPath

			var paths []searchPath
			if f.explain {
				paths = append(results[*result], stepPath)
			}

			results[*result] = paths
		}
	}

	return results, nil
}

func (f *ObjectSearch) findPermission(params *checkParams) (searchResults, error) {
	return nil, nil
}
