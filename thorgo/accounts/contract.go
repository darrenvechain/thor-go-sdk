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
func (c *Contract) Call(method string, value interface{}, args ...interface{}) error {
	packed, err := c.abi.Pack(method, args...)
	if err != nil {
		return fmt.Errorf("failed to pack method %s: %w", method, err)
	}
	clause := transaction.NewClause(&c.address).WithData(packed).WithValue(big.NewInt(0))
	request := client.InspectRequest{
		Clauses: []*transaction.Clause{clause},
	}
	response, err := c.client.Inspect(request)
	if err != nil {
		return fmt.Errorf("failed to inspect contract: %w", err)
	}
	inspection := response[0]
	if inspection.Reverted {
		return errors.New("contract call reverted")
	}
	if inspection.VmError != "" {
		return errors.New(inspection.VmError)
	}
	decoded, err := hexutil.Decode(inspection.Data)
	if err != nil {
		return fmt.Errorf("failed to decode data: %w", err)
	}
	err = c.abi.UnpackIntoInterface(value, method, decoded)
	if err != nil {
		return fmt.Errorf("failed to unpack method %s: %w", method, err)
	}
	return nil
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

// EventCriteria returns criteria that can be used to query contract events.
// The matchers must be provided in the same order as the event inputs.
// Pass nil for values that should be ignored.
func (c *Contract) EventCriteria(name string, matchers ...interface{}) (client.EventCriteria, error) {
	ev, ok := c.abi.Events[name]
	if !ok {
		return client.EventCriteria{}, fmt.Errorf("event %s not found", name)
	}
	criteria := client.EventCriteria{
		Address: &c.address,
		Topic0:  &ev.ID,
	}

	for i := range ev.Inputs {
		if i >= len(matchers) {
			break
		}
		if matchers[i] == nil {
			continue
		}
		if !ev.Inputs[i].Indexed {
			return client.EventCriteria{}, errors.New("can't match non-indexed event inputs")
		}
		topics, err := abi.MakeTopics(
			[]interface{}{matchers[i]},
		)
		if err != nil {
			return client.EventCriteria{}, err
		}

		switch i + 1 {
		case 1:
			criteria.Topic1 = &topics[0][0]
		case 2:
			criteria.Topic2 = &topics[0][0]
		case 3:
			criteria.Topic3 = &topics[0][0]
		case 4:
			criteria.Topic4 = &topics[0][0]
		}
	}

	return criteria, nil
}

type Event struct {
	Name string
	Args map[string]interface{}
}

func (c *Contract) DecodeEvents(events []client.EventLog) ([]Event, error) {
	var decoded []Event
	for _, ev := range events {
		if len(ev.Topics) < 2 {
			continue
		}

		eventABI, err := c.abi.EventByID(ev.Topics[0])
		if err != nil {
			continue
		}

		var indexed abi.Arguments
		for _, arg := range eventABI.Inputs {
			if arg.Indexed {
				indexed = append(indexed, arg)
			}
		}

		values := make(map[string]interface{})
		err = abi.ParseTopicsIntoMap(values, indexed, ev.Topics[1:])
		if err != nil {
			return nil, err
		}

		data, err := hexutil.Decode(ev.Data)
		if err != nil {
			return nil, err
		}
		err = eventABI.Inputs.UnpackIntoMap(values, data)
		if err != nil {
			return nil, err
		}

		decoded = append(decoded, Event{
			Name: eventABI.Name,
			Args: values,
		})
	}
	return decoded, nil
}
