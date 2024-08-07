package client

import (
	"math/big"
	"testing"

	"github.com/darrenvechain/thor-go-sdk/solo"
	"github.com/darrenvechain/thor-go-sdk/transaction"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
)

func TestClient_SendTransaction(t *testing.T) {
	account1 := solo.Keys()[0]
	account2 := solo.Keys()[1]
	account2Addr := crypto.PubkeyToAddress(account2.PublicKey)

	vetClause := transaction.NewClause(&account2Addr).
		WithValue(big.NewInt(1000))

	txBody := new(transaction.Builder).
		Gas(3_000_000).
		GasPriceCoef(255).
		ChainTag(client.ChainTag()).
		Expiration(100000000).
		BlockRef(transaction.NewBlockRef(0)).
		Nonce(transaction.Nonce()).
		Clause(vetClause).
		Build()

	signingHash := txBody.SigningHash()
	signature, err := crypto.Sign(signingHash[:], account1)
	assert.NoError(t, err)

	signedTx := txBody.WithSignature(signature)

	res, err := client.SendTransaction(signedTx)
	assert.NoError(t, err)
	assert.Equal(t, signedTx.ID().String(), res.ID.String())

	tx, err := client.PendingTransaction(signedTx.ID())
	assert.NoError(t, err)
	assert.Equal(t, signedTx.ID().String(), tx.ID.String())
}
