package v2

import (
	"io"
	"log"
	"strings"

	dsc "github.com/aserto-dev/go-directory/aserto/directory/common/v2"
	dsr "github.com/aserto-dev/go-directory/aserto/directory/reader/v2"
	dsw "github.com/aserto-dev/go-directory/aserto/directory/writer/v2"
	"gopkg.in/yaml.v2"

	"github.com/aserto-dev/azm"
)

var (
	_ azm.Loader = &ModelV2{}
	_ azm.Saver  = &ModelV2{}
	_ azm.Reader = &ModelV2{}
	_ azm.Writer = &ModelV2{}
)

const (
	schemaVersion int    = 2
	modelName     string = "m2"
	union         string = "union"
	perms         string = "permissions"
)

type ModelV2 struct{}
type Manifest map[string]ObjectRelation
type ObjectRelation map[string]Relation
type Relation map[string][]string

func Model() *ModelV2 {
	return &ModelV2{}
}

func (*ModelV2) Read(r io.Reader) (*azm.Model, error) {
	manifest := make(Manifest, 0)

	dec := yaml.NewDecoder(r)
	dec.SetStrict(true)
	if err := dec.Decode(&manifest); err != nil {
		return nil, err
	}

	return convert(manifest)
}

func (*ModelV2) Write(w io.Writer, m *azm.Model) error {
	return nil
}

func (*ModelV2) Load(r dsr.ReaderClient) (*azm.Model, error) {
	return nil, nil
}

func (*ModelV2) Save(w dsw.WriterClient, m *azm.Model) error {
	return nil
}

func convert(m Manifest) (*azm.Model, error) {
	var err error

	model := azm.New(modelName, schemaVersion)

	model, err = m.objectTypes(model)
	if err != nil {
		return nil, err
	}

	model, err = m.relations(model)
	if err != nil {
		return nil, err
	}

	model, err = m.permissions(model)
	if err != nil {
		return nil, err
	}

	return model, nil
}

func (m Manifest) objectTypes(model *azm.Model) (*azm.Model, error) {
	for objectType := range m {
		log.Printf("o:%s\n", objectType)

		if _, err := model.SetObjectType(&azm.ObjectType{
			Name: objectType,
		}); err != nil {
			return nil, err
		}
	}

	return model, nil
}

func (m Manifest) relations(model *azm.Model) (*azm.Model, error) {
	for objectType, objectRelation := range m {
		objType, err := model.SetObjectType(&azm.ObjectType{Name: objectType})
		if err != nil {
			return nil, err
		}

		for relationType, v := range objectRelation {
			if len(v[union]) == 0 && len(v[perms]) > 0 {
				continue
			}

			rt := &azm.RelationType{
				Name: relationType,
			}

			if len(v[union]) == 0 {
				rt.Definition = ""
				rt.Operator = azm.None
				rt.Relations = []string{}
			} else {
				rt.Definition = strings.Join(v[union], " | ")
				rt.Operator = azm.Union
				rt.Relations = v[union]
			}

			log.Printf("r:%s#%s u:[%s]\n", objectType, relationType, strings.Join(v[union], ","))
			if _, err := objType.SetRelationType(rt); err != nil {
				return nil, err
			}
		}
		if _, err := model.SetObjectType(&azm.ObjectType{Name: objectType}); err != nil {
			return nil, err
		}

	}

	return model, nil
}

func (m Manifest) permissions(model *azm.Model) (*azm.Model, error) {
	for objectType, objectRelation := range m {
		objType, err := model.SetObjectType(&azm.ObjectType{Name: objectType})
		if err != nil {
			return nil, err
		}

		for relationType, v := range objectRelation {
			if len(v[perms]) == 0 {
				continue
			}

			pt := &azm.PermissionType{
				Name:        relationType,
				Definition:  strings.Join(v[perms], " | "),
				Operator:    azm.Union,
				Permissions: v[perms],
			}

			log.Printf("r:%s#%s p:[%s]\n", objectType, relationType, strings.Join(v[perms], ","))
			if _, err := objType.SetPermissionType(pt); err != nil {
				return nil, err
			}
		}

		if _, err := model.SetObjectType(&azm.ObjectType{Name: objectType}); err != nil {
			return nil, err
		}
	}

	return model, nil
}

func (*ModelV2) Update(objectTypes []*dsc.ObjectType, relationTypes []*dsc.RelationType) (*azm.Model, error) {
	manifest := make(Manifest, 0)

	for _, objType := range objectTypes {
		objRel := make(ObjectRelation, 0)

		for _, relType := range relationTypes {
			if relType.ObjectType != objType.Name {
				continue
			}

			relation := make(Relation, 0)
			relation[union] = relType.Unions
			relation[perms] = relType.Permissions
			objRel[relType.Name] = relation
		}

		manifest[objType.Name] = objRel
	}

	return convert(manifest)
}
