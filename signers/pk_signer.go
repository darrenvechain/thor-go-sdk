package signers

import (
	"crypto/ecdsa"

	"github.com/darrenvechain/thor-go-sdk/thorgo"
	"github.com/darrenvechain/thor-go-sdk/thorgo/transactions"
	"github.com/darrenvechain/thor-go-sdk/transaction"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

type PKSigner struct {
	key  *ecdsa.PrivateKey
	thor *thorgo.Thor
}

func FromPK(key *ecdsa.PrivateKey, thor *thorgo.Thor) *PKSigner {
	return &PKSigner{key: key, thor: thor}
}

func GeneratePK(thor *thorgo.Thor) (*PKSigner, error) {
	key, err := crypto.GenerateKey()
	if err != nil {
		return nil, err
	}

	return &PKSigner{key: key, thor: thor}, nil
}

func (s *PKSigner) Address() common.Address {
	return crypto.PubkeyToAddress(s.key.PublicKey)
}

func (s *PKSigner) Sign(data []byte) ([]byte, error) {
	return crypto.Sign(data, s.key)
}

func (s *PKSigner) PublicKey() *ecdsa.PublicKey {
	return &s.key.PublicKey
}

func (s *PKSigner) SignTransaction(tx *transaction.Transaction) (*transaction.Transaction, error) {
	signingHash := tx.SigningHash()
	signature, err := s.Sign(signingHash.Bytes())
	if err != nil {
		return nil, err
	}

	return tx.WithSignature(signature), nil
}

func (s *PKSigner) SendClauses(clauses []*transaction.Clause) (*transactions.Visitor, error) {
	tx, err := s.thor.TxBuilder(clauses, s.Address()).Build()
	if err != nil {
		return nil, err
	}

	return s.SendTransaction(tx)
}

func (s *PKSigner) SendTransaction(tx *transaction.Transaction) (*transactions.Visitor, error) {
	signingHash := tx.SigningHash()
	signature, err := s.Sign(signingHash.Bytes())
	if err != nil {
		return nil, err
	}

	signedTx := tx.WithSignature(signature)
	res, err := s.thor.Client().SendTransaction(signedTx)
	if err != nil {
		return nil, err
	}

	return s.thor.Transaction(res.ID), nil
}
