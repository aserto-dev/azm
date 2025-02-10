package graph_test

import (
	"strings"
	"testing"

	azmgraph "github.com/aserto-dev/azm/graph"
	"github.com/aserto-dev/azm/mempool"
	v3 "github.com/aserto-dev/azm/v3"
	rq "github.com/stretchr/testify/require"
)

func TestComputedSet(t *testing.T) {
	require := rq.New(t)

	m, err := v3.LoadFile("./computed_set.yaml")
	require.NoError(err)
	require.NotNil(m)

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
