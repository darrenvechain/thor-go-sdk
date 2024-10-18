package accounts

import (
	"bytes"
	"errors"
	"fmt"
	"math/big"

	"github.com/darrenvechain/thorgo/client"
	"github.com/darrenvechain/thorgo/crypto/transaction"
	"github.com/darrenvechain/thorgo/transactions"
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

// NewContract creates a new contract instance.
func NewContract(
	client *client.Client,
	address common.Address,
	abi *abi.ABI,
) *Contract {
	return &Contract{client: client, Address: address, ABI: abi, revision: nil}
}

// NewContractAt creates a new contract instance at a specific revision. It should be used to query historical contract states.
func NewContractAt(
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

// DecodeCall decodes the result of a contract call, for example, decoding a clause's 'data'.
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

// Send executes a transaction with a single clause.
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

// EventCriteria generates criteria to query contract events by name.
// Matchers correspond to event input parameters and must be in the same order as the event's inputs.
// Use nil for any event input you want to ignore.
//
// For example, consider the following event:
//
//	event Transfer(address indexed from, address indexed to, uint256 value);
//
// To filter events based on the 'to' address while ignoring the 'from' address and 'value', you can pass nil for those values:
//
//	to := common.HexToAddress("0x87AA2B76f29583E4A9095DBb6029A9C41994E25B")
//	criteria, err := contract.EventCriteria("Transfer", nil, &to)
//
// Returns an EventCriteria object and any error encountered.
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

// DecodeEvents parses logs into a slice of decoded events.
// The logs are typically obtained from a contract's filtered events.
//
// For example, consider the Solidity event:
//
//	event Transfer(address indexed from, address indexed to, uint256 value);
//
// To retrieve and decode "Transfer" events where the 'to' address matches a given value, you would:
//
//	to := common.HexToAddress("0x87AA2B76f29583E4A9095DBb6029A9C41994E25B")
//	logs, _ := client.FilterEvents(contract.EventCriteria("Transfer", nil, &to))
//	events, _ := contract.DecodeEvents(logs)
//
// Once decoded, you can iterate over the events and access their name and arguments:
//
//	for _, event := range events {
//	  fmt.Println(event.Name, event.Args)
//	}
//
// This function returns a slice of decoded event objects and any error encountered.
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
