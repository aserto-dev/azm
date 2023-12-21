package model_test

import (
	"bytes"
	"encoding/json"
	"os"
	"testing"

	"github.com/aserto-dev/azm/model"
	"github.com/aserto-dev/azm/types"
	v3 "github.com/aserto-dev/azm/v3"
	"github.com/aserto-dev/go-directory/pkg/derr"
	"github.com/hashicorp/go-multierror"
	"github.com/nsf/jsondiff"
	stretch "github.com/stretchr/testify/require"
)

var m1 = model.Model{
	Version: 2,
	Objects: map[types.ObjectName]*types.Object{
		types.ObjectName("user"): {},
		types.ObjectName("group"): {
			Relations: map[types.RelationName]*types.Relation{
				types.RelationName("member"): {
					Union: []*types.RelationRef{
						{Object: types.ObjectName("user")},
						{Object: types.ObjectName("group"), Relation: types.RelationName("member")},
					},
				},
			},
		},

		types.ObjectName("folder"): {
			Relations: map[types.RelationName]*types.Relation{
				types.RelationName("owner"): {
					Union: []*types.RelationRef{
						{Object: types.ObjectName("user")},
					},
				},
			},
			Permissions: map[types.RelationName]*types.Permission{
				types.RelationName("read"): {
					Union: []*types.PermissionTerm{{RelOrPerm: "owner"}},
				},
			},
		},
		types.ObjectName("document"): {
			Relations: map[types.RelationName]*types.Relation{
				types.RelationName("parent_folder"): {
					Union: []*types.RelationRef{{Object: types.ObjectName("folder")}},
				},
				types.RelationName("writer"): {
					Union: []*types.RelationRef{{Object: types.ObjectName("user")}},
				},
				types.RelationName("reader"): {
					Union: []*types.RelationRef{
						{Object: types.ObjectName("user")},
						{Object: types.ObjectName("user"), Relation: "*"},
					},
				},
			},
			Permissions: map[types.RelationName]*types.Permission{
				types.RelationName("edit"): {
					Union: []*types.PermissionTerm{{RelOrPerm: "writer"}},
				},
				types.RelationName("view"): {
					Union: []*types.PermissionTerm{
						{RelOrPerm: "reader"},
						{RelOrPerm: "writer"},
					},
				},
				types.RelationName("read_and_write"): {
					Intersection: []*types.PermissionTerm{
						{RelOrPerm: "reader"},
						{RelOrPerm: "writer"},
					},
				},
				types.RelationName("can_only_read"): {
					Exclusion: &types.ExclusionPermission{
						Include: &types.PermissionTerm{RelOrPerm: "reader"},
						Exclude: &types.PermissionTerm{RelOrPerm: "writer"},
					},
				},
				types.RelationName("read"): {
					Union: []*types.PermissionTerm{{Base: "parent_folder", RelOrPerm: "read"}},
				},
			},
		},
	},
}

func TestProgrammaticModel(t *testing.T) {
	stretch.NoError(t, m1.Validate())

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
		Objects: map[types.ObjectName]*types.Object{
			types.ObjectName("new_user"): {},
			types.ObjectName("group"): {
				Relations: map[types.RelationName]*types.Relation{
					types.RelationName("member"): {
						Union: []*types.RelationRef{
							{Object: types.ObjectName("new_user")},
						},
					},
				},
			},
			types.ObjectName("folder"): {
				Relations: map[types.RelationName]*types.Relation{
					types.RelationName("owner"): {
						Union: []*types.RelationRef{{Object: types.ObjectName("new_user")}},
					},
					types.RelationName("viewer"): {
						Union: []*types.RelationRef{{Object: types.ObjectName("new_user")}},
					},
				},
				Permissions: map[types.RelationName]*types.Permission{
					types.RelationName("read"): {
						Union: []*types.PermissionTerm{{RelOrPerm: "owner"}},
					},
				},
			},
			types.ObjectName("document"): {
				Relations: map[types.RelationName]*types.Relation{
					types.RelationName("writer"): {
						Union: []*types.RelationRef{{Object: types.ObjectName("new_user")}},
					},
					types.RelationName("reader"): {
						Union: []*types.RelationRef{
							{Object: types.ObjectName("new_user")},
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
	stretch.Equal(t, diffm1m3.Added.Objects[0], types.ObjectName("new_user"))
	stretch.Equal(t, len(diffm1m3.Removed.Objects), 1)
	stretch.Equal(t, diffm1m3.Removed.Objects[0], types.ObjectName("user"))

	stretch.Equal(t, len(diffm1m3.Added.Relations), 3)
	stretch.Equal(t, diffm1m3.Added.Relations["folder"], map[types.RelationName][]string{"owner": {"new_user"}, "viewer": {}})
	stretch.Equal(t, len(diffm1m3.Removed.Relations), 3)
	stretch.Equal(t, diffm1m3.Removed.Relations["document"], map[types.RelationName][]string{"parent_folder": {}, "reader": {"user", "user:*"}, "writer": {"user"}})
}

func TestGraph(t *testing.T) {
	m := model.Model{
		Version: 2,
		Objects: map[types.ObjectName]*types.Object{
			types.ObjectName("user"): {
				Relations: map[types.RelationName]*types.Relation{
					types.RelationName("rel_name"): {
						Union: []*types.RelationRef{{Object: types.ObjectName("ext_obj")}},
					},
				},
			},
			types.ObjectName("ext_obj"): {},
			types.ObjectName("group"): {
				Relations: map[types.RelationName]*types.Relation{
					types.RelationName("member"): {
						Union: []*types.RelationRef{
							{Object: types.ObjectName("user")},
							{Object: types.ObjectName("group"), Relation: types.RelationName("member")},
						},
					},
				},
			},
			types.ObjectName("folder"): {
				Relations: map[types.RelationName]*types.Relation{
					types.RelationName("owner"): {
						Union: []*types.RelationRef{{Object: types.ObjectName("user")}},
					},
				},
			},
			types.ObjectName("document"): {
				Relations: map[types.RelationName]*types.Relation{
					types.RelationName("parent_folder"): {
						Union: []*types.RelationRef{{Object: types.ObjectName("folder")}},
					},
					types.RelationName("writer"): {
						Union: []*types.RelationRef{{Object: types.ObjectName("user")}},
					},
					types.RelationName("reader"): {
						Union: []*types.RelationRef{
							{Object: types.ObjectName("user")},
							{Object: types.ObjectName("user"), Relation: "*"},
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

func TestValidation(t *testing.T) {
	tests := []struct {
		name           string
		manifest       string
		expectedErrors []string
	}{
		{
			"valid manifest",
			"./testdata/valid.yaml",
			[]string{},
		},
		{
			"relation/permission collision",
			"./testdata/rel_perm_collision.yaml",
			[]string{
				"permission name 'file:writer' conflicts with 'file:writer' relation",
				"relation 'file:bad' has no definition",
			},
		},
		{
			"relations to undefined targets",
			"./testdata/undefined_rel_targets.yaml",
			[]string{
				"relation 'file:owner' references undefined object type 'person'",
				"relation 'file:reader' references undefined object type 'team'",
				"relation 'file:reader' references undefined object type 'project'",
				"relation 'file:writer' references undefined object type 'team'",
				"relation 'file:admin' references undefined relation type 'group#admin'",
			},
		},
		{
			"permissions to undefined targets",
			"./testdata/undefined_perm_targets.yaml",
			[]string{
				"permission 'folder:read' references undefined relation type 'folder:parent'",
				"permission 'folder:view' references undefined relation or permission 'folder:viewer'",
				"permission 'folder:view' references undefined relation or permission 'folder:guest'",
				"permission 'folder:write' references undefined relation or permission 'folder:editor'",
			},
		},
		{
			"cyclic relation definitions",
			"./testdata/invalid_cycles.yaml",
			[]string{
				"relation 'team:member' is circular and does not resolve to any object types",
				"relation 'team:owner' is circular and does not resolve to any object types",
				"relation 'project:owner' is circular and does not resolve to any object types",
			},
		},
		{
			"permissions with invalid targets",
			"./testdata/invalid_perms.yaml",
			[]string{
				"permission 'file:write' references 'owner->write', which can resolve to undefined relation or permission 'user:write'",
				"permission 'file:update' references 'parent->write', which can resolve to undefined relation or permission 'folder:write'",
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

			if len(test.expectedErrors) == 0 {
				assert.NoError(err)
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

			for _, expected := range test.expectedErrors {
				assert.ErrorContains(merr, expected)
			}

			// verify that we got the expected number of errors.
			assert.Len(merr.Errors, len(test.expectedErrors))
		})
	}
}

func TestResolution(t *testing.T) {
	assert := stretch.New(t)
	m, err := loadManifest("./testdata/valid.yaml")
	assert.NoError(err)

	// Relations
	assert.Equal([]types.ObjectName{"user"}, m.Objects["team"].Relations["owner"].SubjectTypes)
	assert.Equal([]types.ObjectName{"team"}, m.Objects["group"].Relations["owner"].SubjectTypes)
	assert.Equal([]types.ObjectName{"group"}, m.Objects["group"].Relations["parent"].SubjectTypes)

	// - order-agnostic set comparison: a subset of equal length.
	assert.Len(m.Objects["team"].Relations["member"].SubjectTypes, 2)
	assert.Subset(m.Objects["team"].Relations["member"].SubjectTypes, []types.ObjectName{"user", "team"})

	assert.Len(m.Objects["group"].Relations["member"].SubjectTypes, 2)
	assert.Subset(m.Objects["group"].Relations["member"].SubjectTypes, []types.ObjectName{"user", "team"})

	assert.Len(m.Objects["group"].Relations["manager"].SubjectTypes, 2)
	assert.Subset(m.Objects["group"].Relations["manager"].SubjectTypes, []types.ObjectName{"user", "team"})

	// Permissions
	assert.Len(m.Objects["group"].Permissions["manage"].SubjectTypes, 2)
	assert.Subset(m.Objects["group"].Permissions["manage"].SubjectTypes, []types.ObjectName{"user", "team"})

	assert.Len(m.Objects["group"].Permissions["delete"].SubjectTypes, 2)
	assert.Subset(m.Objects["group"].Permissions["delete"].SubjectTypes, []types.ObjectName{"user", "team"})

	assert.Equal([]types.ObjectName{"team"}, m.Objects["group"].Permissions["purge"].SubjectTypes)

}

func loadManifest(path string) (*model.Model, error) {
	r, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer r.Close()

	return v3.Load(r)
}
