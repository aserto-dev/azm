package model_test

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/aserto-dev/azm/model"
	"github.com/nsf/jsondiff"
	stretch "github.com/stretchr/testify/require"
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
	stretch.NoError(t, err)

	w, err := os.Create("./model_test.json")
	stretch.NoError(t, err)
	defer w.Close()

	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "  ")
	if err := enc.Encode(m1); err != nil {
		stretch.NoError(t, err)
	}

	b2, err := os.ReadFile("./model.json")
	stretch.NoError(t, err)

	m2 := model.Model{}
	if err := json.Unmarshal(b2, &m2); err != nil {
		stretch.NoError(t, err)
	}

	opts := jsondiff.DefaultJSONOptions()
	if diff, str := jsondiff.Compare(b1, b2, &opts); diff != jsondiff.FullMatch {
		stretch.Equal(t, jsondiff.FullMatch, diff, "diff: %s", str)
	}
}

func TestModel(t *testing.T) {
	buf, err := os.ReadFile("./model.json")
	stretch.NoError(t, err)

	m := model.Model{}
	if err := json.Unmarshal(buf, &m); err != nil {
		stretch.NoError(t, err)
	}

	b1, err := json.Marshal(m)
	stretch.NoError(t, err)

	var m2 model.Model
	if err := json.Unmarshal(b1, &m2); err != nil {
		stretch.NoError(t, err)
	}

	b2, err := json.Marshal(m2)
	stretch.NoError(t, err)

	opts := jsondiff.DefaultConsoleOptions()
	if diff, str := jsondiff.Compare(buf, b1, &opts); diff != jsondiff.FullMatch {
		stretch.Equal(t, diff, jsondiff.FullMatch, "diff: %s", str)
	}

	if diff, str := jsondiff.Compare(buf, b2, &opts); diff != jsondiff.FullMatch {
		stretch.Equal(t, diff, jsondiff.FullMatch, "diff: %s", str)
	}

	if diff, str := jsondiff.Compare(b1, b2, &opts); diff != jsondiff.FullMatch {
		stretch.Equal(t, diff, jsondiff.FullMatch, "diff: %s", str)
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

	diffM1WithM2 := m1.Diff(&m2)
	stretch.Equal(t, len(diffM1WithM2.Added.Objects), 0)
	stretch.Equal(t, len(diffM1WithM2.Added.Relations), 0)
	stretch.Equal(t, len(diffM1WithM2.Removed.Objects), 4)
	stretch.Equal(t, len(diffM1WithM2.Removed.Relations), 0)

	diffM1WithNill := m1.Diff(nil)
	stretch.Equal(t, len(diffM1WithNill.Added.Objects), 0)
	stretch.Equal(t, len(diffM1WithNill.Removed.Objects), 4)
	stretch.Equal(t, len(diffM1WithNill.Added.Relations), 0)
	stretch.Equal(t, len(diffM1WithNill.Removed.Relations), 0)

	var m4 *model.Model
	diffNilWithM1 := m4.Diff(&m1)
	stretch.Equal(t, len(diffNilWithM1.Added.Objects), 4)
	stretch.Equal(t, len(diffNilWithM1.Removed.Objects), 0)
	stretch.Equal(t, len(diffNilWithM1.Added.Relations), 0)
	stretch.Equal(t, len(diffNilWithM1.Removed.Relations), 0)

	diffNilWithNil := m4.Diff(nil)
	stretch.Equal(t, len(diffNilWithNil.Added.Objects), 0)
	stretch.Equal(t, len(diffNilWithNil.Removed.Objects), 0)
	stretch.Equal(t, len(diffNilWithNil.Added.Relations), 0)
	stretch.Equal(t, len(diffNilWithNil.Removed.Relations), 0)

	diffM1WithM3 := m1.Diff(&m3)
	stretch.Equal(t, len(diffM1WithM3.Added.Objects), 1)
	stretch.Equal(t, diffM1WithM3.Added.Objects[0], model.ObjectName("new_user"))
	stretch.Equal(t, len(diffM1WithM3.Removed.Objects), 1)
	stretch.Equal(t, diffM1WithM3.Removed.Objects[0], model.ObjectName("user"))

	stretch.Equal(t, len(diffM1WithM3.Added.Relations), 1)
	stretch.Equal(t, diffM1WithM3.Added.Relations["folder"], []model.RelationName{"viewer"})
	stretch.Equal(t, len(diffM1WithM3.Removed.Relations), 1)
	stretch.Equal(t, diffM1WithM3.Removed.Relations["document"], []model.RelationName{"parent_folder"})
}

func TestGraph(t *testing.T) {
	graph := m1.GetGraph()

	traversal := graph.TraverseGraph("document")
	require.Equal(t, len(traversal), 5)
	traversal = graph.TraverseGraph("group")
	require.Equal(t, len(traversal), 3)
	require.Equal(t, traversal, []string{"group", "member", "user"})
}
