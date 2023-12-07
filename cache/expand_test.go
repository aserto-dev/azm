package cache_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/aserto-dev/azm/cache"
	"github.com/aserto-dev/azm/model"
	v2 "github.com/aserto-dev/azm/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// load model cache from serialized model file.
func loadModelCache(t *testing.T, path string) *cache.Cache {
	r, err := os.Open(path)
	require.NoError(t, err)
	defer r.Close()

	var mc model.Model
	dec := json.NewDecoder(r)
	if err := dec.Decode(&mc); err != nil {
		require.NoError(t, err)
	}

	return cache.New(&mc)
}

// helper to regenerate the serialized cache from a manifest.
func loadFromManifest(t *testing.T, path string) *cache.Cache { // nolint:unused
	r, err := os.Open(path)
	require.NoError(t, err)
	defer r.Close()

	m, err := v2.Load(r)
	require.NoError(t, err)

	cachefile := strings.TrimSuffix(path, filepath.Ext(path)) + ".json"
	w, err := os.Create(cachefile)
	require.NoError(t, err)
	defer w.Close()

	require.NoError(t, m.Write(w))

	return cache.New(m)
}

func TestExpandRelation(t *testing.T) {
	mc := loadModelCache(t, "./expand_test.json")

	// tenant:directory-reader does not exist, results should be an empty array.
	relations := mc.ExpandRelation("tenant", "directory-reader")
	assert.Len(t, relations, 0)

	// system:directory-store-writer is not union-ed with any other relations,
	// results should be an single element array, of the requested relation.
	relations = mc.ExpandRelation("system", "directory-store-writer")
	assert.Len(t, relations, 1)
	assert.Contains(t, relations, model.RelationName("directory-store-writer"))

	// tenant:directory-client-reader is union-ed by directory-client-writer.
	relations = mc.ExpandRelation("tenant", "directory-client-reader")
	assert.Len(t, relations, 2)
	assert.Contains(t, relations, model.RelationName("directory-client-reader"))
	assert.Contains(t, relations, model.RelationName("directory-client-writer"))

	// tenant:viewer is union-ed by owner, admin, member.
	relations = mc.ExpandRelation("tenant", "viewer")
	assert.Len(t, relations, 4)
	assert.Contains(t, relations, model.RelationName("viewer"))
	assert.Contains(t, relations, model.RelationName("owner"))
	assert.Contains(t, relations, model.RelationName("admin"))
	assert.Contains(t, relations, model.RelationName("member"))

	// tenant:member is union-ed by owner, admin.
	relations = mc.ExpandRelation("tenant", "member")
	assert.Len(t, relations, 3)
	assert.Contains(t, relations, model.RelationName("owner"))
	assert.Contains(t, relations, model.RelationName("admin"))
	assert.Contains(t, relations, model.RelationName("member"))

	// tenant:admin is union-ed by owner.
	relations = mc.ExpandRelation("tenant", "admin")
	assert.Len(t, relations, 2)
	assert.Contains(t, relations, model.RelationName("owner"))
	assert.Contains(t, relations, model.RelationName("admin"))

	// tenant:owner is not union-ed by any other relation
	relations = mc.ExpandRelation("tenant", "owner")
	assert.Len(t, relations, 1)
	assert.Contains(t, relations, model.RelationName("owner"))

	// tenant:none-relation none-relation is a none existing relation
	relations = mc.ExpandRelation("tenant", "non-relation")
	assert.Len(t, relations, 0)

	// none-tenant:none-relation, none-tenant is a none existing tenant
	relations = mc.ExpandRelation("none-tenant", "non-relation")
	assert.Len(t, relations, 0)
}

func TestExpandPermission(t *testing.T) {
	mc := loadModelCache(t, "./expand_test.json")

	relations := mc.ExpandPermission("tenant", "none-permission")
	assert.Len(t, relations, 0)

	relations = mc.ExpandPermission("none-tenant", "none-permission")
	assert.Len(t, relations, 0)

	relations = mc.ExpandPermission("tenant", "aserto.directory.writer.v2.Writer.SetObject")
	assert.Len(t, relations, 3)
	assert.Contains(t, relations, model.RelationName("owner"))
	assert.Contains(t, relations, model.RelationName("admin"))
	assert.Contains(t, relations, model.RelationName("directory-client-writer"))

	relations = mc.ExpandPermission("tenant", "aserto.directory.reader.v2.Reader.GetObject")
	assert.Len(t, relations, 6)
	assert.Contains(t, relations, model.RelationName("owner"))
	assert.Contains(t, relations, model.RelationName("admin"))
	assert.Contains(t, relations, model.RelationName("member"))
	assert.Contains(t, relations, model.RelationName("viewer"))
	assert.Contains(t, relations, model.RelationName("directory-client-reader"))
	assert.Contains(t, relations, model.RelationName("directory-client-writer"))

	relations = mc.ExpandPermission("tenant", "aserto.tenant.onboarding.v1.Onboarding.ClaimTenant")
	assert.Len(t, relations, 1)
	assert.Contains(t, relations, model.RelationName("owner"))
}
