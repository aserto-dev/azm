package v3

import (
	"io"

	"github.com/aserto-dev/azm/model"

	"gopkg.in/yaml.v3"
)

func Load(r io.Reader) (*model.Model, error) {
	manifest := Manifest{}
	dec := yaml.NewDecoder(r)
	dec.KnownFields(true)

	if err := dec.Decode(&manifest); err != nil {
		return nil, err
	}

	m := model.Model{
		Version: model.ModelVersion,
		Objects: map[model.ObjectName]*model.Object{},
	}

	for otn, ot := range manifest.ObjectTypes {

		relations := map[model.RelationName][]*model.Relation{}

		for rn, rd := range ot.Relations {

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

		for pn, pd := range ot.Permissions {
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

		m.Objects[model.ObjectName(otn)] = &model.Object{
			Relations:   relations,
			Permissions: permissions,
		}
	}

	return &m, nil
}
