package accounts

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

func TestContract_AsClause(t *testing.T) {
	receiver, err := txmanager.GeneratePK()
	assert.NoError(t, err)

	// transfer clause
	clause, err := vtho.AsClause("transfer", receiver.Address(), big.NewInt(1000))
	assert.NoError(t, err)
	assert.Equal(t, clause.Value(), big.NewInt(0))
	assert.Equal(t, clause.To().Hex(), vthoAddr.Hex())
}

func TestContract_Send(t *testing.T) {
	receiver, err := txmanager.GeneratePK()
	assert.NoError(t, err)

	tx, err := vtho.Send(account1, "transfer", receiver.Address(), big.NewInt(1000))
	assert.NoError(t, err)

	receipt, err := tx.Wait()
	assert.NoError(t, err)
	assert.False(t, receipt.Reverted)
}

func TestContract_EventCriteria(t *testing.T) {
	receiver, err := txmanager.GeneratePK()
	assert.NoError(t, err)

	tx, err := vtho.Send(account1, "transfer", receiver.Address(), big.NewInt(1000))
	assert.NoError(t, err)

	receipt, _ := tx.Wait()
	assert.False(t, receipt.Reverted)

	// event criteria - match the newly created receiver
	criteria, err := vtho.EventCriteria("Transfer", nil, receiver.Address())
	assert.NoError(t, err)

	// fetch events
	transfers, err := events.New(thorClient, []client.EventCriteria{criteria}).Apply(0, 100)
	assert.NoError(t, err)

	// decode events
	decodedEvs, err := vtho.DecodeEvents(*transfers)
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
