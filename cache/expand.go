package cache

import (
	"github.com/aserto-dev/azm/model"
	"github.com/samber/lo"
)

// ExpandRelation, returns list of relations which are a union of the given relation.
// For example, when a writer relation inherits reader, the expansion of a reader = reader + writer.
func (c *Cache) ExpandRelation(on model.ObjectName, rn model.RelationName) []model.RelationName {
	c.mtx.RLock()
	defer c.mtx.RUnlock()

	results := []model.RelationName{}

	// starting object type and relation must exist in order to be expanded.
	if o, ok := c.model.Objects[on]; !ok {
		return results
	} else if _, ok := o.Relations[rn]; !ok {
		return results
	}

	// get relation set for given object:relation.
	r := c.model.Objects[on].Relations[rn]

	// include given permission in result set
	results = append(results, rn)

	// iterate through relation set, determine if it "unions" with the given relation.
	for _, rt := range r.Union {
		switch {
		case rt.Subject != nil && rt.Subject.Object == on:
			results = append(results, rt.Subject.Relation)
		case rt.Direct != nil:
			results = append(results, c.ExpandRelation(on, model.RelationName(rt.Direct.Object))...)
		}
	}

	return lo.Uniq(results)
}

// ExpandPermission returns list of relations which cover the given permission for the given object type.
func (c *Cache) ExpandPermission(on model.ObjectName, pn model.PermissionName) []model.RelationName {
	c.mtx.RLock()
	defer c.mtx.RUnlock()

	norm, _ := model.NormalizeIdentifier(string(pn))
	pn = model.PermissionName(norm)

	results := []model.RelationName{}

	// starting object type and permission must exist in order to be expanded.
	o, ok := c.model.Objects[on]
	if !ok {
		return results
	}
	if _, ok := o.Permissions[pn]; !ok {
		return results
	}

	p := c.model.Objects[on].Permissions[pn]

	results = append(results, c.expandUnion(o, p.Union...)...)

	for _, rn := range results {
		results = append(results, c.ExpandRelation(on, rn)...)
	}

	return lo.Uniq(results)
}

// convert union []string to []model.RelationName.
func (c *Cache) expandUnion(o *model.Object, u ...*model.PermissionRef) []model.RelationName {
	result := []model.RelationName{}
	for _, ref := range u {
		if ref.Base != "" {
			panic("expandUnion: arrow permissions not supported yet")
		}
		rn := model.RelationName(ref.RelOrPerm)
		result = append(result, rn)

		exp := lo.FilterMap(o.Relations[rn].Union, func(r *model.RelationTerm, _ int) (*model.PermissionRef, bool) {
			if r.Direct == nil {
				return &model.PermissionRef{}, false
			}
			_, ok := o.Relations[model.RelationName(r.Direct.Object)]
			return &model.PermissionRef{RelOrPerm: string(r.Direct.Object)}, ok

		})
		result = append(result, c.expandUnion(o, exp...)...)
	}
	return result
}
