package transactions

import (
	"github.com/darrenvechain/thor-go-sdk/client"
	"github.com/darrenvechain/thor-go-sdk/solo"
	"github.com/ethereum/go-ethereum/common"
	"math/big"
	"testing"

	"github.com/darrenvechain/thor-go-sdk/transaction"
	"github.com/stretchr/testify/assert"
)

var (
	account1      = solo.Signers()[0]
	account2      = solo.Signers()[1]
	account3      = solo.Signers()[2]
	vthoAddr      = common.HexToAddress("0x0000000000000000000000000000456e65726779")
	thorClient, _ = client.FromURL(solo.URL)
)

func TestDelegation(t *testing.T) {
	to := account3.Address()
	clause := transaction.NewClause(&to).WithValue(big.NewInt(1000))
	tx, err := NewBuilder(thorClient, []*transaction.Clause{clause}, account1.Address()).
		Delegate().
		Build()
	assert.NoError(t, err)

	delegatorSig, err := account2.DelegateTransaction(tx, account1.Address())
	assert.NoError(t, err)
	tx, err = account1.SignDelegated(tx, delegatorSig)
	assert.NoError(t, err)

	res, err := thorClient.SendTransaction(tx)
	assert.NoError(t, err)

	receipt, err := New(thorClient, res.ID).Wait()
	assert.NoError(t, err)
	assert.False(t, receipt.Reverted)

	// Check if the transaction was delegated
	assert.Equal(t, receipt.GasPayer, account2.Address())
	assert.Equal(t, receipt.Meta.TxOrigin, account1.Address())
}
