package v2

import (
	"io"

	"github.com/aserto-dev/azm"
	"github.com/aserto-dev/azm/model"
	"github.com/aserto-dev/azm/types"

	"gopkg.in/yaml.v3"
)

func Load(r io.Reader) (*model.Model, error) {
	manifest := Manifest{}
	dec := yaml.NewDecoder(r)
	dec.KnownFields(true)

	m := model.Model{
		Version: model.ModelVersion,
		Objects: map[types.ObjectName]*types.Object{},
	}

	if err := dec.Decode(&manifest); err != nil {
		if err == io.EOF {
			return &m, nil
		}
		return nil, err
	}

	for objName, obj := range manifest {
		on := types.ObjectName(objName)

		// create object type if not exists
		if _, ok := m.Objects[on]; !ok {
			m.Objects[on] = &types.Object{
				Relations:   map[types.RelationName]*types.Relation{},
				Permissions: map[types.RelationName]*types.Permission{},
			}
		}

		// get object type instance.
		o := m.Objects[on]

		// create all relation instances
		for relName := range obj {
			if _, ok := o.Relations[types.RelationName(relName)]; !ok {
				o.Relations[types.RelationName(relName)] = &types.Relation{Union: []*types.RelationRef{}}
			}
		}

		for relName, rel := range obj {
			// create a subject relation for each union-ed relation, using the same object type.
			for _, v := range rel.Union {
				rs, ok := o.Relations[types.RelationName(v)]
				if !ok {
					return nil, azm.ErrRelationNotFound.Msg(v)
				}

				rs.Union = append(rs.Union, &types.RelationRef{
					Object:   on,
					Relation: types.RelationName(relName),
				})

				o.Relations[types.RelationName(v)] = rs
			}

			for _, v := range rel.Perms {

				norm, _ := types.NormalizeIdentifier(v)
				pn := types.RelationName(norm)

				// if permission does not exist, create permission definition.
				if pd, ok := o.Permissions[pn]; !ok {
					p := &types.Permission{
						Union: []*types.PermissionTerm{{RelOrPerm: types.RelationName(relName)}},
					}
					o.Permissions[pn] = p
				} else {
					pd.Union = append(pd.Union, &types.PermissionTerm{RelOrPerm: types.RelationName(relName)})
					o.Permissions[pn] = pd
				}
			}
		}
	}

	return &m, nil
}
