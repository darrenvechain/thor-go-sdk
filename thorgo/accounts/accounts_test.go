package accounts_test

import (
	"testing"

	"github.com/darrenvechain/thor-go-sdk/builtins"
	"github.com/darrenvechain/thor-go-sdk/client"
	"github.com/darrenvechain/thor-go-sdk/solo"
	"github.com/darrenvechain/thor-go-sdk/thorgo"
	"github.com/darrenvechain/thor-go-sdk/thorgo/accounts"
	"github.com/darrenvechain/thor-go-sdk/txmanager"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
)

var (
	thorClient, _ = client.FromURL(solo.URL)
	thor          = thorgo.FromClient(thorClient)
	vtho          = builtins.VTHO
	vthoContract  = vtho.Load(thor)
	account1      = txmanager.FromPK(solo.Keys()[0], thor)
)

// TestGetAccount fetches a thor solo account and checks if the balance and energy are greater than 0
func TestGetAccount(t *testing.T) {
	acc, err := accounts.New(thorClient, account1.Address()).Get()

	assert.NoError(t, err, "Account.httpGet should not return an error")
	assert.NotNil(t, acc, "Account.httpGet should return an account")

	assert.Greater(t, acc.Balance.Uint64(), uint64(0))
	assert.Greater(t, acc.Energy.Uint64(), uint64(0))
	assert.False(t, acc.HasCode)
}

// TestGetAccountForRevision fetches a thor solo account for the genesis block
// and checks if the balance and energy are greater than 0
func TestGetAccountForRevision(t *testing.T) {
	acc, err := accounts.New(thorClient, account1.Address()).Revision(solo.GenesisID()).Get()

	assert.NoError(t, err, "Account.httpGet should not return an error")
	assert.NotNil(t, acc, "Account.httpGet should return an account")

	assert.Greater(t, acc.Balance.Uint64(), uint64(0))
	assert.Greater(t, acc.Energy.Uint64(), uint64(0))
	assert.False(t, acc.HasCode)
}

// TestGetCode fetches the code of the VTHO contract and checks if the code length is greater than 2 (0x)
func TestGetCode(t *testing.T) {
	vtho, err := accounts.New(thorClient, vtho.Address).Code()

	assert.NoError(t, err, "Account.Code should not return an error")
	assert.NotNil(t, vtho, "Account.Code should return a code")
	assert.Greater(t, len(vtho.Code), 2)
}

// TestGetCodeForRevision fetches the code of the VTHO contract for the genesis block
func TestGetCodeForRevision(t *testing.T) {
	vtho, err := accounts.New(thorClient, vtho.Address).Revision(solo.GenesisID()).Code()

	assert.NoError(t, err, "Account.Code should not return an error")
	assert.NotNil(t, vtho, "Account.Code should return a code")
	assert.Greater(t, len(vtho.Code), 2)
}

// TestGetStorage fetches a storage position of the VTHO contract and checks if the value is empty
func TestGetStorage(t *testing.T) {
	storage, err := accounts.New(thorClient, vtho.Address).Storage(common.Hash{})

	assert.NoError(t, err, "Account.Storage should not return an error")
	assert.NotNil(t, storage, "Account.Storage should return a storage")
	assert.Equal(t, common.Hash{}.Hex(), storage.Value)
}

// TestGetStorageForRevision fetches a storage position of the VTHO contract for the genesis block
func TestGetStorageForRevision(t *testing.T) {
	storage, err := accounts.New(thorClient, vtho.Address).Revision(solo.GenesisID()).Storage(common.Hash{})

	assert.NoError(t, err, "Account.Storage should not return an error")
	assert.NotNil(t, storage, "Account.Storage should return a storage")
	assert.Equal(t, common.Hash{}.Hex(), storage.Value)
}
