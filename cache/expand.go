package cache

import (
	"github.com/aserto-dev/azm/model"
	v2 "github.com/aserto-dev/azm/v2"
	"github.com/samber/lo"
)

// ExpandRelation, returns list of relations which are a union of the given relation.
// For example, when a writer relation inherits reader, the expansion of a reader = reader + writer.
func (c *Cache) ExpandRelation(otn model.ObjectName, rtn model.RelationName) []model.RelationName {
	c.mtx.RLock()
	defer c.mtx.RUnlock()

	results := []model.RelationName{}

	// starting object type and relation must exist in order to be expanded.
	if ot, ok := c.model.Objects[otn]; !ok {
		return results
	} else if _, ok := ot.Relations[rtn]; !ok {
		return results
	}

	// get relation set for given object:relation.
	rs := c.model.Objects[otn].Relations[rtn]

	// include given permission in result set
	results = append(results, rtn)

	// iterate through relation set, determine if it "unions" with the given relation.
	for _, r := range rs {
		if r.Subject != nil && r.Subject.Object == otn {
			results = append(results, r.Subject.Relation)
		}
	}

	return lo.Uniq(results)
}

// ExpandPermission returns list of relations which cover the given permission for the given object type.
func (c *Cache) ExpandPermission(on model.ObjectName, pn model.PermissionName) []model.RelationName {
	c.mtx.RLock()
	defer c.mtx.RUnlock()

	pn = model.PermissionName(v2.NormalizePermission(string(pn)))

	results := []model.RelationName{}

	// starting object type and permission must exist in order to be expanded.
	if ot, ok := c.model.Objects[on]; !ok {
		return results
	} else if _, ok := ot.Permissions[pn]; !ok {
		return results
	}

	p := c.model.Objects[on].Permissions[pn]

	results = append(results, expandUnion(p.Union)...)

	for _, rn := range results {
		results = append(results, c.ExpandRelation(on, rn)...)
	}

	return lo.Uniq(results)
}

// convert union []string to []model.RelationName.
func expandUnion(u []string) []model.RelationName {
	result := []model.RelationName{}
	for _, v := range u {
		result = append(result, model.RelationName(v))
	}
	return result
}
