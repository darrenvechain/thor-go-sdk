package blocks

import (
	"fmt"
	"time"

	"github.com/darrenvechain/thor-go-sdk/client"
	"github.com/ethereum/go-ethereum/common"
)

type Blocks struct {
	client *client.Client
}

func New(c *client.Client) *Blocks {
	return &Blocks{client: c}
}

func (b *Blocks) ByID(id common.Hash) (*client.Block, error) {
	return b.client.Block(id.Hex())
}

func (b *Blocks) Best() (*client.Block, error) {
	return b.client.Block("best")
}

func (b *Blocks) Finalized() (*client.Block, error) {
	return b.client.Block("finalized")
}

func (b *Blocks) ByNumber(number uint64) (*client.Block, error) {
	return b.client.Block(fmt.Sprintf("%d", number))
}

func (b *Blocks) Expanded(revision string) (*client.ExpandedBlock, error) {
	return b.client.ExpandedBlock(revision)
}

func (b *Blocks) WaitForNext() (*client.Block, error) {
	currentBlock, err := b.client.Block("best")
	if err != nil {
		return nil, err
	}

	for i := 0; i < 60; i++ {
		nextBlock, err := b.client.Block(fmt.Sprintf("%d", currentBlock.Number+1))
		if err == nil {
			return nextBlock, nil
		}
		time.Sleep(1 * time.Second / 2)
	}

	return nil, fmt.Errorf("timed out waiting for next block")
}
