package model

import (
	"fmt"

	"github.com/aserto-dev/go-directory/pkg/derr"
	set "github.com/deckarep/golang-set/v2"
	"github.com/hashicorp/go-multierror"
	"github.com/samber/lo"
)

type validator struct {
	*Model
	opts *validationOptions
}

func newValidator(m *Model, opts *validationOptions) *validator {
	return &validator{Model: m, opts: opts}
}

func (v *validator) validate() error {
	if !v.opts.skipNameValidation {
		// validate that all object/relation/permission names are valid identifiers.
		if err := v.validateNames(); err != nil {
			return derr.ErrInvalidArgument.Err(err)
		}
	}

	// Pass 1 (syntax): ensure no name collisions and all relations reference existing objects/relations.
	if err := v.validateReferences(); err != nil {
		return derr.ErrInvalidArgument.Err(err)
	}

	// Pass 2: resolve all relations to a set of possible subject types.
	if err := v.resolveRelations(); err != nil {
		return derr.ErrInvalidArgument.Err(err)
	}

	// Pass 3: validate all arrow operators in permissions. This requires that all relations have already been resolved.
	if err := v.validatePermissions(); err != nil {
		return derr.ErrInvalidArgument.Err(err)
	}

	// Pass 4: resolve all permissions to a set of possible subject types.
	if err := v.resolvePermissions(); err != nil {
		return derr.ErrInvalidArgument.Err(err)
	}

	return nil
}

func (v *validator) validateNames() error {
	var errs error

	for on, o := range v.Objects {
		if !on.Valid() {
			errs = multierror.Append(errs, derr.ErrInvalidObjectType.Msgf("invalid type name '%s': %s", on, msgInvalidIdentifier))
		}

		for rn := range o.Relations {
			if !rn.Valid() {
				errs = multierror.Append(errs, derr.ErrInvalidRelationType.Msgf("invalid relation name '%s': %s", rn, msgInvalidIdentifier))
			}
		}

		for pn := range o.Permissions {
			if !pn.Valid() {
				errs = multierror.Append(errs, derr.ErrInvalidPermission.Msgf("invalid permission name '%s': %s", pn, msgInvalidIdentifier))
			}
		}
	}

	return errs
}

func (v *validator) validateReferences() error {
	var errs error

	for on, o := range v.Objects {
		validatePerms := true
		if err := v.validateUniqueNames(on, o); err != nil {
			errs = multierror.Append(errs, err)
			validatePerms = false
		}

		if err := v.validateObjectRels(on, o); err != nil {
			errs = multierror.Append(errs, err)
		}

		if validatePerms {
			if err := v.validateObjectPerms(on, o); err != nil {
				errs = multierror.Append(errs, err)
			}
		}
	}

	return errs
}

func (v *validator) validateUniqueNames(on ObjectName, o *Object) error {
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

func (v *validator) validateObjectRels(on ObjectName, o *Object) error {
	var errs error
	for rn, rs := range o.Relations {
		for _, r := range rs.Union {
			if r.Assignment() == RelationAssignmentUnknown {
				errs = multierror.Append(errs, derr.ErrInvalidRelationType.Msgf(
					"relation '%s:%s' has no definition", on, rn),
				)
				continue
			}

			o := v.Objects[r.Object]
			if o == nil {
				errs = multierror.Append(errs, derr.ErrInvalidRelationType.Msgf(
					"relation '%s:%s' references undefined object type '%s'", on, rn, r.Object),
				)
				continue
			}

			if r.IsSubject() {
				if _, ok := o.Relations[r.Relation]; !ok {
					errs = multierror.Append(errs, derr.ErrInvalidRelationType.Msgf(
						"relation '%s:%s' references undefined relation type '%s#%s'", on, rn, r.Object, r.Relation),
					)
				}
			}
		}
	}

	return errs
}

func (v *validator) validateObjectPerms(on ObjectName, o *Object) error {
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
			if term == nil {
				errs = multierror.Append(errs, derr.ErrInvalidPermission.Msgf(
					"permission '%s:%s' has an empty term", on, pn),
				)
				continue
			}
			switch {
			case term.IsArrow():
				// this is an arrow operator.
				// validate that the base relation exists on this object type.
				// at this stage we don't yet resolve the relation to a set of subject types.
				if !o.HasRelOrPerm(term.Base) {
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

func (v *validator) validatePermissions() error {
	var errs error
	for on, o := range v.Objects {
		for pn, p := range o.Permissions {
			if err := v.validatePermission(on, pn, p); err != nil {
				errs = multierror.Append(errs, err)
			}
		}
	}
	return errs
}

func (v *validator) validatePermission(on ObjectName, pn RelationName, p *Permission) error {
	o := v.Objects[on]

	var errs error
	for _, term := range p.Terms() {
		if term.IsArrow() {
			// given a reference base->rel_or_perm, validate that all object types that `base` can resolve to
			// have a permission or relation named `rel_or_perm`.
			if o.HasPermission(term.Base) {
				if !v.opts.allowPermissionInArrowBase {
					errs = multierror.Append(errs, derr.ErrInvalidPermission.Msgf(
						"permission '%s:%s' references permission '%s', which is not allowed in arrow base. only relations can be used.", on, pn, term.Base),
					)
				}
				continue
			}
			r := o.Relations[term.Base]
			for _, st := range r.SubjectTypes {
				if !v.Objects[st].HasRelOrPerm(term.RelOrPerm) {
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

func (v *validator) resolveRelations() error {
	var errs error
	for on, o := range v.Objects {
		for rn, r := range o.Relations {
			seen := set.NewSet(RelationRef{Object: on, Relation: rn})
			subs, intermediates := v.resolveRelation(r, seen)
			switch len(subs) {
			case 0:
				errs = multierror.Append(errs, derr.ErrInvalidRelationType.Msgf(
					"relation '%s:%s' is circular and does not resolve to any object types", on, rn),
				)
			default:
				r.SubjectTypes = subs
				r.Intermediates = intermediates
			}
		}
	}

	return errs
}

func (v *validator) resolveRelation(r *Relation, seen relSet) ([]ObjectName, RelationRefs) {
	if len(r.SubjectTypes) > 0 {
		// already resolved
		return r.SubjectTypes, r.Intermediates
	}

	subjectTypes := set.NewSet[ObjectName]()
	intermediateTypes := set.NewSet[RelationRef]()
	for _, rr := range r.Union {
		if rr.IsSubject() {
			intermediateTypes.Add(*rr)
			if !seen.Contains(*rr) {
				seen.Add(*rr)
				subs, intermediates := v.resolveRelation(v.Objects[rr.Object].Relations[rr.Relation], seen)
				subjectTypes.Append(subs...)
				intermediateTypes.Append(intermediates...)
			}
		} else {
			subjectTypes.Add(rr.Object)
		}
	}
	return subjectTypes.ToSlice(), intermediateTypes.ToSlice()
}

func (v *validator) resolvePermissions() error {
	var errs error

	seen := set.NewSet[RelationRef]()
	for on, o := range v.Objects {
		for pn := range o.Permissions {
			subjs, _ := v.resolvePermission(&RelationRef{on, pn}, seen)
			if subjs.IsEmpty() {
				errs = multierror.Append(errs, derr.ErrInvalidPermission.Msgf(
					"permission '%s:%s' cannot be satisfied by any type", on, pn),
				)
			}
		}
	}

	return errs
}

func (v *validator) resolvePermission(ref *RelationRef, seen relSet) (objSet, relSet) {
	p, ok := v.Objects[ref.Object].Permissions[ref.Relation]
	if !ok {
		// No such permission. Most likely a bug in the model inversion logic.
		// Return empty sets which result in a validation error.
		return set.NewSet[ObjectName](), set.NewSet[RelationRef]()
	}

	if len(p.SubjectTypes) > 0 {
		// already resolved
		return set.NewSet(p.SubjectTypes...), set.NewSet(p.Intermediates...)
	}

	if seen.Contains(*ref) {
		// cycle detected
		return set.NewSet[ObjectName](), set.NewSet[RelationRef]()
	}
	seen.Add(*ref)

	for _, term := range p.Terms() {
		term.SubjectTypes, term.Intermediates = v.resolvePermissionTerm(ref.Object, term, seen)
	}

	// filter out terms that have no subject types. They represent cycles that are still being resolved.
	resolvedTerms := lo.Filter(p.Terms(), func(term *PermissionTerm, _ int) bool {
		return len(term.SubjectTypes) > 0
	})

	var (
		subjTypes     objSet
		intermediates relSet
	)

	switch {
	case p.IsUnion():
		subjTypes = set.NewSet(lo.FlatMap(resolvedTerms, func(term *PermissionTerm, _ int) []ObjectName {
			return term.SubjectTypes
		})...)
		intermediates = set.NewSet(lo.FlatMap(resolvedTerms, func(term *PermissionTerm, _ int) []RelationRef {
			return term.Intermediates
		})...)

	case p.IsIntersection():
		subjTypes = lo.Reduce(resolvedTerms, func(acc objSet, term *PermissionTerm, i int) objSet {
			subjs := set.NewSet(term.SubjectTypes...)

			if i == 0 {
				return subjs
			}

			return acc.Intersect(subjs)

		}, nil)
		intermediates = lo.Reduce(resolvedTerms, func(acc relSet, term *PermissionTerm, i int) relSet {
			subjs := set.NewSet(term.Intermediates...)

			if i == 0 {
				return subjs
			}

			return acc.Intersect(subjs)

		}, nil)

	case p.IsExclusion():
		subjTypes = set.NewSet(p.Exclusion.Include.SubjectTypes...)
		intermediates = set.NewSet(p.Exclusion.Include.Intermediates...)
	}

	p.SubjectTypes = subjTypes.ToSlice()
	p.Intermediates = intermediates.ToSlice()

	return subjTypes, intermediates
}

func (v *validator) resolvePermissionTerm(on ObjectName, term *PermissionTerm, seen relSet) ([]ObjectName, RelationRefs) {
	var refs set.Set[RelationRef]
	intermediates := set.NewSet[RelationRef]()

	switch {
	case term.IsArrow():
		o := v.Objects[on]
		var (
			sts []ObjectName
		)
		if o.HasRelation(term.Base) {
			sts = o.Relations[term.Base].SubjectTypes
			intermediates.Append(o.Relations[term.Base].Intermediates...)
		} else {
			types, interims := v.resolvePermission(&RelationRef{Object: on, Relation: term.Base}, seen)
			sts = types.ToSlice()
			intermediates.Append(interims.ToSlice()...)
		}
		refs = set.NewSet(lo.Map(sts, func(st ObjectName, _ int) RelationRef {
			return RelationRef{Object: st, Relation: term.RelOrPerm}
		})...)

	default:
		refs = set.NewSet(RelationRef{Object: on, Relation: term.RelOrPerm})
	}

	subjectTypes := set.NewSet[ObjectName]()
	for ref := range refs.Iter() {
		o := v.Objects[ref.Object]

		if o.HasRelation(ref.Relation) {
			// Relations are already resolved to a set of subject types.
			subjectTypes.Append(o.Relations[ref.Relation].SubjectTypes...)
			intermediates.Append(o.Relations[ref.Relation].Intermediates...)
			continue
		}

		sts, interims := v.resolvePermission(&ref, seen)
		subjectTypes = subjectTypes.Union(sts)
		intermediates.Append(interims.ToSlice()...)
	}

	return subjectTypes.ToSlice(), intermediates.ToSlice()
}
