package azm_test

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/aserto-dev/azm"
	v2 "github.com/aserto-dev/azm/v2"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v2"
)

func TestJSON(t *testing.T) {
	r, err := os.Open("./v2/test/test.yaml")
	require.NoError(t, err)
	require.NotNil(t, r)
	defer r.Close()

	m1, err := v2.Model().Read(r)
	require.NoError(t, err)
	require.NotNil(t, m1)

	buf, err := json.Marshal(m1)
	if err != nil {
		require.NoError(t, err)
	}

	var m2 azm.Model
	if err := json.Unmarshal(buf, &m2); err != nil {
		require.NoError(t, err)
	}
}

func TestYAML(t *testing.T) {
	r, err := os.Open("./v2/test/test.yaml")
	require.NoError(t, err)
	require.NotNil(t, r)
	defer r.Close()

	m1, err := v2.Model().Read(r)
	require.NoError(t, err)
	require.NotNil(t, m1)

	buf, err := yaml.Marshal(m1)
	if err != nil {
		require.NoError(t, err)
	}
	var m2 azm.Model
	if err := yaml.Unmarshal(buf, &m2); err != nil {
		require.NoError(t, err)
	}
}
