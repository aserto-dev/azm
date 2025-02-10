package graph_test

import (
	"strings"
	"testing"

	azmgraph "github.com/aserto-dev/azm/graph"
	"github.com/aserto-dev/azm/mempool"
	"github.com/aserto-dev/azm/model"
	v3 "github.com/aserto-dev/azm/v3"
	dsc "github.com/aserto-dev/go-directory/aserto/directory/common/v3"
	"github.com/samber/lo"
	rq "github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func TestComputedSet(t *testing.T) {
	require := rq.New(t)

	m, err := v3.LoadFile("./computed_set.yaml")
	require.NoError(err)

	tests := []struct {
		check    string
		expected bool
	}{
		{"resource:album#can_view@identity:zappa", true},
		{"resource:album#can_view@user:frank", true},

		{"component:guitar#can_repair@identity:zappa", true},
		{"component:guitar#can_repair@user:frank", true},
		{"component:coil#can_repair@identity:zappa", false},
		{"component:coil#can_repair@user:frank", false},
		{"component:pickup#can_repair@identity:zappa", false},
		{"component:pickup#can_repair@user:frank", false},

		{"component:guitar#can_repair@identity:duncan", true},
		{"component:guitar#can_repair@user:seymour", true},
		{"component:pickup#can_repair@identity:duncan", true},
		{"component:pickup#can_repair@user:seymour", true},
		{"component:coil#can_repair@identity:duncan", true},
		{"component:coil#can_repair@user:seymour", true},
		{"component:magnet#can_repair@identity:duncan", false},
		{"component:magnet#can_repair@user:seymour", false},
	}

	pool := mempool.NewRelationsPool()

	for _, test := range tests {
		t.Run(test.check, func(tt *testing.T) {
			require := rq.New(tt)

			checker := azmgraph.NewCheck(m, checkReq(test.check, true), csRels.GetRelations, pool)

			res, err := checker.Check()
			require.NoError(err)
			tt.Log("trace:\n", strings.Join(checker.Trace(), "\n"))
			require.Equal(test.expected, res)
		})
	}
}

func TestComputedSetSearchSubjects(t *testing.T) {
	require := rq.New(t)
	m, err := v3.LoadFile("./computed_set.yaml")
	require.NoError(err)
	require.NotNil(m)

	tests := []searchTest{
		{"user:frank#identifier@identity:?", []object{{"identity", "zappa"}}},
		{"resource:album#can_view@user:?", []object{{"user", "frank"}}},
		{"resource:album#can_view@identity:?", []object{{"identity", "zappa"}}},
		{"component:guitar#can_repair@user:?", []object{{"user", "seymour"}, {"user", "frank"}}},
		{"component:guitar#can_repair@identity:?", []object{{"identity", "duncan"}, {"identity", "zappa"}}},
		{"component:guitar#can_repair@group:?#member", []object{{"group", "guitarists"}}},
		{"component:pickup#can_repair@identity:?", []object{{"identity", "duncan"}}},
		{"component:pickup#can_repair@user:?", []object{{"user", "seymour"}}},
	}

	pool := mempool.NewRelationsPool()
	for _, test := range tests {
		t.Run(test.search, func(tt *testing.T) {
			require := rq.New(tt)
			subjSearch, err := azmgraph.NewSubjectSearch(m, graphReq(test.search), csRels.GetRelations, pool)
			require.NoError(err)

			res, err := subjSearch.Search()
			require.NoError(err)
			tt.Logf("explanation: +%v\n", res.Explanation.AsMap())
			tt.Logf("trace: +%v\n", res.Trace)

			subjects := lo.Map(res.Results, func(s *dsc.ObjectIdentifier, _ int) object {
				return object{
					Type: model.ObjectName(s.ObjectType),
					ID:   model.ObjectID(s.ObjectId),
				}
			})

			for _, e := range test.expected {
				require.Contains(subjects, e)
			}

			require.Equal(len(test.expected), len(subjects), subjects)
		})
	}
}

func TestComputedSetSearchObjects(t *testing.T) {
	t.Skip("FIXME: this doesn't work yet")
	require := rq.New(t)
	m, err := v3.LoadFile("./computed_set.yaml")
	require.NoError(err)
	require.NotNil(m)

	im := m.Invert()
	mnfst := manifest(im)

	b, err := yaml.Marshal(mnfst)
	require.NoError(err)

	t.Logf("inverted model:\n%s\n", b)

	require.NoError(
		im.Validate(model.SkipNameValidation, model.AllowPermissionInArrowBase),
	)

	tests := []searchTest{}

	pool := mempool.NewRelationsPool()
	for _, test := range tests {
		t.Run(test.search, func(tt *testing.T) {
			require := rq.New(tt)

			objSearch, err := azmgraph.NewObjectSearch(m, graphReq(test.search), csRels.GetRelations, pool)
			require.NoError(err)

			res, err := objSearch.Search()
			require.NoError(err)
			tt.Logf("explanation: +%v\n", res.Explanation.AsMap())
			tt.Logf("trace: +%v\n", res.Trace)

			objects := lo.Map(res.Results, func(s *dsc.ObjectIdentifier, _ int) object {
				return object{
					Type: model.ObjectName(s.ObjectType),
					ID:   model.ObjectID(s.ObjectId),
				}
			})

			for _, e := range test.expected {
				require.Contains(objects, e)
			}

			require.Equal(len(test.expected), len(objects), objects)
		})
	}
}

var csRels = NewRelationsReader(
	"user:frank#identifier@identity:zappa",
	"group:guitarists#member@user:frank",
	"group:musicians#member@group:guitarists#member",
	"resource:album#viewer@group:musicians#member",

	"user:seymour#identifier@identity:duncan",
	"component:coil#maintainer@user:seymour",
	"component:string#maintainer@group:guitarists#member",
	"component:pickup#part@component:magnet",
	"component:pickup#part@component:coil",
	"component:guitar#part@component:pickup#part",
	"component:guitar#part@component:string",
)
