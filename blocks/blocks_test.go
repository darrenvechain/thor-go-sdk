package blocks

import (
	"testing"

	"github.com/darrenvechain/thorgo/client"
	"github.com/darrenvechain/thorgo/solo"
	"github.com/stretchr/testify/assert"
)

var (
	thorClient, _ = client.FromURL(solo.URL)
	blocks        = New(thorClient)
)

// TestGetBestBlock fetches the best block from the network
func TestGetBestBlock(t *testing.T) {
	block, err := blocks.Best()
	assert.NoError(t, err)
	assert.NotNil(t, block)
}

// TestGetBlockByNumber fetches a block by its number
func TestGetBlockByNumber(t *testing.T) {
	block, err := blocks.ByNumber(0)
	assert.NoError(t, err)
	assert.NotNil(t, block)
}

// TestGetBlockByID fetches a block by its ID
func TestGetBlockByID(t *testing.T) {
	block, err := blocks.ByID(solo.GenesisID())
	assert.NoError(t, err)
	assert.NotNil(t, block)
}

// TestGetFinalizedBlock fetches the finalized block from the network
func TestGetFinalizedBlock(t *testing.T) {
	block, err := blocks.Finalized()
	assert.NoError(t, err)
	assert.NotNil(t, block)
}

// TestGetExpandedBlock fetches a block where all the transactions are expanded
// It accepts a revision, which can be a block ID, block number, "best" or "finalized"
func TestGetExpandedBlock(t *testing.T) {
	block, err := blocks.Expanded(solo.GenesisID().Hex())
	assert.NoError(t, err)
	assert.NotNil(t, block)
}

// TestWaitForNextBlock waits for the next block to be produced
func TestWaitForNextBlock(t *testing.T) {
	block, err := blocks.Ticker()
	assert.NoError(t, err)
	assert.NotNil(t, block)
}
