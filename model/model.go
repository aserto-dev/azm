package model

const ModelVersion int = 1

type Model struct {
	Version int                    `json:"version"`
	Objects map[ObjectName]*Object `json:"types"`
}

type ObjectName string
type RelationName string
type PermissionName string

type Object struct {
	Relations   map[RelationName][]*Relation   `json:"relations,omitempty"`
	Permissions map[PermissionName]*Permission `json:"permissions,omitempty"`
}

type Relation struct {
	Direct   ObjectName       `json:"direct,omitempty"`
	Subject  *SubjectRelation `json:"subject,omitempty"`
	Wildcard ObjectName       `json:"wildcard,omitempty"`
}

type SubjectRelation struct {
	Object   ObjectName   `json:"object,omitempty"`
	Relation RelationName `json:"relation,omitempty"`
}

type Permission struct {
	Union        []string             `json:"union,omitempty"`
	Intersection []string             `json:"intersection,omitempty"`
	Exclusion    *ExclusionPermission `json:"exclusion,omitempty"`
	Arrow        *ArrowPermission     `json:"arrow,omitempty"`
}

type ExclusionPermission struct {
	Base     string `json:"base,omitempty"`
	Subtract string `json:"subtract,omitempty"`
}

type ArrowPermission struct {
	Relation   string `json:"relation,omitempty"`
	Permission string `json:"permission,omitempty"`
}
