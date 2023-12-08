package cache

import (
	"sort"

	"github.com/aserto-dev/azm/model"
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
func (c *Cache) GetObjectTypes() (ObjectTypeSlice, error) {
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

	on := model.ObjectName(objectType)
	rn := model.RelationName(relation)

	o, ok := c.model.Objects[on]
	if !ok {
		return &dsc2.RelationType{}, derr.ErrObjectTypeNotFound.Msg(objectType)
	}

	if _, ok := o.Relations[rn]; !ok {
		return &dsc2.RelationType{}, derr.ErrRelationNotFound.Msg(objectType + ":" + relation)
	}

	return &dsc2.RelationType{
		ObjectType:  objectType,
		Name:        relation,
		DisplayName: objectType + ":" + relation,
		Ordinal:     0,
		Status:      0,
		Unions:      c.getRelationUnions(o, on, rn),
		Permissions: c.getRelationPermissions(o, rn),
		CreatedAt:   timestamppb.Now(),
		UpdatedAt:   timestamppb.Now(),
		Hash:        "",
	}, nil
}

func (*Cache) getRelationPermissions(o *model.Object, rn model.RelationName) []string {
	permissions := []string{}
	for pn, p := range o.Permissions {
		union := lo.Map(p.Union, func(r *model.PermissionRef, _ int) string {
			if r.Base != "" {
				panic("arrow permissions not supported yet")
			}
			return r.RelOrPerm
		})
		if lo.Contains(union, string(rn)) {
			permissions = append(permissions, string(pn))
		}
	}
	return permissions
}

func (*Cache) getRelationUnions(o *model.Object, on model.ObjectName, rn model.RelationName) []string {
	unions := []string{}
	for name, r := range o.Relations {
		for _, rt := range r.Union {
			if rt.IsSubject() && rt.Object == on && rt.Relation == rn {
				unions = append(unions, string(name))
			}
		}
	}
	return unions
}

// GetRelationTypes, v2 backwards-compatibility accessor function, returns list of v2 RelationType instances, optionally filtered by by object type.
func (c *Cache) GetRelationTypes(objectType string) (RelationTypeSlice, error) {
	c.mtx.RLock()
	defer c.mtx.RUnlock()

	results := []*dsc2.RelationType{}

	objectTypes := c.model.Objects
	if objectType != "" {
		if o, ok := c.model.Objects[model.ObjectName(objectType)]; !ok {
			return results, derr.ErrObjectTypeNotFound.Msg(objectType)
		} else {
			objectTypes = map[model.ObjectName]*model.Object{
				model.ObjectName(objectType): o,
			}
		}
	}

	for on, o := range objectTypes {
		for rn := range o.Relations {

			results = append(results, &dsc2.RelationType{
				ObjectType:  string(on),
				Name:        string(rn),
				DisplayName: string(on) + ":" + string(rn),
				Ordinal:     0,
				Status:      0,
				Unions:      c.getRelationUnions(o, on, rn),
				Permissions: c.getRelationPermissions(o, rn),
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

	norm, _ := model.NormalizeIdentifier(permission)
	pn := model.PermissionName(norm)

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
func (c *Cache) GetPermissions() (PermissionSlice, error) {
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
