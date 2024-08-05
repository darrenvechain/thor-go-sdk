package internal

import (
	"math/big"
	"testing"

	"github.com/darrenvechain/thor-go-sdk/transaction"
	"github.com/darrenvechain/thor-go-sdk/txmanager"
	"github.com/stretchr/testify/assert"
)

// TestPKSigner demonstrates ease the ease of sending a transaction using a private key signer
func TestPKSigner(t *testing.T) {
	signer := txmanager.FromPK(account1)

	vetClause := transaction.NewClause(&account2Addr).WithValue(big.NewInt(1000))

	tx, err := thor.TxBuilder([]*transaction.Clause{vetClause}, signer.Address()).Send(signer)
	assert.NoError(t, err)

	receipt, err := tx.Wait()
	assert.NoError(t, err)
	assert.False(t, receipt.Reverted)
}
