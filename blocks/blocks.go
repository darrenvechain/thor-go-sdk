package blocks

import (
	"fmt"
	"sync/atomic"
	"time"

	"github.com/darrenvechain/thorgo/client"
	"github.com/ethereum/go-ethereum/common"
)

type Blocks struct {
	client *client.Client
	best   atomic.Value
}

func New(c *client.Client) *Blocks {
	return &Blocks{client: c}
}

// ByID returns the block by the given ID.
func (b *Blocks) ByID(id common.Hash) (*client.Block, error) {
	return b.client.Block(id.Hex())
}

// Best returns the latest block on chain.
func (b *Blocks) Best() (block *client.Block, err error) {
	// Load the best block from the cache.
	if best, ok := b.best.Load().(*client.Block); ok {
		// Convert the timestamp to UTC time.
		bestTime := time.Unix(best.Timestamp, 0).UTC()
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

// Finalized returns the finalized block.
func (b *Blocks) Finalized() (*client.Block, error) {
	return b.client.Block("finalized")
}

// Justified returns the justified block.
func (b *Blocks) Justified() (*client.Block, error) {
	return b.client.Block("justified")
}

// ByNumber returns the block by the given number.
func (b *Blocks) ByNumber(number uint64) (*client.Block, error) {
	return b.client.Block(fmt.Sprintf("%d", number))
}

// Expanded returns the expanded block information.
// This includes the transactions and receipts.
func (b *Blocks) Expanded(revision string) (*client.ExpandedBlock, error) {
	return b.client.ExpandedBlock(revision)
}

// Ticker waits for the next block to be produced
// Returns the next block
func (b *Blocks) Ticker() (*client.Block, error) {
	best, err := b.Best()
	if err != nil {
		return nil, err
	}

	// Sleep until the current block + 10 seconds
	predictedTime := time.Unix(best.Timestamp, 0).Add(10 * time.Second)
	time.Sleep(time.Until(predictedTime))

	ticker := time.NewTicker(1 * time.Second)
	timeout := time.NewTimer(30 * time.Second)

	for {
		select {
		case <-ticker.C:
			nextBlock, err := b.client.Block(fmt.Sprintf("%d", best.Number+1))
			if err == nil {
				return nextBlock, nil
			}
		case <-timeout.C:
			return nil, fmt.Errorf("timed out waiting for next block")
		}
	}
}
