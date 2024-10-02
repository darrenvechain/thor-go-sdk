package accounts_test

import (
	"math/big"
	"testing"

	"github.com/darrenvechain/thor-go-sdk/client"
	"github.com/darrenvechain/thor-go-sdk/thorgo/events"
	"github.com/darrenvechain/thor-go-sdk/txmanager"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
)

func TestContract_Call(t *testing.T) {
	// name
	var name string
	err := vthoContract.Call("name", &name)
	assert.NoError(t, err)
	assert.Equal(t, "VeThor", name)

	// symbol
	var symbol string
	err = vthoContract.Call("symbol", &symbol)
	assert.NoError(t, err)
	assert.Equal(t, "VTHO", symbol)

	// decimals
	var decimals uint8
	err = vthoContract.Call("decimals", &decimals)
	assert.NoError(t, err)
	assert.Equal(t, uint8(18), decimals)
}

func TestContract_DecodeCall(t *testing.T) {
	packed, err := vtho.ABI.Pack("balanceOf", account1.Address())
	assert.NoError(t, err)

	balance := new(big.Int)
	err = vthoContract.DecodeCall(packed, &balance)
	assert.NoError(t, err)
	assert.Greater(t, balance.Uint64(), uint64(0))
}

func TestContract_AsClause(t *testing.T) {
	receiver, err := txmanager.GeneratePK(thor)
	assert.NoError(t, err)

	// transfer clause
	clause, err := vthoContract.AsClause("transfer", receiver.Address(), big.NewInt(1000))
	assert.NoError(t, err)
	assert.Equal(t, clause.Value(), big.NewInt(0))
	assert.Equal(t, clause.To().Hex(), vtho.Address.Hex())
}

func TestContract_Send(t *testing.T) {
	receiver, err := txmanager.GeneratePK(thor)
	assert.NoError(t, err)

	tx, err := vthoContract.Send(account1, "transfer", receiver.Address(), big.NewInt(1000))
	assert.NoError(t, err)

	receipt, err := tx.Wait()
	assert.NoError(t, err)
	assert.False(t, receipt.Reverted)
}

func TestContract_EventCriteria(t *testing.T) {
	receiver, err := txmanager.GeneratePK(thor)
	assert.NoError(t, err)

	tx, err := vthoContract.Send(account1, "transfer", receiver.Address(), big.NewInt(1000))
	assert.NoError(t, err)

	receipt, _ := tx.Wait()
	assert.False(t, receipt.Reverted)

	// event criteria - match the newly created receiver
	criteria, err := vthoContract.EventCriteria("Transfer", nil, receiver.Address())
	assert.NoError(t, err)

	// fetch events
	transfers, err := events.New(thorClient, []client.EventCriteria{criteria}).Apply(0, 100)
	assert.NoError(t, err)

	// decode events
	decodedEvs, err := vthoContract.DecodeEvents(transfers)
	assert.NoError(t, err)

	ev := decodedEvs[0]
	assert.Equal(t, "Transfer", ev.Name)
	assert.NotNil(t, ev.Args["from"])
	assert.NotNil(t, ev.Args["to"])
	assert.NotNil(t, ev.Args["value"])
	assert.IsType(t, common.Address{}, ev.Args["from"])
	assert.IsType(t, common.Address{}, ev.Args["to"])
	assert.IsType(t, &big.Int{}, ev.Args["value"])
}
