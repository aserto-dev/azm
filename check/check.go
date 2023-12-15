package check

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

type RelationReader func(*dsc.Relation) ([]*dsc.Relation, error)

type Checker struct {
	m       *model.Model
	params  *checkParams
	getRels RelationReader

	memo checkMemo
}

func New(m *model.Model, req *dsr.CheckRequest, reader RelationReader) *Checker {
	return &Checker{
		m: m,
		params: &checkParams{
			ot:  model.ObjectName(req.ObjectType),
			oid: ObjectID(req.ObjectId),
			rel: model.RelationName(req.Relation),
			st:  model.ObjectName(req.SubjectType),
			sid: ObjectID(req.SubjectId),
		},
		getRels: reader,
		memo:    checkMemo{},
	}
}

func (c *Checker) Check() (bool, error) {
	o := c.m.Objects[c.params.ot]
	if o == nil {
		return false, derr.ErrObjectTypeNotFound.Msg(c.params.ot.String())
	}

	if !o.HasRelOrPerm(c.params.rel) {
		return false, derr.ErrRelationNotFound.Msg(c.params.rel.String())
	}

	return c.check(c.params)
}

type checkStatus int

const (
	checkStatusUnknown checkStatus = iota
	checkStatusPending
	checkStatusTrue
	checkStatusFalse
)

// checkMemo tracks pending checks to detect cycles, and caches the results of completed checks.
type checkMemo map[checkParams]checkStatus

// MarkVisited returns the status of a check. If the check has not been visited, it is marked as pending.
func (m checkMemo) MarkVisited(params *checkParams) checkStatus {
	status, ok := m[*params]
	if !ok {
		m[*params] = checkStatusPending
		status = checkStatusUnknown
	}
	return status
}

// MarkComplete records the result of a check.
func (m checkMemo) MarkComplete(params *checkParams, checkResult bool) {
	status := checkStatusFalse
	if checkResult {
		status = checkStatusTrue
	}
	m[*params] = status
}

type checkParams struct {
	ot  model.ObjectName
	oid ObjectID
	rel model.RelationName
	st  model.ObjectName
	sid ObjectID
}

func (p *checkParams) String() string {
	return fmt.Sprintf("%s:%s#%s@%s:%s", p.ot, p.oid, p.rel, p.st, p.sid)
}

func (c *Checker) check(params *checkParams) (bool, error) {
	prior := c.memo.MarkVisited(params)
	switch prior {
	case checkStatusPending:
		// We have a cycle.
		return false, nil
	case checkStatusTrue, checkStatusFalse:
		// We already checked this relation.
		return prior == checkStatusTrue, nil
	case checkStatusUnknown:
		// this is the first time we're running this check.
	}

	o := c.m.Objects[params.ot]

	var (
		result bool
		err    error
	)
	if o.HasRelation(params.rel) {
		result, err = c.checkRelation(params)
	} else {
		result, err = c.checkPermission(params)
	}

	c.memo.MarkComplete(params, result)

	return result, err
}

func (c *Checker) checkRelation(params *checkParams) (bool, error) {
	r := c.m.Objects[params.ot].Relations[params.rel]
	steps := c.stepRelation(r, params.st)

	for _, step := range steps {
		req := &dsc.Relation{
			ObjectType:  params.ot.String(),
			ObjectId:    params.oid.String(),
			Relation:    params.rel.String(),
			SubjectType: step.Object.String(),
		}

		switch {
		case step.IsWildcard():
			req.SubjectId = "*"
		case step.IsSubject():
			req.SubjectRelation = step.Relation.String()
		}

		rels, err := c.getRels(req)
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

		case step.IsWildcard():
			if len(rels) > 0 {
				// We have a wildcard match.
				return true, nil
			}

		case step.IsSubject():
			for _, rel := range rels {
				if ok, err := c.check(&checkParams{
					ot:  step.Object,
					oid: ObjectID(rel.SubjectId),
					rel: step.Relation,
					st:  params.st,
					sid: params.sid,
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

func (c *Checker) checkPermission(params *checkParams) (bool, error) {
	p := c.m.Objects[params.ot].Permissions[params.rel]

	if !lo.Contains(p.SubjectTypes, params.st) {
		// The subject type cannot have this permission.
		return false, nil
	}

	terms := p.Terms()
	termChecks := make([][]*checkParams, 0, len(terms))
	for _, pt := range terms {
		// expand arrow operators.
		expanded, err := c.expandTerm(pt, params)
		if err != nil {
			return false, err
		}
		termChecks = append(termChecks, expanded)
	}

	switch {
	case p.IsUnion():
		return c.checkAny(termChecks)
	case p.IsIntersection():
		return c.checkAll(termChecks)
	case p.IsExclusion():
		include, err := c.checkAny(termChecks[:1])
		switch {
		case err != nil:
			return false, err
		case !include:
			// Short-circuit: The include term is false, so the permission is false.
			return false, nil
		}

		exclude, err := c.checkAny(termChecks[1:])
		if err != nil {
			return false, err
		}

		return !exclude, nil
	}

	return false, derr.ErrUnknown.Msg("unknown permission operator")
}

func (c *Checker) expandTerm(pt *model.PermissionTerm, params *checkParams) ([]*checkParams, error) {
	if pt.IsArrow() {
		// Resolve the base of the arrow.
		rels, err := c.getRels(&dsc.Relation{
			ObjectType: params.ot.String(),
			ObjectId:   params.oid.String(),
			Relation:   pt.Base.String(),
		})
		if err != nil {
			return []*checkParams{}, err
		}

		expanded := lo.Map(rels, func(rel *dsc.Relation, _ int) *checkParams {
			return &checkParams{
				ot:  model.ObjectName(rel.SubjectType),
				oid: ObjectID(rel.SubjectId),
				rel: pt.RelOrPerm,
				st:  params.st,
				sid: params.sid,
			}
		})

		return expanded, nil
	}

	return []*checkParams{{ot: params.ot, oid: params.oid, rel: pt.RelOrPerm, st: params.st, sid: params.sid}}, nil
}

func (c *Checker) checkAny(checks [][]*checkParams) (bool, error) {
	for _, check := range checks {
		var (
			ok  bool
			err error
		)

		switch len(check) {
		case 0:
			ok, err = false, nil
		case 1:
			ok, err = c.check(check[0])
		default:
			ok, err = c.checkAny(lo.Map(check, func(params *checkParams, _ int) []*checkParams {
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

func (c *Checker) checkAll(checks [][]*checkParams) (bool, error) {
	for _, check := range checks {
		// if the base of an arrow operator resolves to multiple objects (e.g. multiple "parents")
		// then a match on any of them is sufficient.
		ok, err := c.checkAny([][]*checkParams{check})
		if err != nil {
			return false, err
		}
		if !ok {
			return false, nil
		}
	}
	return true, nil
}

func (c *Checker) stepRelation(r *model.Relation, subjs ...model.ObjectName) []*model.RelationRef {
	steps := lo.FilterMap(r.Union, func(rt *model.RelationTerm, _ int) (*model.RelationRef, bool) {
		if rt.IsDirect() || rt.IsWildcard() {
			// include direct or wildcard with the expected types.
			return rt.RelationRef, len(subjs) == 0 || lo.Contains(subjs, rt.Object)
		}

		// include subject relations that can resolve to the expected types.
		return rt.RelationRef, len(subjs) == 0 || len(lo.Intersect(c.m.Objects[rt.Object].Relations[rt.Relation].SubjectTypes, subjs)) > 0
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
