package model

import (
	"bytes"
	"encoding/json"
	"io"
	"time"
)

const ModelVersion int = 1

type Model struct {
	Version  int                    `json:"version"`
	Objects  map[ObjectName]*Object `json:"types"`
	Metadata *Metadata              `json:"metadata"`
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

type Metadata struct {
	UpdatedAt time.Time `json:"updated_at"`
	ETag      string    `json:"etag"`
}

type Diff struct {
	Added   Changes
	Removed Changes
}

type Changes struct {
	Objects   []ObjectName
	Relations map[ObjectName][]RelationName
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
	// m - newmodel => deletions
	deleted := m.subtract(newModel)

	return &Diff{Added: *added, Removed: *deleted}
}

func (m *Model) subtract(newModel *Model) *Changes {
	chgs := &Changes{
		Objects:   make([]ObjectName, 0),
		Relations: make(map[ObjectName][]RelationName),
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
			for relname := range obj.Relations {
				if newModel.Objects[objName].Relations[relname] == nil {
					chgs.Relations[objName] = append(chgs.Relations[objName], relname)
				}
			}
		}
	}

	return chgs
}
