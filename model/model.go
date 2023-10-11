package model

import (
	"bytes"
	"encoding/json"
	"io"
	"time"

	"github.com/aserto-dev/azm/graph"
)

const ModelVersion int = 1

type Model struct {
	Version  int                    `json:"version"`
	Objects  map[ObjectName]*Object `json:"types"`
	Metadata *Metadata              `json:"metadata"`
}

type ObjectName Identifier
type RelationName Identifier
type PermissionName Identifier

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

func (m *Model) GetGraph() *graph.Graph {
	grph := graph.NewGraph()
	for objectName := range m.Objects {
		grph.AddNode(string(objectName))
	}
	for objectName, obj := range m.Objects {
		for relName, rel := range obj.Relations {
			for _, rl := range rel {
				if string(rl.Direct) != "" {
					grph.AddEdge(string(objectName), string(rl.Direct), string(relName))
				} else if rl.Subject != nil {
					grph.AddEdge(string(objectName), string(rl.Subject.Object), string(relName))
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
	// newModel - m => additions
	added := newModel.subtract(m)
	// m - newModel => deletions
	deleted := m.subtract(newModel)

	return &Diff{Added: *added, Removed: *deleted}
}

func (m *Model) subtract(newModel *Model) *Changes {
	changes := &Changes{
		Objects:   make([]ObjectName, 0),
		Relations: make(map[ObjectName][]RelationName),
	}

	if m == nil {
		return changes
	}

	if newModel == nil {
		for objName := range m.Objects {
			changes.Objects = append(changes.Objects, objName)
		}
		return changes
	}

	for objName, obj := range m.Objects {
		if newModel.Objects[objName] == nil {
			changes.Objects = append(changes.Objects, objName)
		} else {
			for relName := range obj.Relations {
				if newModel.Objects[objName].Relations[relName] == nil {
					changes.Relations[objName] = append(changes.Relations[objName], relName)
				}
			}
		}
	}

	return changes
}
