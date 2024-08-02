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
