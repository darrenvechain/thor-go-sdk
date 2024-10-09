package hdwallet

import (
	"crypto/ecdsa"
	"errors"

	"github.com/btcsuite/btcd/btcutil/hdkeychain"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

// PathVET is the root path to which custom derivation endpoints are appended.
// As such, the first account will be at m/44'/818'/0'/0, the second
// at m/44'/818'/0'/1, etc.
var PathVET = accounts.DerivationPath{0x80000000 + 44, 0x80000000 + 818, 0x80000000 + 0, 0}
var PathETH = accounts.DefaultRootDerivationPath

// Wallet is the underlying wallet struct.
type Wallet struct {
	masterKey *hdkeychain.ExtendedKey
	seed      []byte
	path      accounts.DerivationPath
}

// FromSeed generates a wallet from a BIP-39 seed.
func FromSeed(seed []byte) (*Wallet, error) {
	return FromSeedAt(seed, PathVET)
}

// FromSeedAt generates a wallet from a BIP-39 seed and a specific derivation path.
func FromSeedAt(seed []byte, path accounts.DerivationPath) (*Wallet, error) {
	masterKey, err := hdkeychain.NewMaster(seed, &chaincfg.MainNetParams)
	if err != nil {
		return nil, err
	}

	return &Wallet{
		masterKey: masterKey,
		seed:      seed,
		path:      path,
	}, nil
}

// FromMnemonic generates a wallet using the PathVET derivation path.
func FromMnemonic(mnemonic string) (*Wallet, error) {
	return FromMnemonicAt(mnemonic, PathVET)
}

// FromMnemonicAt generates a wallet from a BIP-39 mnemonic and a specific derivation path.
func FromMnemonicAt(mnemonic string, path accounts.DerivationPath) (*Wallet, error) {
	seed, err := NewSeedFromMnemonic(mnemonic)
	if err != nil {
		return nil, err
	}

	masterKey, err := hdkeychain.NewMaster(seed, &chaincfg.MainNetParams)
	if err != nil {
		return nil, err
	}

	return &Wallet{
		masterKey: masterKey,
		seed:      seed,
		path:      path,
	}, nil
}

// Derive returns a new wallet derived from the root seed and master key at the given path.
func (w *Wallet) Derive(path accounts.DerivationPath) *Wallet {
	return &Wallet{
		masterKey: w.masterKey,
		seed:      w.seed,
		path:      path,
	}
}

// Child returns a child of the current wallet at the given index.
func (w *Wallet) Child(index uint32) *Wallet {
	path := append(w.path, index)
	return w.Derive(path)
}

// GetPrivateKey returns the ECDSA private key of the account.
func (w *Wallet) GetPrivateKey() (*ecdsa.PrivateKey, error) {
	return w.derivePrivateKey(w.path)
}

// MustGetPrivateKey returns the ECDSA private key of the account.
func (w *Wallet) MustGetPrivateKey() *ecdsa.PrivateKey {
	privateKey, err := w.GetPrivateKey()
	if err != nil {
		panic(err)
	}
	return privateKey
}

// GetPublicKey returns the ECDSA public key of the account.
func (w *Wallet) GetPublicKey() (*ecdsa.PublicKey, error) {
	return w.derivePublicKey(w.path)
}

// MustGetPublicKey returns the ECDSA public key of the account.
func (w *Wallet) MustGetPublicKey() *ecdsa.PublicKey {
	publicKey, err := w.GetPublicKey()
	if err != nil {
		panic(err)
	}
	return publicKey
}

// GetAddress returns the address of the current master key.
func (w *Wallet) GetAddress() (common.Address, error) {
	return w.deriveAddress(w.path)
}

// MustGetAddress returns the address of the current master key.
func (w *Wallet) MustGetAddress() common.Address {
	addr, err := w.GetAddress()
	if err != nil {
		panic(err)
	}
	return addr
}

// DerivePrivateKey derives the private key of the derivation path.
func (w *Wallet) derivePrivateKey(path accounts.DerivationPath) (*ecdsa.PrivateKey, error) {
	var err error
	key := w.masterKey
	for _, n := range path {
		key, err = key.Derive(n)
		if err != nil {
			return nil, err
		}
	}

	privateKey, err := key.ECPrivKey()
	if err != nil {
		return nil, err
	}
	privateKeyECDSA := privateKey.ToECDSA()

	return privateKeyECDSA, nil
}

// DerivePublicKey derives the public key of the derivation path.
func (w *Wallet) derivePublicKey(path accounts.DerivationPath) (*ecdsa.PublicKey, error) {
	privateKeyECDSA, err := w.derivePrivateKey(path)
	if err != nil {
		return nil, err
	}

	publicKey := privateKeyECDSA.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, errors.New("failed to get public key")
	}

	return publicKeyECDSA, nil
}

// DeriveAddress derives the account address of the derivation path.
func (w *Wallet) deriveAddress(path accounts.DerivationPath) (common.Address, error) {
	publicKeyECDSA, err := w.derivePublicKey(path)
	if err != nil {
		return common.Address{}, err
	}

	address := crypto.PubkeyToAddress(*publicKeyECDSA)
	return address, nil
}
