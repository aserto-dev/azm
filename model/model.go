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
	"github.com/samber/lo"
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
	BaseTypes []ObjectName `json:"base_types,omitempty"`
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
//   - Within an object, a permission cannot share the same name as a relation.
//   - Direct relations must reference existing objects .
//   - Wildcard relations must reference existing objects.
//   - Subject relations must reference existing object#relation pairs.
//   - Arrow permissions (relation->rel_or_perm) must reference existing relations/permissions.
func (m *Model) Validate() error {
	if err := m.validateUniqueNames(); err != nil {
		// if there are name collisions, we can't validate relations or permissions.
		return err
	}

	var errs error
	for on, o := range m.Objects {
		if err := m.validateObjectRels(on, o); err != nil {
			errs = multierror.Append(errs, err)

			// if there are relation errors, we can't validate permissions.
			continue
		}

		if err := m.validateObjectPerms(on, o); err != nil {
			errs = multierror.Append(errs, err)
		}
	}

	if errs != nil {
		return derr.ErrInvalidArgument.Err(errs)
	}
	return nil
}

func (m *Model) validateUniqueNames() error {
	var errs error

	for on, o := range m.Objects {
		rels := lo.Map(lo.Keys(o.Relations), func(rn RelationName, _ int) string {
			return string(rn)
		})
		perms := lo.Map(lo.Keys(o.Permissions), func(pn PermissionName, _ int) string {
			return string(pn)
		})

		rpCollisions := lo.Intersect(rels, perms)
		for _, collision := range rpCollisions {
			errs = multierror.Append(errs, derr.ErrInvalidPermission.Msgf(
				"permission name '%[1]s:%[2]s' conflicts with '%[1]s:%[2]s' relation", on, collision),
			)
		}
	}

	if errs != nil {
		return derr.ErrInvalidArgument.Err(errs)
	}
	return nil
}

func (m *Model) validateObjectRels(on ObjectName, o *Object) error {
	var errs error
	for rn, rs := range o.Relations {
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

	return errs
}

func (m *Model) validateObjectPerms(on ObjectName, o *Object) error {
	var errs error
	for pn, p := range o.Permissions {
		var refs []*RelationRef

		switch {
		case p.Union != nil:
			refs = p.Union
		case p.Intersection != nil:
			refs = p.Intersection
		case p.Exclusion != nil:
			refs = []*RelationRef{p.Exclusion.Include, p.Exclusion.Exclude}
		}

		for _, ref := range refs {
			var bases []ObjectName
			switch ref.Base {
			case "":
				bases = []ObjectName{on}
			default:
				rel := o.Relations[ref.Base]
				if rel == nil {
					errs = multierror.Append(errs, derr.ErrInvalidPermission.Msgf(
						"permission '%s:%s' references undefined relation type '%s:%s'", on, pn, on, ref.Base),
					)
					continue
				}

				for _, r := range rel {
					bases = lo.Uniq(append(bases, m.relationSubjectTypes(r)...))
				}
			}

			for _, base := range bases {
				_, foundRelation := m.Objects[base].Relations[RelationName(ref.RelOrPerm)]
				_, foundPermission := m.Objects[base].Permissions[PermissionName(ref.RelOrPerm)]
				if !(foundRelation || foundPermission) {
					switch base {
					case on:
						errs = multierror.Append(errs, derr.ErrInvalidPermission.Msgf(
							"permission '%s:%s' references undefined relation or permission '%s:%s'", on, pn, base, ref.RelOrPerm),
						)
					default:
						arrow := fmt.Sprintf("%s->%s", ref.Base, ref.RelOrPerm)
						errs = multierror.Append(errs, derr.ErrInvalidPermission.Msgf(
							"permission '%s:%s' references '%s', which can resolve to undefined relation or permission '%s:%s' ",
							on, pn, arrow, base, ref.RelOrPerm,
						))
					}

					continue
				}
			}

			ref.BaseTypes = bases
		}
	}

	return errs
}

// relationSubjectTypes returns a list of all object types that can be the subject of the given relation.
func (m *Model) relationSubjectTypes(r *Relation) []ObjectName {
	switch {
	case r.Direct != "":
		return []ObjectName{r.Direct}
	case r.Wildcard != "":
		return []ObjectName{r.Wildcard}
	case r.Subject != nil:
		subjRel := m.Objects[r.Subject.Object].Relations[r.Subject.Relation]
		return lo.Uniq(
			lo.FlatMap(subjRel, func(r *Relation, _ int) []ObjectName {
				return m.relationSubjectTypes(r)
			}),
		)
	}

	return []ObjectName{}
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
