package check

import (
	"fmt"
	"sort"

	"github.com/aserto-dev/azm/model"
	dsc "github.com/aserto-dev/go-directory/aserto/directory/common/v3"
	dsr "github.com/aserto-dev/go-directory/aserto/directory/reader/v3"
	"github.com/aserto-dev/go-directory/pkg/derr"

	"github.com/samber/lo"
)

type ObjectID = model.ObjectID

type RelationReader func(*dsc.Relation) ([]*dsc.Relation, error)

type Checker struct {
	m       *model.Model
	params  *CheckParams
	getRels RelationReader

	memo *checkMemo
}

func New(m *model.Model, req *dsr.CheckRequest, reader RelationReader) *Checker {
	return &Checker{
		m: m,
		params: &CheckParams{
			OT:  model.ObjectName(req.ObjectType),
			OID: ObjectID(req.ObjectId),
			Rel: model.RelationName(req.Relation),
			ST:  model.ObjectName(req.SubjectType),
			SID: ObjectID(req.SubjectId),
		},
		getRels: reader,
		memo:    newCheckMemo(req.Trace),
	}
}

func NewGraph(m *model.Model, req *dsr.GetGraphRequest, reader RelationReader) *Checker {
	return &Checker{
		m: m,
		params: &CheckParams{
			OT:   model.ObjectName(req.ObjectType),
			OID:  ObjectID(req.ObjectId),
			Rel:  model.RelationName(req.Relation),
			ST:   model.ObjectName(req.SubjectType),
			SID:  ObjectID(req.SubjectId),
			SRel: model.RelationName(req.SubjectRelation),
		},
		getRels: reader,
		memo:    newCheckMemo(false),
	}
}

func (c *Checker) Check() (bool, error) {
	o := c.m.Objects[c.params.OT]
	if o == nil {
		return false, derr.ErrObjectTypeNotFound.Msg(c.params.OT.String())
	}

	if !o.HasRelOrPerm(c.params.Rel) {
		return false, derr.ErrRelationNotFound.Msg(c.params.Rel.String())
	}

	results, err := c.check(c.params, firstMatch)
	if err != nil {
		return false, err
	}

	return results.status() == checkStatusTrue, nil
}

func (c *Checker) Search() (CheckResults, error) {
	o := c.m.Objects[c.params.OT]
	if o == nil {
		return nil, derr.ErrObjectTypeNotFound.Msg(c.params.OT.String())
	}

	checks := []*CheckParams{}

	switch c.params.Rel {
	case "":
		for _, rel := range lo.Union(lo.Keys(o.Relations), lo.Keys(o.Permissions)) {
			check := *c.params
			check.Rel = rel
			checks = append(checks, &check)
		}

	default:
		if !o.HasRelOrPerm(c.params.Rel) {
			return nil, derr.ErrRelationNotFound.Msg(c.params.Rel.String())
		}

		checks = append(checks, c.params)
	}

	results := CheckResults{}
	for _, check := range checks {
		res, err := c.check(check, allPaths)
		if err != nil {
			return nil, err
		}

		results = results.append(res)
	}

	return results, nil
}

func (c *Checker) Trace() []string {
	return c.memo.Trace()
}

type CheckParams struct {
	OT   model.ObjectName
	OID  ObjectID
	Rel  model.RelationName
	ST   model.ObjectName
	SID  ObjectID
	SRel model.RelationName
}

func (p *CheckParams) String() string {
	return fmt.Sprintf("%s:%s#%s@%s:%s", p.OT, p.OID, p.Rel, p.ST, p.SID)
}

func (p *CheckParams) IsMatch(rel *dsc.Relation) bool {
	return (p.SID == "" || p.SID == ObjectID(rel.SubjectId)) &&
		(p.OID == "" || p.OID == ObjectID(rel.ObjectId))
}

func (p *CheckParams) AsRelation() *dsc.Relation {
	return &dsc.Relation{
		ObjectType:      p.OT.String(),
		ObjectId:        p.OID.String(),
		Relation:        p.Rel.String(),
		SubjectType:     p.ST.String(),
		SubjectId:       p.SID.String(),
		SubjectRelation: p.SRel.String(),
	}
}

func (c *Checker) check(params *CheckParams, paths checkPaths) (CheckResults, error) {
	prior := c.memo.MarkVisited(params)
	switch prior {
	case checkStatusPending:
		// We have a cycle.
		return CheckResults{}, nil
	case checkStatusTrue, checkStatusFalse:
		// We already checked this relation.
		return c.memo.Results(params), nil
	case checkStatusUnknown:
		// this is the first time we're running this check.
	}

	o := c.m.Objects[params.OT]

	var (
		results CheckResults
		err     error
	)
	if o.HasRelation(params.Rel) {
		results, err = c.checkRelation(params, paths)
	} else {
		results, err = c.checkPermission(params, paths)
	}

	c.memo.MarkComplete(params, results)

	return results, err
}

type checkPaths bool

const (
	firstMatch checkPaths = false
	allPaths   checkPaths = true
)

func (c *Checker) checkRelation(params *CheckParams, paths checkPaths) (CheckResults, error) {
	r := c.m.Objects[params.OT].Relations[params.Rel]
	steps := c.stepRelation(r, params.ST)

	results := CheckResults{}

	for _, step := range steps {
		var (
			res CheckResults
			err error
		)
		switch {
		case step.IsDirect():
			res, err = c.checkDirect(step, params, paths)
		case step.IsWildcard():
			res, err = c.checkWildcard(step, params, paths)
		case step.IsSubject():
			res, err = c.checkSubject(step, params, paths)
		}

		if err != nil {
			return results, err
		}

		if len(res) > 0 {
			results = results.append(res)
			if paths == firstMatch {
				break
			}
		}
	}

	return results, nil
}

func (c *Checker) stepRelation(r *model.Relation, subjs ...model.ObjectName) []*model.RelationRef {
	steps := lo.FilterMap(r.Union, func(rr *model.RelationRef, _ int) (*model.RelationRef, bool) {
		if rr.IsDirect() || rr.IsWildcard() {
			// include direct or wildcard with the expected types.
			return rr, len(subjs) == 0 || lo.Contains(subjs, rr.Object)
		}

		// include subject relations that match or can resolve to the expected types.
		include := len(subjs) == 0 ||
			len(lo.Intersect(c.m.Objects[rr.Object].Relations[rr.Relation].SubjectTypes, subjs)) > 0 ||
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

func (c *Checker) checkDirect(step *model.RelationRef, params *CheckParams, paths checkPaths) (CheckResults, error) {
	req := &dsc.Relation{
		ObjectType:  params.OT.String(),
		ObjectId:    params.OID.String(),
		Relation:    params.Rel.String(),
		SubjectType: step.Object.String(),
		SubjectId:   params.SID.String(),
	}

	results := CheckResults{}

	rels, err := c.getRels(req)
	if err != nil {
		return results, err
	}

	for _, rel := range rels {
		if params.IsMatch(rel) {
			results = results.addResult(rel)

			if paths == firstMatch {
				return results, nil
			}
		}
	}

	return results, nil
}

func (c *Checker) checkWildcard(step *model.RelationRef, params *CheckParams, paths checkPaths) (CheckResults, error) {
	req := &dsc.Relation{
		ObjectType:  params.OT.String(),
		ObjectId:    params.OID.String(),
		Relation:    params.Rel.String(),
		SubjectType: step.Object.String(),
		SubjectId:   "*",
	}

	results := CheckResults{}

	rels, err := c.getRels(req)
	if err != nil {
		return results, err
	}

	if len(rels) > 0 {
		// We have a wildcard match.
		results = results.addResult(rels...)
	}

	return results, nil
}

func (c *Checker) checkSubject(step *model.RelationRef, params *CheckParams, paths checkPaths) (CheckResults, error) {

	results := CheckResults{}

	if params.OID == "" {
		check := &CheckParams{
			OT:   step.Object,
			Rel:  step.Relation,
			ST:   c.params.ST,
			SID:  c.params.SID,
			SRel: c.params.SRel,
		}
		subjects := c.memo.Results(check)
		switch subjects.status() {
		case checkStatusUnknown:
			subs, err := c.check(check, paths)
			if err != nil {
				return results, err
			} else {
				subjects = subs
			}

		case checkStatusPending:
			// We found a cycle and need resolve the relations and
			// start unwinding the recursion.
			rels, err := c.getRels(check.AsRelation())
			if err != nil {
				return results, err
			}

			subjects = subjects.addResult(rels...)
			subjects = append(subjects, CheckParams{OT: params.ST, OID: params.SID})
		}

		for _, sub := range subjects {
			rels, err := c.getRels(&dsc.Relation{
				ObjectType:      params.OT.String(),
				Relation:        params.Rel.String(),
				SubjectType:     sub.OT.String(),
				SubjectId:       sub.OID.String(),
				SubjectRelation: params.SRel.String(),
			})
			if err != nil {
				return results, err
			}

			results = results.addResult(rels...)
		}

		return results, nil
	}

	req := &dsc.Relation{
		ObjectType:      params.OT.String(),
		ObjectId:        params.OID.String(),
		Relation:        params.Rel.String(),
		SubjectType:     step.Object.String(),
		SubjectRelation: step.Relation.String(),
	}
	rels, err := c.getRels(req)
	if err != nil {
		return results, err
	}

	for _, rel := range rels {
		if params.SRel == model.RelationName(rel.SubjectRelation) && params.ST == model.ObjectName(rel.SubjectType) {
			// We're searching for a subject relation (not a Check call) and we have a match.
			results = results.addResult(rel)
		}

		check := &CheckParams{
			OT:  step.Object,
			OID: ObjectID(rel.SubjectId),
			Rel: step.Relation,
			ST:  params.ST,
			SID: params.SID,
		}

		res, err := c.check(check, paths)
		if err != nil {
			return results, err
		}

		if len(res) > 0 {
			results = results.append(res)
			if paths == firstMatch {
				break
			}
		}
	}

	return results, nil
}

func (c *Checker) checkPermission(params *CheckParams, paths checkPaths) (CheckResults, error) {
	p := c.m.Objects[params.OT].Permissions[params.Rel]

	results := CheckResults{}

	subjTypes := []model.ObjectName{}
	if params.SRel == "" {
		subjTypes = append(subjTypes, params.ST)
	} else {
		subjTypes = c.m.Objects[params.ST].Relations[params.SRel].SubjectTypes
	}

	if len(lo.Intersect(subjTypes, p.SubjectTypes)) == 0 {
		// The subject type cannot have this permission.
		return results, nil
	}

	terms := p.Terms()
	termChecks := make([][]*CheckParams, 0, len(terms))
	for _, pt := range terms {
		// expand arrow operators.
		expanded, err := c.expandTerm(pt, params)
		if err != nil {
			return results, err
		}
		termChecks = append(termChecks, expanded)
	}

	switch {
	case p.IsUnion():
		return c.checkAny(termChecks, paths)
	case p.IsIntersection():
		return c.checkAll(termChecks, paths)
	case p.IsExclusion():
		include, err := c.checkAny(termChecks[:1], paths)
		switch {
		case err != nil:
			return results, err
		case len(include) == 0:
			// Short-circuit: The include term is false, so the permission is false.
			return results, nil
		}

		exclude, err := c.checkAny(termChecks[1:], paths)
		if err != nil {
			return results, err
		}

		results, _ := lo.Difference(include, exclude)
		return results, nil
	}

	return results, derr.ErrUnknown.Msg("unknown permission operator")
}

func (c *Checker) expandTerm(pt *model.PermissionTerm, params *CheckParams) ([]*CheckParams, error) {
	if pt.IsArrow() {
		// Resolve the base of the arrow.
		rels, err := c.getRels(&dsc.Relation{
			ObjectType: params.OT.String(),
			ObjectId:   params.OID.String(),
			Relation:   pt.Base.String(),
		})
		if err != nil {
			return []*CheckParams{}, err
		}

		expanded := lo.Map(rels, func(rel *dsc.Relation, _ int) *CheckParams {
			return &CheckParams{
				OT:   model.ObjectName(rel.SubjectType),
				OID:  ObjectID(rel.SubjectId),
				Rel:  pt.RelOrPerm,
				ST:   params.ST,
				SID:  params.SID,
				SRel: params.SRel,
			}
		})

		return expanded, nil
	}

	return []*CheckParams{{OT: params.OT, OID: params.OID, Rel: pt.RelOrPerm, ST: params.ST, SID: params.SID, SRel: params.SRel}}, nil
}

func (c *Checker) checkAny(checks [][]*CheckParams, paths checkPaths) (CheckResults, error) {
	results := CheckResults{}

	for _, check := range checks {
		var (
			res CheckResults
			err error
		)

		switch len(check) {
		case 0:
			res, err = CheckResults{}, nil
		case 1:
			res, err = c.check(check[0], paths)
		default:
			res, err = c.checkAny(lo.Map(check, func(params *CheckParams, _ int) []*CheckParams {
				return []*CheckParams{params}
			}), paths)
		}

		if err != nil {
			return res, err
		}

		if len(res) > 0 && paths == firstMatch {
			return res, nil
		}

		results = results.append(res)
	}

	return results, nil
}

func (c *Checker) checkAll(checks [][]*CheckParams, paths checkPaths) (CheckResults, error) {
	results := []CheckResults{}

	for _, check := range checks {
		// if the base of an arrow operator resolves to multiple objects (e.g. multiple "parents")
		// then a match on any of them is sufficient.
		result, err := c.checkAny([][]*CheckParams{check}, paths)
		if err != nil {
			return CheckResults{}, err
		}
		if len(result) == 0 {
			return result, nil
		}

		results = append(results, result)
	}

	intersection := lo.Reduce(results, func(agg CheckResults, item CheckResults, i int) CheckResults {
		if i == 0 {
			return item
		}

		return lo.Intersect(agg, item)
	}, CheckResults{})

	return intersection, nil
}
