package walk

import (
	"fmt"
	"sort"
	"strings"

	"github.com/aserto-dev/azm/model"
	dsc "github.com/aserto-dev/go-directory/aserto/directory/common/v3"
	dsr "github.com/aserto-dev/go-directory/aserto/directory/reader/v3"
	"github.com/aserto-dev/go-directory/pkg/derr"

	"github.com/samber/lo"
)

type ObjectID string

func (id ObjectID) String() string {
	return string(id)
}

type Object struct {
	Type model.ObjectName
	ID   ObjectID
}

type RelationReader func(*dsc.Relation) ([]*dsc.Relation, error)

type Walker struct {
	m       *model.Model
	obj     *Object
	rel     model.RelationName
	subj    *Object
	getRels RelationReader

	memo checkMemo
}

func New(m *model.Model, req *dsr.CheckRequest, reader RelationReader) *Walker {
	return &Walker{
		m:       m,
		obj:     &Object{Type: model.ObjectName(req.ObjectType), ID: ObjectID(req.ObjectId)},
		rel:     model.RelationName(req.Relation),
		subj:    &Object{Type: model.ObjectName(req.SubjectType), ID: ObjectID(req.SubjectId)},
		getRels: reader,
		memo:    checkMemo{},
	}
}

func (w *Walker) Check() (bool, error) {
	o := w.m.Objects[w.obj.Type]
	if o == nil {
		return false, derr.ErrObjectTypeNotFound.Msg(w.obj.Type.String())
	}

	if !o.HasRelOrPerm(w.rel) {
		return false, derr.ErrRelationNotFound.Msg(w.rel.String())
	}

	params := &checkParams{
		ot:       w.obj.Type,
		oid:      w.obj.ID,
		relation: w.rel,
		st:       w.subj.Type,
		sid:      w.subj.ID,
	}

	return w.check(params)
}

type checkStatus int

const (
	checkStatusUnknown checkStatus = iota
	checkStatusPending
	checkStatusTrue
	checkStatusFalse
)

type checkMemo map[checkParams]checkStatus

type checkParams struct {
	ot       model.ObjectName
	oid      ObjectID
	relation model.RelationName
	st       model.ObjectName
	sid      ObjectID
}

func (p *checkParams) String() string {
	return fmt.Sprintf("%s:%s#%s@%s:%s", p.ot, p.oid, p.relation, p.st, p.sid)
}

func (w *Walker) check(params *checkParams) (bool, error) {
	prior := w.markVisited(params)
	switch prior {
	case checkStatusPending:
		// We have a cycle.
		return false, nil
	case checkStatusTrue, checkStatusFalse:
		// We already checked this relation.
		return prior == checkStatusTrue, nil
	}

	o := w.m.Objects[params.ot]

	var (
		result bool
		err    error
	)
	if o.HasRelation(params.relation) {
		result, err = w.checkRelation(params)
	} else {
		result, err = w.checkPermission(params)
	}

	status := checkStatusFalse
	if err == nil && result {
		status = checkStatusTrue
	}
	w.markComplete(params, status)

	return result, err
}

func (w *Walker) markVisited(params *checkParams) checkStatus {
	status, ok := w.memo[*params]
	if !ok {
		w.memo[*params] = checkStatusPending
		status = checkStatusUnknown
	}
	return status
}

func (w *Walker) markComplete(params *checkParams, status checkStatus) {
	w.memo[*params] = status
}

func (w *Walker) checkRelation(params *checkParams) (bool, error) {
	r := w.m.Objects[params.ot].Relations[params.relation]
	steps := w.stepRelation(r, params.st)

	for _, step := range steps {
		if step.IsWildcard() && step.Object == params.st {
			// We have a wildcard match.
			return true, nil
		}

		rels, err := w.getRels(&dsc.Relation{
			ObjectType:      params.ot.String(),
			ObjectId:        params.oid.String(),
			Relation:        params.relation.String(),
			SubjectType:     step.Object.String(),
			SubjectRelation: step.Relation.String(),
		})
		if err != nil {
			return false, err
		}
		switch {
		case step.IsDirect():
			for _, rel := range rels {
				if rel.SubjectId == params.sid.String() {
					return true, nil
				}
			}
		case step.IsSubject():
			for _, rel := range rels {
				if ok, err := w.check(&checkParams{
					ot:       step.Object,
					oid:      ObjectID(rel.SubjectId),
					relation: step.Relation,
					st:       params.st,
					sid:      params.sid,
				}); err != nil {
					return false, err
				} else if ok {
					return true, nil
				}
			}
		}
	}

	return false, nil
}

func (w *Walker) checkPermission(params *checkParams) (bool, error) {
	p := w.m.Objects[params.ot].Permissions[params.relation]

	if !lo.Contains(p.SubjectTypes, params.st) {
		// The subject type cannot have this permission.
		return false, nil
	}

	terms := p.Terms()
	termChecks := make([][]*checkParams, 0, len(terms))
	for _, pt := range terms {
		// expand arrow operators.
		expanded, err := w.expandTerm(pt, params)
		if err != nil {
			return false, err
		}
		termChecks = append(termChecks, expanded)
	}

	switch {
	case p.IsUnion():
		return w.checkAny(termChecks)
	case p.IsIntersection():
		return w.checkAll(termChecks)
	case p.IsExclusion():
		include, err := w.checkAny(termChecks[:1])
		switch {
		case err != nil:
			return false, err
		case !include:
			// Short-circuit: The include term is false, so the permission is false.
			return false, nil
		}

		exclude, err := w.checkAny(termChecks[1:])
		if err != nil {
			return false, err
		}

		return !exclude, nil
	}

	return false, derr.ErrUnknown.Msg("unknown permission operator")
}

func (w *Walker) expandTerm(pt *model.PermissionTerm, params *checkParams) ([]*checkParams, error) {
	if pt.IsArrow() {
		// Resolve the base of the arrow.
		rels, err := w.getRels(&dsc.Relation{
			ObjectType: params.ot.String(),
			ObjectId:   params.oid.String(),
			Relation:   pt.Base.String(),
		})
		if err != nil {
			return []*checkParams{}, err
		}

		expanded := lo.Map(rels, func(rel *dsc.Relation, _ int) *checkParams {
			return &checkParams{
				ot:       model.ObjectName(rel.SubjectType),
				oid:      ObjectID(rel.SubjectId),
				relation: pt.RelOrPerm,
				st:       params.st,
				sid:      params.sid,
			}
		})

		return expanded, nil
	}

	return []*checkParams{{ot: params.ot, oid: params.oid, relation: pt.RelOrPerm, st: params.st, sid: params.sid}}, nil
}

func (w *Walker) checkAny(checks [][]*checkParams) (bool, error) {
	for _, check := range checks {
		var (
			ok  bool
			err error
		)

		switch len(check) {
		case 0:
			ok, err = false, nil
		case 1:
			ok, err = w.check(check[0])
		default:
			ok, err = w.checkAny(lo.Map(check, func(params *checkParams, _ int) []*checkParams {
				return []*checkParams{params}
			}))
		}

		if err != nil {
			return false, err
		}

		fmt.Printf(
			"checkAny: %s - %v\n",
			strings.Join(
				lo.Map(check, func(p *checkParams, _ int) string {
					return p.String()
				}),
				", ",
			), ok)
		if ok {
			return true, nil
		}
	}

	return false, nil
}

func (w *Walker) checkAll(checks [][]*checkParams) (bool, error) {
	for _, check := range checks {
		// if the base of an arrow operator resolves to multiple objects (e.g. multiple "parents")
		// then a match on any of them is sufficient.
		ok, err := w.checkAny([][]*checkParams{check})
		if err != nil {
			return false, err
		}
		if !ok {
			return false, nil
		}
	}
	return true, nil
}

func (w *Walker) stepRelation(r *model.Relation, subjs ...model.ObjectName) []*model.RelationRef {
	steps := lo.FilterMap(r.Union, func(rt *model.RelationTerm, _ int) (*model.RelationRef, bool) {
		if rt.IsDirect() || rt.IsWildcard() {
			// include direct or wildcard with the expected types.
			return rt.RelationRef, len(subjs) == 0 || lo.Contains(subjs, rt.Object)
		}

		// include subject relations that can resolve to the expected types.
		return rt.RelationRef, len(subjs) == 0 || len(lo.Intersect(w.m.Objects[rt.Object].Relations[rt.Relation].SubjectTypes, subjs)) > 0
	})

	sort.Slice(steps, func(i, j int) bool {
		switch {
		case steps[i].Assignment() > steps[j].Assignment():
			// Wildcard < Subjetc < Direct
			return true
		case steps[i].Assignment() == steps[j].Assignment():
			return steps[i].String() < steps[j].String()
		default:
			return false
		}
	})

	return steps
}
