package model_test

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/aserto-dev/azm/model"
	"github.com/nsf/jsondiff"
	"github.com/stretchr/testify/require"
)

var m1 = model.Model{
	Version: 1,
	Objects: map[model.ObjectName]*model.Object{
		model.ObjectName("user"): {},
		model.ObjectName("group"): {
			Relations: map[model.RelationName][]*model.Relation{
				model.RelationName("member"): {
					&model.Relation{Direct: model.ObjectName("user")},
					&model.Relation{Subject: &model.SubjectRelation{
						Object:   model.ObjectName("group"),
						Relation: model.RelationName("member"),
					}},
				},
			},
		},
		model.ObjectName("folder"): {
			Relations: map[model.RelationName][]*model.Relation{
				model.RelationName("owner"): {
					&model.Relation{Direct: model.ObjectName("user")},
				},
			},
			Permissions: map[model.PermissionName]*model.Permission{
				model.PermissionName("read"): {
					Union: []string{"owner"},
				},
			},
		},
		model.ObjectName("document"): {
			Relations: map[model.RelationName][]*model.Relation{
				model.RelationName("parent_folder"): {
					{Direct: model.ObjectName("folder")},
				},
				model.RelationName("writer"): {
					{Direct: model.ObjectName("user")},
				},
				model.RelationName("reader"): {
					{Direct: model.ObjectName("user")},
					{Wildcard: model.ObjectName("user")},
				},
			},
			Permissions: map[model.PermissionName]*model.Permission{
				model.PermissionName("edit"): {
					Union: []string{"writer"},
				},
				model.PermissionName("view"): {
					Union: []string{"reader", "writer"},
				},
				model.PermissionName("read_and_write"): {
					Intersection: []string{"reader", "writer"},
				},
				model.PermissionName("can_only_read"): {
					Exclusion: &model.ExclusionPermission{
						Base:     "reader",
						Subtract: "writer",
					},
				},
				model.PermissionName("read"): {
					Arrow: &model.ArrowPermission{
						Relation:   "parent_folder",
						Permission: "read",
					},
				},
			},
		},
	},
}

func TestProgrammaticModel(t *testing.T) {
	b1, err := json.Marshal(m1)
	require.NoError(t, err)

	w, err := os.Create("./model_test.json")
	require.NoError(t, err)
	defer w.Close()

	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "  ")
	if err := enc.Encode(m1); err != nil {
		require.NoError(t, err)
	}

	b2, err := os.ReadFile("./model.json")
	require.NoError(t, err)

	m2 := model.Model{}
	if err := json.Unmarshal(b2, &m2); err != nil {
		require.NoError(t, err)
	}

	opts := jsondiff.DefaultJSONOptions()
	if diff, str := jsondiff.Compare(b1, b2, &opts); diff != jsondiff.FullMatch {
		require.Equal(t, jsondiff.FullMatch, diff, "diff: %s", str)
	}
}

func TestModel(t *testing.T) {
	buf, err := os.ReadFile("./model.json")
	require.NoError(t, err)

	m := model.Model{}
	if err := json.Unmarshal(buf, &m); err != nil {
		require.NoError(t, err)
	}

	b1, err := json.Marshal(m)
	require.NoError(t, err)

	var m2 model.Model
	if err := json.Unmarshal(b1, &m2); err != nil {
		require.NoError(t, err)
	}

	b2, err := json.Marshal(m2)
	require.NoError(t, err)

	opts := jsondiff.DefaultConsoleOptions()
	if diff, str := jsondiff.Compare(buf, b1, &opts); diff != jsondiff.FullMatch {
		require.Equal(t, diff, jsondiff.FullMatch, "diff: %s", str)
	}

	if diff, str := jsondiff.Compare(buf, b2, &opts); diff != jsondiff.FullMatch {
		require.Equal(t, diff, jsondiff.FullMatch, "diff: %s", str)
	}

	if diff, str := jsondiff.Compare(b1, b2, &opts); diff != jsondiff.FullMatch {
		require.Equal(t, diff, jsondiff.FullMatch, "diff: %s", str)
	}
}

func TestDiff(t *testing.T) {
	m2 := model.Model{
		Version: 1,
		Objects: nil,
	}

	m3 := model.Model{
		Version: 1,
		Objects: map[model.ObjectName]*model.Object{
			model.ObjectName("new_user"): {},
			model.ObjectName("group"): {
				Relations: map[model.RelationName][]*model.Relation{
					model.RelationName("member"): {
						&model.Relation{Direct: model.ObjectName("new_user")},
						&model.Relation{Subject: &model.SubjectRelation{
							Object:   model.ObjectName("group"),
							Relation: model.RelationName("member"),
						}},
					},
				},
			},
			model.ObjectName("folder"): {
				Relations: map[model.RelationName][]*model.Relation{
					model.RelationName("owner"): {
						&model.Relation{Direct: model.ObjectName("new_user")},
					},
					model.RelationName("viewer"): {
						&model.Relation{Direct: model.ObjectName("new_user")},
					},
				},
				Permissions: map[model.PermissionName]*model.Permission{
					model.PermissionName("read"): {
						Union: []string{"owner"},
					},
				},
			},
			model.ObjectName("document"): {
				Relations: map[model.RelationName][]*model.Relation{
					model.RelationName("writer"): {
						{Direct: model.ObjectName("new_user")},
					},
					model.RelationName("reader"): {
						{Direct: model.ObjectName("new_user")},
						{Wildcard: model.ObjectName("new_user")},
					},
				},
			},
		},
	}

	diffm1m2 := m1.Diff(&m2)
	require.Equal(t, len(diffm1m2.Added.Objects), 0)
	require.Equal(t, len(diffm1m2.Added.Relations), 0)
	require.Equal(t, len(diffm1m2.Removed.Objects), 4)
	require.Equal(t, len(diffm1m2.Removed.Relations), 0)

	diffm1nill := m1.Diff(nil)
	require.Equal(t, len(diffm1nill.Added.Objects), 0)
	require.Equal(t, len(diffm1nill.Removed.Objects), 4)
	require.Equal(t, len(diffm1nill.Added.Relations), 0)
	require.Equal(t, len(diffm1nill.Removed.Relations), 0)

	var m4 *model.Model
	diffnillm1 := m4.Diff(&m1)
	require.Equal(t, len(diffnillm1.Added.Objects), 4)
	require.Equal(t, len(diffnillm1.Removed.Objects), 0)
	require.Equal(t, len(diffnillm1.Added.Relations), 0)
	require.Equal(t, len(diffnillm1.Removed.Relations), 0)

	diffBetweenNills := m4.Diff(nil)
	require.Equal(t, len(diffBetweenNills.Added.Objects), 0)
	require.Equal(t, len(diffBetweenNills.Removed.Objects), 0)
	require.Equal(t, len(diffBetweenNills.Added.Relations), 0)
	require.Equal(t, len(diffBetweenNills.Removed.Relations), 0)

	diffm1m3 := m1.Diff(&m3)
	require.Equal(t, len(diffm1m3.Added.Objects), 1)
	require.Equal(t, diffm1m3.Added.Objects[0], model.ObjectName("new_user"))
	require.Equal(t, len(diffm1m3.Removed.Objects), 1)
	require.Equal(t, diffm1m3.Removed.Objects[0], model.ObjectName("user"))

	require.Equal(t, len(diffm1m3.Added.Relations), 1)
	require.Equal(t, diffm1m3.Added.Relations["folder"], []model.RelationName{"viewer"})
	require.Equal(t, len(diffm1m3.Removed.Relations), 1)
	require.Equal(t, diffm1m3.Removed.Relations["document"], []model.RelationName{"parent_folder"})
}
