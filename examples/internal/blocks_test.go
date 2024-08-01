package internal

import (
	"github.com/darrenvechain/thor-go-sdk/solo"
	"github.com/stretchr/testify/assert"
	"testing"
)

// TestGetBestBlock fetches the best block from the network
func TestGetBestBlock(t *testing.T) {
	block, err := thor.Blocks.Best()
	assert.NoError(t, err)
	assert.NotNil(t, block)
}

// TestGetBlockByNumber fetches a block by its number
func TestGetBlockByNumber(t *testing.T) {
	block, err := thor.Blocks.ByNumber(0)
	assert.NoError(t, err)
	assert.NotNil(t, block)
}

// TestGetBlockByID fetches a block by its ID
func TestGetBlockByID(t *testing.T) {
	block, err := thor.Blocks.ByID(solo.GenesisID())
	assert.NoError(t, err)
	assert.NotNil(t, block)
}

// TestGetFinalizedBlock fetches the finalized block from the network
func TestGetFinalizedBlock(t *testing.T) {
	block, err := thor.Blocks.Finalized()
	assert.NoError(t, err)
	assert.NotNil(t, block)
}

// TestGetExpandedBlock fetches a block where all the transactions are expanded
// It accepts a revision, which can be a block ID, block number, "best" or "finalized"
func TestGetExpandedBlock(t *testing.T) {
	block, err := thor.Blocks.Expanded(solo.GenesisID().Hex())
	assert.NoError(t, err)
	assert.NotNil(t, block)
}

// TestWaitForNextBlock waits for the next block to be produced
func TestWaitForNextBlock(t *testing.T) {
	block, err := thor.Blocks.WaitForNext()
	assert.NoError(t, err)
	assert.NotNil(t, block)
}
