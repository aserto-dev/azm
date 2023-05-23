package v2_test

import (
	"encoding/json"
	"os"
	"testing"

	v2 "github.com/aserto-dev/azm/v2"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v2"
)

func TestLoadModelV2FromManifest(t *testing.T) {
	r, err := os.Open("./test/test.yaml")
	require.NoError(t, err)
	require.NotNil(t, r)
	defer r.Close()

	m, err := v2.Model().Read(r)
	require.NoError(t, err)
	require.NotNil(t, m)
}

func TestModelV2YAML(t *testing.T) {
	r, err := os.Open("./test/test.yaml")
	require.NoError(t, err)
	require.NotNil(t, r)
	defer r.Close()

	m, err := v2.Model().Read(r)
	require.NoError(t, err)
	require.NotNil(t, m)

	if err := yaml.NewEncoder(os.Stdout).Encode(m); err != nil {
		require.NoError(t, err)
	}
}

func TestModelV2JSON(t *testing.T) {
	r, err := os.Open("./test/test.yaml")
	require.NoError(t, err)
	require.NotNil(t, r)
	defer r.Close()

	m, err := v2.Model().Read(r)
	require.NoError(t, err)
	require.NotNil(t, m)

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	enc.SetEscapeHTML(true)
	if err := enc.Encode(m); err != nil {
		require.NoError(t, err)
	}
}

func TestResolvePermission(t *testing.T) {
	r, err := os.Open("./test/test.yaml")
	require.NoError(t, err)
	require.NotNil(t, r)
	defer r.Close()

	m, err := v2.Model().Read(r)
	require.NoError(t, err)
	require.NotNil(t, m)

	relations, err := m.ResolvePermission("database", "read")
	require.NoError(t, err)
	require.NotNil(t, relations)
}

func TestResolveRelation(t *testing.T) {
	r, err := os.Open("./test/test.yaml")
	require.NoError(t, err)
	require.NotNil(t, r)
	defer r.Close()

	m, err := v2.Model().Read(r)
	require.NoError(t, err)
	require.NotNil(t, m)

	relations, err := m.ResolveRelation("database", "reader")
	require.NoError(t, err)
	require.NotNil(t, relations)
}
