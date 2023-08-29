package cache

import (
	"sort"

	"github.com/aserto-dev/azm/model"
	v2 "github.com/aserto-dev/azm/v2"
	dsc2 "github.com/aserto-dev/go-directory/aserto/directory/common/v2"
	"github.com/aserto-dev/go-directory/pkg/derr"
	"github.com/samber/lo"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// GetObjectType, v2 backwards-compatibility accessor function, returns v2 ObjectType by name.
func (c *Cache) GetObjectType(objectType string) (*dsc2.ObjectType, error) {
	c.mtx.RLock()
	defer c.mtx.RUnlock()

	if _, ok := c.model.Objects[model.ObjectName(objectType)]; ok {
		return &dsc2.ObjectType{
			Name:        objectType,
			DisplayName: title(objectType),
			IsSubject:   false,
			Ordinal:     0,
			Status:      0,
			Schema:      &structpb.Struct{},
			CreatedAt:   timestamppb.Now(),
			UpdatedAt:   timestamppb.Now(),
			Hash:        "",
		}, nil
	}

	return &dsc2.ObjectType{}, derr.ErrObjectTypeNotFound.Msg(objectType)
}

// GetObjectTypes, v2 backwards-compatibility accessor function, returns list of v2.ObjectType instances.
func (c *Cache) GetObjectTypes() ([]*dsc2.ObjectType, error) {
	c.mtx.RLock()
	defer c.mtx.RUnlock()

	results := []*dsc2.ObjectType{}

	for objectType := range c.model.Objects {
		results = append(results, &dsc2.ObjectType{
			Name:        string(objectType),
			DisplayName: title(string(objectType)),
			IsSubject:   false,
			Ordinal:     0,
			Status:      0,
			Schema:      &structpb.Struct{},
			CreatedAt:   timestamppb.Now(),
			UpdatedAt:   timestamppb.Now(),
			Hash:        "",
		})
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Name < results[j].Name
	})

	return results, nil
}

// GetRelationType, v2 backwards-compatibility accessor function, returns v2 RelationType by object type and relation name.
func (c *Cache) GetRelationType(objectType, relation string) (*dsc2.RelationType, error) {
	c.mtx.RLock()
	defer c.mtx.RUnlock()

	otn := model.ObjectName(objectType)
	rtn := model.RelationName(relation)

	ot, ok := c.model.Objects[otn]
	if !ok {
		return &dsc2.RelationType{}, derr.ErrObjectTypeNotFound.Msg(objectType)
	}

	if _, ok := ot.Relations[rtn]; !ok {
		return &dsc2.RelationType{}, derr.ErrRelationNotFound.Msg(objectType + ":" + relation)
	}

	return &dsc2.RelationType{
		ObjectType:  objectType,
		Name:        relation,
		DisplayName: objectType + ":" + relation,
		Ordinal:     0,
		Status:      0,
		Unions:      c.getRelationUnions(ot, otn, rtn),
		Permissions: c.getRelationPermissions(ot, rtn),
		CreatedAt:   timestamppb.Now(),
		UpdatedAt:   timestamppb.Now(),
		Hash:        "",
	}, nil
}

func (*Cache) getRelationPermissions(ot *model.Object, rtn model.RelationName) []string {
	permissions := []string{}
	for pn, p := range ot.Permissions {
		if lo.Contains(p.Union, string(rtn)) {
			permissions = append(permissions, string(pn))
		}
	}
	return permissions
}

func (*Cache) getRelationUnions(ot *model.Object, otn model.ObjectName, rtn model.RelationName) []string {
	unions := []string{}
	for rn, rs := range ot.Relations {
		for _, r := range rs {
			if r.Subject != nil && r.Subject.Object == otn && r.Subject.Relation == rtn {
				unions = append(unions, string(rn))
			}
		}
	}
	return unions
}

// GetRelationTypes, v2 backwards-compatibility accessor function, returns list of v2 RelationType instances, optionally filtered by by object type.
func (c *Cache) GetRelationTypes(objectType string) ([]*dsc2.RelationType, error) {
	c.mtx.RLock()
	defer c.mtx.RUnlock()

	results := []*dsc2.RelationType{}

	objectTypes := c.model.Objects
	if objectType != "" {
		if ot, ok := c.model.Objects[model.ObjectName(objectType)]; !ok {
			return results, derr.ErrObjectTypeNotFound.Msg(objectType)
		} else {
			objectTypes = map[model.ObjectName]*model.Object{
				model.ObjectName(objectType): ot,
			}
		}
	}

	for otn, ot := range objectTypes {
		for rn := range ot.Relations {

			results = append(results, &dsc2.RelationType{
				ObjectType:  string(otn),
				Name:        string(rn),
				DisplayName: string(otn) + ":" + string(rn),
				Ordinal:     0,
				Status:      0,
				Unions:      c.getRelationUnions(ot, otn, rn),
				Permissions: c.getRelationPermissions(ot, rn),
				CreatedAt:   timestamppb.Now(),
				UpdatedAt:   timestamppb.Now(),
				Hash:        "",
			})
		}
	}

	sort.Slice(results, func(i, j int) bool {
		iv, jv := results[i], results[j]
		switch {
		case iv.ObjectType != jv.ObjectType:
			return iv.ObjectType < jv.ObjectType
		default:
			return iv.Name < jv.Name
		}
	})

	return results, nil
}

// GetPermission, v2 backwards-compatibility accessor function, returns v2 Permission by permission name.
func (c *Cache) GetPermission(permission string) (*dsc2.Permission, error) {
	c.mtx.RLock()
	defer c.mtx.RUnlock()

	pn := model.PermissionName(v2.NormalizePermission(permission))

	for _, o := range c.model.Objects {
		if _, ok := o.Permissions[pn]; ok {
			return &dsc2.Permission{
				Name:        string(pn),
				DisplayName: string(pn),
				CreatedAt:   timestamppb.Now(),
				UpdatedAt:   timestamppb.Now(),
				Hash:        "",
			}, nil
		}
	}

	return &dsc2.Permission{}, derr.ErrPermissionNotFound.Msg(permission)
}

// GetPermissions, v2 backwards-compatibility accessor function, returns list of v2 Permission instances.
func (c *Cache) GetPermissions() ([]*dsc2.Permission, error) {
	c.mtx.RLock()
	defer c.mtx.RUnlock()

	aggregator := []*dsc2.Permission{}

	ts := timestamppb.Now()
	for _, o := range c.model.Objects {
		for pn := range o.Permissions {
			aggregator = append(aggregator, &dsc2.Permission{
				Name:        string(pn),
				DisplayName: string(pn),
				CreatedAt:   ts,
				UpdatedAt:   ts,
				Hash:        "",
			})
		}
	}

	results := lo.UniqBy(aggregator, func(i *dsc2.Permission) string {
		return i.Name
	})

	sort.Slice(results, func(i, j int) bool {
		return results[i].Name < results[j].Name
	})

	return results, nil
}

func title(s string) string {
	return cases.Title(language.AmericanEnglish, cases.NoLower).String(s)
}
