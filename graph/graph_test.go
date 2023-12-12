package graph_test

import (
	"os"
	"testing"

	v3 "github.com/aserto-dev/azm/v3"
	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/require"
)

func TestGraph(t *testing.T) {
	r, err := os.Open("../walk/walk_test.yaml")
	require.NoError(t, err)

	m, err := v3.Load(r)
	require.NoError(t, err)
	require.NotNil(t, m)
	g := m.GetGraph()
	require.NotNil(t, g)
	spew.Dump(g)
}
