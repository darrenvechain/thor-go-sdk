package internal

import (
	"github.com/darrenvechain/thor-go-sdk/solo"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"testing"
)

// TestGetAccount fetches a thor solo account and checks if the balance and energy are greater than 0
func TestGetAccount(t *testing.T) {
	acc, err := thor.Account(account1Addr).Get()

	assert.NoError(t, err, "Account.Get should not return an error")
	assert.NotNil(t, acc, "Account.Get should return an account")

	assert.Greater(t, acc.Balance.Uint64(), uint64(0))
	assert.Greater(t, acc.Energy.Uint64(), uint64(0))
	assert.False(t, acc.HasCode)
}

// TestGetAccountForRevision fetches a thor solo account for the genesis block
// and checks if the balance and energy are greater than 0
func TestGetAccountForRevision(t *testing.T) {
	acc, err := thor.Account(account1Addr).Revision(solo.GenesisID()).Get()

	assert.NoError(t, err, "Account.Get should not return an error")
	assert.NotNil(t, acc, "Account.Get should return an account")

	assert.Greater(t, acc.Balance.Uint64(), uint64(0))
	assert.Greater(t, acc.Energy.Uint64(), uint64(0))
	assert.False(t, acc.HasCode)
}

// TestGetCode fetches the code of the VTHO contract and checks if the code length is greater than 2 (0x)
func TestGetCode(t *testing.T) {
	vtho, err := thor.Account(vthoAddr).Code()

	assert.NoError(t, err, "Account.Code should not return an error")
	assert.NotNil(t, vtho, "Account.Code should return a code")
	assert.Greater(t, len(vtho.Code), 2)
}

// TestGetCodeForRevision fetches the code of the VTHO contract for the genesis block
func TestGetCodeForRevision(t *testing.T) {
	vtho, err := thor.Account(vthoAddr).Revision(solo.GenesisID()).Code()

	assert.NoError(t, err, "Account.Code should not return an error")
	assert.NotNil(t, vtho, "Account.Code should return a code")
	assert.Greater(t, len(vtho.Code), 2)
}

// TestGetStorage fetches a storage position of the VTHO contract and checks if the value is empty
func TestGetStorage(t *testing.T) {
	storage, err := thor.Account(vthoAddr).Storage(common.Hash{})

	assert.NoError(t, err, "Account.Storage should not return an error")
	assert.NotNil(t, storage, "Account.Storage should return a storage")
	assert.Equal(t, common.Hash{}.Hex(), storage.Value)
}

// TestGetStorageForRevision fetches a storage position of the VTHO contract for the genesis block
func TestGetStorageForRevision(t *testing.T) {
	storage, err := thor.Account(vthoAddr).Revision(solo.GenesisID()).Storage(common.Hash{})

	assert.NoError(t, err, "Account.Storage should not return an error")
	assert.NotNil(t, storage, "Account.Storage should return a storage")
	assert.Equal(t, common.Hash{}.Hex(), storage.Value)
}
