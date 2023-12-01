package v3

import (
	"io"

	"github.com/aserto-dev/azm/model"
	"github.com/aserto-dev/azm/parser"
	"github.com/samber/lo"

	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
)

func Load(r io.Reader) (*model.Model, error) {
	m := model.Model{
		Version: model.ModelVersion,
		Objects: map[model.ObjectName]*model.Object{},
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

		relations := lo.MapEntries(o.Relations, func(rn RelationName, rd string) (model.RelationName, []*model.Relation) {
			log.Debug().Str("object", string(on)).Str("relation", string(rn)).Msg("loading relation")

			return model.RelationName(rn), parser.ParseRelation(rd)
		})

		permissions := lo.MapEntries(o.Permissions, func(pn PermissionName, pd string) (model.PermissionName, *model.Permission) {
			log.Debug().Str("object", string(on)).Str("permission", string(pn)).Msg("loading permission")

			return model.PermissionName(pn), parser.ParsePermission(pd)
		})

		m.Objects[model.ObjectName(on)] = &model.Object{
			Relations:   relations,
			Permissions: permissions,
		}
	}

	return &m, m.Validate()
}
