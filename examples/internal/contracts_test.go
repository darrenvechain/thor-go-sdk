package internal

import (
	"math/big"
	"strings"
	"testing"

	"github.com/darrenvechain/thor-go-sdk/thorgo/accounts"
	"github.com/darrenvechain/thor-go-sdk/txmanager"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/stretchr/testify/assert"
)

var (
	vthoABI abi.ABI
	vtho    *accounts.Contract
)

func init() {
	vthoABI, _ = abi.JSON(strings.NewReader(erc20ABI))
	vtho = thor.Account(vthoAddr).Contract(vthoABI)
}

func TestContractRead(t *testing.T) {
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

func TestContractClause(t *testing.T) {
	// construct accounts
	receiver, err := txmanager.GeneratePK()
	assert.NoError(t, err)

	// balanceOf - 1
	result, err := vtho.Call("balanceOf", account1Addr)
	assert.NoError(t, err)
	values, _ := result.([]interface{})

	// transfer 60% of the balance to account2
	amount, ok := values[0].(*big.Int)
	assert.True(t, ok)
	transferAmount := amount.Mul(amount, big.NewInt(60)).Div(amount, big.NewInt(100))

	// transfer clause
	clause, err := vtho.AsClause("transfer", receiver.Address(), transferAmount)
	assert.NoError(t, err)
	assert.NoError(t, err)

	// decoded
	decoded, err := vthoABI.Unpack("transfer", clause.Data())
	assert.NoError(t, err)
	assert.Equal(t, receiver.Address().Hex(), decoded[0].(string))
	assert.Equal(t, transferAmount, decoded[1].(*big.Int))
}

func TestContractSend(t *testing.T) {
	signer := txmanager.FromPK(account1)
	receiver, err := txmanager.GeneratePK()
	assert.NoError(t, err)

	tx, err := vtho.Send(signer, "transfer", receiver.Address(), big.NewInt(1000))
	assert.NoError(t, err)

	receipt, err := tx.Wait()
	assert.NoError(t, err)

	assert.False(t, receipt.Reverted)
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
                "name":"to",
                "type":"address"
            },
            {
                "name":"value",
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
                "name":"owner",
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
    },
    {
        "anonymous":false,
        "inputs":[
            {
                "indexed":true,
                "name":"from",
                "type":"address"
            },
            {
                "indexed":true,
                "name":"to",
                "type":"address"
            },
            {
                "indexed":false,
                "name":"value",
                "type":"uint256"
            }
        ],
        "name":"Transfer",
        "type":"event"
    }
]`
