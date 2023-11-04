package v3

import (
	"strconv"
	"strings"

	"github.com/aserto-dev/azm"
	"gopkg.in/yaml.v3"
)

const SupportedSchemaVersion int = 3

const (
	UnionIdentifier        string = "|"
	IntersectionIdentifier string = "&"
	ExclusionIdentifier    string = " - "
	RelationIdentifier     string = "#"
	WildcardIdentifier     string = ":*"
	ArrowIdentifier        string = "->"
)

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
	Relations   map[RelationName]RelationDefinition   `yaml:"relations,omitempty"`
	Permissions map[PermissionName]PermissionOperator `yaml:"permissions,omitempty"`
}

type RelationName string

type RelationDefinition struct {
	Definition []interface {
		isRelationDefinition()
	} `yaml:"definition"`
}

type DirectRelation struct {
	ObjectType string `yaml:"direct_relation"`
}

func (*DirectRelation) isRelationDefinition() {}

type SubjectRelation struct {
	ObjectType string `yaml:"object_relation"`
	Relation   string `yaml:"subject_relation"`
}

func (*SubjectRelation) isRelationDefinition() {}

type WildcardRelation struct {
	ObjectType string `yaml:"wildcard_relation"`
}

func (*WildcardRelation) isRelationDefinition() {}

type PermissionName string

type PermissionOperator struct {
	Operator interface {
		isPermissionOperator()
	} `yaml:"operator"`
}

type UnionOperator struct {
	Union []string `yaml:"union"`
}

func (*UnionOperator) isPermissionOperator() {}

type IntersectionOperator struct {
	Intersection []string `yaml:"intersect"`
}

func (*IntersectionOperator) isPermissionOperator() {}

type ExclusionOperator struct {
	Base     string `yaml:"base"`
	Subtract string `yaml:"subtract"`
}

func (*ExclusionOperator) isPermissionOperator() {}

type ArrowOperator struct {
	Relation   string `yaml:"relation"`
	Permission string `yaml:"permission"`
}

func (*ArrowOperator) isPermissionOperator() {}

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

func (r *RelationDefinition) UnmarshalYAML(value *yaml.Node) error {
	s := strings.Split(value.Value, UnionIdentifier)
	for _, v := range s {
		switch {
		// subject relation
		case strings.Contains(v, RelationIdentifier):
			sr := strings.Split(v, RelationIdentifier)
			r.Definition = append(r.Definition, &SubjectRelation{
				ObjectType: strings.TrimSpace(sr[0]),
				Relation:   strings.TrimSpace(sr[1]),
			})
		// wildcard relation
		case strings.Contains(v, WildcardIdentifier):
			wc := strings.Split(v, WildcardIdentifier)
			r.Definition = append(r.Definition, &WildcardRelation{
				ObjectType: strings.TrimSpace(wc[0]),
			})
		// direct relation
		default:
			r.Definition = append(r.Definition, &DirectRelation{
				ObjectType: strings.TrimSpace(v),
			})
		}
	}

	return nil
}

func (p *PermissionOperator) UnmarshalYAML(value *yaml.Node) error {
	switch {
	// union (OR)
	case strings.Contains(value.Value, UnionIdentifier):
		s := strings.Split(value.Value, UnionIdentifier)
		union := []string{}
		for _, v := range s {
			union = append(union, strings.TrimSpace(v))
		}
		*p = PermissionOperator{
			Operator: &UnionOperator{
				Union: union,
			},
		}
	// intersection (AND)
	case strings.Contains(value.Value, IntersectionIdentifier):
		s := strings.Split(value.Value, IntersectionIdentifier)
		intersect := []string{}
		for _, v := range s {
			intersect = append(intersect, strings.TrimSpace(v))
		}
		*p = PermissionOperator{
			Operator: &IntersectionOperator{
				Intersection: intersect,
			},
		}
	// arrow
	case strings.Contains(value.Value, ArrowIdentifier):
		s := strings.Split(value.Value, ArrowIdentifier)
		*p = PermissionOperator{
			Operator: &ArrowOperator{
				Relation:   strings.TrimSpace(s[0]),
				Permission: strings.TrimSpace(s[1]),
			},
		}
	// exclusion (NOT)
	case strings.Contains(value.Value, ExclusionIdentifier):
		s := strings.Split(value.Value, ExclusionIdentifier)
		*p = PermissionOperator{
			Operator: &ExclusionOperator{
				Base:     strings.TrimSpace(s[0]),
				Subtract: strings.TrimSpace(s[1]),
			},
		}
	// default union of one
	default:
		*p = PermissionOperator{
			Operator: &UnionOperator{
				Union: []string{strings.TrimSpace(value.Value)},
			},
		}
	}

	return nil
}
