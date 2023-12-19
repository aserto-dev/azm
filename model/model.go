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

func New(r io.Reader) (*Model, error) {
	m := Model{}
	dec := json.NewDecoder(r)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&m); err != nil {
		return nil, err
	}
	return &m, nil
}

type ObjectID string

func (id ObjectID) String() string {
	return string(id)
}

func (id ObjectID) IsWildcard() bool {
	return id == "*"
}

type relation struct {
	on  ObjectName
	oid ObjectID
	rn  RelationName
	sn  ObjectName
	sid ObjectID
	srn RelationName
}

func (r *relation) String() string {
	srn := ""
	if r.srn != "" {
		srn = "#" + r.srn.String()
	}

	return fmt.Sprintf("%s:%s#%s@%s:%s%s", r.on, r.oid, r.rn, r.sn, r.sid, srn)
}

type objSet set.Set[ObjectName]
type relSet set.Set[RelationRef]

func (m *Model) GetGraph() *graph.Graph {
	grph := graph.NewGraph()
	for objectName := range m.Objects {
		grph.AddNode(string(objectName))
	}
	for objectName, obj := range m.Objects {
		for relName, rel := range obj.Relations {
			for _, rl := range rel.Union {
				if rl.IsDirect() {
					grph.AddEdge(string(objectName), string(rl.Object), string(relName))
				} else if rl.IsSubject() {
					grph.AddEdge(string(objectName), string(rl.Object), string(relName))
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

	// Pass 3: validate all arrow operators in permissions. This requires that all relations have already been resolved.
	if err := m.validatePermissions(); err != nil {
		return derr.ErrInvalidArgument.Err(err)
	}

	// Pass 4: resolve all permissions to a set of possible subject types.
	if err := m.resolvePermissions(); err != nil {
		return derr.ErrInvalidArgument.Err(err)
	}

	return nil
}

func (m *Model) ValidateRelation(on ObjectName, oid ObjectID, rn RelationName, sn ObjectName, sid ObjectID, srn RelationName) error {
	rel := &relation{on, oid, rn, sn, sid, srn}

	if oid.IsWildcard() {
		return derr.ErrInvalidRelation.Msgf("[%s] object id cannot be a wildcard", rel)
	}

	o := m.Objects[on]
	if o == nil {
		return derr.ErrInvalidRelation.Err(derr.ErrObjectTypeNotFound.Msgf("%s", on)).Msgf("[%s]", rel)
	}

	r := o.Relations[rn]
	if r == nil {
		return derr.ErrInvalidRelation.Err(derr.ErrRelationTypeNotFound.Msgf("%s:%s", on, rn)).Msgf("[%s]", rel)
	}

	// Find all valid assignments for the given subject type.
	refs := lo.Filter(r.Union, func(rr *RelationRef, _ int) bool {
		return rr.Object == rel.sn
	})

	if len(refs) == 0 {
		return derr.ErrInvalidRelation.Msgf("[%s] subject type '%s' is not valid for relation '%s:%s'", rel, rel.sn, on, rn)
	}

	if rel.sid.IsWildcard() {
		// Wildcard assignment.
		if rel.srn != "" {
			return derr.ErrInvalidRelation.Msgf("[%s] wildcard assignment cannot include subject relation", rel)
		}

		if !lo.ContainsBy(refs, func(rr *RelationRef) bool { return rr.IsWildcard() }) {
			return derr.ErrInvalidRelation.Msgf(
				"[%s] wildcard assignment of '%s' are not allowed on relation '%s:%s'",
				rel, sn, on, rn,
			)
		}
	}

	assignment := RelationRef{Object: sn, Relation: srn}
	if !lo.ContainsBy(refs, func(rr *RelationRef) bool { return *rr == assignment }) {
		return derr.ErrInvalidRelation.Msgf("[%s] invalid assignment", rel)
	}

	return nil
}

func (m *Model) Diff(newModel *Model) *diff.Diff {
	// newmodel - m => additions
	added := newModel.subtract(m)
	// m - newModel => deletions
	deleted := m.subtract(newModel)

	return &diff.Diff{Added: *added, Removed: *deleted}
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
