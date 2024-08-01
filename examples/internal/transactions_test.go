package internal

import (
	"github.com/darrenvechain/thor-go-sdk/transaction"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
	"math/big"
	"testing"
)

// TestTransactions demonstrates how to build, sign, send, and wait for a transaction
func TestTransactions(t *testing.T) {
	// build a transaction
	vetClause := transaction.NewClause(&account2Addr).WithValue(big.NewInt(1000))
	unsigned, err := thor.TxBuilder([]*transaction.Clause{vetClause}, account1Addr).Build()
	assert.NoError(t, err)

	// sign it
	signature, err := crypto.Sign(unsigned.SigningHash().Bytes(), account1)
	assert.NoError(t, err)
	signed := unsigned.WithSignature(signature)

	// send it
	res, err := thor.Client().SendTransaction(signed)
	assert.NoError(t, err)

	tx := thor.Transaction(res.ID)

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
