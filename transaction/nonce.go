package transaction

import (
	"crypto/rand"
	"encoding/binary"
)

// Nonce generates a random nonce using the crypto/rand package.
func Nonce() uint64 {
	var nonce [8]byte
	_, _ = rand.Read(nonce[:])
	return binary.BigEndian.Uint64(nonce[:])
}
