# Thor GO SDK
The Thor GO SDK is a Golang library designed to provide an easy and intuitive way to interact with the VechainThor blockchain. It simplifies blockchain interactions, making it straightforward for developers to build and manage applications on VechainThor.

## Key Features
- **Easy-to-Use Interface**: Provides a simple and accessible API for VechainThor interactions.
- **Blockchain Interaction**: Facilitates transactions, smart contract interactions, and more.
- **Golang Support**: Leverages the power and efficiency of Go for blockchain development.

## Note on Geth
The Thor GO SDK is built on top of the latest version of [geth](https://github.com/ethereum/go-ethereum). Familiarity with the Geth repository is encouraged, particularly when working with Application Binary Interfaces (ABIs), cryptographic operations (hashing, signing, and managing private keys), and other low-level blockchain functions. Understanding these elements can help in effectively utilizing the SDK and troubleshooting any related issues.

## Installation
To install the Thor GO SDK, run the following command:

```bash
go get github.com/darrenvechain/thor-go-sdk
```

## Examples

```golang
package main

import (
	"fmt"
	"github.com/darrenvechain/thor-go-sdk/thorgo"
	"github.com/darrenvechain/thor-go-sdk/solo"
	"github.com/ethereum/go-ethereum/common"
)

func main() {
	// Create a new client
	thor, err := thorgo.FromURL(solo.URL)

	// Get an accounts balance
	acc, err := thor.Account(common.HexToAddress("0x0000000000000000000000000000456e6570")).Get()
	fmt.Println(acc.Balance)
}
```

The following test files demonstrate how to interact with the VechainThor blockchain using the Thor GO SDK:

- [Accounts](./examples/internal/accounts_test.go) - Fetch account details, code and storage.
- [Blocks](./examples/internal/blocks_test.go) - Fetch blocks by revisions.
- [Contracts](./examples/internal/contracts_test.go) - Read state and create transaction clauses.
- [Events](./examples/internal/log_events_test.go) - Query on chain events.
- [Transfers](./examples/internal/log_transfers_test.go) - Query on chain VET transfers.
- [Transactions](./examples/internal/transactions_test.go) - Send transactions and query transaction details.
- [Signers](./examples/internal/private_key_signer_test.go) - Simplify transaction building, signing, sending and retrieving.
