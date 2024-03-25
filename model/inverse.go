package model

import (
	"fmt"

	"github.com/samber/lo"
)

const (
	ObjectNameSeparator       = "^"
	SubjectRelationSeparator  = "#"
	GeneratedPermissionPrefix = "$"
)

func (m *Model) Invert() *Model {
	return newInverter(m).invert()
}

type inverter struct {
	m     *Model
	im    *Model
	subst map[RelationName]RelationName
}

func newInverter(m *Model) *inverter {
	return &inverter{
		m: m,
		im: &Model{
			Version:  m.Version,
			Objects:  lo.MapValues(m.Objects, func(o *Object, _ ObjectName) *Object { return NewObject() }),
			Metadata: m.Metadata,
		},
		subst: map[RelationName]RelationName{},
	}
}

func (i *inverter) invert() *Model {
	// invert all relations before inverting permissions.
	// this is necessary to create synthetic permissions for subject relations.
	// these permissions are stored in the substitution map (i.subst) and used in inverted permissions.
	for on, o := range i.m.Objects {
		for rn, r := range o.Relations {
			i.invertRelation(on, rn, r)
		}
	}

	for _, o := range i.im.Objects {
		i.applySubstitutions(o)
	}

	for on, o := range i.m.Objects {
		for pn, p := range o.Permissions {
			i.invertPermission(on, pn, o, p)
		}
	}

	return i.im
}

func (i *inverter) invertRelation(on ObjectName, rn RelationName, r *Relation) {
	unionObjs := lo.Associate(r.Union, func(rr *RelationRef) (ObjectName, bool) { return rr.Object, true })

	for _, rr := range r.Union {
		irn := InverseRelation(on, rn)
		isrn := InverseRelation(on, rn, rr.Relation)
		i.im.Objects[rr.Object].Relations[isrn] = &Relation{Union: []*RelationRef{{Object: on}}}
		if rr.IsSubject() {
			// add a synthetic permission to reverse the expansion of the subject relation
			ipn := PermForRel(irn)
			srel := i.m.Objects[rr.Object].Relations[rr.Relation]

			for _, subj := range srel.AllRefs() {
				p := permissionOrNew(i.im.Objects[subj.Object], ipn, permissionKindUnion)
				i.addSubstitution(irn, ipn)
				if _, ok := unionObjs[subj.Object]; ok {
					p.AddTerm(&PermissionTerm{RelOrPerm: InverseRelation(on, rn, subj.Relation)})
				}
				rel := InverseRelation(rr.Object, rr.Relation)
				p.AddTerm(&PermissionTerm{Base: rel, RelOrPerm: ipn})
			}
		}
	}
}

func (i *inverter) applySubstitutions(o *Object) {
	for pn, p := range o.Permissions {
		for _, pt := range p.Terms() {
			if !pt.IsArrow() && PermForRel(pt.RelOrPerm) == pn {
				continue
			}
			pt.Base = i.sub(pt.Base)
			pt.RelOrPerm = i.sub(pt.RelOrPerm)
		}
	}
}

func (i *inverter) invertPermission(on ObjectName, pn RelationName, o *Object, p *Permission) {
	kind := kindOf(p)
	ipn := InverseRelation(on, pn)

	for _, pt := range p.Terms() {
		switch {
		case pt.IsArrow():
			baseRel := o.Relations[pt.Base]
			itip := i.irelSub(on, pt.Base)

			for _, subj := range p.SubjectTypes {
				ip := permissionOrNew(i.im.Objects[subj], ipn, kind)
				for _, baseSubj := range baseRel.AllTypes() {
					ip.AddTerm(&PermissionTerm{Base: i.irelSub(baseSubj, pt.RelOrPerm), RelOrPerm: itip})

					// create a subject relation to expand the recursive permission
					r := relationOrNew(i.im.Objects[baseSubj], itip)
					r.AddRef(&RelationRef{Object: baseSubj, Relation: itip})
				}
			}

			for _, rr := range p.Intermediates {
				isrpn := InverseRelation(on, pn, rr.Relation)
				ip := permissionOrNew(i.im.Objects[rr.Object], isrpn, kind)
				for _, baseSubj := range baseRel.AllTypes() {
					term := &PermissionTerm{Base: i.sub(InverseRelation(baseSubj, pt.RelOrPerm, rr.Relation)), RelOrPerm: itip}
					ip.AddTerm(term)
				}
			}

		case o.HasRelation(pt.RelOrPerm):
			r := o.Relations[pt.RelOrPerm]
			for _, subj := range r.SubjectTypes {
				ip := permissionOrNew(i.im.Objects[subj], ipn, kind)
				ip.AddTerm(&PermissionTerm{RelOrPerm: i.irelSub(on, pt.RelOrPerm)})
			}

			for _, rr := range r.Intermediates {
				sripn := InverseRelation(on, pn, rr.Relation)
				ip := permissionOrNew(i.im.Objects[rr.Object], sripn, kind)
				ip.AddTerm(&PermissionTerm{RelOrPerm: i.irelSub(on, pt.RelOrPerm)})
			}

		case o.HasPermission(pt.RelOrPerm):
			for _, subj := range subjs(o, pt.RelOrPerm) {
				ip := permissionOrNew(i.im.Objects[subj], ipn, kind)
				ip.AddTerm(&PermissionTerm{RelOrPerm: i.irelSub(on, pt.RelOrPerm)})
			}
		}
	}
}

func (i *inverter) irelSub(on ObjectName, rn RelationName) RelationName {
	return i.sub(InverseRelation(on, rn))
}

func (i *inverter) sub(rn RelationName) RelationName {
	if pn, ok := i.subst[rn]; ok {
		return pn
	}

	return rn
}

func (i *inverter) addSubstitution(rn, pn RelationName) {
	i.subst[rn] = pn
}

type permissionKind int

const (
	permissionKindUnion permissionKind = iota
	permissionKindIntersection
	permissionKindExclusion
)

func kindOf(p *Permission) permissionKind {
	switch {
	case p.IsUnion():
		return permissionKindUnion
	case p.IsIntersection():
		return permissionKindIntersection
	case p.IsExclusion():
		return permissionKindExclusion
	}

	panic("unknown permission kind")
}

func relationOrNew(o *Object, rn RelationName) *Relation {
	r := o.Relations[rn]
	if r != nil {
		return r
	}

	r = &Relation{}
	o.Relations[rn] = r

	return r
}

func permissionOrNew(o *Object, pn RelationName, kind permissionKind) *Permission {
	p := o.Permissions[pn]
	if p != nil {
		return p
	}

	p = &Permission{}
	terms := PermissionTerms{}
	switch kind {
	case permissionKindUnion:
		p.Union = terms
	case permissionKindIntersection:
		p.Intersection = terms
	case permissionKindExclusion:
		p.Exclusion = &ExclusionPermission{}
	}

	o.Permissions[pn] = p

	return p
}

func InverseRelation(on ObjectName, rn RelationName, srn ...RelationName) RelationName {
	irn := RelationName(fmt.Sprintf("%s%s%s", on, ObjectNameSeparator, rn))

	switch {
	case len(srn) == 0 || srn[0] == "":
		return irn
	case len(srn) == 1:
		return RelationName(fmt.Sprintf("%s%s%s", irn, SubjectRelationSeparator, srn[0]))
	default:
		panic("too many subject relations")
	}
}

func PermForRel(rn RelationName) RelationName {
	return RelationName(fmt.Sprintf("%s%s", GeneratedPermissionPrefix, rn))
}

func subjs(o *Object, rn RelationName) []ObjectName {
	if o.HasRelation(rn) {
		return o.Relations[rn].AllTypes()
	}

	return o.Permissions[rn].SubjectTypes
}
