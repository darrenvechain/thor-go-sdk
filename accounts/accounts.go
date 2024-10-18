package accounts

import (
	"math/big"

	"github.com/darrenvechain/thorgo/client"
	"github.com/darrenvechain/thorgo/crypto/transaction"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

type Visitor struct {
	client   *client.Client
	account  common.Address
	revision *common.Hash
}

func New(c *client.Client, account common.Address) *Visitor {
	return &Visitor{client: c, account: account}
}

// Revision sets the optional revision for the API calls.
func (a *Visitor) Revision(revision common.Hash) *Visitor {
	a.revision = &revision
	return a
}

// Get fetches the account information for the given address.
func (a *Visitor) Get() (*client.Account, error) {
	if a.revision == nil {
		return a.client.Account(a.account)
	}
	return a.client.AccountAt(a.account, *a.revision)
}

// Code fetches the byte code of the contract at the given address.
func (a *Visitor) Code() (*client.AccountCode, error) {
	if a.revision == nil {
		return a.client.AccountCode(a.account)
	}

	return a.client.AccountCodeAt(a.account, *a.revision)
}

// Storage fetches the storage value for the given key.
func (a *Visitor) Storage(key common.Hash) (*client.AccountStorage, error) {
	if a.revision == nil {
		return a.client.AccountStorage(a.account, key)
	}

	return a.client.AccountStorageAt(a.account, key, *a.revision)
}

// Call executes a read-only contract call.
func (a *Visitor) Call(calldata []byte) (*client.InspectResponse, error) {
	clause := transaction.NewClause(&a.account).WithData(calldata).WithValue(big.NewInt(0))

	request := client.InspectRequest{
		Clauses: []*transaction.Clause{clause},
	}
	var (
		inspection []client.InspectResponse
		err        error
	)

	if a.revision == nil {
		inspection, err = a.client.Inspect(request)
	} else {
		inspection, err = a.client.InspectAt(request, *a.revision)
	}

	if err != nil {
		return nil, err
	}

	return &inspection[0], nil
}

// Contract returns a new Contract instance.
func (a *Visitor) Contract(abi *abi.ABI) *Contract {
	return NewContractAt(a.client, a.account, abi, a.revision)
}
