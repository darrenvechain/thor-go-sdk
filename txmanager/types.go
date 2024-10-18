package txmanager

import (
	"github.com/darrenvechain/thorgo/crypto/transaction"
	"github.com/ethereum/go-ethereum/common"
)

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
