package internal

import (
	"github.com/darrenvechain/thor-go-sdk/signers"
	"github.com/darrenvechain/thor-go-sdk/transaction"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"math/big"
	"strings"
	"testing"
)

func TestContractRead(t *testing.T) {
	vthoABI, err := abi.JSON(strings.NewReader(erc20ABI))
	assert.NoError(t, err)
	vtho := thor.Account(common.HexToAddress("0x0000000000000000000000000000456e65726779")).Contract(vthoABI)

	// name
	result, err := vtho.Call("name")
	assert.NoError(t, err)
	values, _ := result.([]interface{})
	assert.Equal(t, "VeThor", values[0])

	// symbol
	result, err = vtho.Call("symbol")
	values, _ = result.([]interface{})
	assert.Equal(t, "VTHO", values[0])

	// decimals
	result, err = vtho.Call("decimals")
	values, _ = result.([]interface{})
	assert.Equal(t, uint8(18), values[0])
}

func TestContractWrite(t *testing.T) {
	// construct contract
	vthoABI, err := abi.JSON(strings.NewReader(erc20ABI))
	assert.NoError(t, err)
	vtho := thor.Account(common.HexToAddress("0x0000000000000000000000000000456e65726779")).Contract(vthoABI)

	// construct accounts
	sender := signers.FromPK(account1, thor)
	receiver, err := signers.GeneratePK(thor)
	assert.NoError(t, err)

	// balanceOf - 1
	result, err := vtho.Call("balanceOf", account1Addr)
	assert.NoError(t, err)
	values, _ := result.([]interface{})

	// transfer half 60% of the balance to account2
	amount, ok := values[0].(*big.Int)
	assert.True(t, ok)
	transferAmount := amount.Mul(amount, big.NewInt(60)).Div(amount, big.NewInt(100))

	// transfer clause
	clause, err := vtho.AsClause("transfer", receiver.Address(), transferAmount)
	assert.NoError(t, err)
	tx, err := sender.SendClauses([]*transaction.Clause{clause})
	assert.NoError(t, err)

	// wait for the receipt
	receipt, err := tx.Wait()
	assert.NoError(t, err)
	assert.False(t, receipt.Reverted)

	// balanceOf - 2
	result, err = vtho.Call("balanceOf", receiver.Address())
	assert.NoError(t, err)
	values, _ = result.([]interface{})
	assert.Equal(t, transferAmount, values[0])
}

const erc20ABI = `[
	{
		"constant":true,
		"inputs":[],
		"name":"name",
		"outputs":[
			{
				"name":"",	
				"type":"string"
			}
		],
		"payable":false,
		"stateMutability":"view",
		"type":"function"
	},
	{
		"constant":true,
		"inputs":[],
		"name":"symbol",
		"outputs":[
			{
				"name":"",
				"type":"string"
			}
		],	
		"payable":false,
		"stateMutability":"view",
		"type":"function"
	},
	{
		"constant":true,
		"inputs":[],
		"name":"decimals",	
		"outputs":[	
			{
				"name":"",	
				"type":"uint8"
			}	
		],		
		"payable":false,	
		"stateMutability":"view",
		"type":"function"	
	},
	{
		"constant":false,
		"inputs":[
			{
				"name":"_to",
				"type":"address"
			},
			{
				"name":"_value",
				"type":"uint256"
			}
		],
		"name":"transfer",
		"outputs":[
			{
				"name":"",
				"type":"bool"
			}
		],
		"payable":false,
		"stateMutability":"nonpayable",
		"type":"function"
	},
	{
		"constant":true,
		"name":"balanceOf",
		"inputs":[
			{
				"name":"_owner",
				"type":"address"
			}
		],
		"outputs":[
			{
				"name":"",
				"type":"uint256"
			}
		],
		"payable":false,
		"stateMutability":"view",
		"type":"function"
	}
]`
