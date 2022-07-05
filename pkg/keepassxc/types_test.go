package keepassxc_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/MarkusFreitag/keepassxc-go/pkg/keepassxc"
)

func TestTypePassword(t *testing.T) {
	pass := keepassxc.Password("supersecret")
	require.Equal(t, "*****", pass.String())
	require.Equal(t, "supersecret", pass.Plaintext())
}

func TestTypeBoolString(t *testing.T) {
	data := struct {
		Bool keepassxc.BoolString
	}{}
	err := json.Unmarshal([]byte(`{"bool": "true"}`), &data)
	require.Nil(t, err)
	require.True(t, bool(data.Bool))

	data = struct {
		Bool keepassxc.BoolString
	}{}
	err = json.Unmarshal([]byte(`{"bool": "false"}`), &data)
	require.Nil(t, err)
	require.False(t, bool(data.Bool))
}

func TestTypeFields(t *testing.T) {
	fields := make(keepassxc.Fields, 0)
	require.Equal(t, "", fields.String())
	fields = append(fields, "first")
	require.Equal(t, "first", fields.String())
	fields = append(fields, "second", "third")
	require.Equal(t, "first,second,third", fields.String())
}
