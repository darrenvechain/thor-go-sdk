package transactions_test

import (
	"math/big"
	"testing"

	"github.com/darrenvechain/thorgo"
	"github.com/darrenvechain/thorgo/client"
	"github.com/darrenvechain/thorgo/crypto/tx"
	"github.com/darrenvechain/thorgo/solo"
	"github.com/darrenvechain/thorgo/transactions"
	"github.com/darrenvechain/thorgo/txmanager"
	"github.com/stretchr/testify/assert"
)

var (
	thorClient, _ = client.FromURL(solo.URL)
	thor          = thorgo.FromClient(thorClient)
	account1      = txmanager.FromPK(solo.Keys()[0], thor)
	account2      = txmanager.FromPK(solo.Keys()[1], thor)
)

// TestTransactions demonstrates how to build, sign, send, and wait for a transaction
func TestTransactions(t *testing.T) {
	// build a transaction
	to := account2.Address()
	vetClause := tx.NewClause(&to).WithValue(big.NewInt(1000))
	unsigned, err := transactions.NewTransactor(thorClient, []*tx.Clause{vetClause}).Build(account1.Address())
	assert.NoError(t, err)

	// sign it
	signature, err := account1.SignTransaction(unsigned)
	assert.NoError(t, err)
	signed := unsigned.WithSignature(signature)

	// send it
	res, err := thorClient.SendTransaction(signed)
	assert.NoError(t, err)

	tx := transactions.New(thorClient, res.ID)

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
