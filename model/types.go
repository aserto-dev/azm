package model

import "fmt"

type ObjectName Identifier
type RelationName Identifier
type PermissionName Identifier

func (on ObjectName) String() string {
	return string(on)
}

func (rn RelationName) String() string {
	return string(rn)
}

func (pn PermissionName) String() string {
	return string(pn)
}

func (pn PermissionName) RN() RelationName {
	return RelationName(pn)
}

type Object struct {
	Relations   map[RelationName]*Relation     `json:"relations,omitempty"`
	Permissions map[PermissionName]*Permission `json:"permissions,omitempty"`
}

func (o *Object) HasRelation(name RelationName) bool {
	if o == nil {
		return false
	}

	return o.Relations[name] != nil
}

func (o *Object) HasPermission(name PermissionName) bool {
	if o == nil {
		return false
	}
	return o.Permissions[name] != nil
}

func (o *Object) HasRelOrPerm(name string) bool {
	return o.HasRelation(RelationName(name)) || o.HasPermission(PermissionName(name))
}

type Relation struct {
	Union        []*RelationTerm `json:"union,omitempty"`
	SubjectTypes []ObjectName    `json:"subject_types,omitempty"`
}

type RelationTerm struct {
	*RelationRef
	SubjectTypes []ObjectName `json:"subject_types,omitempty"`
}

type RelationRef struct {
	Object   ObjectName   `json:"object,omitempty"`
	Relation RelationName `json:"relation,omitempty"`
}

type RelationAssignment int

const (
	RelationAssignmentUnknown RelationAssignment = iota
	RelationAssignmentDirect
	RelationAssignmentSubject
	RelationAssignmentWildcard
)

func NewRelationRef(on ObjectName, rn RelationName) *RelationRef {
	return &RelationRef{Object: on, Relation: rn}
}

func (rr *RelationRef) String() string {
	if rr.Relation == "" {
		return string(rr.Object)
	}
	return fmt.Sprintf("%s:%s", rr.Object, rr.Relation)
}

func (rr *RelationRef) Assignment() RelationAssignment {
	if rr == nil {
		return RelationAssignmentUnknown
	}

	switch {
	case rr.Relation == "*":
		return RelationAssignmentWildcard
	case rr.Relation != "":
		return RelationAssignmentSubject
	case rr.Object != "":
		return RelationAssignmentDirect
	default:
		return RelationAssignmentUnknown
	}
}

func (rr *RelationRef) IsDirect() bool {
	return rr.Assignment() == RelationAssignmentDirect
}

func (rr *RelationRef) IsWildcard() bool {
	return rr.Assignment() == RelationAssignmentWildcard
}

func (rr *RelationRef) IsSubject() bool {
	return rr.Assignment() == RelationAssignmentSubject
}

type Permission struct {
	Union        []*PermissionRef     `json:"union,omitempty"`
	Intersection []*PermissionRef     `json:"intersection,omitempty"`
	Exclusion    *ExclusionPermission `json:"exclusion,omitempty"`
}

func (p *Permission) Refs() []*PermissionRef {
	var refs []*PermissionRef

	switch {
	case p.Union != nil:
		refs = p.Union
	case p.Intersection != nil:
		refs = p.Intersection
	case p.Exclusion != nil:
		refs = []*PermissionRef{p.Exclusion.Include, p.Exclusion.Exclude}
	}

	return refs
}

type PermissionRef struct {
	Base      RelationName `json:"base,omitempty"`
	RelOrPerm string       `json:"rel_or_perm"`
	BaseTypes []ObjectName `json:"base_types,omitempty"`
}

type ExclusionPermission struct {
	Include *PermissionRef `json:"include,omitempty"`
	Exclude *PermissionRef `json:"exclude,omitempty"`
}

type ArrowPermission struct {
	Relation   string `json:"relation,omitempty"`
	Permission string `json:"permission,omitempty"`
}
