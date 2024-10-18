package txmanager

import (
	"crypto/ecdsa"

	"github.com/darrenvechain/thorgo"
	"github.com/darrenvechain/thorgo/crypto/tx"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

type PKManager struct {
	key  *ecdsa.PrivateKey
	thor *thorgo.Thor
}

func FromPK(key *ecdsa.PrivateKey, thor *thorgo.Thor) *PKManager {
	return &PKManager{key: key, thor: thor}
}

func GeneratePK(thor *thorgo.Thor) (*PKManager, error) {
	key, err := crypto.GenerateKey()
	if err != nil {
		return nil, err
	}

	return &PKManager{key: key, thor: thor}, nil
}

func (p *PKManager) Address() (addr common.Address) {
	return crypto.PubkeyToAddress(p.key.PublicKey)
}

func (p *PKManager) PublicKey() *ecdsa.PublicKey {
	return &p.key.PublicKey
}

func (p *PKManager) SendClauses(clauses []*tx.Clause) (common.Hash, error) {
	tx, err := p.thor.Transactor(clauses).Build(p.Address())
	if err != nil {
		return common.Hash{}, err
	}
	signature, err := p.SignTransaction(tx)
	if err != nil {
		return common.Hash{}, err
	}
	res, err := p.thor.Client.SendTransaction(tx.WithSignature(signature))
	if err != nil {
		return common.Hash{}, err
	}
	return res.ID, nil
}

func (p *PKManager) SignTransaction(tx *tx.Transaction) ([]byte, error) {
	signature, err := crypto.Sign(tx.SigningHash().Bytes(), p.key)
	if err != nil {
		return nil, err
	}

	return signature, nil
}
