package transactions_test

import (
	"math/big"
	"testing"

	"github.com/darrenvechain/thor-go-sdk/client"
	"github.com/darrenvechain/thor-go-sdk/solo"
	"github.com/darrenvechain/thor-go-sdk/thorgo"
	"github.com/darrenvechain/thor-go-sdk/thorgo/transactions"
	"github.com/darrenvechain/thor-go-sdk/transaction"
	"github.com/darrenvechain/thor-go-sdk/txmanager"
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
	vetClause := transaction.NewClause(&to).WithValue(big.NewInt(1000))
	unsigned, err := transactions.NewBuilder(thorClient, []*transaction.Clause{vetClause}, account1.Address()).Build()
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
