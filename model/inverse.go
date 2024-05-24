package model

import (
	"fmt"

	set "github.com/deckarep/golang-set/v2"
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
		isrn := InverseRelation(on, rn, rr.Relation)
		i.im.Objects[rr.Object].Relations[isrn] = &Relation{Union: []*RelationRef{{Object: on}}}
		if rr.IsSubject() {
			// add a synthetic permission to reverse the expansion of the subject relation
			srel := i.m.Objects[rr.Object].Relations[rr.Relation]

			for _, subj := range srel.AllRefs() {
				ipr := InverseRelation(on, rn, subj.Relation)
				ipn := PermForRel(ipr)
				p := permissionOrNew(i.im.Objects[subj.Object], ipn, permissionKindUnion)
				i.addSubstitution(ipr, ipn)
				if _, ok := unionObjs[subj.Object]; ok {
					p.AddTerm(&PermissionTerm{RelOrPerm: InverseRelation(on, rn, subj.Relation)})
				}
				p.AddTerm(&PermissionTerm{Base: InverseRelation(rr.Object, rr.Relation, subj.Relation), RelOrPerm: isrn})
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
	var typeSet relSet

	switch {
	case p.IsUnion():
		typeSet = set.NewSet(p.Types()...)
	case p.IsIntersection():
		typeSet = lo.Reduce(p.Terms(), func(acc relSet, pt *PermissionTerm, i int) relSet {
			s := set.NewSet(pt.Types()...)
			if i == 0 {
				return s
			}
			if s.IsEmpty() {
				return acc
			}
			return acc.Intersect(s)
		}, nil)
	case p.IsExclusion():
		typeSet = set.NewSet(p.Exclusion.Include.Types()...)
	}

	for _, pt := range p.Terms() {
		newTermInverter(i, on, pn, o, p, pt, typeSet).invert()
	}
}

func (i *inverter) irelSub(on ObjectName, rn RelationName, srn ...RelationName) RelationName {
	return i.sub(InverseRelation(on, rn, srn...))
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

type termInverter struct {
	inv      *inverter
	objName  ObjectName
	permName RelationName
	obj      *Object
	perm     *Permission
	term     *PermissionTerm
	typeSet  relSet
	kind     permissionKind
}

func newTermInverter(i *inverter, on ObjectName, pn RelationName, o *Object, p *Permission, pt *PermissionTerm, typeSet relSet) *termInverter {
	return &termInverter{
		inv:      i,
		objName:  on,
		permName: pn,
		obj:      o,
		perm:     p,
		term:     pt,
		typeSet:  typeSet,
		kind:     kindOf(p),
	}
}

func (ti *termInverter) invert() {
	switch {
	case ti.term.IsArrow():
		// create a subject relation to expand the recursive permission
		ti.invertArrow()

	case ti.obj.HasRelation(ti.term.RelOrPerm):
		ti.invertRelation()

	case ti.obj.HasPermission(ti.term.RelOrPerm):
		ti.invertPermission()
	}
}

func (ti *termInverter) invertArrow() {
	itip := ti.inv.irelSub(ti.objName, ti.term.Base)
	baseRel := ti.obj.Relations[ti.term.Base]

	typeRefs := set.NewSet(ti.perm.Types()...)
	typeRefs.Intersect(ti.typeSet).Each(func(rr RelationRef) bool {
		iName := InverseRelation(ti.objName, ti.permName, rr.Relation)
		iPerm := permissionOrNew(ti.inv.im.Objects[rr.Object], iName, ti.kind)
		for _, baseRR := range baseRel.Types() {
			term := &PermissionTerm{Base: ti.inv.sub(InverseRelation(baseRR.Object, ti.term.RelOrPerm, rr.Relation)), RelOrPerm: itip}
			iPerm.AddTerm(term)

			rel := relationOrNew(ti.inv.im.Objects[baseRR.Object], itip)
			rel.AddRef(&RelationRef{Object: baseRR.Object, Relation: itip})
		}
		return false // resume iteration
	})
}

func (ti *termInverter) invertRelation() {
	typeRefs := set.NewSet(ti.obj.Relations[ti.term.RelOrPerm].Types()...)
	typeRefs.Intersect(ti.typeSet).Each(func(rr RelationRef) bool {
		iName := InverseRelation(ti.objName, ti.permName, rr.Relation)
		ip := permissionOrNew(ti.inv.im.Objects[rr.Object], iName, ti.kind)
		ip.AddTerm(&PermissionTerm{RelOrPerm: ti.inv.irelSub(ti.objName, ti.term.RelOrPerm, rr.Relation)})
		return false // resume iteration
	})
}

func (ti *termInverter) invertPermission() {
	typeRefs := set.NewSet(types(ti.obj, ti.term.RelOrPerm)...)
	typeRefs.Intersect(ti.typeSet).Each(func(rr RelationRef) bool {
		iName := InverseRelation(ti.objName, ti.permName, rr.Relation)
		ip := permissionOrNew(ti.inv.im.Objects[rr.Object], iName, ti.kind)
		ip.AddTerm(&PermissionTerm{RelOrPerm: ti.inv.irelSub(ti.objName, ti.term.RelOrPerm, rr.Relation)})
		return false // resume iteration
	})
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

func types(o *Object, rn RelationName) RelationRefs {
	if o.HasRelation(rn) {
		return o.Relations[rn].Types()
	}

	return o.Permissions[rn].Types()
}
