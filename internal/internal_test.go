package internal_test

import (
	"testing"

	"github.com/MarkusFreitag/keepassxc-go/internal"
	"github.com/kevinburke/nacl"
	"github.com/stretchr/testify/require"
)

func TestNaclNonceEncoding(t *testing.T) {
	nonce := nacl.NewNonce()

	encodedNonce := internal.NaclNonceToB64(nonce)
	decodedNonce := internal.B64ToNaclNonce(encodedNonce)

	require.Equal(t, nonce, decodedNonce)
}

func TestNaclKeyEncoding(t *testing.T) {
	key := nacl.NewKey()

	encodedKey := internal.NaclKeyToB64(key)
	decodedKey := internal.B64ToNaclKey(encodedKey)

	require.Equal(t, key, decodedKey)
}
