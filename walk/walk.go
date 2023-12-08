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
	stats   *Stats
	obj     *Object
	rel     model.RelationName
	subj    *Object
	getRels RelationReader

	visited map[model.RelationRef]ObjectIDs
}

func New(m *model.Model, stats *Stats, req *dsr.CheckRequest, reader RelationReader) *Walker {
	return &Walker{
		m:       m,
		stats:   stats,
		obj:     &Object{Type: model.ObjectName(req.ObjectType), ID: ObjectID(req.ObjectId)},
		rel:     model.RelationName(req.Relation),
		subj:    &Object{Type: model.ObjectName(req.SubjectType), ID: ObjectID(req.SubjectId)},
		getRels: reader,
		visited: map[model.RelationRef]ObjectIDs{},
	}
}

func (w *Walker) Check() (bool, error) {
	o := w.m.Objects[model.ObjectName(w.obj.Type)]
	if o == nil {
		return false, derr.ErrObjectTypeNotFound.Msg(w.obj.Type.String())
	}

	if o.HasRelation(w.rel) {
		return w.checkRelation(w.obj, w.rel, w.subj)
	}

	if o.HasPermission(model.PermissionName(w.rel)) {
		return w.checkPermission()
	}

	return false, derr.ErrRelationNotFound.Msg(w.rel.String())
}

func (w *Walker) checkRelation(
	object *Object,
	relation model.RelationName,
	subject *Object,
) (bool, error) {
	rr := &model.RelationRef{Object: object.Type, Relation: relation}

	visited := w.visited[*rr]
	if visited == nil {
		visited = set.NewSet[ObjectID]()
		w.visited[*rr] = visited
	}

	objID := ObjectID(object.ID)
	if visited.Contains(objID) {
		// We already checked this relation and it didn't resolve to the subject.
		return false, nil
	}

	visited.Add(ObjectID(object.ID))

	steps := w.step(rr, subject.Type)

	for _, step := range steps {
		if step.IsWildcard() && step.Object == subject.Type {
			// We have a wildcard match.
			return true, nil
		}

		rels, err := w.getRels(&dsc.Relation{
			ObjectType:      object.Type.String(),
			ObjectId:        object.ID.String(),
			Relation:        relation.String(),
			SubjectType:     step.Object.String(),
			SubjectRelation: step.Relation.String(),
		})
		if err != nil {
			return false, err
		}
		switch {
		case step.IsDirect():
			for _, rel := range rels {
				if rel.SubjectId == subject.ID.String() {
					return true, nil
				}
			}
		case step.IsSubject():
			for _, rel := range rels {
				if ok, err := w.checkRelation(&Object{Type: step.Object, ID: ObjectID(rel.SubjectId)}, step.Relation, subject); err != nil {
					return false, err
				} else if ok {
					return true, nil
				}
			}
		}
	}

	return false, nil
}

func (w *Walker) checkPermission() (bool, error) {
	return false, nil
}

func (w *Walker) step(rr *model.RelationRef, subjs ...model.ObjectName) []*model.RelationRef {
	if rr == nil {
		return []*model.RelationRef{}
	}

	o := w.m.Objects[rr.Object]
	if o == nil {
		return []*model.RelationRef{}
	}

	r := o.Relations[rr.Relation]
	if r == nil {
		return []*model.RelationRef{}
	}

	steps := lo.FilterMap(r.Union, func(rt *model.RelationTerm, _ int) (*model.RelationRef, bool) {
		if lo.Contains([]model.RelationAssignment{
			model.RelationAssignmentDirect, model.RelationAssignmentWildcard,
		}, rt.Assignment()) {
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
