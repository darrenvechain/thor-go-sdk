package accounts

import (
	"github.com/darrenvechain/thor-go-sdk/client"
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

// Revision sets the revision for the API call. Optional.
func (a *Visitor) Revision(revision common.Hash) *Visitor {
	a.revision = &revision
	return a
}

// Get fetches the account information for the given address.
func (a *Visitor) Get() (*client.Account, error) {
	if a.revision == nil {
		return a.client.Account(a.account)
	}
	return a.client.AccountForRevision(a.account, *a.revision)
}

// Code fetches the byte code of the contract at the given address.
func (a *Visitor) Code() (*client.AccountCode, error) {
	if a.revision == nil {
		return a.client.AccountCode(a.account)
	}

	return a.client.AccountCodeForRevision(a.account, *a.revision)
}

// Storage fetches the storage value for the given key.
func (a *Visitor) Storage(key common.Hash) (*client.AccountStorage, error) {
	if a.revision == nil {
		return a.client.AccountStorage(a.account, key)
	}

	return a.client.AccountStorageForRevision(a.account, key, *a.revision)
}

// Contract returns a new Contract instance.
func (a *Visitor) Contract(abi abi.ABI) *Contract {
	return NewContract(a.client, a.account, abi, a.revision)
}
