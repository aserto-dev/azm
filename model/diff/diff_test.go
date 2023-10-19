package diff_test

import (
	"context"
	"testing"

	"github.com/aserto-dev/azm/model/diff"
	"github.com/aserto-dev/go-directory/pkg/derr"
	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestValidateDiffNoDeletion(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockDirectoryValidator := diff.NewMockDirectoryValidator(ctrl)

	dif := diff.Diff{Removed: diff.Changes{}, Added: diff.Changes{}}
	err := dif.Validate(context.Background(), mockDirectoryValidator)

	require.NoError(t, err)
}

func TestValidateDiffWithObjectTypeDeletion(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockDirectoryValidator := diff.NewMockDirectoryValidator(ctrl)
	objType := "user"
	bCtx := context.Background()

	dif := diff.Diff{Removed: diff.Changes{Objects: []string{objType}}, Added: diff.Changes{}}

	mockDirectoryValidator.EXPECT().HasObjectInstances(bCtx, objType).Return(false, nil)
	err := dif.Validate(bCtx, mockDirectoryValidator)

	require.NoError(t, err)
}

func TestValidateDiffWith2ObjectTypeDeletion(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockDirectoryValidator := diff.NewMockDirectoryValidator(ctrl)
	objTypes := []string{"user", "member"}
	bCtx := context.Background()

	dif := diff.Diff{Removed: diff.Changes{Objects: objTypes}, Added: diff.Changes{}}

	mockDirectoryValidator.EXPECT().HasObjectInstances(bCtx, objTypes[0]).Return(false, nil)
	mockDirectoryValidator.EXPECT().HasObjectInstances(bCtx, objTypes[1]).Return(true, nil)
	err := dif.Validate(bCtx, mockDirectoryValidator)

	require.Error(t, err)
	require.Contains(t, err.Error(), derr.ErrObjectTypeInUse.Message)
}

func TestValidateDiffWithRelationTypeDeletion(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockDirectoryValidator := diff.NewMockDirectoryValidator(ctrl)
	objTypes := []string{"user", "member"}
	relationTypes := map[string][]string{"folder": {"parent_folder"}}
	bCtx := context.Background()

	dif := diff.Diff{Removed: diff.Changes{Objects: objTypes, Relations: relationTypes}, Added: diff.Changes{}}

	mockDirectoryValidator.EXPECT().HasObjectInstances(bCtx, objTypes[0]).Return(false, nil)
	mockDirectoryValidator.EXPECT().HasObjectInstances(bCtx, objTypes[1]).Return(false, nil)
	mockDirectoryValidator.EXPECT().HasRelationInstances(bCtx, "folder", relationTypes["folder"][0]).Return(true, nil)
	err := dif.Validate(bCtx, mockDirectoryValidator)

	require.Error(t, err)
	require.Contains(t, err.Error(), derr.ErrRelationTypeInUse.Message)
}
