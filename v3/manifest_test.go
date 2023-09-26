package v3_test

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/aserto-dev/azm/model"
	v3 "github.com/aserto-dev/azm/v3"
	"github.com/nsf/jsondiff"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func TestManifestUnmarshal(t *testing.T) {
	buf, err := os.ReadFile("./manifest.yaml")
	require.NoError(t, err)

	manifest := v3.Manifest{}
	if err := yaml.Unmarshal(buf, &manifest); err != nil {
		require.NoError(t, err)
	}

	enc := yaml.NewEncoder(os.Stderr)
	enc.SetIndent(2)
	if err := enc.Encode(&manifest); err != nil {
		require.NoError(t, err)
	}
}

func TestLoadModel(t *testing.T) {
	r, err := os.Open("./manifest.yaml")
	require.NoError(t, err)

	m1, err := v3.Load(r)
	require.NoError(t, err)
	require.NotNil(t, m1)

	b1, err := json.Marshal(m1)
	require.NoError(t, err)

	b2, err := os.ReadFile("../model/model.json")
	require.NoError(t, err)

	m2 := model.Model{}
	if err := json.Unmarshal(b2, &m2); err != nil {
		require.NoError(t, err)
	}

	opts := jsondiff.DefaultJSONOptions()
	if diff, str := jsondiff.Compare(b1, b2, &opts); diff != jsondiff.FullMatch {
		require.Equal(t, jsondiff.FullMatch, diff, "diff: %s", str)
	}
}
