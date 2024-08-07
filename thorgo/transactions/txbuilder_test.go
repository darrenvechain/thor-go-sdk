package transactions

import (
	"math/big"
	"testing"

	"github.com/darrenvechain/thor-go-sdk/solo"
	"github.com/darrenvechain/thor-go-sdk/transaction"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
)

func TestContractClause(t *testing.T) {
	// account 1
	account1 := solo.Keys()[0]
	account1Addr := crypto.PubkeyToAddress(account1.PublicKey)

	// account 2
	account2 := solo.Keys()[1]
	account2Addr := crypto.PubkeyToAddress(account2.PublicKey)

	// transfer clause
	clause := transaction.NewClause(&account2Addr).WithData([]byte{}).WithValue(big.NewInt(1000))
	txbuilder := NewBuilder(thorClient, []*transaction.Clause{clause}, account1Addr)

	// simulation
	simulation, err := txbuilder.Simulate()
	assert.NoError(t, err)
	assert.False(t, simulation.Reverted())

	// build
	tx, err := txbuilder.Build()
	assert.NoError(t, err)

	// sign
	signingHash := tx.SigningHash()
	signature, _ := crypto.Sign(signingHash.Bytes(), account1)
	signedTx := tx.WithSignature(signature)

	// send
	res, err := thorClient.SendTransaction(signedTx)
	assert.NoError(t, err)
	assert.NotNil(t, res)
}
