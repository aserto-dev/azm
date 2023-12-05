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
	set "github.com/deckarep/golang-set/v2"
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
	Relations   map[RelationName]*Relation     `json:"relations,omitempty"`
	Permissions map[PermissionName]*Permission `json:"permissions,omitempty"`
}

func (o *Object) HasRelOrPerm(name string) bool {
	if o == nil {
		return false
	}
	if _, foundRelation := o.Relations[RelationName(name)]; foundRelation {
		return true
	}
	_, foundPermission := o.Permissions[PermissionName(name)]
	return foundPermission
}

type Relation struct {
	Union        []*RelationTerm `json:"union,omitempty"`
	SubjectTypes []ObjectName    `json:"subject_types,omitempty"`
}

type RelationTerm struct {
	Direct   ObjectName       `json:"direct,omitempty"`
	Subject  *SubjectRelation `json:"subject,omitempty"`
	Wildcard ObjectName       `json:"wildcard,omitempty"`
}

type RelationRef struct {
	Object   ObjectName   `json:"object,omitempty"`
	Relation RelationName `json:"relation,omitempty"`
}

type SubjectRelation struct {
	*RelationRef
	SubjectTypes []ObjectName `json:"subject_types,omitempty"`
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
			for _, rl := range rel.Union {
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
//   - Within an object, a permission cannot share the same name as a relation.
//   - Direct relations must reference existing objects .
//   - Wildcard relations must reference existing objects.
//   - Subject relations must reference existing object#relation pairs.
//   - Arrow permissions (relation->rel_or_perm) must reference existing relations/permissions.
func (m *Model) Validate() error {
	// Pass 1 (syntax): ensure no name collisions and all relations reference existing objects/relations.
	if err := m.validateReferences(); err != nil {
		return derr.ErrInvalidArgument.Err(err)
	}

	// Pass 2: resolve all relations to a set of possible subject types.
	if err := m.resolveRelations(); err != nil {
		return derr.ErrInvalidArgument.Err(err)
	}

	if err := m.resolvePermissions(); err != nil {
		return derr.ErrInvalidArgument.Err(err)
	}

	return nil
}

func (m *Model) validateReferences() error {
	var errs error

	for on, o := range m.Objects {
		validatePerms := true
		if err := m.validateUniqueNames(on, o); err != nil {
			errs = multierror.Append(errs, err)
			validatePerms = false
		}

		if err := m.validateObjectRels(on, o); err != nil {
			errs = multierror.Append(errs, err)
		}

		if validatePerms {
			if err := m.validateObjectPerms(on, o); err != nil {
				errs = multierror.Append(errs, err)
			}
		}
	}

	return errs
}

func (m *Model) validateUniqueNames(on ObjectName, o *Object) error {
	rels := lo.Map(lo.Keys(o.Relations), func(rn RelationName, _ int) string {
		return string(rn)
	})
	perms := lo.Map(lo.Keys(o.Permissions), func(pn PermissionName, _ int) string {
		return string(pn)
	})

	rpCollisions := lo.Intersect(rels, perms)

	var errs error
	for _, collision := range rpCollisions {
		errs = multierror.Append(errs, derr.ErrInvalidPermission.Msgf(
			"permission name '%[1]s:%[2]s' conflicts with '%[1]s:%[2]s' relation", on, collision),
		)
	}

	return errs
}

func (m *Model) validateObjectRels(on ObjectName, o *Object) error {
	var errs error
	for rn, rs := range o.Relations {
		for _, r := range rs.Union {
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
			default:
				errs = multierror.Append(errs, derr.ErrInvalidRelation.Msgf(
					"relation '%s:%s' has no definition", on, rn),
				)
			}
		}
	}

	return errs
}

func (m *Model) validateObjectPerms(on ObjectName, o *Object) error {
	var errs error
	for pn, p := range o.Permissions {
		refs := p.Refs()
		if len(refs) == 0 {
			errs = multierror.Append(errs, derr.ErrInvalidPermission.Msgf(
				"permission '%s:%s' has no definition", on, pn),
			)
			continue
		}

		for _, ref := range refs {
			if ref.Base != "" {
				// validate that the base relation exists on this object type.
				// at this stage we don't yet resolve the relation to a set of subject types.
				if _, hasRelation := o.Relations[ref.Base]; !hasRelation {
					errs = multierror.Append(errs, derr.ErrInvalidPermission.Msgf(
						"permission '%s:%s' references undefined relation type '%s:%s'", on, pn, on, ref.Base),
					)
				}
			}
		}
	}

	return errs
}

func (m *Model) resolveRelations() error {
	var errs error
	for on, o := range m.Objects {
		for rn, r := range o.Relations {
			seen := set.NewSet(RelationRef{Object: on, Relation: rn})
			subs := m.resolveRelation(r, seen)
			switch len(subs) {
			case 0:
				errs = multierror.Append(errs, derr.ErrInvalidRelation.Msgf(
					"relation '%s:%s' is circular and does not resolve to any object types", on, rn),
				)
			default:
				r.SubjectTypes = lo.Uniq(subs)
			}
		}
	}

	return errs
}

type RelSet set.Set[RelationRef]

func (m *Model) resolveRelation(r *Relation, seen RelSet) []ObjectName {
	if len(r.SubjectTypes) > 0 {
		// already resolved
		return r.SubjectTypes
	}

	subjectTypes := []ObjectName{}
	for _, rt := range r.Union {
		switch {
		case rt.Direct != "":
			subjectTypes = append(subjectTypes, rt.Direct)
		case rt.Wildcard != "":
			subjectTypes = append(subjectTypes, rt.Wildcard)
		case rt.Subject != nil:
			if !seen.Contains(*rt.Subject.RelationRef) {
				seen.Add(*rt.Subject.RelationRef)
				subjectTypes = append(subjectTypes, m.resolveRelation(
					m.Objects[rt.Subject.Object].Relations[rt.Subject.Relation],
					seen)...,
				)
			}

		}
	}
	return subjectTypes
}

func (m *Model) resolvePermissions() error {
	var errs error

	for on, o := range m.Objects {
		for pn, p := range o.Permissions {
			if err := m.resolvePermission(on, pn, p); err != nil {
				errs = multierror.Append(errs, err)
			}
		}
	}

	return errs
}

func (m *Model) resolvePermission(on ObjectName, pn PermissionName, p *Permission) error {
	var errs error
	for _, ref := range p.Refs() {
		bases := []ObjectName{}
		switch ref.Base {
		case "":
			bases = append(bases, on)
		default:
			// relations are already resolved at this point.
			bases = append(bases, m.Objects[on].Relations[ref.Base].SubjectTypes...)
		}

		for _, base := range bases {
			if !m.Objects[base].HasRelOrPerm(ref.RelOrPerm) {
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
	}

	return errs
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
