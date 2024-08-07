package thorgo

import (
	"github.com/darrenvechain/thor-go-sdk/client"
	"github.com/darrenvechain/thor-go-sdk/thorgo/accounts"
	"github.com/darrenvechain/thor-go-sdk/thorgo/blocks"
	"github.com/darrenvechain/thor-go-sdk/thorgo/events"
	"github.com/darrenvechain/thor-go-sdk/thorgo/transactions"
	"github.com/darrenvechain/thor-go-sdk/thorgo/transfers"
	"github.com/darrenvechain/thor-go-sdk/transaction"
	"github.com/ethereum/go-ethereum/common"
)

type Thor struct {
	client *client.Client
	Blocks *blocks.Blocks
}

func FromURL(url string) (*Thor, error) {
	c, err := client.FromURL(url)
	if err != nil {
		return nil, err
	}

	return &Thor{client: c, Blocks: blocks.New(c)}, nil
}

func FromClient(c *client.Client) *Thor {
	return &Thor{client: c, Blocks: blocks.New(c)}
}

func (t *Thor) Account(address common.Address) *accounts.Visitor {
	return accounts.New(t.client, address)
}

func (t *Thor) Transaction(hash common.Hash) *transactions.Visitor {
	return transactions.New(t.client, hash)
}

func (t *Thor) TxBuilder(clauses []*transaction.Clause, caller common.Address) *transactions.Builder {
	return transactions.NewBuilder(t.client, clauses, caller)
}

func (t *Thor) Events(criteria []client.EventCriteria) *events.Filter {
	return events.New(t.client, criteria)
}

func (t *Thor) Transfers(criteria []client.TransferCriteria) *transfers.Filter {
	return transfers.New(t.client, criteria)
}

func (t *Thor) Client() *client.Client {
	return t.client
}
