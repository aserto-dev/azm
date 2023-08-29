package cache_test

import (
	"testing"

	dsc2 "github.com/aserto-dev/go-directory/aserto/directory/common/v2"
	"github.com/aserto-dev/go-directory/pkg/derr"
	"github.com/stretchr/testify/assert"
)

func TestGetObjectTypeV2(t *testing.T) {
	mc := loadModelCache(t, "./metadata_test.json")

	objectType, err := mc.GetObjectType("")
	assert.Error(t, err, derr.ErrObjectTypeNotFound)
	assert.Equal(t, &dsc2.ObjectType{}, objectType)

	objectType, err = mc.GetObjectType("foobar")
	assert.Error(t, err, derr.ErrObjectTypeNotFound)
	assert.Equal(t, &dsc2.ObjectType{}, objectType)

	objectType, err = mc.GetObjectType("tenant")
	assert.NoError(t, err)
	assert.Equal(t, objectType.Name, "tenant")
}

func TestGetObjectTypesV2(t *testing.T) {
	mc := loadModelCache(t, "./metadata_test.json")

	objectTypes, err := mc.GetObjectTypes()
	assert.NoError(t, err)
	assert.Len(t, objectTypes, 11)
}

func TestGetRelationTypeV2(t *testing.T) {
	mc := loadModelCache(t, "./metadata_test.json")

	relationType, err := mc.GetRelationType("", "")
	assert.Error(t, err, derr.ErrRelationTypeNotFound)
	assert.Equal(t, &dsc2.RelationType{}, relationType)

	relationType, err = mc.GetRelationType("tenant", "")
	assert.Error(t, err, derr.ErrObjectTypeNotFound)
	assert.Equal(t, &dsc2.RelationType{}, relationType)

	relationType, err = mc.GetRelationType("tenant", "admin")
	assert.NoError(t, err)
	assert.Equal(t, relationType.ObjectType, "tenant")
	assert.Equal(t, relationType.Name, "admin")
	assert.Len(t, relationType.Unions, 2)
	assert.Contains(t, relationType.Unions, "member")
	assert.Contains(t, relationType.Unions, "viewer")
	assert.Len(t, relationType.Permissions, 171)
}

func TestGetRelationTypesV2(t *testing.T) {
	mc := loadModelCache(t, "./metadata_test.json")

	relationTypes, err := mc.GetRelationTypes("foobar")
	assert.Error(t, err, derr.ErrObjectTypeNotFound)
	assert.Len(t, relationTypes, 0)

	relationTypes, err = mc.GetRelationTypes("")
	assert.NoError(t, err)
	assert.Len(t, relationTypes, 23)

	relationTypes, err = mc.GetRelationTypes("tenant")
	assert.NoError(t, err)
	assert.Len(t, relationTypes, 10)
}

func TestGetPermissionV2(t *testing.T) {
	mc := loadModelCache(t, "./metadata_test.json")

	permission, err := mc.GetPermission("")
	assert.Error(t, err, derr.ErrPermissionNotFound)
	assert.Equal(t, &dsc2.Permission{}, permission)

	permission, err = mc.GetPermission("foo.bar")
	assert.Error(t, err, derr.ErrPermissionNotFound)
	assert.Equal(t, &dsc2.Permission{}, permission)

	permission, err = mc.GetPermission("aserto.discovery.policy.v2.Discovery.OPAInstanceDiscovery")
	assert.NoError(t, err)
	assert.Equal(t, permission.Name, "aserto_discovery_policy_v2_discovery_opainstancediscovery")
}

func TestGetPermissionsV2(t *testing.T) {
	mc := loadModelCache(t, "./metadata_test.json")

	permissions, err := mc.GetPermissions()
	assert.NoError(t, err)
	assert.Len(t, permissions, 233)
}
