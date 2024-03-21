package graph_test

import (
	"strings"
	"testing"

	"github.com/aserto-dev/azm/graph"
	"github.com/aserto-dev/azm/model"
	v3 "github.com/aserto-dev/azm/v3"
	dsc "github.com/aserto-dev/go-directory/aserto/directory/common/v3"
	"github.com/samber/lo"
	req "github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func TestInversion(t *testing.T) {
	require := req.New(t)

	m, err := v3.Load(strings.NewReader(folderModel))
	require.NoError(err)
	require.NotNil(m)

	im := m.Invert()
	require.NotNil(im)

	mnfst, err := manifest(im)
	require.NoError(err)

	b, err := yaml.Marshal(mnfst)
	require.NoError(err)

	t.Logf("inverted model:\n%s\n", b)
	// require.Fail("inverted model")
}

func TestReverseSearch(t *testing.T) {
	require := req.New(t)

	m, err := v3.Load(strings.NewReader(folderModel))
	require.NoError(err)
	require.NotNil(m)

	rm, err := v3.Load(strings.NewReader(inverseModel))
	require.NoError(err)
	require.NotNil(rm)

	folderIsOwner := rm.Objects["user"].Permissions["folder_is_owner"]
	require.NotNil(folderIsOwner)

	require.Equal([]model.ObjectName{"folder"}, folderIsOwner.SubjectTypes)

	rels := NewRelationsReader(
		"folder:leaf#parent@folder:branch",
		"folder:branch#parent@folder:root",
		"folder:root#owner@group:admins",
		"doc:leaf_doc#parent@folder:leaf",
		"doc:root_doc#parent@folder:root",
		"group:admins#member@user:admin",
		"group:admins#member@group:superusers#member",
		"group:superusers#member@user:su",
	)

	for _, test := range reverseSearchTests {
		t.Run(test.search, func(tt *testing.T) {
			assert := req.New(tt)

			req := invertedGraphReq(test.search)
			s, err := graph.NewSubjectSearch(rm, req, reverseLookup(rm, rels.GetRelations))
			assert.NoError(err)

			res, err := s.Search()
			assert.NoError(err)

			tt.Logf("request: +%v\n", req)
			tt.Logf("explanation: +%v\n", res.Explanation)
			tt.Logf("trace: +%v\n", res.Trace)

			objects := lo.Map(res.Results, func(s *dsc.ObjectIdentifier, _ int) object {
				return object{
					Type: model.ObjectName(s.ObjectType),
					ID:   model.ObjectID(s.ObjectId),
				}
			})

			for _, e := range test.expected {
				assert.Contains(objects, e)
			}

			assert.Equal(len(test.expected), len(objects), objects)
		})
	}

}

func manifest(m *model.Model) (*v3.Manifest, error) {
	mnfst := v3.Manifest{
		ModelInfo: &v3.ModelInfo{Version: v3.SchemaVersion(v3.SupportedSchemaVersion)},
		ObjectTypes: lo.MapEntries(m.Objects, func(on model.ObjectName, o *model.Object) (v3.ObjectTypeName, *v3.ObjectType) {
			return v3.ObjectTypeName(on), &v3.ObjectType{
				Relations: lo.MapEntries(o.Relations, func(rn model.RelationName, r *model.Relation) (v3.RelationName, string) {
					return v3.RelationName(rn), strings.Join(
						lo.Map(r.Union, func(rr *model.RelationRef, _ int) string {
							return rr.String()
						}),
						" | ",
					)
				}),
				Permissions: lo.MapEntries(o.Permissions, func(pn model.RelationName, p *model.Permission) (v3.PermissionName, string) {
					name := v3.PermissionName(pn)
					var (
						terms    []*model.PermissionTerm
						operator string
					)
					switch {
					case p.IsUnion():
						terms = p.Union
						operator = " | "
					case p.IsIntersection():
						terms = p.Intersection
						operator = " & "
					case p.IsExclusion():
						terms = []*model.PermissionTerm{p.Exclusion.Include, p.Exclusion.Exclude}
						operator = " - "
					}

					return name, strings.Join(lo.Map(terms, func(pt *model.PermissionTerm, _ int) string {
						return pt.String()
					}), operator)
				}),
			}
		}),
	}

	return &mnfst, nil
}

var reverseSearchTests = []searchTest{
	{"folder:?#is_owner@user:admin", []object{{"folder", "leaf"}, {"folder", "branch"}, {"folder", "root"}}},
	{"doc:?#is_owner@user:admin", []object{{"doc", "leaf_doc"}, {"doc", "root_doc"}}},
	{"folder:?#is_owner@user:su", []object{{"folder", "leaf"}, {"folder", "branch"}, {"folder", "root"}}},
	{"folder:?#is_owner@group:superusers#member", []object{{"folder", "leaf"}, {"folder", "branch"}, {"folder", "root"}}},
	{"doc:?#is_owner@group:admins#member", []object{{"doc", "leaf_doc"}, {"doc", "root_doc"}}},
	{"group:?#member@user:admin", []object{{"group", "admins"}}},
	// {"group:?#member@user:su", []object{{"group", "admins"}, {"group", "superusers"}}},
}

const folderModel = `
model:
  version: 3

types:
  user: {}

  group:
    relations:
      member: user | group#member

  folder:
    relations:
      parent: folder
      owner: user | group#member
    permissions:
      is_owner: owner | parent->is_owner

  doc:
    relations:
      parent: folder
      owner: user | group#member
      reader: user | group#member
    permissions:
      is_owner: owner | parent->is_owner
      can_read: reader | is_owner
`

const inverseModel = `
model:
  version: 3

types:
  doc: {}

  folder:
    relations:
      doc_parent: doc | folder#doc_parent
      folder_parent: folder | folder#folder_parent

  group:
    relations:
      group_member: group
      folder_owner: folder
      doc_owner: doc
      doc_reader: doc
    permissions:
      folder_is_owner: folder_owner | group_member->folder_is_owner | folder_owner->folder_parent
      doc_is_owner: doc_owner | group_member->doc_is_owner | folder_owner->doc_parent
      doc_can_read: doc_reader | group_member->doc_can_read | doc_is_owner
      r_group_member: group_member | group_member->r_group_member

  user:
    relations:
      group_member: group
      folder_owner: folder
      doc_owner: doc
      doc_reader: doc
    permissions:
      folder_is_owner: folder_owner | group_member->folder_is_owner | folder_is_owner->folder_parent
      doc_is_owner: doc_owner | group_member->doc_is_owner | folder_is_owner->doc_parent
      doc_can_read: doc_reader | group_member->doc_can_read | doc_is_owner
      r_group_member: group_member | group_member->r_group_member
`

func reverseLookup(_ *model.Model, reader graph.RelationReader) graph.RelationReader {
	return func(r *dsc.Relation) ([]*dsc.Relation, error) {
		x := strings.SplitN(r.Relation, "_", 2)

		rr := &dsc.Relation{
			ObjectType:  r.SubjectType,
			ObjectId:    r.SubjectId,
			Relation:    x[1],
			SubjectType: r.ObjectType,
			SubjectId:   r.ObjectId,
		}

		res, err := reader(rr)
		if err != nil {
			return nil, err
		}

		return lo.Map(res, func(r *dsc.Relation, _ int) *dsc.Relation {
			return &dsc.Relation{
				ObjectType:  r.SubjectType,
				ObjectId:    r.SubjectId,
				Relation:    r.Relation,
				SubjectType: r.ObjectType,
				SubjectId:   r.ObjectId,
			}
		}), nil
	}
}
