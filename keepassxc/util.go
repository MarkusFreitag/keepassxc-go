package keepassxc

import (
	"encoding/base64"

	"github.com/kevinburke/nacl"
)

func NaclNonceToB64(nonce nacl.Nonce) string {
	return base64.StdEncoding.EncodeToString((*nonce)[:])
}

func B64ToNaclNonce(b64Nonce string) nacl.Nonce {
	decoded, err := base64.StdEncoding.DecodeString(b64Nonce)
	if err != nil {
		panic(err)
	}
	nonce := new([nacl.NonceSize]byte)
	copy(nonce[:], decoded)
	return nonce
}

func NaclKeyToB64(key nacl.Key) string {
	return base64.StdEncoding.EncodeToString((*key)[:])
}

func B64ToNaclKey(b64Key string) nacl.Key {
	decoded, err := base64.StdEncoding.DecodeString(b64Key)
	if err != nil {
		panic(err)
	}
	key := new([nacl.KeySize]byte)
	copy(key[:], decoded)
	return key
}
