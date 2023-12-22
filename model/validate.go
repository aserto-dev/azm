package model

import (
	"fmt"

	"github.com/aserto-dev/go-directory/pkg/derr"
	set "github.com/deckarep/golang-set/v2"
	"github.com/hashicorp/go-multierror"
	"github.com/samber/lo"
)

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

func (m *Model) validateUniqueNames(on ObjectName, o *Object) error {
	rels := lo.Map(lo.Keys(o.Relations), func(rn RelationName, _ int) string {
		return string(rn)
	})
	perms := lo.Map(lo.Keys(o.Permissions), func(pn RelationName, _ int) string {
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

func (m *Model) validateObjectRels(on ObjectName, o *Object) error {
	var errs error
	for rn, rs := range o.Relations {
		for _, r := range rs.Union {
			if r.Assignment() == RelationAssignmentUnknown {
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

func (m *Model) validateObjectPerms(on ObjectName, o *Object) error {
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

func (m *Model) validatePermission(on ObjectName, pn RelationName, p *Permission) error {
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

func (m *Model) resolveRelations() error {
	var errs error
	for on, o := range m.Objects {
		for rn, r := range o.Relations {
			seen := set.NewSet(RelationRef{Object: on, Relation: rn})
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

func (m *Model) resolveRelation(r *Relation, seen relSet) []ObjectName {
	if len(r.SubjectTypes) > 0 {
		// already resolved
		return r.SubjectTypes
	}

	subjectTypes := set.NewSet[ObjectName]()
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

func (m *Model) resolvePermissions() error {
	var errs error

	seen := set.NewSet[RelationRef]()
	for on, o := range m.Objects {
		for pn := range o.Permissions {
			subjs := m.resolvePermission(&RelationRef{on, pn}, seen)
			if subjs.IsEmpty() {
				errs = multierror.Append(errs, derr.ErrInvalidPermission.Msgf(
					"permission '%s:%s' cannot be satisfied by any type", on, pn),
				)
			}
		}
	}

	return errs
}

func (m *Model) resolvePermission(ref *RelationRef, seen relSet) objSet {
	p := m.Objects[ref.Object].Permissions[ref.Relation]

	if len(p.SubjectTypes) > 0 {
		// already resolved
		return set.NewSet(p.SubjectTypes...)
	}

	if seen.Contains(*ref) {
		// cycle detected
		return set.NewSet[ObjectName]()
	}
	seen.Add(*ref)

	for _, term := range p.Terms() {
		term.SubjectTypes = m.resolvePermissionTerm(ref.Object, term, seen)
	}

	// filter out terms that have no subject types. They represent cycles that are still being resolved.
	resolvedTerms := lo.Filter(p.Terms(), func(term *PermissionTerm, _ int) bool {
		return len(term.SubjectTypes) > 0
	})

	var subjTypes objSet

	switch {
	case p.IsUnion():
		subjTypes = set.NewSet(lo.FlatMap(resolvedTerms, func(term *PermissionTerm, _ int) []ObjectName {
			return term.SubjectTypes
		})...)

	case p.IsIntersection():
		subjTypes = lo.Reduce(resolvedTerms, func(acc objSet, term *PermissionTerm, i int) objSet {
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

func (m *Model) resolvePermissionTerm(on ObjectName, term *PermissionTerm, seen relSet) []ObjectName {
	var refs set.Set[RelationRef]

	switch {
	case term.IsArrow():
		sts := m.Objects[on].Relations[term.Base].SubjectTypes
		refs = set.NewSet(lo.Map(sts, func(st ObjectName, _ int) RelationRef {
			return RelationRef{Object: st, Relation: term.RelOrPerm}
		})...)

	default:
		refs = set.NewSet(RelationRef{Object: on, Relation: term.RelOrPerm})
	}

	subjectTypes := set.NewSet[ObjectName]()
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
