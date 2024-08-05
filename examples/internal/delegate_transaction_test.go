package internal

import (
	"github.com/darrenvechain/thor-go-sdk/transaction"
	"github.com/darrenvechain/thor-go-sdk/txmanager"
	"github.com/stretchr/testify/assert"
	"math/big"
	"testing"
)

func TestDelegation(t *testing.T) {
	origin := txmanager.FromPK(account1)
	delegator := txmanager.FromPK(account2)

	clause := transaction.NewClause(&account3Addr).WithValue(big.NewInt(1000))
	tx, err := thor.TxBuilder([]*transaction.Clause{clause}, origin.Address()).
		Delegate().
		Build()
	assert.NoError(t, err)

	delegatorSig, err := delegator.DelegateTransaction(tx, origin.Address())
	assert.NoError(t, err)
	tx, err = origin.SignDelegated(tx, delegatorSig)
	assert.NoError(t, err)

	res, err := thor.Client().SendTransaction(tx)
	assert.NoError(t, err)

	receipt, err := thor.Transaction(res.ID).Wait()
	assert.NoError(t, err)
	assert.False(t, receipt.Reverted)

	// Check if the transaction was delegated
	assert.Equal(t, receipt.GasPayer, delegator.Address())
	assert.Equal(t, receipt.Meta.TxOrigin, origin.Address())
}
