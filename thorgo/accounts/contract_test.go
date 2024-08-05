package accounts

import (
	"math/big"
	"strings"
	"testing"

	"github.com/darrenvechain/thor-go-sdk/client"
	"github.com/darrenvechain/thor-go-sdk/solo"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
)

var (
	thorClient *client.Client
	vthoAddr   = common.HexToAddress("0x0000000000000000000000000000456E65726779")
	vtho       *Contract
)

func init() {
	var err error
	thorClient, err = client.FromURL(solo.URL)
	if err != nil {
		panic(err)
	}
	vthoABI, _ := abi.JSON(strings.NewReader(erc20ABI))
	vtho = New(thorClient, vthoAddr).Contract(vthoABI)
}

func TestContractCall(t *testing.T) {
	// name
	result, err := vtho.Call("name")
	assert.NoError(t, err)
	values, _ := result.([]interface{})
	assert.Equal(t, "VeThor", values[0])

	// symbol
	result, _ = vtho.Call("symbol")
	values, _ = result.([]interface{})
	assert.Equal(t, "VTHO", values[0])

	// decimals
	result, _ = vtho.Call("decimals")
	values, _ = result.([]interface{})
	assert.Equal(t, uint8(18), values[0])
}

func TestContractClause(t *testing.T) {
	// address 2
	account2, _ := solo.Key(1)
	account2Addr := crypto.PubkeyToAddress(account2.PublicKey)

	// transfer clause
	clause, err := vtho.AsClause("transfer", account2Addr, big.NewInt(1000))
	assert.NoError(t, err)
	assert.NotNil(t, clause)
	assert.Equal(t, clause.Value(), big.NewInt(0))
	assert.Equal(t, clause.To().Hex(), vthoAddr.Hex())
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
	}
]`
