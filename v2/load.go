package v2

import (
	"io"

	"github.com/aserto-dev/azm"
	"github.com/aserto-dev/azm/model"
	"github.com/aserto-dev/azm/types"

	"gopkg.in/yaml.v3"
)

type ObjectName = types.ObjectName
type Object = types.Object
type RelationName = types.RelationName
type Permission = types.Permission
type RelationRef = types.RelationRef
type PermissionTerm = types.PermissionTerm

func Load(r io.Reader) (*model.Model, error) {
	manifest := Manifest{}
	dec := yaml.NewDecoder(r)
	dec.KnownFields(true)

	m := model.Model{
		Version: model.ModelVersion,
		Objects: map[ObjectName]*Object{},
	}

	if err := dec.Decode(&manifest); err != nil {
		if err == io.EOF {
			return &m, nil
		}
		return nil, err
	}

	for objName, obj := range manifest {
		on := ObjectName(objName)

		// create object type if not exists
		if _, ok := m.Objects[on]; !ok {
			m.Objects[on] = &Object{
				Relations:   map[RelationName]*types.Relation{},
				Permissions: map[RelationName]*Permission{},
			}
		}

		// get object type instance.
		o := m.Objects[on]

		// create all relation instances
		for relName := range obj {
			if _, ok := o.Relations[RelationName(relName)]; !ok {
				o.Relations[RelationName(relName)] = &types.Relation{Union: []*RelationRef{}}
			}
		}

		for relName, rel := range obj {
			// create a subject relation for each union-ed relation, using the same object type.
			for _, v := range rel.Union {
				rs, ok := o.Relations[RelationName(v)]
				if !ok {
					return nil, azm.ErrRelationNotFound.Msg(v)
				}

				rs.Union = append(rs.Union, &RelationRef{
					Object:   on,
					Relation: RelationName(relName),
				})

				o.Relations[RelationName(v)] = rs
			}

			for _, v := range rel.Perms {

				norm, _ := types.NormalizeIdentifier(v)
				pn := RelationName(norm)

				// if permission does not exist, create permission definition.
				if pd, ok := o.Permissions[pn]; !ok {
					p := &Permission{
						Union: []*PermissionTerm{{RelOrPerm: RelationName(relName)}},
					}
					o.Permissions[pn] = p
				} else {
					pd.Union = append(pd.Union, &PermissionTerm{RelOrPerm: RelationName(relName)})
					o.Permissions[pn] = pd
				}
			}
		}
	}

	return &m, nil
}
