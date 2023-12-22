package v3

import (
	"io"

	"github.com/aserto-dev/azm/model"
	"github.com/aserto-dev/azm/parser"
	"github.com/aserto-dev/azm/types"
	"github.com/samber/lo"

	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
)

type Object = types.Object
type ObjectName = types.ObjectName
type RelationRef = types.RelationRef
type Relation = types.Relation
type Permission = types.Permission

func Load(r io.Reader) (*model.Model, error) {
	m := model.Model{
		Version: model.ModelVersion,
		Objects: map[ObjectName]*Object{},
	}

	dec := yaml.NewDecoder(r)
	dec.KnownFields(true)

	manifest := Manifest{}
	if err := dec.Decode(&manifest); err != nil {
		if err == io.EOF {
			return &m, nil
		}
		return nil, err
	}

	for on, o := range manifest.ObjectTypes {
		log.Debug().Str("object", string(on)).Msg("loading object")

		relationTerms := lo.MapEntries(o.Relations, func(rn RelationName, rd string) (types.RelationName, []*RelationRef) {
			log.Debug().Str("object", string(on)).Str("relation", string(rn)).Msg("loading relation")

			return types.RelationName(rn), parser.ParseRelation(rd)
		})

		relations := lo.MapEntries(relationTerms, func(rn types.RelationName, rts []*RelationRef) (types.RelationName, *Relation) {
			return rn, &Relation{Union: rts}
		})

		permissions := lo.MapEntries(o.Permissions, func(pn PermissionName, pd string) (types.RelationName, *Permission) {
			log.Debug().Str("object", string(on)).Str("permission", string(pn)).Msg("loading permission")

			return types.RelationName(pn), parser.ParsePermission(pd)
		})

		m.Objects[ObjectName(on)] = &Object{
			Relations:   relations,
			Permissions: permissions,
		}
	}

	return &m, m.Validate()
}
