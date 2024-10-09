package accounts

import (
	"bytes"
	"errors"
	"fmt"
	"math/big"

	"github.com/darrenvechain/thor-go-sdk/client"
	"github.com/darrenvechain/thor-go-sdk/crypto/transaction"
	"github.com/darrenvechain/thor-go-sdk/thorgo/transactions"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

// Contract represents a smart contract on the blockchain.
type Contract struct {
	client   *client.Client
	revision *common.Hash
	ABI      *abi.ABI
	Address  common.Address
}

func NewContract(
	client *client.Client,
	address common.Address,
	abi *abi.ABI,
	revision *common.Hash,
) *Contract {
	return &Contract{client: client, Address: address, ABI: abi, revision: revision}
}

// Call executes a read-only contract call.
func (c *Contract) Call(method string, value interface{}, args ...interface{}) error {
	packed, err := c.ABI.Pack(method, args...)
	if err != nil {
		return fmt.Errorf("failed to pack method %s: %w", method, err)
	}
	clause := transaction.NewClause(&c.Address).WithData(packed).WithValue(big.NewInt(0))
	request := client.InspectRequest{
		Clauses: []*transaction.Clause{clause},
	}
	var response []client.InspectResponse
	if c.revision == nil {
		response, err = c.client.Inspect(request)
	} else {
		response, err = c.client.InspectAt(request, *c.revision)
	}
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
	err = c.ABI.UnpackIntoInterface(value, method, decoded)
	if err != nil {
		return fmt.Errorf("failed to unpack method %s: %w", method, err)
	}
	return nil
}

// DecodeCall decodes the result of a contract call.
// The data must include the method signature.
func (c *Contract) DecodeCall(data []byte, value interface{}) error {
	var method string
	for name, m := range c.ABI.Methods {
		if len(data) >= 4 && bytes.Equal(data[:4], m.ID) {
			method = name
			break
		}
	}

	if method == "" {
		return errors.New("method signature not found")
	}

	data = data[4:]

	err := c.ABI.UnpackIntoInterface(value, method, data)
	if err != nil {
		return fmt.Errorf("failed to unpack method %s: %w", method, err)
	}
	return nil
}

// AsClause returns a transaction clause for the given method and arguments.
func (c *Contract) AsClause(method string, args ...interface{}) (*transaction.Clause, error) {
	packed, err := c.ABI.Pack(method, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to pack method %s: %w", method, err)
	}
	return transaction.NewClause(&c.Address).WithData(packed).WithValue(big.NewInt(0)), nil
}

type TxManager interface {
	SendClauses(clauses []*transaction.Clause) (common.Hash, error)
}

// Send executes a single clause
func (c *Contract) Send(manager TxManager, method string, args ...interface{}) (*transactions.Visitor, error) {
	clause, err := c.AsClause(method, args...)
	if err != nil {
		return &transactions.Visitor{}, fmt.Errorf("failed to pack method %s: %w", method, err)
	}
	txId, err := manager.SendClauses([]*transaction.Clause{clause})
	if err != nil {
		return &transactions.Visitor{}, fmt.Errorf("failed to send transaction: %w", err)
	}
	return transactions.New(c.client, txId), nil
}

// EventCriteria returns criteria that can be used to query contract events.
// The matchers must be provided in the same order as the event inputs.
// Pass nil for values that should be ignored.
func (c *Contract) EventCriteria(name string, matchers ...interface{}) (client.EventCriteria, error) {
	ev, ok := c.ABI.Events[name]
	if !ok {
		return client.EventCriteria{}, fmt.Errorf("event %s not found", name)
	}
	criteria := client.EventCriteria{
		Address: &c.Address,
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
	Log  client.EventLog
}

func (c *Contract) DecodeEvents(logs []client.EventLog) ([]Event, error) {
	var decoded []Event
	for _, log := range logs {
		if len(log.Topics) < 2 {
			continue
		}

		eventABI, err := c.ABI.EventByID(log.Topics[0])
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
		err = abi.ParseTopicsIntoMap(values, indexed, log.Topics[1:])
		if err != nil {
			return nil, err
		}

		data, err := hexutil.Decode(log.Data)
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
			Log:  log,
		})
	}
	return decoded, nil
}
