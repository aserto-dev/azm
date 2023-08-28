package model_test

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/aserto-dev/azm/model"
	"github.com/nsf/jsondiff"
	"github.com/stretchr/testify/require"
)

func TestProgrammaticModel(t *testing.T) {
	m1 := model.Model{
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
		require.Equal(t, jsondiff.FullMatch, diff)
		fmt.Fprintf(os.Stdout, "differences %s", str)
	}
}

func TestModel(t *testing.T) {
	buf, err := os.ReadFile("./model.json")
	require.NoError(t, err)

	m1 := model.Model{}
	if err := json.Unmarshal(buf, &m1); err != nil {
		require.NoError(t, err)
	}

	b1, err := json.Marshal(m1)
	require.NoError(t, err)

	var m2 model.Model
	if err := json.Unmarshal(b1, &m2); err != nil {
		require.NoError(t, err)
	}

	b2, err := json.Marshal(m2)
	require.NoError(t, err)

	opts := jsondiff.DefaultConsoleOptions()
	if diff, _ := jsondiff.Compare(buf, b1, &opts); diff != jsondiff.FullMatch {
		require.Equal(t, diff, jsondiff.FullMatch)
	}

	if diff, _ := jsondiff.Compare(buf, b2, &opts); diff != jsondiff.FullMatch {
		require.Equal(t, diff, jsondiff.FullMatch)
	}

	if diff, _ := jsondiff.Compare(b1, b2, &opts); diff != jsondiff.FullMatch {
		require.Equal(t, diff, jsondiff.FullMatch)
	}
}
