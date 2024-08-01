package internal

import (
	"github.com/darrenvechain/thor-go-sdk/signers"
	"github.com/darrenvechain/thor-go-sdk/transaction"
	"github.com/stretchr/testify/assert"
	"math/big"
	"testing"
)

// TestPKSigner demonstrates ease the ease of sending a transaction using a private key signer
func TestPKSigner(t *testing.T) {
	signer := signers.FromPK(account1, thor)

	vetClause := transaction.NewClause(&account2Addr).WithValue(big.NewInt(1000))

	tx, err := signer.SendClauses([]*transaction.Clause{vetClause})
	assert.NoError(t, err)

	receipt, err := tx.Wait()
	assert.NoError(t, err)
	assert.False(t, receipt.Reverted)
}
