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
					Union: []*model.RelationRef{{RelOrPerm: "owner"}},
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
					Union: []*model.RelationRef{{RelOrPerm: "writer"}},
				},
				model.PermissionName("view"): {
					Union: []*model.RelationRef{{RelOrPerm: "reader"}, {RelOrPerm: "writer"}},
				},
				model.PermissionName("read_and_write"): {
					Intersection: []*model.RelationRef{{RelOrPerm: "reader"}, {RelOrPerm: "writer"}},
				},
				model.PermissionName("can_only_read"): {
					Exclusion: &model.ExclusionPermission{
						Include: &model.RelationRef{RelOrPerm: "reader"},
						Exclude: &model.RelationRef{RelOrPerm: "writer"},
					},
				},
				model.PermissionName("read"): {
					Union: []*model.RelationRef{{Base: "parent_folder", RelOrPerm: "read"}},
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

	diffm1m3 := m1.Diff(&m3)
	stretch.Equal(t, len(diffm1m3.Added.Objects), 1)
	stretch.Equal(t, diffm1m3.Added.Objects[0], "new_user")
	stretch.Equal(t, len(diffm1m3.Removed.Objects), 1)
	stretch.Equal(t, diffm1m3.Removed.Objects[0], "user")

	stretch.Equal(t, len(diffm1m3.Added.Relations), 1)
	stretch.Equal(t, diffm1m3.Added.Relations["folder"], []string{"viewer"})
	stretch.Equal(t, len(diffm1m3.Removed.Relations), 1)
	stretch.Equal(t, diffm1m3.Removed.Relations["document"], []string{"parent_folder"})
}

func TestGraph(t *testing.T) {
	m := model.Model{
		Version: 1,
		Objects: map[model.ObjectName]*model.Object{
			model.ObjectName("user"): {
				Relations: map[model.RelationName][]*model.Relation{
					model.RelationName("rel_name"): {
						&model.Relation{Direct: model.ObjectName("ext_obj")},
					},
				},
			},
			model.ObjectName("ext_obj"): {},
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
			},
		},
	}

	docExtObjResults := [][]string{
		{"document", "writer", "user", "rel_name", "ext_obj"},
		{"document", "reader", "user", "rel_name", "ext_obj"},
		{"document", "parent_folder", "folder", "owner", "user", "rel_name", "ext_obj"},
	}

	docUserResults := [][]string{
		{"document", "writer", "user"},
		{"document", "reader", "user"},
		{"document", "parent_folder", "folder", "owner", "user"},
	}

	groupExtObjResults := [][]string{
		{"group", "member", "group", "member", "user", "rel_name", "ext_obj"},
		{"group", "member", "user", "rel_name", "ext_obj"},
	}

	graph := m.GetGraph()

	search := graph.FindPaths("document", "ext_obj")
	stretch.Equal(t, len(search), 3)
	for _, expected := range docExtObjResults {
		stretch.Contains(t, search, expected)
	}

	search = graph.FindPaths("document", "user")
	stretch.Equal(t, len(search), 3)
	for _, expected := range docUserResults {
		stretch.Contains(t, search, expected)
	}

	search = graph.FindPaths("group", "ext_obj")
	stretch.Equal(t, len(search), 2)
	for _, expected := range groupExtObjResults {
		stretch.Contains(t, search, expected)
	}
}
