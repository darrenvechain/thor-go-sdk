package transactions

import (
	"fmt"
	"time"

	"github.com/darrenvechain/thorgo/blocks"
	"github.com/darrenvechain/thorgo/client"
	"github.com/ethereum/go-ethereum/common"
)

type Visitor struct {
	client *client.Client
	hash   common.Hash
	blocks *blocks.Blocks
}

func New(client *client.Client, hash common.Hash) *Visitor {
	return &Visitor{client: client, hash: hash, blocks: blocks.New(client)}
}

func (v *Visitor) ID() common.Hash {
	return v.hash
}

// Get fetches the transaction by its hash. This includes the clauses, but not the outputs.
func (v *Visitor) Get() (*client.Transaction, error) {
	return v.client.Transaction(v.hash)
}

// Receipt fetches the transaction receipt by its hash. This includes the outputs.
func (v *Visitor) Receipt() (*client.TransactionReceipt, error) {
	return v.client.TransactionReceipt(v.hash)
}

// Raw fetches the raw transaction by its hash.
func (v *Visitor) Raw() (*client.RawTransaction, error) {
	return v.client.RawTransaction(v.hash)
}

// Pending includes the transaction in the pending pool when querying for a transaction.
func (v *Visitor) Pending() (*client.Transaction, error) {
	return v.client.PendingTransaction(v.hash)
}

// Wait for the transaction to be included in a block.
// It will wait for ~6 blocks to be produced.
func (v *Visitor) Wait() (*client.TransactionReceipt, error) {
	return v.WaitFor(time.Minute)
}

// WaitFor the transaction to be included in a block.
// It will wait for the given duration.
// If the transaction is not included in a block within the duration, it will return an error.
func (v *Visitor) WaitFor(duration time.Duration) (*client.TransactionReceipt, error) {
	receipt, err := v.client.TransactionReceipt(v.hash)
	if err == nil {
		return receipt, nil
	}

	timeout := time.After(duration)

	for {
		select {
		case <-timeout:
			return nil, fmt.Errorf("timed out waiting for the tx receipt %s", v.hash.String())
		default:
			_, err := v.blocks.Ticker()
			if err != nil {
				time.Sleep(1 * time.Second)
			}
			receipt, err = v.client.TransactionReceipt(v.hash)
			if err == nil {
				return receipt, nil
			}
		}
	}
}
