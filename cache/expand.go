package cache

import (
	"github.com/aserto-dev/azm/types"
	"github.com/samber/lo"
)

// ExpandRelation, returns list of relations which are a union of the given relation.
// For example, when a writer relation inherits reader, the expansion of a reader = reader + writer.
func (c *Cache) ExpandRelation(on types.ObjectName, rn types.RelationName) []types.RelationName {
	c.mtx.RLock()
	defer c.mtx.RUnlock()

	results := []types.RelationName{}

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
		case rt.IsSubject() && rt.Object == on:
			results = append(results, rt.Relation)
		case rt.IsDirect():
			results = append(results, c.ExpandRelation(on, types.RelationName(rt.Object))...)
		}
	}

	return lo.Uniq(results)
}

// ExpandPermission returns list of relations which cover the given permission for the given object type.
func (c *Cache) ExpandPermission(on types.ObjectName, pn types.RelationName) []types.RelationName {
	c.mtx.RLock()
	defer c.mtx.RUnlock()

	norm, _ := types.NormalizeIdentifier(string(pn))
	pn = types.RelationName(norm)

	results := []types.RelationName{}

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
func (c *Cache) expandUnion(o *types.Object, u ...*types.PermissionTerm) []types.RelationName {
	result := []types.RelationName{}
	for _, ref := range u {
		if ref.IsArrow() {
			continue
		}

		result = append(result, ref.RelOrPerm)
		exp := lo.FilterMap(o.Relations[ref.RelOrPerm].Union, func(r *types.RelationRef, _ int) (*types.PermissionTerm, bool) {
			if !r.IsDirect() {
				return &types.PermissionTerm{}, false
			}
			_, ok := o.Relations[types.RelationName(r.Object)]
			return &types.PermissionTerm{RelOrPerm: types.RelationName(r.Object)}, ok

		})
		result = append(result, c.expandUnion(o, exp...)...)
	}
	return result
}
