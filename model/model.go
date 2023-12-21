package model

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/aserto-dev/azm/graph"
	"github.com/aserto-dev/azm/types"
	"github.com/aserto-dev/go-directory/pkg/derr"
	set "github.com/deckarep/golang-set/v2"
	"github.com/hashicorp/go-multierror"
	"github.com/samber/lo"
)

const ModelVersion int = 2

type Model struct {
	Version  int                                `json:"version"`
	Objects  map[types.ObjectName]*types.Object `json:"types"`
	Metadata *Metadata                          `json:"metadata"`
}

type Metadata struct {
	UpdatedAt time.Time `json:"updated_at"`
	ETag      string    `json:"etag"`
}

func New(r io.Reader) (*Model, error) {
	m := Model{}
	dec := json.NewDecoder(r)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&m); err != nil {
		return nil, err
	}
	return &m, nil
}

type objSet set.Set[types.ObjectName]
type relSet set.Set[types.RelationRef]

func (m *Model) GetGraph() *graph.Graph {
	grph := graph.NewGraph()
	for objectName := range m.Objects {
		grph.AddNode(string(objectName))
	}
	for objectName, obj := range m.Objects {
		for relName, rel := range obj.Relations {
			for _, rl := range rel.Union {
				if rl.IsDirect() {
					grph.AddEdge(string(objectName), string(rl.Object), string(relName))
				} else if rl.IsSubject() {
					grph.AddEdge(string(objectName), string(rl.Object), string(relName))
				}
			}
		}
	}

	return grph
}

func (m *Model) Reader() (io.Reader, error) {
	b := bytes.Buffer{}
	enc := json.NewEncoder(&b)
	if err := enc.Encode(m); err != nil {
		return nil, err
	}
	return bytes.NewReader(b.Bytes()), nil
}

func (m *Model) Write(w io.Writer) error {
	enc := json.NewEncoder(w)
	return enc.Encode(m)
}

func (m *Model) Diff(newModel *Model) *Diff {
	// newmodel - m => additions
	added := newModel.subtract(m)
	// m - newModel => deletions
	deleted := m.subtract(newModel)

	return &Diff{Added: *added, Removed: *deleted}
}

// Validate enforces the model's internal consistency.
//
// It enforces the following rules:
//   - Within an object, a permission cannot share the same name as a relation.
//   - Direct relations must reference existing objects .
//   - Wildcard relations must reference existing objects.
//   - Subject relations must reference existing object#relation pairs.
//   - Arrow permissions (relation->rel_or_perm) must reference existing relations/permissions.
func (m *Model) Validate() error {
	// Pass 1 (syntax): ensure no name collisions and all relations reference existing objects/relations.
	if err := m.validateReferences(); err != nil {
		return derr.ErrInvalidArgument.Err(err)
	}

	// Pass 2: resolve all relations to a set of possible subject types.
	if err := m.resolveRelations(); err != nil {
		return derr.ErrInvalidArgument.Err(err)
	}

	// Pass 3: validate all arrow operators in permissions. This requires that all relations have already been resolved.
	if err := m.validatePermissions(); err != nil {
		return derr.ErrInvalidArgument.Err(err)
	}

	// Pass 4: resolve all permissions to a set of possible subject types.
	if err := m.resolvePermissions(); err != nil {
		return derr.ErrInvalidArgument.Err(err)
	}

	return nil
}

func (m *Model) validateReferences() error {
	var errs error

	for on, o := range m.Objects {
		validatePerms := true
		if err := m.validateUniqueNames(on, o); err != nil {
			errs = multierror.Append(errs, err)
			validatePerms = false
		}

		if err := m.validateObjectRels(on, o); err != nil {
			errs = multierror.Append(errs, err)
		}

		if validatePerms {
			if err := m.validateObjectPerms(on, o); err != nil {
				errs = multierror.Append(errs, err)
			}
		}
	}

	return errs
}

func (m *Model) validateUniqueNames(on types.ObjectName, o *types.Object) error {
	rels := lo.Map(lo.Keys(o.Relations), func(rn types.RelationName, _ int) string {
		return string(rn)
	})
	perms := lo.Map(lo.Keys(o.Permissions), func(pn types.RelationName, _ int) string {
		return string(pn)
	})

	rpCollisions := lo.Intersect(rels, perms)

	var errs error
	for _, collision := range rpCollisions {
		errs = multierror.Append(errs, derr.ErrInvalidPermission.Msgf(
			"permission name '%[1]s:%[2]s' conflicts with '%[1]s:%[2]s' relation", on, collision),
		)
	}

	return errs
}

func (m *Model) validateObjectRels(on types.ObjectName, o *types.Object) error {
	var errs error
	for rn, rs := range o.Relations {
		for _, r := range rs.Union {
			if r.Assignment() == types.RelationAssignmentUnknown {
				errs = multierror.Append(errs, derr.ErrInvalidRelation.Msgf(
					"relation '%s:%s' has no definition", on, rn),
				)
				continue
			}

			o := m.Objects[r.Object]
			if o == nil {
				errs = multierror.Append(errs, derr.ErrInvalidRelation.Msgf(
					"relation '%s:%s' references undefined object type '%s'", on, rn, r.Object),
				)
				continue
			}

			if r.IsSubject() {
				if _, ok := o.Relations[r.Relation]; !ok {
					errs = multierror.Append(errs, derr.ErrInvalidRelation.Msgf(
						"relation '%s:%s' references undefined relation type '%s#%s'", on, rn, r.Object, r.Relation),
					)
				}
			}
		}
	}

	return errs
}

func (m *Model) validateObjectPerms(on types.ObjectName, o *types.Object) error {
	var errs error
	for pn, p := range o.Permissions {
		terms := p.Terms()
		if len(terms) == 0 {
			errs = multierror.Append(errs, derr.ErrInvalidPermission.Msgf(
				"permission '%s:%s' has no definition", on, pn),
			)
			continue
		}

		for _, term := range terms {
			switch {
			case term.IsArrow():
				// this is an arrow operator.
				// validate that the base relation exists on this object type.
				// at this stage we don't yet resolve the relation to a set of subject types.
				if _, hasRelation := o.Relations[term.Base]; !hasRelation {
					errs = multierror.Append(errs, derr.ErrInvalidPermission.Msgf(
						"permission '%s:%s' references undefined relation type '%s:%s'", on, pn, on, term.Base),
					)
				}

			default:
				// validate that the relation/permission exists on this object type.
				if !o.HasRelOrPerm(term.RelOrPerm) {
					errs = multierror.Append(errs, derr.ErrInvalidPermission.Msgf(
						"permission '%s:%s' references undefined relation or permission '%s:%s'", on, pn, on, term.RelOrPerm),
					)
				}
			}

		}
	}

	return errs
}

func (m *Model) resolveRelations() error {
	var errs error
	for on, o := range m.Objects {
		for rn, r := range o.Relations {
			seen := set.NewSet(types.RelationRef{Object: on, Relation: rn})
			subs := m.resolveRelation(r, seen)
			switch len(subs) {
			case 0:
				errs = multierror.Append(errs, derr.ErrInvalidRelation.Msgf(
					"relation '%s:%s' is circular and does not resolve to any object types", on, rn),
				)
			default:
				r.SubjectTypes = subs
			}
		}
	}

	return errs
}

func (m *Model) resolveRelation(r *types.Relation, seen relSet) []types.ObjectName {
	if len(r.SubjectTypes) > 0 {
		// already resolved
		return r.SubjectTypes
	}

	subjectTypes := set.NewSet[types.ObjectName]()
	for _, rr := range r.Union {
		switch {
		case rr.IsSubject():
			if !seen.Contains(*rr) {
				seen.Add(*rr)
				subjectTypes.Append(m.resolveRelation(m.Objects[rr.Object].Relations[rr.Relation], seen)...)
			}
		default:
			subjectTypes.Add(rr.Object)
		}
	}
	return subjectTypes.ToSlice()
}

func (m *Model) validatePermissions() error {
	var errs error
	for on, o := range m.Objects {
		for pn, p := range o.Permissions {
			if err := m.validatePermission(on, pn, p); err != nil {
				errs = multierror.Append(errs, err)
			}
		}
	}
	return errs
}

func (m *Model) validatePermission(on types.ObjectName, pn types.RelationName, p *types.Permission) error {
	o := m.Objects[on]

	var errs error
	for _, term := range p.Terms() {
		if term.IsArrow() {
			// given a reference base->rel_or_perm, validate that all object types that `base` can resolve to
			// have a permission or relation named `rel_or_perm`.
			r := o.Relations[term.Base]
			for _, st := range r.SubjectTypes {
				if !m.Objects[st].HasRelOrPerm(term.RelOrPerm) {
					arrow := fmt.Sprintf("%s->%s", term.Base, term.RelOrPerm)
					errs = multierror.Append(errs, derr.ErrInvalidPermission.Msgf(
						"permission '%s:%s' references '%s', which can resolve to undefined relation or permission '%s:%s' ",
						on, pn, arrow, st, term.RelOrPerm,
					))
				}
			}

		}
	}

	return errs
}

func (m *Model) resolvePermissions() error {
	var errs error

	seen := set.NewSet[types.RelationRef]()
	for on, o := range m.Objects {
		for pn := range o.Permissions {
			subjs := m.resolvePermission(&types.RelationRef{Object: on, Relation: pn}, seen)
			if subjs.IsEmpty() {
				errs = multierror.Append(errs, derr.ErrInvalidPermission.Msgf(
					"permission '%s:%s' cannot be satisfied by any type", on, pn),
				)
			}
		}
	}

	return errs
}

func (m *Model) resolvePermission(ref *types.RelationRef, seen relSet) objSet {
	p := m.Objects[ref.Object].Permissions[ref.Relation]

	if len(p.SubjectTypes) > 0 {
		// already resolved
		return set.NewSet(p.SubjectTypes...)
	}

	if seen.Contains(*ref) {
		// cycle detected
		return set.NewSet[types.ObjectName]()
	}
	seen.Add(*ref)

	for _, term := range p.Terms() {
		term.SubjectTypes = m.resolvePermissionTerm(ref.Object, term, seen)
	}

	// filter out terms that have no subject types. They represent cycles that are still being resolved.
	resolvedTerms := lo.Filter(p.Terms(), func(term *types.PermissionTerm, _ int) bool {
		return len(term.SubjectTypes) > 0
	})

	var subjTypes objSet

	switch {
	case p.IsUnion():
		subjTypes = set.NewSet(lo.FlatMap(resolvedTerms, func(term *types.PermissionTerm, _ int) []types.ObjectName {
			return term.SubjectTypes
		})...)

	case p.IsIntersection():
		subjTypes = lo.Reduce(resolvedTerms, func(acc objSet, term *types.PermissionTerm, i int) objSet {
			subjs := set.NewSet(term.SubjectTypes...)

			if i == 0 {
				return subjs
			}

			return acc.Intersect(subjs)

		}, nil)

	case p.IsExclusion():
		subjTypes = set.NewSet(p.Exclusion.Include.SubjectTypes...)
	}

	p.SubjectTypes = subjTypes.ToSlice()

	return subjTypes
}

func (m *Model) resolvePermissionTerm(on types.ObjectName, term *types.PermissionTerm, seen relSet) []types.ObjectName {
	var refs set.Set[types.RelationRef]

	switch {
	case term.IsArrow():
		sts := m.Objects[on].Relations[term.Base].SubjectTypes
		refs = set.NewSet(lo.Map(sts, func(st types.ObjectName, _ int) types.RelationRef {
			return types.RelationRef{Object: st, Relation: term.RelOrPerm}
		})...)

	default:
		refs = set.NewSet(types.RelationRef{Object: on, Relation: term.RelOrPerm})
	}

	subjectTypes := set.NewSet[types.ObjectName]()
	for ref := range refs.Iter() {
		o := m.Objects[ref.Object]

		if o.HasRelation(ref.Relation) {
			// Relations are already resolved to a set of subject types.
			subjectTypes.Append(o.Relations[ref.Relation].SubjectTypes...)
			continue
		}

		subjectTypes = subjectTypes.Union(m.resolvePermission(&ref, seen))
	}

	return subjectTypes.ToSlice()
}

func (m *Model) subtract(newModel *Model) *Changes {
	chgs := &Changes{
		Objects:   make([]types.ObjectName, 0),
		Relations: make(map[types.ObjectName]map[types.RelationName][]string),
	}

	if m == nil {
		return chgs
	}

	if newModel == nil {
		for objName := range m.Objects {
			chgs.Objects = append(chgs.Objects, objName)
		}
		return chgs
	}

	for objName, obj := range m.Objects {
		if newModel.Objects[objName] == nil {
			chgs.Objects = append(chgs.Objects, objName)
		} else {
			for relname, rel := range obj.Relations {
				if newModel.Objects[objName].Relations[relname] == nil {
					if chgs.Relations[objName] == nil {
						chgs.Relations[objName] = make(map[types.RelationName][]string, 0)
					}
					chgs.Relations[objName][relname] = []string{}
					continue
				}
				relDiff := substractArray(rel.Union, newModel.Objects[objName].Relations[relname].Union)
				if len(relDiff) > 0 {
					if chgs.Relations[objName] == nil {
						chgs.Relations[objName] = make(map[types.RelationName][]string, 0)
					}
					chgs.Relations[objName][relname] = relDiff
				}
			}
		}
	}

	return chgs
}

func substractArray(arr1, arr2 []*types.RelationRef) []string {
	result := []string{}
	for _, elem := range arr1 {
		found := false
		for _, elem2 := range arr2 {
			if elem == elem2 {
				found = true
				break
			}
		}
		if !found {
			result = append(result, elem.String())
		}
	}
	return result
}
