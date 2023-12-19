package cache

import (
	"sync"

	"github.com/aserto-dev/azm/model"
	"github.com/aserto-dev/azm/model/diff"
	dsc "github.com/aserto-dev/go-directory/aserto/directory/common/v3"
)

type Cache struct {
	model *model.Model
	mtx   sync.RWMutex
}

// New, create new model cache instance.
func New(m *model.Model) *Cache {
	return &Cache{
		model: m,
		mtx:   sync.RWMutex{},
	}
}

// UpdateModel, swaps the cache model instance.
func (c *Cache) UpdateModel(m *model.Model) error {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	c.model = m
	return nil
}

// Returns a diff struct resulted between the old and the new model.
func (c *Cache) Diff(other *model.Model) *diff.Diff {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	return c.model.Diff(other)
}

// ObjectExists, checks if given object type name exists in the model cache.
func (c *Cache) ObjectExists(on model.ObjectName) bool {
	c.mtx.RLock()
	defer c.mtx.RUnlock()

	_, ok := c.model.Objects[on]
	return ok
}

// RelationExists, checks if given relation type, for the given object type, exists in the model cache.
func (c *Cache) RelationExists(on model.ObjectName, rn model.RelationName) bool {
	c.mtx.RLock()
	defer c.mtx.RUnlock()

	if obj, ok := c.model.Objects[on]; ok {
		_, ok := obj.Relations[rn]
		return ok
	}
	return false
}

// PermissionExists, checks if given permission, for the given object type, exists in the model cache.
func (c *Cache) PermissionExists(on model.ObjectName, pn model.RelationName) bool {
	c.mtx.RLock()
	defer c.mtx.RUnlock()

	if obj, ok := c.model.Objects[on]; ok {
		_, ok := obj.Permissions[pn]
		return ok
	}
	return false
}

func (c *Cache) Metadata() *model.Metadata {
	c.mtx.RLock()
	defer c.mtx.RUnlock()
	return c.model.Metadata
}

func (c *Cache) ValidateRelation(relation *dsc.Relation) error {
	c.mtx.RLock()
	defer c.mtx.RUnlock()

	return c.model.ValidateRelation(
		model.ObjectName(relation.ObjectType),
		model.ObjectID(relation.ObjectId),
		model.RelationName(relation.Relation),
		model.ObjectName(relation.SubjectType),
		model.ObjectID(relation.SubjectId),
		model.RelationName(relation.SubjectRelation),
	)
}
