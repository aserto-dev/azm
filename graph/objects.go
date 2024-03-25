package graph

import (
	"strings"

	"github.com/aserto-dev/azm/model"
	dsc "github.com/aserto-dev/go-directory/aserto/directory/common/v3"
	dsr "github.com/aserto-dev/go-directory/aserto/directory/reader/v3"
	"github.com/samber/lo"
)

type ObjectSearch struct {
	subjectSearch  *SubjectSearch
	wildcardSearch *SubjectSearch
}

func NewObjectSearch(m *model.Model, req *dsr.GetGraphRequest, reader RelationReader) (*ObjectSearch, error) {
	params := searchParams(req)
	if err := validate(m, params); err != nil {
		return nil, err
	}

	im := m.Invert()
	if err := im.Validate(); err != nil {
		// TODO: we should persist the inverted model instead of computing it on the fly.
		panic(err)
	}

	iParams := invertGetGraphRequest(im, req)

	return &ObjectSearch{
		subjectSearch: &SubjectSearch{graphSearch{
			m:       im,
			params:  iParams,
			getRels: invertedRelationReader(im, reader),
			memo:    newSearchMemo(req.Trace),
			explain: req.Explain,
		}},
		wildcardSearch: &SubjectSearch{graphSearch{
			m:       im,
			params:  wildcardParams(iParams),
			getRels: invertedRelationReader(im, reader),
			memo:    newSearchMemo(req.Trace),
			explain: req.Explain,
		}},
	}, nil
}

func (s *ObjectSearch) Search() (*dsr.GetGraphResponse, error) {
	resp := &dsr.GetGraphResponse{}

	res, err := s.subjectSearch.search(s.subjectSearch.params)
	if err != nil {
		return resp, err
	}

	wildRes, err := s.wildcardSearch.search(s.wildcardSearch.params)
	if err != nil {
		return resp, err
	}

	for obj, paths := range wildRes {
		res[obj] = append(res[obj], paths...)
	}

	memo := s.subjectSearch.memo
	memo.history = append(memo.history, s.wildcardSearch.memo.history...)

	resp.Results = res.Subjects()

	if s.subjectSearch.explain {
		resp.Explanation = res.Explain()
	}

	resp.Trace = memo.Trace()

	return resp, nil
}

func invertGetGraphRequest(im *model.Model, req *dsr.GetGraphRequest) *relation {
	rel := model.InverseRelation(model.ObjectName(req.ObjectType), model.RelationName(req.Relation))
	relPerm := model.PermForRel(rel)
	if im.Objects[model.ObjectName(req.SubjectType)].HasPermission(relPerm) {
		rel = relPerm
	} else if req.SubjectRelation != "" {
		rel = model.InverseRelation(model.ObjectName(req.ObjectType), model.RelationName(req.Relation), model.RelationName(req.SubjectRelation))
	}

	iReq := &relation{
		ot:  model.ObjectName(req.SubjectType),
		oid: ObjectID(req.SubjectId),
		rel: rel,
		st:  model.ObjectName(req.ObjectType),
		sid: ObjectID(req.ObjectId),
		// TODO: what do we do with subject relations
		// srel: model.RelationName(req.SubjectRelation),
	}

	o := im.Objects[iReq.ot]
	srPerm := model.GeneratedPermissionPrefix + iReq.rel
	if o.HasRelation(iReq.rel) && o.HasPermission(srPerm) {
		iReq.rel = srPerm
	}

	return iReq
}

func wildcardParams(params *relation) *relation {
	wildcard := *params
	wildcard.oid = "*"
	return &wildcard
}

func invertedRelationReader(m *model.Model, reader RelationReader) RelationReader {
	return func(r *dsc.Relation) ([]*dsc.Relation, error) {
		objSplit := strings.SplitN(r.Relation, model.ObjectNameSeparator, 2)
		obj := model.ObjectName(objSplit[0])
		relSplit := strings.SplitN(objSplit[1], model.SubjectRelationSeparator, 2)
		rel := relSplit[0]
		srel := ""
		if len(relSplit) > 1 {
			srel = relSplit[1]
		}

		perm := model.PermForRel(model.RelationName(rel))
		if m.Objects[obj].HasPermission(perm) {
			rel = perm.String()
		}

		rr := &dsc.Relation{
			ObjectType:      r.SubjectType,
			ObjectId:        r.SubjectId,
			Relation:        rel,
			SubjectType:     r.ObjectType,
			SubjectId:       r.ObjectId,
			SubjectRelation: srel,
		}

		res, err := reader(rr)
		if err != nil {
			return nil, err
		}

		return lo.Map(res, func(r *dsc.Relation, _ int) *dsc.Relation {
			return &dsc.Relation{
				ObjectType:  r.SubjectType,
				ObjectId:    r.SubjectId,
				Relation:    r.Relation,
				SubjectType: r.ObjectType,
				SubjectId:   r.ObjectId,
			}
		}), nil
	}
}
