package v3

import (
	"io"

	"github.com/aserto-dev/azm/model"
	"github.com/aserto-dev/azm/parser"

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

		relations := make(map[model.RelationName][]*model.Relation, len(o.Relations))

		for rn, rd := range o.Relations {
			log.Debug().Str("object", string(on)).Str("relation", string(rn)).Msg("loading relation")

			relations[model.RelationName(rn)] = parser.ParseRelation(rd)
		}

		permissions := make(map[model.PermissionName]*model.Permission, len(o.Permissions))

		for pn, pd := range o.Permissions {
			log.Debug().Str("object", string(on)).Str("permission", string(pn)).Msg("loading permission")

			permissions[model.PermissionName(pn)] = parser.ParsePermission(pd)
		}

		m.Objects[model.ObjectName(on)] = &model.Object{
			Relations:   relations,
			Permissions: permissions,
		}
	}

	return &m, nil
}
