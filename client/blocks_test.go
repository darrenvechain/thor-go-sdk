package client

import (
	"testing"

	"github.com/darrenvechain/thor-go-sdk/solo"
	"github.com/stretchr/testify/assert"
)

func TestClient_Block(t *testing.T) {
	block, err := client.Block("1")
	assert.NoError(t, err)
	assert.NotNil(t, block)
}

func TestClient_BestBlock(t *testing.T) {
	block, err := client.BestBlock()
	assert.NoError(t, err)
	assert.NotNil(t, block)
}

func TestClient_GenesisBlock(t *testing.T) {
	block := client.GenesisBlock()
	assert.NotNil(t, block)
}

func TestClient_ExpandedBlock(t *testing.T) {
	block, err := client.ExpandedBlock("0")
	assert.NoError(t, err)
	assert.NotNil(t, block)
}

func TestClient_ExpandedBlockWithTxs(t *testing.T) {
	c, err := FromURL("https://mainnet.vechain.org")
	assert.NoError(t, err)

	blk, err := c.ExpandedBlock("0x0125fb07988ff3c36b261b5f7227688c1c0473c4873825ac299bc256ea991b0f")
	assert.NoError(t, err)

	assert.NotNil(t, blk)
	assert.NotNil(t, blk.Transactions)
	assert.Greater(t, len(blk.Transactions), 0)
}

func TestClient_ChainTag(t *testing.T) {
	chainTag := client.ChainTag()
	assert.Equal(t, solo.ChainTag(), chainTag)
}

func TestClient_BlockRef(t *testing.T) {
	genesis, err := client.Block("0")
	assert.NoError(t, err)
	assert.NotNil(t, genesis)
	assert.Equal(t, genesis.BlockRef().Number(), uint32(0))
}
