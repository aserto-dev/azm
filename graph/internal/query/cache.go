package query

import dsc "github.com/aserto-dev/go-directory/aserto/directory/common/v3"

type Relation struct {
	Load
	Path
}

func (r *Relation) Identifier(out *dsc.RelationIdentifier) {
	out.ObjectType = r.OT.String()
	out.ObjectId = r.OID.String()
	out.Relation = r.RT.String()
	out.SubjectType = r.ST.String()
	out.SubjectId = r.SID.String()
	out.SubjectRelation = r.SRT.String()
}

type Cache struct {
	sets  map[Relation]ObjSet
	calls map[Relation]ObjSet
}

func (c *Cache) LookupSet(rel *Relation) (ObjSet, bool) {
	if c.sets == nil {
		return nil, false
	}
	set, ok := c.sets[*rel]
	return set, ok
}

func (c *Cache) StoreSet(rel *Relation, set ObjSet) {
	if c.sets == nil {
		c.sets = make(map[Relation]ObjSet)
	}
	c.sets[*rel] = set
}

func (c *Cache) LookupCall(rel *Relation) (ObjSet, bool) {
	if c.calls == nil {
		return nil, false
	}
	set, ok := c.calls[*rel]
	return set, ok
}

func (c *Cache) StoreCall(rel *Relation, set ObjSet) {
	if c.calls == nil {
		c.calls = make(map[Relation]ObjSet)
	}
	c.calls[*rel] = set
}
