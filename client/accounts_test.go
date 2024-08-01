package client

import (
	"github.com/darrenvechain/thor-go-sdk/hex"
	"github.com/darrenvechain/thor-go-sdk/solo"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"math/big"
	"strings"
	"testing"
)

func TestClient_Account(t *testing.T) {
	acc, err := client.Account(common.HexToAddress("0xd1d37b8913563fC25BC5bB2E669eB3dBC6b87762"))
	assert.NoError(t, err)
	assert.Equal(t, hex.Int{Int: big.NewInt(0)}.Int64(), acc.Balance.Int64())
	assert.Equal(t, hex.Int{Int: big.NewInt(0)}.Int64(), acc.Energy.Int64())
	assert.False(t, acc.HasCode)
}

func TestClient_AccountForRevision(t *testing.T) {
	acc, err := client.AccountForRevision(
		common.HexToAddress("0xd1d37b8913563fC25BC5bB2E669eB3dBC6b87762"),
		solo.GenesisID(),
	)

	assert.NoError(t, err)
	assert.Equal(t, hex.Int{Int: big.NewInt(0)}.Int64(), acc.Balance.Int64())
	assert.Equal(t, hex.Int{Int: big.NewInt(0)}.Int64(), acc.Energy.Int64())
	assert.False(t, acc.HasCode)
}

func TestClient_AccountCode(t *testing.T) {
	res, err := client.AccountCode(common.HexToAddress("0x0000000000000000000000000000456E65726779"))
	assert.NoError(t, err)
	assert.Greater(t, len(res.Code), 2)
}

func TestClient_AccountCodeForRevision(t *testing.T) {
	res, err := client.AccountCodeForRevision(
		common.HexToAddress("0x0000000000000000000000000000456E65726779"),
		solo.GenesisID(),
	)
	assert.NoError(t, err)
	assert.Greater(t, len(res.Code), 2)
}

func TestClient_AccountStorage(t *testing.T) {
	res, err := client.AccountStorage(
		common.HexToAddress("0x0000000000000000000000000000456E65726779"),
		common.HexToHash(strings.Repeat("0", 64)),
	)

	assert.NoError(t, err)
	assert.Greater(t, len(res.Value), 2)
}

func TestClient_AccountStorageForRevision(t *testing.T) {
	res, err := client.AccountStorageForRevision(
		common.HexToAddress("0x0000000000000000000000000000456E65726779"),
		common.HexToHash(strings.Repeat("0", 64)),
		solo.GenesisID(),
	)

	assert.NoError(t, err)
	assert.Greater(t, len(res.Value), 2)
}
