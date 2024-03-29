package keystore_test

import (
	"os"
	"testing"

	monkey "github.com/agiledragon/gomonkey/v2"
	"github.com/stretchr/testify/require"

	"github.com/MarkusFreitag/keepassxc-go/pkg/keystore"
)

func TestKeystore(t *testing.T) {
	fakeUserConfigDir, err := os.MkdirTemp("", "fakeUserConfigDir")
	require.Nil(t, err)
	require.NotEmpty(t, fakeUserConfigDir)
	defer os.RemoveAll(fakeUserConfigDir)

	patch := monkey.ApplyFunc(os.UserConfigDir, func() (string, error) {
		return fakeUserConfigDir, nil
	})
	defer patch.Reset()

	require.Nil(t, new(keystore.Profile).NaclKey())

	store, err := keystore.Load()
	require.Nil(t, err)
	require.NotNil(t, store)
	require.Equal(t, 0, len(store.Profiles))

	profile, err := store.Get("abc")
	require.Nil(t, profile)
	require.NotNil(t, err)
	require.Equal(t, keystore.ErrEmptyKeystore, err)

	expectedProfile := &keystore.Profile{Name: "abc", Key: "secretkeysecretkeysecretkeysecretkey"}
	err = store.Add(expectedProfile)
	require.Nil(t, err)

	err = store.Add(expectedProfile)
	require.NotNil(t, err)
	require.Equal(t, "profile named 'abc' already exists", err.Error())

	profile, err = store.Get("abc")
	require.Nil(t, err)
	require.NotNil(t, profile)
	require.Equal(t, *expectedProfile, *profile)

	profile, err = store.Get("")
	require.Nil(t, err)
	require.NotNil(t, profile)
	require.Equal(t, *expectedProfile, *profile)

	err = store.Save()
	require.Nil(t, err)

	store, err = keystore.Load()
	require.Nil(t, err)
	require.NotNil(t, store)
	require.Equal(t, "abc", store.Default)
	require.Equal(t, 1, len(store.Profiles))
	require.Equal(t, "abc", store.Profiles[0].Name)
	require.Equal(t, "secretkeysecretkeysecretkeysecretkey", store.Profiles[0].Key)
	require.NotNil(t, store.Profiles[0].NaclKey())

	profile, err = store.DefaultProfile()
	require.Nil(t, err)
	require.NotNil(t, profile)
	require.Equal(t, "abc", profile.Name)

	expectedProfile = &keystore.Profile{Name: "def", Key: "secretkeysecretkeysecretkeysecretkey"}
	err = store.Add(expectedProfile)
	require.Nil(t, err)

	profile, err = store.Get("")
	require.NotNil(t, err)
	require.Equal(t, keystore.ErrToManyProfiles, err)
	require.Nil(t, profile)

	profile, err = store.Get("def")
	require.Nil(t, err)
	require.NotNil(t, profile)
	require.Equal(t, *expectedProfile, *profile)

	profile, err = store.DefaultProfile()
	require.Nil(t, err)
	require.NotNil(t, profile)
	require.Equal(t, "abc", profile.Name)

	store.Default = ""

	err = store.Save()
	require.Nil(t, err)

	store, err = keystore.Load()
	require.Nil(t, err)
	require.NotNil(t, store)
	require.Empty(t, "", store.Default)

	profile, err = store.DefaultProfile()
	require.Nil(t, profile)
	require.NotNil(t, err)
	require.Equal(t, keystore.ErrToManyProfiles, err)

	store.Default = "abc"
	require.Nil(t, store.Save())

	store, err = keystore.Load()
	require.Nil(t, err)
	require.NotNil(t, store)
	require.Equal(t, "abc", store.Default)

	store.Default = "moo"
	err = store.Save()
	require.Nil(t, err)

	store, err = keystore.Load()
	require.Nil(t, store)
	require.NotNil(t, err)
	require.Equal(t, keystore.ErrDefaultProfileDoesNotExist, err)

	store = new(keystore.Keystore)
	profile, err = store.DefaultProfile()
	require.Nil(t, err)
	require.NotNil(t, profile)
	require.Empty(t, profile.Name)
}
