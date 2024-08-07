package transactions

import (
	"math/big"
	"testing"

	"github.com/darrenvechain/thor-go-sdk/transaction"
	"github.com/stretchr/testify/assert"
)

// TestTransactions demonstrates how to build, sign, send, and wait for a transaction
func TestTransactions(t *testing.T) {
	// build a transaction
	to := account2.Address()
	vetClause := transaction.NewClause(&to).WithValue(big.NewInt(1000))
	unsigned, err := NewBuilder(thorClient, []*transaction.Clause{vetClause}, account1.Address()).Build()
	assert.NoError(t, err)

	// sign it
	signed, err := account1.SignTransaction(unsigned)
	assert.NoError(t, err)

	// send it
	res, err := thorClient.SendTransaction(signed)
	assert.NoError(t, err)

	tx := New(thorClient, res.ID)

	// fetch the pending transaction
	pending, err := tx.Pending()
	assert.NoError(t, err)
	assert.NotNil(t, pending)

	// wait for the receipt
	receipt, err := tx.Wait()
	assert.NoError(t, err)
	assert.False(t, receipt.Reverted)

	// raw tx
	raw, err := tx.Raw()
	assert.NoError(t, err)
	assert.NotNil(t, raw)
}
