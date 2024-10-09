package txmanager

import (
	"github.com/darrenvechain/thor-go-sdk/crypto/transaction"
	"github.com/ethereum/go-ethereum/common"
)

type Options struct {
	ChainTag   uint8
	BlockRef   string
	Expiration uint
	Gas        uint64
	GasPrice   uint8
	Nonce      uint64
	DependsOn  common.Hash
}

// Manager represents a transaction manager. It is used to send transactions to the blockchain
type Manager interface {
	SendClauses(clauses []*transaction.Clause) (common.Hash, error)
}

// Signer is used for signing transactions
type Signer interface {
	Address() common.Address
	SignTransaction(tx *transaction.Transaction) ([]byte, error)
}

// Delegator handles the payment of transaction fees
type Delegator interface {
	Delegate(tx *transaction.Transaction, origin common.Address) ([]byte, error)
}

type DelegateRequest struct {
	Origin string `json:"origin"`
	Raw    string `json:"raw"`
}

type DelegateResponse struct {
	Signature string `json:"signature"`
}
