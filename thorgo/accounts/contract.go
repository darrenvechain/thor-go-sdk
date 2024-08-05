package accounts

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/darrenvechain/thor-go-sdk/client"
	"github.com/darrenvechain/thor-go-sdk/thorgo/transactions"
	"github.com/darrenvechain/thor-go-sdk/transaction"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

// Contract represents a smart contract on the blockchain.
type Contract struct {
	client   *client.Client
	address  common.Address
	revision *common.Hash
	abi      abi.ABI
}

func NewContract(
	client *client.Client,
	address common.Address,
	abi abi.ABI,
	revision *common.Hash,
) *Contract {
	return &Contract{client: client, address: address, abi: abi, revision: revision}
}

// Call executes a read-only contract call.
func (c *Contract) Call(method string, args ...interface{}) (interface{}, error) {
	packed, err := c.abi.Pack(method, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to pack method %s: %w", method, err)
	}
	clause := transaction.NewClause(&c.address).WithData(packed).WithValue(big.NewInt(0))
	request := client.InspectRequest{
		Clauses: []*transaction.Clause{clause},
	}
	response, err := c.client.Inspect(request)
	if err != nil {
		return nil, fmt.Errorf("failed to inspect contract: %w", err)
	}
	inspection := (*response)[0]
	if inspection.Reverted {
		return nil, errors.New("contract call reverted")
	}
	if inspection.VmError != "" {
		return nil, errors.New(inspection.VmError)
	}
	decoded, err := hexutil.Decode(inspection.Data)
	if err != nil {
		return nil, fmt.Errorf("failed to decode data: %w", err)
	}
	result, err := c.abi.Unpack(method, decoded)
	if err != nil {
		return nil, fmt.Errorf("failed to unpack method %s: %w", method, err)
	}
	return result, nil
}

// AsClause returns a transaction clause for the given method and arguments.
func (c *Contract) AsClause(method string, args ...interface{}) (*transaction.Clause, error) {
	packed, err := c.abi.Pack(method, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to pack method %s: %w", method, err)
	}
	return transaction.NewClause(&c.address).WithData(packed).WithValue(big.NewInt(0)), nil
}

type TxSigner interface {
	SignTransaction(tx *transaction.Transaction) (*transaction.Transaction, error)
	Address() common.Address
}

// Send executes a single clause
func (c *Contract) Send(signer TxSigner, method string, args ...interface{}) (*transactions.Visitor, error) {
	clause, err := c.AsClause(method, args...)
	if err != nil {
		return &transactions.Visitor{}, fmt.Errorf("failed to pack method %s: %w", method, err)
	}
	tx, err := transactions.NewBuilder(c.client, []*transaction.Clause{clause}, signer.Address()).Build()
	if err != nil {
		return &transactions.Visitor{}, fmt.Errorf("failed to build transaction: %w", err)
	}
	tx, err = signer.SignTransaction(tx)
	if err != nil {
		return &transactions.Visitor{}, fmt.Errorf("failed to sign transaction: %w", err)
	}
	res, err := c.client.SendTransaction(tx)
	if err != nil {
		return &transactions.Visitor{}, fmt.Errorf("failed to send transaction: %w", err)
	}
	return transactions.New(c.client, res.ID), nil
}

func (c *Contract) Event(name string) (abi.Event, bool) {
	ev, ok := c.abi.Events[name]
	return ev, ok
}

func (c *Contract) Method(name string) (abi.Method, bool) {
	m, ok := c.abi.Methods[name]
	return m, ok
}
