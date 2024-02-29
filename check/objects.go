package check

import (
	"github.com/aserto-dev/azm/model"
	dsc "github.com/aserto-dev/go-directory/aserto/directory/common/v3"
	dsr "github.com/aserto-dev/go-directory/aserto/directory/reader/v3"
	"github.com/aserto-dev/go-directory/pkg/derr"
	"github.com/samber/lo"
)

type ObjectSearch struct {
	graphSearch
}

func NewObjectSearch(m *model.Model, req *dsr.GetGraphRequest, reader RelationReader) *ObjectSearch {
	return &ObjectSearch{graphSearch{
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

func (s *ObjectSearch) Search() (*dsr.GetGraphResponse, error) {
	resp := &dsr.GetGraphResponse{}

	if err := s.validate(); err != nil {
		return resp, err
	}

	res, err := s.search(s.params)
	if err != nil {
		return resp, err
	}

	resp.Results = res.Objects()

	if s.explain {
		resp.Explanation = res.Explain()
	}

	resp.Trace = s.memo.Trace()

	return resp, nil
}

func (s *ObjectSearch) search(params *relation) (searchResults, error) {
	status := s.memo.MarkVisited(params)
	switch status {
	case searchStatusComplete:
		return s.memo.visited[*params], nil
	case searchStatusPending:
		return nil, derr.ErrCycleDetected
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
		return nil, derr.ErrNotImplemented.Msg("search on permissions is not supported")
	}

	s.memo.MarkComplete(params, results)

	return results, err
}

func (s *ObjectSearch) searchRelation(params *relation) (searchResults, error) {
	results := searchResults{}

	r := s.m.Objects[params.ot].Relations[params.rel]
	steps := s.m.StepRelation(r, params.st)

	for _, step := range steps {
		var (
			res searchResults
			err error
		)
		switch {
		case step.IsDirect(), step.IsWildcard():
			res, err = s.findNeighbor(step, params)
		case step.IsSubject():
			res, err = s.searchSubjectSet(step, params)
		}

		if err != nil {
			return results, err
		}

		results = lo.Assign(results, res)
	}

	return results, nil
}

func (s *ObjectSearch) findNeighbor(step *model.RelationRef, params *relation) (searchResults, error) {
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

	rels, err := s.getRels(req)
	if err != nil {
		return results, err
	}

	for _, rel := range rels {
		if rel.SubjectId != "*" && params.sid != ObjectID(rel.SubjectId) {
			continue
		}

		edge := relation{
			ot:  model.ObjectName(rel.ObjectType),
			oid: ObjectID(rel.ObjectId),
			rel: model.RelationName(rel.Relation),
			st:  model.ObjectName(rel.SubjectType),
			sid: ObjectID(rel.SubjectId),
		}

		obj := edge.object()

		var path []searchPath
		if s.explain {
			path = append(results[*obj], searchPath{&edge}) //nolint: gocritic
		}

		results[*obj] = path
	}

	return results, nil
}

func (s *ObjectSearch) searchSubjectSet(step *model.RelationRef, params *relation) (searchResults, error) {
	expansion := &relation{
		ot:   step.Object,
		rel:  step.Relation,
		st:   s.params.st,
		sid:  s.params.sid,
		srel: s.params.srel,
	}

	subjSet := searchResults{}

	switch s.memo.Status(expansion) {
	case searchStatusNew:
		set, err := s.search(expansion)
		if err != nil {
			return nil, err
		}

		if expansion.rel == expansion.srel {
			set[object{expansion.ot, expansion.sid}] = nil
		}

		subjSet = set

	case searchStatusPending:
		// This is a recursive structure.
		// Expand the subject set to find all sets that contain it.
		set, err := s.expandSubjectSet(expansion)
		if err != nil {
			return nil, err
		}

		subjSet = set

	case searchStatusComplete:
		subjSet = s.memo.visited[*expansion]
	}

	if *params == *expansion {
		return subjSet, nil
	}

	results := searchResults{}
	for subj, paths := range subjSet {
		search := &dsc.Relation{
			ObjectType:      params.ot.String(),
			Relation:        params.rel.String(),
			SubjectType:     subj.Type.String(),
			SubjectId:       subj.ID.String(),
			SubjectRelation: params.srel.String(),
		}
		objects, err := s.getRels(search)
		if err != nil {
			return nil, err
		}

		for _, rel := range objects {
			edge := relation{
				ot:  model.ObjectName(rel.ObjectType),
				oid: ObjectID(rel.ObjectId),
				rel: model.RelationName(rel.Relation),
				st:  model.ObjectName(rel.SubjectType),
				sid: ObjectID(rel.SubjectId),
			}

			obj := edge.object()

			var resPaths []searchPath
			if s.explain {
				resPaths = append(results[*obj], paths...) //nolint: gocritic
			}

			results[*obj] = resPaths
		}
	}

	return results, nil
}

func (s *ObjectSearch) expandSubjectSet(params *relation) (searchResults, error) {
	results := searchResults{}
	backlog := map[relation]searchPath{*params: nil}

	for len(backlog) > 0 {
		// pop item from backlog
		var (
			cur  relation
			path searchPath
		)

		for k, v := range backlog {
			cur = k
			path = v
			break
		}

		delete(backlog, cur)

		rels, err := s.getRels(cur.toProto())
		if err != nil {
			return nil, err
		}

		for _, rel := range rels {
			result := relationFromProto(rel)
			obj := result.object()
			if _, ok := results[*obj]; ok {
				continue
			}

			step := relation{
				ot:   params.ot,
				rel:  params.rel,
				st:   result.ot,
				sid:  result.oid,
				srel: result.rel,
			}
			stepPath := append(path, result) //nolint: gocritic
			backlog[step] = stepPath

			var paths []searchPath
			if s.explain {
				paths = append(results[*obj], stepPath) //nolint: gocritic
			}

			results[*obj] = paths
		}
	}

	return results, nil
}

// func (f *ObjectSearch) searchPermission(params *relation) (searchResults, error) {
//     p := f.m.Objects[params.ot].Permissions[params.rel]

//     // Check if the subject type can have this permission.
//     subjTypes := []model.ObjectName{}
//     if params.srel == "" {
//         subjTypes = append(subjTypes, params.st)
//     } else {
//         subjTypes = f.m.Objects[params.st].Relations[params.srel].SubjectTypes
//     }
//     if len(lo.Intersect(subjTypes, p.SubjectTypes)) == 0 {
//         // The subject type cannot have this permission.
//         return searchResults{}, nil
//     }

//     termResults := []searchResults{}
//     terms := p.Terms()
//     sort.SliceStable(terms, func(i, j int) bool {
//         return !terms[i].IsArrow() && terms[j].IsArrow()
//     })

//     for _, term := range terms {
//         res, err := f.expandPermissionTerms(params, term)
//         if err != nil {
//             return searchResults{}, err
//         }

//         termResults = append(termResults, res)
//     }

//     switch {
//     case p.IsUnion():
//         return lo.Reduce(termResults, func(agg searchResults, res searchResults, _ int) searchResults {
//             for rel, paths := range res {
//                 agg[rel] = append(agg[rel], paths...)
//             }
//             return agg
//         }, searchResults{}), nil
//     case p.IsIntersection():
//     case p.IsExclusion():

//     }
//     return nil, nil
// }

// func (f *ObjectSearch) expandPermissionTerms(params *relation, term *model.PermissionTerm) (searchResults, error) {
//     if term.IsArrow() {
//         return f.expandArrow(params, term)
//     }

//     search := &relation{
//         ot:   params.ot,
//         oid:  params.oid,
//         rel:  term.RelOrPerm,
//         st:   params.st,
//         sid:  params.sid,
//         srel: params.srel,
//     }
//     return f.search(search)
// }

// func (f *ObjectSearch) expandArrow(params *relation, pt *model.PermissionTerm) (searchResults, error) {
//     baseRel := f.m.Objects[params.ot].Relations[pt.Base]
//     results := searchResults{}

//     for _, baseRR := range baseRel.Union {
//         arrowSearch := &relation{
//             ot:   baseRR.Object,
//             rel:  pt.RelOrPerm,
//             st:   params.st,
//             sid:  params.sid,
//             srel: params.srel,
//         }

//         arrowResults, err := f.search(arrowSearch)
//         if err != nil {
//             return nil, err
//         }

//         for arrowResult, arrowPaths := range arrowResults {
//             baseSearch := &relation{
//                 ot:  params.ot,
//                 rel: pt.Base,
//                 st:  arrowResult.ot,
//                 sid: arrowResult.oid,
//             }
//             baseResults, err := f.search(baseSearch)
//             if err != nil {
//                 return nil, err
//             }

//             for baseResult, basePaths := range baseResults {
//                 result := relation{
//                     ot:   baseResult.ot,
//                     oid:  baseResult.oid,
//                     rel:  params.rel,
//                     st:   params.st,
//                     sid:  params.sid,
//                     srel: params.srel,
//                 }
//                 var paths []searchPath
//                 if f.explain {
//                     paths = append(results[result], append(arrowPaths, basePaths...)...) //nolint: gocritic
//                 }

//                 results[result] = paths
//             }
//         }
//     }

//     return results, nil
// }
