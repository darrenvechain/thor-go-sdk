package txmanager

import (
	"crypto/ecdsa"
	"errors"

	"github.com/darrenvechain/thor-go-sdk/transaction"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

type PKManager struct {
	key *ecdsa.PrivateKey
}

func FromPK(key *ecdsa.PrivateKey) *PKManager {
	return &PKManager{key: key}
}

func GeneratePK() (*PKManager, error) {
	key, err := crypto.GenerateKey()
	if err != nil {
		return nil, err
	}

	return &PKManager{key: key}, nil
}

func (p *PKManager) Address() (addr common.Address) {
	return crypto.PubkeyToAddress(p.key.PublicKey)
}

func (p *PKManager) PublicKey() *ecdsa.PublicKey {
	return &p.key.PublicKey
}

func (p *PKManager) SignTransaction(tx *transaction.Transaction) (*transaction.Transaction, error) {
	var signingHash common.Hash

	if tx.Features().IsDelegated() {
		signingHash = tx.DelegatorSigningHash(p.Address())
	} else {
		signingHash = tx.SigningHash()
	}

	signature, err := crypto.Sign(signingHash.Bytes(), p.key)
	if err != nil {
		return nil, err
	}

	return tx.WithSignature(signature), nil
}

func (p *PKManager) SignDelegated(tx *transaction.Transaction, delegatorSig []byte) (*transaction.Transaction, error) {
	if !tx.Features().IsDelegated() {
		return nil, errors.New("cannot sign non-delegated transaction")
	}

	signature, err := crypto.Sign(tx.SigningHash().Bytes(), p.key)
	if err != nil {
		return nil, err
	}

	// signature = originSig + delegatorSig
	signature = append(signature, delegatorSig...)

	return tx.WithSignature(signature), nil
}

func (p *PKManager) DelegateTransaction(tx *transaction.Transaction, origin common.Address) ([]byte, error) {
	if !tx.Features().IsDelegated() {
		return nil, errors.New("cannot sign non-delegated transaction")
	}

	return crypto.Sign(tx.DelegatorSigningHash(origin).Bytes(), p.key)
}
