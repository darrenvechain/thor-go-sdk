package txmanager

import (
	"math/big"
	"testing"

	"github.com/darrenvechain/thor-go-sdk/transaction"
	"github.com/stretchr/testify/assert"
)

// TestPKSigner demonstrates ease the ease of sending a transaction using a private key signer
func TestPKSigner(t *testing.T) {
	signer, err := GeneratePK()
	assert.NoError(t, err)

	to, err := GeneratePK()
	assert.NoError(t, err)
	toAddr := to.Address()
	vetClause := transaction.NewClause(&toAddr).WithValue(big.NewInt(1000))

	tx := new(transaction.Builder).
		GasPriceCoef(1).
		Gas(100000).
		Clause(vetClause).
		ChainTag(10).
		BlockRef(transaction.NewBlockRef(100)).
		Build()

	signedTx, err := signer.SignTransaction(tx)
	assert.NoError(t, err)
	origin, err := signedTx.Origin()
	assert.NoError(t, err)
	assert.Equal(t, signer.Address(), origin)
}
