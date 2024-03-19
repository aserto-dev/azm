package model

import (
	"fmt"

	"github.com/samber/lo"
)

func (m *Model) Invert() *Model {
	return newInverter(m).invert()
}

type inverter struct {
	m     *Model
	im    *Model
	subst map[ObjectName]map[RelationName]RelationName
}

func newInverter(m *Model) *inverter {
	return &inverter{
		m: m,
		im: &Model{
			Version:  m.Version,
			Objects:  lo.MapValues(m.Objects, func(o *Object, _ ObjectName) *Object { return NewObject() }),
			Metadata: m.Metadata,
		},
		subst: map[ObjectName]map[RelationName]RelationName{},
	}
}

func (i *inverter) invert() *Model {
	for on, o := range i.m.Objects {
		for rn, r := range o.Relations {
			i.invertRelation(on, rn, r)
		}

		for pn, p := range o.Permissions {
			kind := kindOf(p)
			ipn := i.iName(on, pn)

			for _, pt := range p.Terms() {
				switch {
				case pt.IsArrow():
					baseRel := o.Relations[pt.Base]
					itip := i.iName(on, pt.Base)

					for _, subj := range p.SubjectTypes {
						ip := permissionOrNew(i.im.Objects[subj], ipn, kind)
						for _, baseSubj := range baseRel.SubjectTypes {
							ip.AddTerm(&PermissionTerm{Base: i.iName(baseSubj, pt.RelOrPerm), RelOrPerm: itip})

							// create a subject relation to expand the recursive permission
							r := relationOrNew(i.im.Objects[baseSubj], itip)
							r.AddRef(&RelationRef{Object: baseSubj, Relation: itip})
						}
					}

				case o.HasRelation(pt.RelOrPerm):
					for _, rr := range o.Relations[pt.RelOrPerm].Union {
						ip := permissionOrNew(i.im.Objects[rr.Object], ipn, kind)
						ip.AddTerm(&PermissionTerm{RelOrPerm: i.iName(on, pt.RelOrPerm)})

						if rr.IsSubject() {
							term := &PermissionTerm{Base: i.iName(rr.Object, rr.Relation), RelOrPerm: ipn}
							ip.AddTerm(term)

							for _, subj := range subjs(i.m.Objects[rr.Object], rr.Relation) {
								ip = permissionOrNew(i.im.Objects[subj], ipn, kind)
								ip.AddTerm(term)
							}
						}
					}

				case o.HasPermission(pt.RelOrPerm):
					for _, subj := range subjs(o, pt.RelOrPerm) {
						ip := permissionOrNew(i.im.Objects[subj], ipn, kind)
						ip.AddTerm(&PermissionTerm{RelOrPerm: i.iName(on, pt.RelOrPerm)})
					}
				}
			}
		}
	}

	return i.im
}

func (i *inverter) invertRelation(on ObjectName, rn RelationName, r *Relation) {
	for _, rr := range r.Union {
		irn := rrel(on, rn)
		i.im.Objects[rr.Object].Relations[irn] = &Relation{Union: []*RelationRef{{Object: on}}}
		if rr.IsSubject() && rr.Object == on {
			// add a synthetic permission to reverse the expansion of the subject relation
			ipn := rsrel(on, rn)
			for _, subj := range r.Union {
				p := permissionOrNew(i.im.Objects[subj.Object], ipn, permissionKindUnion)
				p.AddTerm(&PermissionTerm{RelOrPerm: irn})
				p.AddTerm(&PermissionTerm{Base: ipn, RelOrPerm: ipn})
				i.addSubstitution(subj.Object, rn, ipn)
			}
		}
	}
}

func (i *inverter) iName(on ObjectName, rn RelationName) RelationName {
	if pn, ok := i.subst[on][rn]; ok {
		return pn
	}

	return rrel(on, rn)
}

func (i *inverter) addSubstitution(on ObjectName, rn, pn RelationName) {
	if i.subst[on] == nil {
		i.subst[on] = map[RelationName]RelationName{}
	}

	i.subst[on][rn] = pn
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

func rrel(on ObjectName, rn RelationName) RelationName {
	return RelationName(fmt.Sprintf("%s_%s", on, rn))
}

func rsrel(on ObjectName, rn RelationName) RelationName {
	return RelationName(fmt.Sprintf("r_%s", rrel(on, rn)))
}

func subjs(o *Object, rn RelationName) []ObjectName {
	if o.HasRelation(rn) {
		return o.Relations[rn].SubjectTypes
	}

	return o.Permissions[rn].SubjectTypes
}
