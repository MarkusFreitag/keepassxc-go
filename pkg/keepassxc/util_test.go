package keepassxc_test

import (
	"testing"

	"github.com/kevinburke/nacl"
	"github.com/stretchr/testify/require"

	"github.com/MarkusFreitag/keepassxc-go/pkg/keepassxc"
)

func TestNaclNonceEncoding(t *testing.T) {
	nonce := nacl.NewNonce()

	encodedNonce := keepassxc.NaclNonceToB64(nonce)
	decodedNonce := keepassxc.B64ToNaclNonce(encodedNonce)

	require.Equal(t, nonce, decodedNonce)
}

func TestNaclKeyEncoding(t *testing.T) {
	key := nacl.NewKey()

	encodedKey := keepassxc.NaclKeyToB64(key)
	decodedKey := keepassxc.B64ToNaclKey(encodedKey)

	require.Equal(t, key, decodedKey)
}
