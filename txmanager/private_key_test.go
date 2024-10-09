package txmanager

import (
	"math/big"
	"testing"

	"github.com/darrenvechain/thor-go-sdk/crypto/transaction"
	"github.com/darrenvechain/thor-go-sdk/thorgo"
	"github.com/stretchr/testify/assert"
)

var (
	thor, _ = thorgo.FromURL("http://localhost:8669")
)

var (
	// PKManager should implement Manager
	_ Manager = &PKManager{}
	// PKManager should implement Signer
	_ Signer = &PKManager{}
)

// TestPKSigner demonstrates ease the ease of sending a transaction using a private key signer
func TestPKSigner(t *testing.T) {
	signer, err := GeneratePK(thor)
	assert.NoError(t, err)

	to, err := GeneratePK(thor)
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

	signature, err := signer.SignTransaction(tx)
	assert.NoError(t, err)
	signedTx := tx.WithSignature(signature)
	origin, err := signedTx.Origin()
	assert.NoError(t, err)
	assert.Equal(t, signer.Address(), origin)
}
