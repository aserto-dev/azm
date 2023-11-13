package v3

import (
	"io"

	"github.com/aserto-dev/azm/model"

	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
)

func Load(r io.Reader) (*model.Model, error) {
	manifest := Manifest{}
	dec := yaml.NewDecoder(r)
	dec.KnownFields(true)

	m := model.Model{
		Version: model.ModelVersion,
		Objects: map[model.ObjectName]*model.Object{},
	}

	if err := dec.Decode(&manifest); err != nil {
		if err == io.EOF {
			return &m, nil
		}
		return nil, err
	}

	for on, o := range manifest.ObjectTypes {
		log.Debug().Str("object", string(on)).Msg("loading object")

		relations := map[model.RelationName][]*model.Relation{}

		if o.Relations == nil {
			o.Relations = map[RelationName]RelationDefinition{}
		}

		for rn, rd := range o.Relations {
			log.Debug().Str("object", string(on)).Str("relation", string(rn)).Msg("loading relation")

			for _, v := range rd.Definition {
				if _, ok := relations[model.RelationName(rn)]; !ok {
					relations[model.RelationName(rn)] = []*model.Relation{}
				}

				rs := relations[model.RelationName(rn)]
				r := &model.Relation{}

				switch x := v.(type) {
				case *DirectRelation:
					r.Direct = model.ObjectName(x.ObjectType)

				case *SubjectRelation:
					r.Subject = &model.SubjectRelation{
						Object:   model.ObjectName(x.ObjectType),
						Relation: model.RelationName(x.Relation),
					}

				case *WildcardRelation:
					r.Wildcard = model.ObjectName(x.ObjectType)
				}

				rs = append(rs, r)
				relations[model.RelationName(rn)] = rs
			}
		}

		permissions := map[model.PermissionName]*model.Permission{}

		if o.Permissions == nil {
			o.Permissions = map[PermissionName]PermissionOperator{}
		}

		for pn, pd := range o.Permissions {
			log.Debug().Str("object", string(on)).Str("permission", string(pn)).Msg("loading permission")

			if _, ok := permissions[model.PermissionName(pn)]; !ok {
				permissions[model.PermissionName(pn)] = &model.Permission{}
			}

			p := permissions[model.PermissionName(pn)]

			switch x := pd.Operator.(type) {
			case *UnionOperator:
				p.Union = x.Union

			case *IntersectionOperator:
				p.Intersection = x.Intersection

			case *ExclusionOperator:
				p.Exclusion = &model.ExclusionPermission{
					Base:     x.Base,
					Subtract: x.Subtract,
				}

			case *ArrowOperator:
				p.Arrow = &model.ArrowPermission{
					Relation:   x.Relation,
					Permission: x.Permission,
				}
			}

			permissions[model.PermissionName(pn)] = p
		}

		m.Objects[model.ObjectName(on)] = &model.Object{
			Relations:   relations,
			Permissions: permissions,
		}
	}

	return &m, nil
}
