package v2

import (
	"io"

	"github.com/aserto-dev/azm"
	"github.com/aserto-dev/azm/model"

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

	for objName, obj := range manifest {
		on := model.ObjectName(objName)

		// create object type if not exists
		if _, ok := m.Objects[on]; !ok {
			m.Objects[on] = &model.Object{
				Relations:   map[model.RelationName][]*model.Relation{},
				Permissions: map[model.PermissionName]*model.Permission{},
			}
		}

		// get object type instance.
		o := m.Objects[on]

		// create all relation instances
		for relName := range obj {
			if _, ok := o.Relations[model.RelationName(relName)]; !ok {
				o.Relations[model.RelationName(relName)] = []*model.Relation{}
			}
		}

		for relName, rel := range obj {
			// create a subject relation for each union-ed relation, using the same object type.
			for _, v := range rel.Union {
				rs, ok := o.Relations[model.RelationName(v)]
				if !ok {
					return nil, azm.ErrRelationNotFound.Msg(v)
				}

				rs = append(rs, &model.Relation{Subject: &model.SubjectRelation{
					Object:   on,
					Relation: model.RelationName(relName),
				}})

				o.Relations[model.RelationName(v)] = rs
			}

			for _, v := range rel.Perms {

				norm, _ := model.NormalizeIdentifier(v)
				pn := model.PermissionName(norm)

				// if permission does not exist, create permission definition.
				if pd, ok := o.Permissions[pn]; !ok {
					p := &model.Permission{
						Union: []string{relName},
					}
					o.Permissions[pn] = p
				} else {
					pd.Union = append(pd.Union, relName)
					o.Permissions[pn] = pd
				}
			}
		}
	}

	return &m, nil
}
