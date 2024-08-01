package transactions

import (
	"fmt"
	"time"

	"github.com/darrenvechain/thor-go-sdk/client"
	"github.com/ethereum/go-ethereum/common"
)

type Visitor struct {
	client *client.Client
	hash   common.Hash
}

func New(client *client.Client, hash common.Hash) *Visitor {
	return &Visitor{client: client, hash: hash}
}

func (t *Visitor) ID() common.Hash {
	return t.hash
}

func (t *Visitor) Get() (*client.Transaction, error) {
	return t.client.Transaction(t.hash)
}

func (t *Visitor) Receipt() (*client.TransactionReceipt, error) {
	return t.client.TransactionReceipt(t.hash)
}

func (t *Visitor) Raw() (*client.RawTransaction, error) {
	return t.client.RawTransaction(t.hash)
}

func (t *Visitor) Pending() (*client.Transaction, error) {
	return t.client.PendingTransaction(t.hash)
}

func (t *Visitor) Wait() (*client.TransactionReceipt, error) {
	var receipt *client.TransactionReceipt
	// loop until the transaction is mined, timeout after 30 seconds
	for i := 0; i < 60; i++ {
		receipt, _ = t.Receipt()
		if receipt != nil {
			return receipt, nil
		}
		time.Sleep(time.Second / 2)
	}
	return nil, fmt.Errorf("timed out waiting for the tx receipt %s", t.hash.String())
}

func (t *Visitor) WaitFor(duration time.Duration) (*client.TransactionReceipt, error) {
	// loop and sleep every second from now until the period
	for i := 0; i < int(duration.Seconds()*2); i++ {
		receipt, _ := t.client.TransactionReceipt(t.hash)
		if receipt != nil {
			return receipt, nil
		}
		time.Sleep(time.Second / 2)
	}
	return nil, fmt.Errorf("timed out waiting for the tx receipt %s", t.hash.String())
}
