package cache_test

import (
	"testing"

	"github.com/aserto-dev/azm/paging"
	dsc2 "github.com/aserto-dev/go-directory/aserto/directory/common/v2"
	"github.com/aserto-dev/go-directory/pkg/derr"
	asserts "github.com/stretchr/testify/assert"
)

func TestGetObjectTypeV2(t *testing.T) {
	assert := asserts.New(t)
	mc := loadModelCache(t, "./metadata_test.json")

	objectType, err := mc.GetObjectType("")
	assert.Error(err, derr.ErrObjectTypeNotFound)
	assert.Equal(&dsc2.ObjectType{}, objectType)

	objectType, err = mc.GetObjectType("foobar")
	assert.Error(err, derr.ErrObjectTypeNotFound)
	assert.Equal(&dsc2.ObjectType{}, objectType)

	objectType, err = mc.GetObjectType("tenant")
	assert.NoError(err)
	assert.Equal(objectType.Name, "tenant")
}

func TestGetObjectTypesV2(t *testing.T) {
	assert := asserts.New(t)
	mc := loadModelCache(t, "./metadata_test.json")

	objectTypes, err := mc.GetObjectTypes()
	assert.NoError(err)
	assert.Len(objectTypes, 11)

	testPagination(t, objectTypes)
}

func TestGetRelationTypeV2(t *testing.T) {
	assert := asserts.New(t)
	mc := loadModelCache(t, "./metadata_test.json")

	relationType, err := mc.GetRelationType("", "")
	assert.Error(err, derr.ErrRelationTypeNotFound)
	assert.Equal(&dsc2.RelationType{}, relationType)

	relationType, err = mc.GetRelationType("tenant", "")
	assert.Error(err, derr.ErrObjectTypeNotFound)
	assert.Equal(&dsc2.RelationType{}, relationType)

	relationType, err = mc.GetRelationType("tenant", "admin")
	assert.NoError(err)
	assert.Equal(relationType.ObjectType, "tenant")
	assert.Equal(relationType.Name, "admin")
	assert.Len(relationType.Unions, 2)
	assert.Contains(relationType.Unions, "member")
	assert.Contains(relationType.Unions, "viewer")
	assert.Len(relationType.Permissions, 171)
}

func TestGetRelationTypesV2(t *testing.T) {
	assert := asserts.New(t)
	mc := loadModelCache(t, "./metadata_test.json")

	relationTypes, err := mc.GetRelationTypes("foobar")
	assert.Error(err, derr.ErrObjectTypeNotFound)
	assert.Len(relationTypes, 0)
	testPagination(t, relationTypes)

	relationTypes, err = mc.GetRelationTypes("")
	assert.NoError(err)
	assert.Len(relationTypes, 23)
	testPagination(t, relationTypes)

	relationTypes, err = mc.GetRelationTypes("tenant")
	assert.NoError(err)
	assert.Len(relationTypes, 10)
	testPagination(t, relationTypes)
}

func TestGetPermissionV2(t *testing.T) {
	assert := asserts.New(t)
	mc := loadModelCache(t, "./metadata_test.json")

	permission, err := mc.GetPermission("")
	assert.Error(err, derr.ErrPermissionNotFound)
	assert.Equal(&dsc2.Permission{}, permission)

	permission, err = mc.GetPermission("foo.bar")
	assert.Error(err, derr.ErrPermissionNotFound)
	assert.Equal(&dsc2.Permission{}, permission)

	permission, err = mc.GetPermission("aserto.discovery.policy.v2.Discovery.OPAInstanceDiscovery")
	assert.NoError(err)
	assert.Equal(permission.Name, "aserto_discovery_policy_v2_discovery_opainstancediscovery")
}

func TestGetPermissionsV2(t *testing.T) {
	assert := asserts.New(t)
	mc := loadModelCache(t, "./metadata_test.json")

	permissions, err := mc.GetPermissions()
	assert.NoError(err)
	assert.Len(permissions, 233)
	testPagination(t, permissions)
}

type Named interface {
	GetName() string
}

type PagingSlice[T any] interface {
	~[]T
	Paginate(page *dsc2.PaginationRequest) (*paging.Result[T], error)
}

func testPagination[T Named, S PagingSlice[T]](t *testing.T, slice S) {
	assert := asserts.New(t)
	for pageSize := 1; pageSize <= len(slice); pageSize++ {
		page := &dsc2.PaginationRequest{Size: int32(pageSize)}

		current := 0
		for {
			result, err := slice.Paginate(page)
			assert.NoError(err)

			for i := 0; i < len(result.Items); i++ {
				assert.Equal(result.Items[i].GetName(), slice[current].GetName())
				current++
			}

			if result.PageInfo.NextToken != "" {
				page.Token = result.PageInfo.NextToken
			} else {
				break
			}
		}
	}
}
