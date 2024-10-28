package v2_test

// import (
// 	"bytes"
// 	"encoding/json"
// 	"os"
// 	"testing"

// 	v2 "github.com/aserto-dev/azm/v2"
// 	"github.com/stretchr/testify/require"
// 	"gopkg.in/yaml.v3"
// )

// func TestLoadManifest(t *testing.T) {
// 	buf, err := os.ReadFile("./manifest.yaml")
// 	require.NoError(t, err)

// 	manifest := v2.Manifest{}
// 	if err := yaml.Unmarshal(buf, &manifest); err != nil {
// 		require.NoError(t, err)
// 	}

// 	enc := yaml.NewEncoder(os.Stderr)
// 	if err := enc.Encode(&manifest); err != nil {
// 		require.NoError(t, err)
// 	}
// }

// func TestLoadModel(t *testing.T) {
// 	r, err := os.Open("./manifest.yaml")
// 	require.NoError(t, err)

// 	model, err := v2.Load(r)
// 	require.NoError(t, err)
// 	require.NotNil(t, model)

// 	w, err := os.Create("./model_test.json")
// 	require.NoError(t, err)

// 	enc := json.NewEncoder(w)
// 	enc.SetIndent("", "  ")
// 	enc.SetEscapeHTML(false)
// 	if err := enc.Encode(model); err != nil {
// 		require.NoError(t, err)
// 	}
// }

// func TestLoadEmptyManifest(t *testing.T) {
// 	r := bytes.NewReader([]byte{})

// 	m1, err := v2.Load(r)
// 	require.NoError(t, err)
// 	require.NotNil(t, m1)

// 	b1, err := json.Marshal(m1)
// 	require.NoError(t, err)
// 	require.NotNil(t, b1)
// }
