package thorgo

import (
	"github.com/darrenvechain/thorgo/accounts"
	"github.com/darrenvechain/thorgo/blocks"
	"github.com/darrenvechain/thorgo/client"
	"github.com/darrenvechain/thorgo/crypto/tx"
	"github.com/darrenvechain/thorgo/events"
	"github.com/darrenvechain/thorgo/transactions"
	"github.com/darrenvechain/thorgo/transfers"
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

// Account can be used to query account information such as balance, code, storage, etc.
// It also provides a way to interact with contracts.
func (t *Thor) Account(address common.Address) *accounts.Visitor {
	return accounts.New(t.Client, address)
}

// Transaction provides utility functions to fetch or wait for transactions and their receipts.
func (t *Thor) Transaction(hash common.Hash) *transactions.Visitor {
	return transactions.New(t.Client, hash)
}

// Transactor creates a new transaction builder which makes it easier to build, simulate, build and send transactions.
func (t *Thor) Transactor(clauses []*tx.Clause) *transactions.Transactor {
	return transactions.NewTransactor(t.Client, clauses)
}

// Events sets up a query builder to fetch smart contract solidity events.
func (t *Thor) Events(criteria []client.EventCriteria) *events.Filter {
	return events.New(t.Client, criteria)
}

// Transfers sets up a query builder to fetch VET transfers.
func (t *Thor) Transfers(criteria []client.TransferCriteria) *transfers.Filter {
	return transfers.New(t.Client, criteria)
}

// Deployer makes it easier to deploy contracts.
func (t *Thor) Deployer(bytecode []byte, abi *abi.ABI) *accounts.Deployer {
	return accounts.NewDeployer(t.Client, bytecode, abi)
}
