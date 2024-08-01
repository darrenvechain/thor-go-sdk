package accounts

import (
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/common/hexutil"

	"github.com/darrenvechain/thor-go-sdk/client"
	"github.com/darrenvechain/thor-go-sdk/transaction"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

// Contract represents a smart contract on the blockchain.
type Contract struct {
	client   *client.Client
	account  common.Address
	revision *common.Hash
	abi      abi.ABI
}

func NewContract(
	client *client.Client,
	address common.Address,
	abi abi.ABI,
	revision *common.Hash,
) *Contract {
	return &Contract{client: client, account: address, abi: abi, revision: revision}
}

// Call executes a read-only contract call.
func (m *Contract) Call(method string, args ...interface{}) (interface{}, error) {
	packed, err := m.abi.Pack(method, args...)
	if err != nil {
		return nil, err
	}
	clause := transaction.NewClause(&m.account).WithData(packed).WithValue(big.NewInt(0))
	request := client.InspectRequest{
		Clauses: []*transaction.Clause{clause},
	}
	response, err := m.client.Inspect(request)
	if err != nil {
		return nil, err
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
		return nil, err
	}
	result, err := m.abi.Unpack(method, decoded)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// AsClause returns a transaction clause for the given method and arguments.
func (m *Contract) AsClause(method string, args ...interface{}) (*transaction.Clause, error) {
	packed, err := m.abi.Pack(method, args...)
	if err != nil {
		return nil, err
	}
	return transaction.NewClause(&m.account).WithData(packed).WithValue(big.NewInt(0)), nil
}
