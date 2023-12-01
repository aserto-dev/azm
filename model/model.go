package model

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/aserto-dev/azm/graph"
	"github.com/aserto-dev/azm/model/diff"
	"github.com/aserto-dev/go-directory/pkg/derr"
	"github.com/hashicorp/go-multierror"
)

const ModelVersion int = 2

type Model struct {
	Version  int                    `json:"version"`
	Objects  map[ObjectName]*Object `json:"types"`
	Metadata *Metadata              `json:"metadata"`
}

type Metadata struct {
	UpdatedAt time.Time `json:"updated_at"`
	ETag      string    `json:"etag"`
}

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
	Union        []*RelationRef       `json:"union,omitempty"`
	Intersection []*RelationRef       `json:"intersection,omitempty"`
	Exclusion    *ExclusionPermission `json:"exclusion,omitempty"`
}

type RelationRef struct {
	Base      RelationName `json:"base,omitempty"`
	RelOrPerm string       `json:"rel_or_perm"`
}

type ExclusionPermission struct {
	Include *RelationRef `json:"include,omitempty"`
	Exclude *RelationRef `json:"exclude,omitempty"`
}

type ArrowPermission struct {
	Relation   string `json:"relation,omitempty"`
	Permission string `json:"permission,omitempty"`
}

type ObjectRelation struct {
	Object   ObjectName   `json:"object"`
	Relation RelationName `json:"relation,omitempty"`
}

func NewObjectRelation(on ObjectName, rn RelationName) *ObjectRelation {
	return &ObjectRelation{Object: on, Relation: rn}
}

func (or ObjectRelation) String() string {
	return fmt.Sprintf("%s:%s", or.Object, or.Relation)
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

func (m *Model) Diff(newModel *Model) *diff.Diff {
	// newmodel - m => additions
	added := newModel.subtract(m)
	// m - newModel => deletions
	deleted := m.subtract(newModel)

	return &diff.Diff{Added: *added, Removed: *deleted}
}

// Validate enforces the model's internal consistency.
//
// It enforces the following rules:
//   - A relation cannot share the same name as an object.
//   - Direct relations must reference an existing object.
//   - Wildcard relations must reference an existing object.
//   - Subject relations must reference an existing object#relation pair.
//   - Arrow permissions (relation->rel_or_perm) must reference existing relations/permissions.
func (m *Model) Validate() error {
	var errs error
	for on, o := range m.Objects {
		for rn, rs := range o.Relations {
			if _, ok := m.Objects[ObjectName(rn)]; ok {
				errs = multierror.Append(errs, derr.ErrInvalidRelation.Msgf("relation name '%s:%s' conflicts with object type '%s'", on, rn, rn))
			}

			for _, r := range rs {
				switch {
				case r.Direct != "":
					if _, ok := m.Objects[r.Direct]; !ok {
						errs = multierror.Append(errs, derr.ErrInvalidRelation.Msgf(
							"relation '%s:%s' references undefined object type '%s'", on, rn, r.Direct),
						)
					}
				case r.Wildcard != "":
					if _, ok := m.Objects[r.Wildcard]; !ok {
						errs = multierror.Append(errs, derr.ErrInvalidRelation.Msgf(
							"relation '%s:%s' references undefined object type '%s'", on, rn, r.Wildcard),
						)
					}
				case r.Subject != nil:
					if _, ok := m.Objects[r.Subject.Object]; !ok {
						errs = multierror.Append(errs, derr.ErrInvalidRelation.Msgf(
							"relation '%s:%s' references undefined object type '%s'", on, rn, r.Subject.Object),
						)
						break
					}

					if _, ok := m.Objects[r.Subject.Object].Relations[r.Subject.Relation]; !ok {
						errs = multierror.Append(errs, derr.ErrInvalidRelation.Msgf(
							"relation '%s:%s' references undefined relation type '%s#%s'", on, rn, r.Subject.Object, r.Subject.Relation),
						)
					}
				}
			}
		}

	}

	if errs != nil {
		return derr.ErrInvalidArgument.Err(errs)
	}
	return nil
}

func (m *Model) subtract(newModel *Model) *diff.Changes {
	chgs := &diff.Changes{
		Objects:   make([]string, 0),
		Relations: make(map[string][]string),
	}

	if m == nil {
		return chgs
	}

	if newModel == nil {
		for objName := range m.Objects {
			chgs.Objects = append(chgs.Objects, string(objName))
		}
		return chgs
	}

	for objName, obj := range m.Objects {
		if newModel.Objects[objName] == nil {
			chgs.Objects = append(chgs.Objects, string(objName))
		} else {
			for relname := range obj.Relations {
				if newModel.Objects[objName].Relations[relname] == nil {
					chgs.Relations[string(objName)] = append(chgs.Relations[string(objName)], string(relname))
				}
			}
		}
	}

	return chgs
}
