package types

import "fmt"

type ObjectName Identifier
type RelationName Identifier

func (on ObjectName) String() string {
	return string(on)
}

func (rn RelationName) String() string {
	return string(rn)
}

type Object struct {
	Relations   map[RelationName]*Relation   `json:"relations,omitempty"`
	Permissions map[RelationName]*Permission `json:"permissions,omitempty"`
}

func (o *Object) HasRelation(name RelationName) bool {
	if o == nil {
		return false
	}

	return o.Relations[name] != nil
}

func (o *Object) HasPermission(name RelationName) bool {
	if o == nil {
		return false
	}
	return o.Permissions[name] != nil
}

func (o *Object) HasRelOrPerm(name RelationName) bool {
	return o.HasRelation(name) || o.HasPermission(name)
}

type Relation struct {
	Union        []*RelationRef `json:"union,omitempty"`
	SubjectTypes []ObjectName   `json:"subject_types,omitempty"`
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
	Union        []*PermissionTerm    `json:"union,omitempty"`
	Intersection []*PermissionTerm    `json:"intersection,omitempty"`
	Exclusion    *ExclusionPermission `json:"exclusion,omitempty"`

	SubjectTypes []ObjectName `json:"subject_types,omitempty"`
}

func (p *Permission) IsUnion() bool {
	return p.Union != nil
}

func (p *Permission) IsIntersection() bool {
	return p.Intersection != nil
}

func (p *Permission) IsExclusion() bool {
	return p.Exclusion != nil
}

func (p *Permission) Terms() []*PermissionTerm {
	var refs []*PermissionTerm

	switch {
	case p.IsUnion():
		refs = p.Union
	case p.IsIntersection():
		refs = p.Intersection
	case p.IsExclusion():
		refs = []*PermissionTerm{p.Exclusion.Include, p.Exclusion.Exclude}
	}

	return refs
}

type PermissionTerm struct {
	Base      RelationName `json:"base,omitempty"`
	RelOrPerm RelationName `json:"rel_or_perm"`

	SubjectTypes []ObjectName `json:"subject_types,omitempty"`
}

func (pr *PermissionTerm) IsArrow() bool {
	return pr.Base != ""
}

type ExclusionPermission struct {
	Include *PermissionTerm `json:"include,omitempty"`
	Exclude *PermissionTerm `json:"exclude,omitempty"`
}

type ArrowPermission struct {
	Relation   string `json:"relation,omitempty"`
	Permission string `json:"permission,omitempty"`
}
