package transactions

import (
	"fmt"
	"time"

	"github.com/darrenvechain/thor-go-sdk/client"
	"github.com/darrenvechain/thor-go-sdk/thorgo/blocks"
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

func (v *Visitor) Get() (*client.Transaction, error) {
	return v.client.Transaction(v.hash)
}

func (v *Visitor) Receipt() (*client.TransactionReceipt, error) {
	return v.client.TransactionReceipt(v.hash)
}

func (v *Visitor) Raw() (*client.RawTransaction, error) {
	return v.client.RawTransaction(v.hash)
}

func (v *Visitor) Pending() (*client.Transaction, error) {
	return v.client.PendingTransaction(v.hash)
}

// Wait waits for the transaction to be included in a block.
// It will wait for 6 blocks to be produced.
func (v *Visitor) Wait() (*client.TransactionReceipt, error) {
	receipt, err := v.client.TransactionReceipt(v.hash)
	if err == nil {
		return receipt, nil
	}

	// loop 6 times (6 blocks * 10s/block = 1min)
	for i := 0; i < 6; i++ {
		v.blocks.WaitForNext()
		receipt, err = v.client.TransactionReceipt(v.hash)
		if err == nil {
			return receipt, nil
		}
	}

	return nil, fmt.Errorf("timed out waiting for the tx receipt %s", v.hash.String())
}

func (v *Visitor) WaitFor(duration time.Duration) (*client.TransactionReceipt, error) {
	// loop and sleep every second from now until the period
	for i := 0; i < int(duration.Seconds()*2); i++ {
		receipt, _ := v.client.TransactionReceipt(v.hash)
		if receipt != nil {
			return receipt, nil
		}
		time.Sleep(time.Second / 2)
	}
	return nil, fmt.Errorf("timed out waiting for the tx receipt %s", v.hash.String())
}
