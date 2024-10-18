package txmanager_test

import (
	"math/big"
	"testing"

	"github.com/darrenvechain/thorgo"
	"github.com/darrenvechain/thorgo/accounts"
	"github.com/darrenvechain/thorgo/crypto/tx"
	"github.com/darrenvechain/thorgo/transactions"
	"github.com/darrenvechain/thorgo/txmanager"
	"github.com/stretchr/testify/assert"
)

var (
	thor, _ = thorgo.FromURL("http://localhost:8669")
)

var (
	// PKManager should implement accounts.TxManager
	_ accounts.TxManager = &txmanager.PKManager{}
	// PKManager should implement transactions.Signer
	_ transactions.Signer = &txmanager.PKManager{}
)

// TestPKSigner demonstrates ease the ease of sending a transaction using a private key signer
func TestPKSigner(t *testing.T) {
	signer, err := txmanager.GeneratePK(thor)
	assert.NoError(t, err)

	to, err := txmanager.GeneratePK(thor)
	assert.NoError(t, err)
	toAddr := to.Address()
	vetClause := tx.NewClause(&toAddr).WithValue(big.NewInt(1000))

	tx := new(tx.Builder).
		GasPriceCoef(1).
		Gas(100000).
		Clause(vetClause).
		ChainTag(10).
		BlockRef(tx.NewBlockRef(100)).
		Build()

	signature, err := signer.SignTransaction(tx)
	assert.NoError(t, err)
	signedTx := tx.WithSignature(signature)
	origin, err := signedTx.Origin()
	assert.NoError(t, err)
	assert.Equal(t, signer.Address(), origin)
}
