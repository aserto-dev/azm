package model_test

import (
	"errors"
	"testing"

	"github.com/aserto-dev/azm/model"
	"github.com/aserto-dev/azm/stats"
	"github.com/aserto-dev/azm/types"
	"github.com/aserto-dev/go-directory/pkg/derr"
	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

var ErrBoom = errors.New("Boom")

func TestValidateDiffNoDeletion(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockInstances := model.NewMockInstances(ctrl)

	dif := model.Diff{Removed: model.Changes{}, Added: model.Changes{}}
	err := dif.Validate(mockInstances)

	require.NoError(t, err)
}

func TestValidateDiffWithObjectTypeDeletion(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockInstances := model.NewMockInstances(ctrl)
	objType := types.ObjectName("user")

	dif := model.Diff{Removed: model.Changes{Objects: []types.ObjectName{objType}}, Added: model.Changes{}}

	mockInstances.EXPECT().GetStats().Return(&stats.Stats{}, nil)
	err := dif.Validate(mockInstances)

	require.NoError(t, err)
}

func TestValidateDiffWith2ObjectTypeDeletion(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockInstances := model.NewMockInstances(ctrl)
	objTypes := []types.ObjectName{"user", "member"}

	dif := model.Diff{Removed: model.Changes{Objects: objTypes}, Added: model.Changes{}}

	mockInstances.EXPECT().GetStats().Return(&stats.Stats{ObjectTypes: stats.ObjectTypes{"user": {ObjCount: 1}}}, nil)
	err := dif.Validate(mockInstances)

	require.Error(t, err)
	require.Contains(t, err.Error(), derr.ErrObjectTypeInUse.Message)
}

func TestValidateDiffWithRelationTypeDeletion(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockInstances := model.NewMockInstances(ctrl)
	objTypes := []types.ObjectName{"user", "member"}
	relationTypes := map[types.ObjectName]map[types.RelationName][]string{"folder": {"parent_folder": []string{"document"}}}

	dif := model.Diff{Removed: model.Changes{Objects: objTypes, Relations: relationTypes}, Added: model.Changes{}}

	mockInstances.EXPECT().GetStats().Return(&stats.Stats{ObjectTypes: stats.ObjectTypes{"user": {ObjCount: 0}, "folder": {Count: 1, Relations: stats.Relations{"parent_folder": {Count: 1, SubjectTypes: stats.SubjectTypes{"document": {Count: 1}}}}}}}, nil)
	err := dif.Validate(mockInstances)

	require.Error(t, err)
	require.Contains(t, err.Error(), derr.ErrRelationTypeInUse.Message)
}

func TestValidateDiffWithObjectInstances(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockInstances := model.NewMockInstances(ctrl)
	objTypes := []types.ObjectName{"user", "member"}
	relationTypes := map[types.ObjectName]map[types.RelationName][]string{"folder": {"parent_folder": []string{"document"}}}

	dif := model.Diff{Removed: model.Changes{Objects: objTypes, Relations: relationTypes}, Added: model.Changes{}}

	mockInstances.EXPECT().GetStats().Return(nil, ErrBoom)
	err := dif.Validate(mockInstances)

	require.Error(t, err)
	require.Contains(t, err.Error(), ErrBoom.Error())
}
