package graph

import (
	"strings"

	"github.com/aserto-dev/azm/mempool"
	"github.com/aserto-dev/azm/model"
	dsc "github.com/aserto-dev/go-directory/aserto/directory/common/v3"
	dsr "github.com/aserto-dev/go-directory/aserto/directory/reader/v3"
	"github.com/aserto-dev/go-directory/pkg/derr"
	"github.com/rs/zerolog/log"
	"github.com/samber/lo"
)

type ObjectSearch struct {
	subjectSearch  *SubjectSearch
	wildcardSearch *SubjectSearch
}

func NewObjectSearch(m *model.Model, req *dsr.GetGraphRequest, reader RelationReader, pool *mempool.RelationsPool) (*ObjectSearch, error) {
	params := searchParams(req)
	if err := validate(m, params); err != nil {
		return nil, err
	}

	im := m.Invert()
	// validate the model but skip name validation. To avoid name collisions, the inverted model
	// uses mangled names that are not valid identifiers.
	if err := im.Validate(model.SkipNameValidation, model.AllowPermissionInArrowBase); err != nil {
		log.Err(err).Interface("req", req).Msg("inverted model is invalid")
		// NOTE: we should persist the inverted model instead of computing it on the fly.
		return nil, derr.ErrUnknown.Msg("internal error: unable to search objects.")
	}

	iParams := invertGetGraphRequest(im, req)

	return &ObjectSearch{
		subjectSearch: &SubjectSearch{graphSearch{
			m:       im,
			params:  iParams,
			getRels: invertedRelationReader(im, reader),
			memo:    newSearchMemo(req.Trace),
			explain: req.Explain,
			pool:    pool,
		}},
		wildcardSearch: &SubjectSearch{graphSearch{
			m:       im,
			params:  wildcardParams(iParams),
			getRels: invertedRelationReader(im, reader),
			memo:    newSearchMemo(req.Trace),
			explain: req.Explain,
			pool:    pool,
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

	res = invertResults(res)

	m := s.subjectSearch.m

	memo := s.subjectSearch.memo
	memo.history = append(memo.history, s.wildcardSearch.memo.history...)
	memo.history = lo.Map(memo.history, func(c *searchCall, _ int) *searchCall {
		return &searchCall{
			relation: uninvertRelation(m, c.relation),
			status:   c.status,
		}
	})

	resp.Results = res.Subjects()

	if s.subjectSearch.explain {
		resp.Explanation, _ = res.Explain()
	}

	resp.Trace = memo.Trace()

	return resp, nil
}

func invertGetGraphRequest(im *model.Model, req *dsr.GetGraphRequest) *relation {
	rel := model.InverseRelation(model.ObjectName(req.ObjectType), model.RelationName(req.Relation), model.RelationName(req.SubjectRelation))
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
	return func(r *dsc.RelationIdentifier, relPool MessagePool[*dsc.RelationIdentifier], out *Relations) error {
		ir := uninvertRelation(m, relationFromProto(r))
		if err := reader(ir.asProto(), relPool, out); err != nil {
			return err
		}

		res := *out
		for i, r := range res {
			res[i] = &dsc.RelationIdentifier{
				ObjectType:  r.SubjectType,
				ObjectId:    r.SubjectId,
				Relation:    r.Relation,
				SubjectType: r.ObjectType,
				SubjectId:   r.ObjectId,
			}
		}

		return nil
	}
}

func uninvertRelation(m *model.Model, r *relation) *relation {
	objSplit := strings.SplitN(r.rel.String(), model.ObjectNameSeparator, 2)
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

	return &relation{
		ot:   r.st,
		oid:  r.sid,
		rel:  model.RelationName(rel),
		st:   r.ot,
		sid:  r.oid,
		srel: model.RelationName(srel),
	}
}

func invertResults(res searchResults) searchResults {
	return lo.MapValues(res, func(paths []searchPath, obj object) []searchPath {
		return lo.Map(paths, func(p searchPath, _ int) searchPath {
			return lo.Map(p, func(r *relation, _ int) *relation {
				return &relation{
					ot:   r.st,
					oid:  r.sid,
					rel:  r.rel,
					st:   r.ot,
					sid:  r.oid,
					srel: r.srel,
				}
			})
		})
	})
}
