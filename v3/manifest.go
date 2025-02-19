package v3

import (
	"strconv"

	"github.com/aserto-dev/azm"
	"gopkg.in/yaml.v3"
)

const SupportedSchemaVersion int = 3

type Manifest struct {
	ModelInfo   *ModelInfo                     `yaml:"model"`
	ObjectTypes map[ObjectTypeName]*ObjectType `yaml:"types"`
}

type SchemaVersion int

type ModelInfo struct {
	Version SchemaVersion `yaml:"version"`
}

type ObjectTypeName string

type ObjectType struct {
	Relations   map[RelationName]string   `yaml:"relations,omitempty"`
	Permissions map[PermissionName]string `yaml:"permissions,omitempty"`
}

type RelationName string
type PermissionName string

func (v *SchemaVersion) UnmarshalYAML(value *yaml.Node) error {
	version, err := strconv.Atoi(value.Value)
	if err != nil {
		return err
	}

	if version != SupportedSchemaVersion {
		return azm.ErrInvalidSchemaVersion.Msgf("%d", version)
	}

	*v = SchemaVersion(version)

	return nil
}
