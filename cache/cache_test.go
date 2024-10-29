package cache_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aserto-dev/azm/cache"
	v3 "github.com/aserto-dev/azm/v3"
	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
)

type (
	ON = cache.ObjectName
	RN = cache.RelationName
)

type testCase struct {
	on       ON
	sn       ON
	sr       []RN
	expected []RN
}

func (t *testCase) name() string {
	name := fmt.Sprintf("%s#?@%s", t.on, t.sn)
	switch len(t.sr) {
	case 0:
		return name
	case 1:
		return fmt.Sprintf("%s#%s", name, t.sr[0])
	default:
		srs := strings.Join(lo.Map(t.sr, func(sr RN, _ int) string { return sr.String() }), "|")
		return fmt.Sprintf("%s#[%s]", name, srs)
	}
}

func TestAssignableRelations(t *testing.T) {
	m, err := v3.LoadFile("./cache_test.yaml")
	require.NoError(t, err)
	require.NotNil(t, m)

	c := cache.New(m)

	tests := []*testCase{
		{"machine", "user", nil, []RN{"owner"}},
		{"machine", "tenant", nil, nil},
		{"group", "group", []RN{"member"}, []RN{"member", "guest"}},
		{"group", "user", nil, []RN{"member", "guest"}},
		{"tenant", "user", nil, []RN{"owner", "admin", "viewer"}},
		{"tenant", "group", nil, nil},
		{"tenant", "machine", nil, nil},
		{"tenant", "machine", []RN{"owner"}, []RN{"log_writer", "data_reader"}},
	}

	for _, test := range tests {
		t.Run(test.name(), func(tt *testing.T) {
			assert := require.New(tt)
			actual, err := c.AssignableRelations(test.on, test.sn, test.sr...)
			assert.NoError(err)
			assert.Equal(len(test.expected), len(actual))
			assert.Subset(test.expected, actual)
		})
	}
}
