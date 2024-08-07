package blocks

import (
	"fmt"
	"sync/atomic"
	"time"

	"github.com/darrenvechain/thor-go-sdk/client"
	"github.com/ethereum/go-ethereum/common"
)

type Blocks struct {
	client *client.Client
	best   atomic.Value
}

func New(c *client.Client) *Blocks {
	return &Blocks{client: c}
}

func (b *Blocks) ByID(id common.Hash) (*client.Block, error) {
	return b.client.Block(id.Hex())
}

func (b *Blocks) Best() (block *client.Block, err error) {
	// Load the best block from the cache.
	if best, ok := b.best.Load().(*client.Block); ok {
		// Convert the timestamp to UTC time.
		bestTime := time.Unix(int64(best.Timestamp), 0).UTC()
		if time.Since(bestTime) < 10*time.Second {
			return best, nil
		}
	}

	block, err = b.client.Block("best")
	if err != nil {
		return nil, err
	}

	b.best.Store(block)
	return block, nil
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
	best, err := b.Best()
	if err != nil {
		return nil, err
	}

	// Sleep until the current block + 10 seconds
	predictedTime := time.Unix(int64(best.Timestamp), 0).Add(10 * time.Second)
	time.Sleep(time.Until(predictedTime))

	for i := 0; i < 40; i++ {
		nextBlock, err := b.client.Block(fmt.Sprintf("%d", best.Number+1))
		if err == nil {
			return nextBlock, nil
		}
		time.Sleep(1 * time.Second / 2)
	}

	return nil, fmt.Errorf("timed out waiting for next block")
}
