package thorgo

import (
	"github.com/darrenvechain/thor-go-sdk/client"
	"github.com/darrenvechain/thor-go-sdk/crypto/transaction"
	"github.com/darrenvechain/thor-go-sdk/thorgo/accounts"
	"github.com/darrenvechain/thor-go-sdk/thorgo/blocks"
	"github.com/darrenvechain/thor-go-sdk/thorgo/events"
	"github.com/darrenvechain/thor-go-sdk/thorgo/transactions"
	"github.com/darrenvechain/thor-go-sdk/thorgo/transfers"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

type Thor struct {
	Blocks *blocks.Blocks
	Client *client.Client
}

func FromURL(url string) (*Thor, error) {
	c, err := client.FromURL(url)
	if err != nil {
		return nil, err
	}

	return &Thor{Client: c, Blocks: blocks.New(c)}, nil
}

func FromClient(c *client.Client) *Thor {
	return &Thor{Client: c, Blocks: blocks.New(c)}
}

func (t *Thor) Account(address common.Address) *accounts.Visitor {
	return accounts.New(t.Client, address)
}

func (t *Thor) Transaction(hash common.Hash) *transactions.Visitor {
	return transactions.New(t.Client, hash)
}

func (t *Thor) Transactor(clauses []*transaction.Clause, caller common.Address) *transactions.Transactor {
	return transactions.NewTransactor(t.Client, clauses, caller)
}

func (t *Thor) Events(criteria []client.EventCriteria) *events.Filter {
	return events.New(t.Client, criteria)
}

func (t *Thor) Transfers(criteria []client.TransferCriteria) *transfers.Filter {
	return transfers.New(t.Client, criteria)
}

// Deployer creates a new contract deployer.
func (t *Thor) Deployer(bytecode []byte, abi *abi.ABI) *accounts.Deployer {
	return accounts.NewDeployer(t.Client, bytecode, abi)
}
