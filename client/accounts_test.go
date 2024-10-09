package client

import (
	"strings"
	"testing"

	"github.com/darrenvechain/thor-go-sdk/solo"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
)

func TestClient_Account(t *testing.T) {
	acc, err := client.Account(common.HexToAddress("0xd1d37b8913563fC25BC5bB2E669eB3dBC6b87762"))
	assert.NoError(t, err)

	assert.Zero(t, acc.Balance.ToInt().Int64())
	assert.Zero(t, acc.Energy.ToInt().Int64())
	assert.False(t, acc.HasCode)
}

func TestClient_AccountAt(t *testing.T) {
	acc, err := client.AccountAt(
		common.HexToAddress("0xd1d37b8913563fC25BC5bB2E669eB3dBC6b87762"),
		solo.GenesisID(),
	)

	assert.NoError(t, err)
	assert.Equal(t, int64(0), acc.Balance.ToInt().Int64())
	assert.Equal(t, int64(0), acc.Energy.ToInt().Int64())
	assert.False(t, acc.HasCode)
}

func TestClient_AccountCode(t *testing.T) {
	res, err := client.AccountCode(common.HexToAddress("0x0000000000000000000000000000456E65726779"))
	assert.NoError(t, err)
	assert.Greater(t, len(res.Code), 2)
}

func TestClient_AccountCodeAt(t *testing.T) {
	res, err := client.AccountCodeAt(
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

func TestClient_AccountStorageAt(t *testing.T) {
	res, err := client.AccountStorageAt(
		common.HexToAddress("0x0000000000000000000000000000456E65726779"),
		common.HexToHash(strings.Repeat("0", 64)),
		solo.GenesisID(),
	)

	assert.NoError(t, err)
	assert.Greater(t, len(res.Value), 2)
}
