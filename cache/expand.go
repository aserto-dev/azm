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
	rs := c.model.Objects[on].Relations[rn]

	// include given permission in result set
	results = append(results, rn)

	// iterate through relation set, determine if it "unions" with the given relation.
	for _, r := range rs {
		switch {
		case r.Subject != nil && r.Subject.Object == on:
			results = append(results, r.Subject.Relation)
		case r.Direct != "":
			results = append(results, c.ExpandRelation(on, model.RelationName(r.Direct))...)
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
func (c *Cache) expandUnion(o *model.Object, u ...string) []model.RelationName {
	result := []model.RelationName{}
	for _, v := range u {
		rn := model.RelationName(v)
		result = append(result, rn)

		exp := lo.FilterMap(o.Relations[rn], func(r *model.Relation, _ int) (string, bool) {
			if r.Direct == "" {
				return "", false
			}
			_, ok := o.Relations[model.RelationName(r.Direct)]
			return string(r.Direct), ok

		})
		result = append(result, c.expandUnion(o, exp...)...)
	}
	return result
}
