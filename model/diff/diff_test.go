package diff_test

import (
	"errors"
	"testing"

	"github.com/aserto-dev/azm/model/diff"
	"github.com/aserto-dev/go-directory/pkg/derr"
	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

var ErrBoom = errors.New("Boom")

func TestValidateDiffNoDeletion(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockInstances := diff.NewMockInstances(ctrl)

	dif := diff.Diff{Removed: diff.Changes{}, Added: diff.Changes{}}
	err := dif.Validate(mockInstances)

	require.NoError(t, err)
}

func TestValidateDiffWithObjectTypeDeletion(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockInstances := diff.NewMockInstances(ctrl)
	objType := "user"

	dif := diff.Diff{Removed: diff.Changes{Objects: []string{objType}}, Added: diff.Changes{}}

	mockInstances.EXPECT().ObjectTypes().Return([]string{}, nil)
	mockInstances.EXPECT().RelationTypes().Return([]*diff.RelationKind{}, nil)
	err := dif.Validate(mockInstances)

	require.NoError(t, err)
}

func TestValidateDiffWith2ObjectTypeDeletion(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockInstances := diff.NewMockInstances(ctrl)
	objTypes := []string{"user", "member"}

	dif := diff.Diff{Removed: diff.Changes{Objects: objTypes}, Added: diff.Changes{}}

	mockInstances.EXPECT().ObjectTypes().Return([]string{"user"}, nil)
	mockInstances.EXPECT().RelationTypes().Return([]*diff.RelationKind{}, nil)
	err := dif.Validate(mockInstances)

	require.Error(t, err)
	require.Contains(t, err.Error(), derr.ErrObjectTypeInUse.Message)
}

func TestValidateDiffWithRelationTypeDeletion(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockInstances := diff.NewMockInstances(ctrl)
	objTypes := []string{"user", "member"}
	relationTypes := map[string][]string{"folder": {"parent_folder"}}

	dif := diff.Diff{Removed: diff.Changes{Objects: objTypes, Relations: relationTypes}, Added: diff.Changes{}}

	mockInstances.EXPECT().ObjectTypes().Return([]string{}, nil)
	mockInstances.EXPECT().RelationTypes().Return([]*diff.RelationKind{{Object: "folder", Relation: "parent_folder"}}, nil)
	err := dif.Validate(mockInstances)

	require.Error(t, err)
	require.Contains(t, err.Error(), derr.ErrRelationTypeInUse.Message)
}

func TestValidateDiffWithObjectInstances(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockInstances := diff.NewMockInstances(ctrl)
	objTypes := []string{"user", "member"}
	relationTypes := map[string][]string{"folder": {"parent_folder"}}

	dif := diff.Diff{Removed: diff.Changes{Objects: objTypes, Relations: relationTypes}, Added: diff.Changes{}}

	mockInstances.EXPECT().ObjectTypes().Return([]string{}, ErrBoom)
	err := dif.Validate(mockInstances)

	require.Error(t, err)
	require.Contains(t, err.Error(), ErrBoom.Error())
}
