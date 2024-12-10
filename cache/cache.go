package cache

import (
	"github.com/aserto-dev/azm/model"
	"github.com/aserto-dev/azm/model/diff"
	stts "github.com/aserto-dev/azm/stats"
	dsc "github.com/aserto-dev/go-directory/aserto/directory/common/v3"
	"github.com/aserto-dev/go-directory/pkg/derr"
	"github.com/samber/lo"
)

type (
	ObjectName   = model.ObjectName
	RelationName = model.RelationName
)

type Cache struct {
	model *model.Model
	// mtx      sync.RWMutex
	// relsPool *mempool.RelationsPool
}

// New, create new model cache instance.
func New(m *model.Model) *Cache {
	return &Cache{
		model: m,
		// mtx:      sync.RWMutex{},
		// relsPool: mempool.NewRelationsPool(),
	}
}

// UpdateModel, swaps the cache model instance.
func (c *Cache) UpdateModel(m *model.Model) error {
	// c.mtx.Lock()
	// defer c.mtx.Unlock()
	c.model = m
	return nil
}

func (c *Cache) CanUpdate(other *model.Model, stats *stts.Stats) error {
	// c.mtx.RLock()
	// defer c.mtx.RUnlock()
	return diff.CanUpdateModel(c.model, other, stats)
}

// ObjectExists, checks if given object type name exists in the model cache.
func (c *Cache) ObjectExists(on ObjectName) bool {
	// c.mtx.RLock()
	// defer c.mtx.RUnlock()

	_, ok := c.model.Objects[on]
	return ok
}

// RelationExists, checks if given relation type, for the given object type, exists in the model cache.
func (c *Cache) RelationExists(on ObjectName, rn RelationName) bool {
	// c.mtx.RLock()
	// defer c.mtx.RUnlock()

	if obj, ok := c.model.Objects[on]; ok {
		_, ok := obj.Relations[rn]
		return ok
	}
	return false
}

// PermissionExists, checks if given permission, for the given object type, exists in the model cache.
func (c *Cache) PermissionExists(on ObjectName, pn RelationName) bool {
	// c.mtx.RLock()
	// defer c.mtx.RUnlock()

	if obj, ok := c.model.Objects[on]; ok {
		_, ok := obj.Permissions[pn]
		return ok
	}
	return false
}

func (c *Cache) Metadata() *model.Metadata {
	// c.mtx.RLock()
	// defer c.mtx.RUnlock()
	return c.model.Metadata
}

func (c *Cache) ValidateRelation(relation *dsc.RelationIdentifier) error {
	// c.mtx.RLock()
	// defer c.mtx.RUnlock()

	return c.model.ValidateRelation(
		ObjectName(relation.ObjectType),
		model.ObjectID(relation.ObjectId),
		RelationName(relation.Relation),
		ObjectName(relation.SubjectType),
		model.ObjectID(relation.SubjectId),
		RelationName(relation.SubjectRelation),
	)
}

// AssignableRelations returns the set of relations that can occur between a given object type
// and a subject type, optionally with a subject relation.
//
// If more than one subject relation is provided, AssignableRelations returns relations that match any
// of the given relations. For example, if the manifest has:
//
// types:
//
//	tenant:
//	  relations:
//	    admin: group#member
//	    guest: group#guest
//
// Then AssignableRelations("tenant", "group", "member", "guest") returns ["admin", "guest"].
func (c *Cache) AssignableRelations(on, sn ObjectName, sr ...RelationName) ([]RelationName, error) {
	if !c.ObjectExists(on) {
		return nil, derr.ErrObjectNotFound.Msg(on.String())
	}
	if !c.ObjectExists(sn) {
		return nil, derr.ErrObjectNotFound.Msg(sn.String())
	}
	for _, srel := range sr {
		if !c.RelationExists(sn, srel) {
			return nil, derr.ErrRelationNotFound.Msgf("%s#%s", sn, sr)
		}
	}

	matches := lo.PickBy(c.model.Objects[on].Relations, func(rn RelationName, r *model.Relation) bool {
		for _, ref := range r.Union {
			if ref.Object != sn {
				// type mismatch
				continue
			}

			switch {
			case ref.IsDirect(), ref.IsWildcard():
				if len(sr) == 0 {
					return true
				}
			case ref.IsSubject() && lo.Contains(sr, ref.Relation):
				return true
			}
		}
		return false
	})

	return lo.Keys(matches), nil
}
