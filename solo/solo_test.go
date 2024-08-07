package solo

import (
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
)

func TestKey(t *testing.T) {
	testCases := []struct{ address string }{
		{"0xf077b491b355E64048cE21E3A6Fc4751eEeA77fa"},
		{"0x435933c8064b4Ae76bE665428e0307eF2cCFBD68"},
		{"0x0F872421Dc479F3c11eDd89512731814D0598dB5"},
		{"0xF370940aBDBd2583bC80bfc19d19bc216C88Ccf0"},
		{"0x99602e4Bbc0503b8ff4432bB1857F916c3653B85"},
		{"0x61E7d0c2B25706bE3485980F39A3a994A8207aCf"},
		{"0x361277D1b27504F36a3b33d3a52d1f8270331b8C"},
		{"0xD7f75A0A1287ab2916848909C8531a0eA9412800"},
		{"0xAbEf6032B9176C186F6BF984f548bdA53349f70a"},
		{"0x865306084235Bf804c8Bba8a8d56890940ca8F0b"},
	}

	for i, tc := range testCases {
		t.Run(tc.address, func(t *testing.T) {
			key := Keys()[i]
			addr := crypto.PubkeyToAddress(key.PublicKey).Hex()
			assert.Equal(t, tc.address, addr)
		})
	}
}
