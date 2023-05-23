package azm

import (
	"fmt"
	"io"
	"strings"

	dsr "github.com/aserto-dev/go-directory/aserto/directory/reader/v2"
	dsw "github.com/aserto-dev/go-directory/aserto/directory/writer/v2"
	"github.com/pkg/errors"
)

// Read model instance from manifest.
type Reader interface {
	Read(io.Reader) (*Model, error)
}

// Write model instance to manifest.
type Writer interface {
	Write(io.Writer, *Model) error
}

// Load model instance from directory store.
type Loader interface {
	Load(dsr.ReaderClient) (*Model, error)
}

// Save model instance to directory store.
type Saver interface {
	Save(dsw.WriterClient, *Model) error
}

type Model struct {
	Model       *ModelInfo    `yaml:"model,omitempty" json:"model,omitempty"`
	ObjectTypes []*ObjectType `yaml:"object_types,omitempty" json:"object_types,omitempty"`
	objects     map[string]int
}

type ModelInfo struct {
	Name          string `yaml:"name,omitempty" json:"name,omitempty"`
	SchemaVersion int    `yaml:"schema_version" json:"schema_version"`
}

type ObjectType struct {
	Name        string            `yaml:"name" json:"name"`
	Relations   []*RelationType   `yaml:"relations,omitempty" json:"relations,omitempty"`
	Permissions []*PermissionType `yaml:"permissions,omitempty" json:"permissions,omitempty"`
	relations   map[string]int
	permissions map[string]int
}

type RelationType struct {
	Name       string   `yaml:"name" json:"name"`
	Definition string   `yaml:"definition,omitempty" json:"definition,omitempty"`
	Operator   Operator `yaml:"operator,omitempty" json:"operator,omitempty"`
	Relations  []string `yaml:"relations,omitempty" json:"relations,omitempty"`
}

type PermissionType struct {
	Name        string   `yaml:"name" json:"name"`                                   // name of assignable relation (permission)
	Definition  string   `yaml:"definition,omitempty" json:"definition,omitempty"`   // manifest definition string (regex)
	Operator    Operator `yaml:"operator,omitempty" json:"operator,omitempty"`       // operator kind
	Permissions []string `yaml:"permissions,omitempty" json:"permissions,omitempty"` // list of permission
}

func New(name string, schemaVersion int) *Model {
	return &Model{
		Model: &ModelInfo{
			Name:          name,
			SchemaVersion: schemaVersion,
		},
		ObjectTypes: []*ObjectType{},
		objects:     map[string]int{},
	}
}

func (m *Model) ObjectTypeExists(name string) bool {
	_, ok := m.objects[name]
	return ok
}

func (o *ObjectType) RelationExists(name string) bool {
	_, ok := o.relations[name]
	return ok
}

func (o *ObjectType) PermissionExists(name string) bool {
	_, ok := o.permissions[name]
	return ok
}

func (m *Model) ResolveRelation(objectType, relation string) ([]string, error) {
	resp := []string{}

	objSlot, ok := m.objects[objectType]
	if !ok {
		return resp, errors.Wrapf(ErrObjectTypeNotFound, "%s", objectType)
	}
	objType := m.ObjectTypes[objSlot]

	relSlot, ok := objType.relations[relation]
	if !ok {
		return resp, errors.Wrapf(ErrRelationTypeNotFound, "%s#%s", objectType, relation)
	}

	relType := objType.Relations[relSlot]

	for _, rel := range objType.Relations {
		for _, r := range rel.Relations {
			if strings.EqualFold(relType.Name, r) {
				resp = append(resp, rel.Name)
			}
		}
	}

	return resp, nil
}

func (m *Model) ResolvePermission(objectType, permission string) ([]string, error) {
	resp := []string{}

	objSlot, ok := m.objects[objectType]
	if !ok {
		return resp, errors.Wrapf(ErrObjectTypeNotFound, "%s", objectType)
	}
	objType := m.ObjectTypes[objSlot]

	permSlot, ok := objType.permissions[permission]
	if !ok {
		return resp, errors.Wrapf(ErrPermissionTypeNotFound, "%s#%s", objectType, permission)
	}

	permType := objType.Permissions[permSlot]

	for _, relType := range objType.Relations {
		for _, r := range relType.Relations {
			if strings.EqualFold(permType.Name, r) {
				resp = append(resp, relType.Name)
			}
		}
	}

	return resp, nil
}

func (m *Model) SetObjectType(i *ObjectType) (*ObjectType, error) {
	if i == nil || i.Name == "" {
		return nil, ErrInvalidObjectType
	}

	if i.Relations == nil {
		i.Relations = []*RelationType{}
	}

	if i.Permissions == nil {
		i.Permissions = []*PermissionType{}
	}

	if i.relations == nil {
		i.relations = map[string]int{}
	}

	if i.permissions == nil {
		i.permissions = map[string]int{}
	}

	slot, ok := m.objects[i.Name]
	if !ok {
		if i.Relations == nil {
			i.Relations = []*RelationType{}
		}

		if i.Permissions == nil {
			i.Permissions = []*PermissionType{}
		}

		slot = len(m.ObjectTypes)
		m.objects[i.Name] = slot
		m.ObjectTypes = append(m.ObjectTypes, i)
	}

	return m.ObjectTypes[slot], nil
}

func (o *ObjectType) SetRelationType(i *RelationType) (*RelationType, error) {
	if i == nil || i.Name == "" {
		return nil, ErrInvalidRelationType
	}

	slot, ok := o.relations[i.Name]
	if !ok {
		if o.Relations == nil {
			o.Relations = []*RelationType{}
		}
		slot = len(o.relations)
		o.relations[i.Name] = slot
		o.Relations = append(o.Relations, i)
	}

	return o.Relations[slot], nil
}

func (o *ObjectType) SetPermissionType(i *PermissionType) (*PermissionType, error) {
	if i == nil || i.Name == "" {
		return nil, ErrInvalidPermissionType
	}

	slot, ok := o.permissions[i.Name]
	if !ok {
		if o.Permissions == nil {
			o.Permissions = []*PermissionType{}
		}
		slot = len(o.permissions)
		o.permissions[i.Name] = slot
		o.Permissions = append(o.Permissions, i)
	}

	return o.Permissions[slot], nil
}

func (m *Model) Expand() {
	for _, ot := range m.ObjectTypes {
		for _, rt := range ot.Relations {
			for _, r := range rt.Relations {
				fmt.Printf("r:%s#%s|%s [%s]\n", ot.Name, rt.Name, r, strings.Join(rt.Relations, " | "))
			}
		}
		for _, pt := range ot.Permissions {
			for _, p := range pt.Permissions {
				fmt.Printf("p:%s#%s|%s [%s]\n", ot.Name, pt.Name, p, strings.Join(pt.Permissions, " | "))
			}
		}
	}
}
