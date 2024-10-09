package accounts

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/darrenvechain/thor-go-sdk/client"
	"github.com/darrenvechain/thor-go-sdk/crypto/transaction"
	"github.com/darrenvechain/thor-go-sdk/thorgo/transactions"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

type Deployer struct {
	client   *client.Client
	bytecode []byte
	abi      *abi.ABI
	value    *big.Int
}

func NewDeployer(client *client.Client, bytecode []byte, abi *abi.ABI) *Deployer {
	return &Deployer{client: client, bytecode: bytecode, abi: abi, value: big.NewInt(0)}
}

func (d *Deployer) WithValue(value *big.Int) *Deployer {
	d.value = value
	return d
}

func (d *Deployer) Deploy(sender TxManager, args ...interface{}) (*Contract, common.Hash, error) {
	contractArgs, err := d.abi.Pack("", args...)
	txID := common.Hash{}
	if err != nil {
		return nil, txID, fmt.Errorf("failed to pack contract arguments: %w", err)
	}
	bytecode := append(d.bytecode, contractArgs...)
	clause := transaction.NewClause(nil).WithData(bytecode).WithValue(d.value)
	txID, err = sender.SendClauses([]*transaction.Clause{clause})
	if err != nil {
		return nil, txID, fmt.Errorf("failed to send contract deployment transaction: %w", err)
	}
	receipt, err := transactions.New(d.client, txID).Wait()
	if err != nil {
		return nil, txID, fmt.Errorf("failed to wait for contract deployment: %w", err)
	}
	if receipt.Reverted {
		return nil, txID, errors.New("contract deployment reverted")
	}

	address := common.HexToAddress(receipt.Outputs[0].ContractAddress)

	return NewContract(d.client, address, d.abi, nil), txID, nil
}
