package v2

import (
	"io"
	"strings"

	"github.com/aserto-dev/azm"
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

	for objName, obj := range manifest {
		otn := model.ObjectName(objName)

		// create object type if not exists
		if _, ok := m.Objects[otn]; !ok {
			m.Objects[otn] = &model.Object{
				Relations:   map[model.RelationName][]*model.Relation{},
				Permissions: map[model.PermissionName]*model.Permission{},
			}
		}

		// get object type instance.
		ot := m.Objects[otn]

		// create all relation instances
		for relName := range obj {
			if _, ok := ot.Relations[model.RelationName(relName)]; !ok {
				ot.Relations[model.RelationName(relName)] = []*model.Relation{}
			}
		}

		for relName, rel := range obj {
			// create a subject relation for each union-ed relation, using the same object type.
			for _, v := range rel.Union {
				rs, ok := ot.Relations[model.RelationName(v)]
				if !ok {
					return nil, azm.ErrRelationNotFound.Msg(v)
				}

				rs = append(rs, &model.Relation{Subject: &model.SubjectRelation{
					Object:   otn,
					Relation: model.RelationName(relName),
				}})

				ot.Relations[model.RelationName(v)] = rs
			}

			for _, v := range rel.Perms {

				pn := model.PermissionName(NormalizePermission(v))

				// if permission does not exist, create permission definition.
				if pd, ok := ot.Permissions[pn]; !ok {
					p := &model.Permission{
						Union: []string{relName},
					}
					ot.Permissions[pn] = p
				} else {
					pd.Union = append(pd.Union, relName)
					ot.Permissions[pn] = pd
				}
			}
		}
	}

	return &m, nil
}

func NormalizePermission(p string) string {
	// normalize permission name
	// NOTE: does not enforce max length to 64 characters.
	return strings.ToLower(strings.ReplaceAll(p, ".", "_"))
}
