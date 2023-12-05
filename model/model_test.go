package model_test

import (
	"bytes"
	"encoding/json"
	"os"
	"testing"

	"github.com/aserto-dev/azm/model"
	v3 "github.com/aserto-dev/azm/v3"
	"github.com/aserto-dev/go-directory/pkg/derr"
	"github.com/hashicorp/go-multierror"
	"github.com/nsf/jsondiff"
	stretch "github.com/stretchr/testify/require"
)

var m1 = model.Model{
	Version: 2,
	Objects: map[model.ObjectName]*model.Object{
		model.ObjectName("user"): {},
		model.ObjectName("group"): {
			Relations: map[model.RelationName]*model.Relation{
				model.RelationName("member"): {
					Union: []*model.RelationTerm{
						{Direct: &model.RelationRef{Object: model.ObjectName("user")}},
						{Subject: &model.SubjectRelation{
							RelationRef: &model.RelationRef{
								Object:   model.ObjectName("group"),
								Relation: model.RelationName("member"),
							},
						}},
					},
					SubjectTypes: []model.ObjectName{"user"},
				},
			},
		},

		model.ObjectName("folder"): {
			Relations: map[model.RelationName]*model.Relation{
				model.RelationName("owner"): {
					Union: []*model.RelationTerm{
						{Direct: &model.RelationRef{Object: model.ObjectName("user")}},
					},
					SubjectTypes: []model.ObjectName{"user"},
				},
			},
			Permissions: map[model.PermissionName]*model.Permission{
				model.PermissionName("read"): {
					Union: []*model.PermissionRef{{RelOrPerm: "owner"}},
				},
			},
		},
		model.ObjectName("document"): {
			Relations: map[model.RelationName]*model.Relation{
				model.RelationName("parent_folder"): {
					Union:        []*model.RelationTerm{{Direct: &model.RelationRef{Object: model.ObjectName("folder")}}},
					SubjectTypes: []model.ObjectName{"folder"},
				},
				model.RelationName("writer"): {
					Union:        []*model.RelationTerm{{Direct: &model.RelationRef{Object: model.ObjectName("user")}}},
					SubjectTypes: []model.ObjectName{"user"},
				},
				model.RelationName("reader"): {
					Union: []*model.RelationTerm{
						{Direct: &model.RelationRef{Object: model.ObjectName("user")}},
						{Wildcard: &model.RelationRef{Object: model.ObjectName("user"), Relation: "*"}},
					},
					SubjectTypes: []model.ObjectName{"user"},
				},
			},
			Permissions: map[model.PermissionName]*model.Permission{
				model.PermissionName("edit"): {
					Union: []*model.PermissionRef{{RelOrPerm: "writer"}},
				},
				model.PermissionName("view"): {
					Union: []*model.PermissionRef{{RelOrPerm: "reader"}, {RelOrPerm: "writer"}},
				},
				model.PermissionName("read_and_write"): {
					Intersection: []*model.PermissionRef{{RelOrPerm: "reader"}, {RelOrPerm: "writer"}},
				},
				model.PermissionName("can_only_read"): {
					Exclusion: &model.ExclusionPermission{
						Include: &model.PermissionRef{RelOrPerm: "reader"},
						Exclude: &model.PermissionRef{RelOrPerm: "writer"},
					},
				},
				model.PermissionName("read"): {
					Union: []*model.PermissionRef{{Base: "parent_folder", RelOrPerm: "read"}},
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

	b2, err := os.ReadFile("./testdata/model.json")
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
	buf, err := os.ReadFile("./testdata/model.json")
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
		Version: 2,
		Objects: nil,
	}

	m3 := model.Model{
		Version: 2,
		Objects: map[model.ObjectName]*model.Object{
			model.ObjectName("new_user"): {},
			model.ObjectName("group"): {
				Relations: map[model.RelationName]*model.Relation{
					model.RelationName("member"): {
						Union: []*model.RelationTerm{
							{Direct: &model.RelationRef{Object: model.ObjectName("new_user")}},
							{Subject: &model.SubjectRelation{
								RelationRef: &model.RelationRef{
									Object:   model.ObjectName("group"),
									Relation: model.RelationName("member"),
								},
								SubjectTypes: []model.ObjectName{},
							}},
						},
					},
				},
			},
			model.ObjectName("folder"): {
				Relations: map[model.RelationName]*model.Relation{
					model.RelationName("owner"): {
						Union: []*model.RelationTerm{{Direct: &model.RelationRef{Object: model.ObjectName("new_user")}}},
					},
					model.RelationName("viewer"): {
						Union: []*model.RelationTerm{{Direct: &model.RelationRef{Object: model.ObjectName("new_user")}}},
					},
				},
				Permissions: map[model.PermissionName]*model.Permission{
					model.PermissionName("read"): {
						Union: []*model.PermissionRef{{RelOrPerm: "owner"}},
					},
				},
			},
			model.ObjectName("document"): {
				Relations: map[model.RelationName]*model.Relation{
					model.RelationName("writer"): {
						Union: []*model.RelationTerm{{Direct: &model.RelationRef{Object: model.ObjectName("new_user")}}},
					},
					model.RelationName("reader"): {
						Union: []*model.RelationTerm{
							{Direct: &model.RelationRef{Object: model.ObjectName("new_user")}},
							{Wildcard: &model.RelationRef{Object: model.ObjectName("new_user"), Relation: "*"}},
						},
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
		Version: 2,
		Objects: map[model.ObjectName]*model.Object{
			model.ObjectName("user"): {
				Relations: map[model.RelationName]*model.Relation{
					model.RelationName("rel_name"): {
						Union: []*model.RelationTerm{{Direct: &model.RelationRef{Object: model.ObjectName("ext_obj")}}},
					},
				},
			},
			model.ObjectName("ext_obj"): {},
			model.ObjectName("group"): {
				Relations: map[model.RelationName]*model.Relation{
					model.RelationName("member"): {
						Union: []*model.RelationTerm{
							{Direct: &model.RelationRef{Object: model.ObjectName("user")}},
							{Subject: &model.SubjectRelation{
								RelationRef: &model.RelationRef{
									Object:   model.ObjectName("group"),
									Relation: model.RelationName("member"),
								},
								SubjectTypes: []model.ObjectName{},
							}},
						},
					},
				},
			},
			model.ObjectName("folder"): {
				Relations: map[model.RelationName]*model.Relation{
					model.RelationName("owner"): {
						Union: []*model.RelationTerm{{Direct: &model.RelationRef{Object: model.ObjectName("user")}}},
					},
				},
			},
			model.ObjectName("document"): {
				Relations: map[model.RelationName]*model.Relation{
					model.RelationName("parent_folder"): {
						Union: []*model.RelationTerm{{Direct: &model.RelationRef{Object: model.ObjectName("folder")}}},
					},
					model.RelationName("writer"): {
						Union: []*model.RelationTerm{{Direct: &model.RelationRef{Object: model.ObjectName("user")}}},
					},
					model.RelationName("reader"): {
						Union: []*model.RelationTerm{
							{Direct: &model.RelationRef{Object: model.ObjectName("user")}},
							{Wildcard: &model.RelationRef{Object: model.ObjectName("user"), Relation: "*"}},
						},
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

type expectedErrors int

func TestValidation(t *testing.T) {
	tests := []struct {
		name           string
		manifest       string
		expectedErrors expectedErrors
		validate       func(error, *stretch.Assertions)
	}{
		{
			"valid manifest",
			"./testdata/valid.yaml",
			expectedErrors(0),
			func(_ error, _ *stretch.Assertions) {},
		},
		{
			"relation/permission collision",
			"./testdata/rel_perm_collision.yaml",
			expectedErrors(2),
			func(err error, assert *stretch.Assertions) {
				assert.ErrorContains(err, "permission name 'file:writer' conflicts with 'file:writer' relation")
				assert.ErrorContains(err, "relation 'file:bad' has no definition")
			},
		},
		{
			"relation to undefined targets",
			"./testdata/undefined_targets.yaml",
			expectedErrors(6),
			func(err error, assert *stretch.Assertions) {
				// relations
				assert.ErrorContains(err, "relation 'file:owner' references undefined object type 'person'")
				assert.ErrorContains(err, "relation 'file:reader' references undefined object type 'team'")
				assert.ErrorContains(err, "relation 'file:reader' references undefined object type 'project'")
				assert.ErrorContains(err, "relation 'file:writer' references undefined object type 'team'")
				assert.ErrorContains(err, "relation 'file:admin' references undefined relation type 'group#admin'")
				assert.ErrorContains(err, "permission name 'file:reader' conflicts with 'file:reader' relation")
			},
		},
		{
			"cyclic relation definitions",
			"./testdata/invalid_cycles.yaml",
			expectedErrors(3),
			func(err error, assert *stretch.Assertions) {
				assert.ErrorContains(err, "relation 'team:member' is circular and does not resolve to any object types")
				assert.ErrorContains(err, "relation 'team:owner' is circular and does not resolve to any object types")
				assert.ErrorContains(err, "relation 'project:owner' is circular and does not resolve to any object types")
			},
		},
		{
			"permissions with invalid targets",
			"./testdata/invalid_perms.yaml",
			expectedErrors(2),
			func(err error, assert *stretch.Assertions) {
				assert.ErrorContains(err, "permission 'file:write' references 'owner->write', which can resolve to undefined relation or permission 'user:write'")
				assert.ErrorContains(err, "permission 'file:write' references undefined relation or permission 'file:editor'")
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			assert := stretch.New(tt)
			m, err := loadManifest(test.manifest)

			// Log the model for debugging purposes.
			var b bytes.Buffer
			enc := json.NewEncoder(&b)
			enc.SetIndent("", "  ")
			assert.NoError(enc.Encode(m))
			tt.Logf("model: %s", b.String())

			if test.expectedErrors == 0 {
				return
			}

			// verify that we got a load error.
			assert.Error(err)
			// verify that the error is of type ErrInvalidArgument
			aerr := derr.ErrInvalidArgument
			assert.ErrorAs(err, &aerr)
			assert.Equal("E20015", aerr.Code)

			merr := aerr.Unwrap().(*multierror.Error)
			assert.NotNil(merr)
			test.validate(merr, assert)

			// verify that we got the expected number of errors.
			assert.Len(merr.Errors, int(test.expectedErrors))
		})
	}
}

func loadManifest(path string) (*model.Model, error) {
	r, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer r.Close()

	return v3.Load(r)
}
