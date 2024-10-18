package txmanager

import (
	"github.com/darrenvechain/thorgo/crypto/tx"
	"github.com/ethereum/go-ethereum/common"
)

// Delegator handles the payment of transaction fees
type Delegator interface {
	Delegate(tx *tx.Transaction, origin common.Address) ([]byte, error)
}

type DelegateRequest struct {
	Origin string `json:"origin"`
	Raw    string `json:"raw"`
}

type DelegateResponse struct {
	Signature string `json:"signature"`
}
