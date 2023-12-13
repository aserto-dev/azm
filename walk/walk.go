package walk

import (
	"sort"

	"github.com/aserto-dev/azm/model"
	dsc "github.com/aserto-dev/go-directory/aserto/directory/common/v3"
	dsr "github.com/aserto-dev/go-directory/aserto/directory/reader/v3"
	"github.com/aserto-dev/go-directory/pkg/derr"

	set "github.com/deckarep/golang-set/v2"
	"github.com/samber/lo"
)

type ObjectID string

func (id ObjectID) String() string {
	return string(id)
}

type ObjectIDs set.Set[ObjectID]

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

	visited map[model.RelationRef]ObjectIDs
}

func New(m *model.Model, req *dsr.CheckRequest, reader RelationReader) *Walker {
	return &Walker{
		m:       m,
		obj:     &Object{Type: model.ObjectName(req.ObjectType), ID: ObjectID(req.ObjectId)},
		rel:     model.RelationName(req.Relation),
		subj:    &Object{Type: model.ObjectName(req.SubjectType), ID: ObjectID(req.SubjectId)},
		getRels: reader,
		visited: map[model.RelationRef]ObjectIDs{},
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
		object:   w.obj,
		relation: w.rel,
		subject:  w.subj,
	}

	return w.check(params)
}

type checkParams struct {
	object   *Object
	relation model.RelationName
	subject  *Object
}

func (w *Walker) check(params *checkParams) (bool, error) {
	if !w.markVisited(params.object, params.relation) {
		// We already checked this relation and it didn't resolve to the subject.
		return false, nil
	}

	o := w.m.Objects[params.object.Type]

	if o.HasRelation(params.relation) {
		return w.checkRelation(params)
	}

	p := o.Permissions[params.relation]
	if !lo.Contains(p.SubjectTypes, params.subject.Type) {
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
		if err != nil {
			return false, err
		}
		if !include {
			// Short-circuit: The include term is false, so the permission is false.
			return false, nil
		}

		exclude, err := w.checkAny(termChecks[1:])
		if err != nil {
			return false, err
		}

		return !exclude, nil
	}

	return false, nil
}

func (w *Walker) markVisited(o *Object, rn model.RelationName) bool {
	rr := &model.RelationRef{Object: o.Type, Relation: rn}

	visited := w.visited[*rr]
	if visited == nil {
		visited = set.NewSet[ObjectID]()
		w.visited[*rr] = visited
	}

	if visited.Contains(o.ID) {
		// already visited
		return false
	}

	visited.Add(o.ID)
	return true
}

func (w *Walker) checkRelation(params *checkParams) (bool, error) {
	r := w.m.Objects[params.object.Type].Relations[params.relation]
	steps := w.stepRelation(r, params.subject.Type)

	for _, step := range steps {
		if step.IsWildcard() && step.Object == params.subject.Type {
			// We have a wildcard match.
			return true, nil
		}

		rels, err := w.getRels(&dsc.Relation{
			ObjectType:      params.object.Type.String(),
			ObjectId:        params.object.ID.String(),
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
				if rel.SubjectId == params.subject.ID.String() {
					return true, nil
				}
			}
		case step.IsSubject():
			for _, rel := range rels {
				if ok, err := w.check(&checkParams{
					object:   &Object{Type: step.Object, ID: ObjectID(rel.SubjectId)},
					relation: step.Relation,
					subject:  params.subject,
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

func (w *Walker) expandTerm(pt *model.PermissionTerm, params *checkParams) ([]*checkParams, error) {
	if pt.IsArrow() {
		// Resolve the base of the arrow.
		rels, err := w.getRels(&dsc.Relation{
			ObjectType:  params.object.Type.String(),
			ObjectId:    params.object.ID.String(),
			Relation:    pt.Base.String(),
			SubjectType: params.subject.Type.String(),
		})
		if err != nil {
			return []*checkParams{}, err
		}

		expanded := lo.Map(rels, func(rel *dsc.Relation, _ int) *checkParams {
			return &checkParams{
				object:   &Object{Type: model.ObjectName(rel.SubjectType), ID: ObjectID(rel.SubjectId)},
				relation: pt.RelOrPerm,
				subject:  params.subject,
			}
		})

		return expanded, nil
	}

	return []*checkParams{{object: params.object, relation: pt.RelOrPerm, subject: params.subject}}, nil
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
