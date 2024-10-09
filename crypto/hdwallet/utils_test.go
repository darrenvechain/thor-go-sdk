package hdwallet

import (
	"crypto/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewMnemonic(t *testing.T) {
	mnemonic, err := NewMnemonic(256)
	assert.NoError(t, err)
	assert.NotEmpty(t, mnemonic)
}

func TestNewSeedFromMnemonic(t *testing.T) {
	seed, err := NewSeedFromMnemonic("denial kitchen pet squirrel other broom bar gas better priority spoil cross")
	assert.NoError(t, err)
	assert.NotEmpty(t, seed)
}

func TestParseDerivationPath(t *testing.T) {
	path, err := ParseDerivationPath("m/44'/818'/0'/0")
	assert.NoError(t, err)
	assert.NotEmpty(t, path)
}

func TestNewSeed(t *testing.T) {
	seed, err := NewSeed()
	assert.NoError(t, err)
	assert.NotEmpty(t, seed)
}

func TestNewEntropy(t *testing.T) {
	entropy, err := NewEntropy(256)
	assert.NoError(t, err)
	assert.NotEmpty(t, entropy)
}

func TestNewMnemonicFromEntropy(t *testing.T) {
	entropy := make([]byte, 128/8)
	rand.Read(entropy)
	mnemonic, err := NewMnemonicFromEntropy(entropy)
	assert.NoError(t, err)
	assert.NotEmpty(t, mnemonic)
}
