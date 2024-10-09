package hdwallet

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	mnemonic    = "denial kitchen pet squirrel other broom bar gas better priority spoil cross"
	rootAccount = "0x0c1A60341E1064bEBB94e8769Bd508b11Ca2a27D"
	account0    = "0xf077b491b355E64048cE21E3A6Fc4751eEeA77fa"
)

func TestFromMnemonic(t *testing.T) {
	wallet, err := FromMnemonic(mnemonic)
	assert.NoError(t, err)

	addr, err := wallet.GetAddress()
	assert.NoError(t, err)
	assert.Equal(t, rootAccount, addr.Hex())
}

func TestFromSeed(t *testing.T) {
	seed, err := NewSeedFromMnemonic(mnemonic)
	assert.NoError(t, err)

	wallet, err := FromSeed(seed)
	assert.NoError(t, err)

	addr, err := wallet.GetAddress()
	assert.NoError(t, err)
	assert.Equal(t, rootAccount, addr.Hex())
}

func TestWallet_Child(t *testing.T) {
	wallet, err := FromMnemonic(mnemonic)
	assert.NoError(t, err)

	child := wallet.Child(0)
	addr, err := child.GetAddress()
	assert.NoError(t, err)
	assert.Equal(t, account0, addr.Hex())
}
