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
)
