package check

import (
	"github.com/aserto-dev/azm/model"
	dsc "github.com/aserto-dev/go-directory/aserto/directory/common/v3"
	dsr "github.com/aserto-dev/go-directory/aserto/directory/reader/v3"
	"github.com/aserto-dev/go-directory/pkg/derr"

	"github.com/samber/lo"
)

type Checker struct {
	m       *model.Model
	params  *relation
	getRels RelationReader

	memo *checkMemo
}

func New(m *model.Model, req *dsr.CheckRequest, reader RelationReader) *Checker {
	return &Checker{
		m: m,
		params: &relation{
			ot:  model.ObjectName(req.ObjectType),
			oid: ObjectID(req.ObjectId),
			rel: model.RelationName(req.Relation),
			st:  model.ObjectName(req.SubjectType),
			sid: ObjectID(req.SubjectId),
		},
		getRels: reader,
		memo:    newCheckMemo(req.Trace),
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

func (c *Checker) Trace() []string {
	return c.memo.Trace()
}

func (c *Checker) check(params *relation) (bool, error) {
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

func (c *Checker) checkRelation(params *relation) (bool, error) {
	r := c.m.Objects[params.ot].Relations[params.rel]
	steps := c.m.StepRelation(r, params.st)

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
				if ok, err := c.check(&relation{
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

func (c *Checker) checkPermission(params *relation) (bool, error) {
	p := c.m.Objects[params.ot].Permissions[params.rel]

	if !lo.Contains(p.SubjectTypes, params.st) {
		// The subject type cannot have this permission.
		return false, nil
	}

	terms := p.Terms()
	termChecks := make([]relations, 0, len(terms))
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

func (c *Checker) expandTerm(pt *model.PermissionTerm, params *relation) (relations, error) {
	if pt.IsArrow() {
		// Resolve the base of the arrow.
		rels, err := c.getRels(&dsc.Relation{
			ObjectType: params.ot.String(),
			ObjectId:   params.oid.String(),
			Relation:   pt.Base.String(),
		})
		if err != nil {
			return relations{}, err
		}

		expanded := lo.Map(rels, func(rel *dsc.Relation, _ int) *relation {
			return &relation{
				ot:  model.ObjectName(rel.SubjectType),
				oid: ObjectID(rel.SubjectId),
				rel: pt.RelOrPerm,
				st:  params.st,
				sid: params.sid,
			}
		})

		return expanded, nil
	}

	return relations{{ot: params.ot, oid: params.oid, rel: pt.RelOrPerm, st: params.st, sid: params.sid}}, nil
}

func (c *Checker) checkAny(checks []relations) (bool, error) {
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
			ok, err = c.checkAny(lo.Map(check, func(params *relation, _ int) relations {
				return relations{params}
			}))
		}

		if err != nil {
			return false, err
		}

		if ok {
			return true, nil
		}
	}

	return false, nil
}

func (c *Checker) checkAll(checks []relations) (bool, error) {
	for _, check := range checks {
		// if the base of an arrow operator resolves to multiple objects (e.g. multiple "parents")
		// then a match on any of them is sufficient.
		ok, err := c.checkAny([]relations{check})
		if err != nil {
			return false, err
		}
		if !ok {
			return false, nil
		}
	}
	return true, nil
}
