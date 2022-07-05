package keystore_test

import (
	"os"
	"testing"

	"bou.ke/monkey"
	"github.com/stretchr/testify/require"

	"github.com/MarkusFreitag/keepassxc-go/pkg/keystore"
)

func TestKeystore(t *testing.T) {
	fakeUserConfigDir, err := os.MkdirTemp("", "fakeUserConfigDir")
	require.Nil(t, err)
	require.NotEmpty(t, fakeUserConfigDir)
	defer os.RemoveAll(fakeUserConfigDir)

	monkey.Patch(os.UserConfigDir, func() (string, error) {
		return fakeUserConfigDir, nil
	})
	defer monkey.Unpatch(os.UserConfigDir)

	store, err := keystore.Load()
	require.Nil(t, err)
	require.NotNil(t, store)
	require.Equal(t, 0, len(store.Profiles))

	profile, err := store.Get("abc")
	require.Nil(t, profile)
	require.NotNil(t, err)
	require.Equal(t, "profile named 'abc' not found", err.Error())

	expectedProfile := keystore.Profile{Name: "abc", Key: "secretkey"}
	err = store.Add(&expectedProfile)
	require.Nil(t, err)

	err = store.Add(&expectedProfile)
	require.NotNil(t, err)
	require.Equal(t, "profile named 'abc' already exists", err.Error())

	profile, err = store.Get("abc")
	require.Nil(t, err)
	require.NotNil(t, profile)
	require.Equal(t, expectedProfile, *profile)

	err = store.Save()
	require.Nil(t, err)

	store, err = keystore.Load()
	require.Nil(t, err)
	require.NotNil(t, store)
	require.Equal(t, 1, len(store.Profiles))
	require.Equal(t, "abc", store.Profiles[0].Name)
	require.Equal(t, "secretkey", store.Profiles[0].Key)
}
