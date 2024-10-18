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
go get github.com/darrenvechain/thorgo
``` 

## Packages

### thorgo 

- `github.com/darrenvechain/thorgo`
- `thorgo` is the primary package in the Thor GO SDK. It provides a high-level interface for interacting with the VechainThor blockchain. This package includes functions for querying account balances, transactions, blocks, and smart contracts. It also supports simulating, building, and sending transactions, as well as interacting with smart contracts for reading and transacting.

### client

- `github.com/darrenvechain/thorgo/client`
- The `client` package provides raw API access to the VechainThor blockchain. It allows developers to query the blockchain directly without the need for higher-level abstractions provided by `thorgo`.

### txmanager

- `github.com/darrenvechain/thorgo/txmanager`
- The `txmanager` package provides a way to sign, send, and delegate transactions.
- The delegation managers can be used to easily delegate transaction gas fees.
- **Note**: The private key implementations in this package are not secure. It is recommended to use a secure key management solution in a production environment.
- To create your own transaction manager or signer, you can implement the `accounts.TxManager` for contract interaction and the `transactions.Signer` to sign transactions.
    
```golang
// github.com/darrenvechain/thorgo/accounts
type TxManager interface {
    SendClauses(clauses []*transaction.Clause) (common.Hash, error)
}
```

```golang
// github.com/darrenvechain/thorgo/transactions
type Signer interface {
  SignTransaction(tx *transaction.Transaction) ([]byte, error)
  Address() common.Address
}
```

### tx

- `github.com/darrenvechain/thorgo/crypto/tx`
- The `tx` package is a copy of the [vechain/thor/tx](https://github.com/vechain/thor/tree/master/tx) package and can be used to build transactions where `thorgo` does not provide the necessary functionality.

### solo

- `github.com/darrenvechain/thorgo/solo`
- The `solo` package provides quick access to Thor solo values for testing and development purposes.

### certificate

- `github.com/darrenvechain/thorgo/crypto/certificate`
- The `certificate` package provides a way to encode, sign, and verify certificates in accordance with [VIP-192](https://github.com/vechain/VIPs/blob/master/vips/VIP-192.md)

### hdwallet

- `github.com/darrenvechain/thorgo/crypto/hdwallet`
- The `hdwallet` package provides a way to generate HD wallets and derive keys from them.

## Examples

### 1: Creating a New Client

```golang
package main

import (
	"fmt"
    
	"github.com/darrenvechain/thorgo"
	"github.com/darrenvechain/thorgo/solo"
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

### 2: Interacting with a contract + Delegated Transaction


<details>
  <summary>Expand</summary>

```golang
package main

import (
    "log/slog"
    "math/big"
    "strings"

    "github.com/darrenvechain/thorgo/solo"
    "github.com/darrenvechain/thorgo"
    "github.com/darrenvechain/thorgo/txmanager"
    "github.com/ethereum/go-ethereum/accounts/abi"
    "github.com/ethereum/go-ethereum/common"
)

func main() {
    thor, _ := thorgo.FromURL("http://localhost:8669")

    // Load a contract
    contractABI, _ := abi.JSON(strings.NewReader(vthoABI))
    vtho := thor.Account(common.HexToAddress("0x0000000000000000000000000000456e65726779")).Contract(&contractABI)

    // Create a delegated transaction manager
    origin := txmanager.FromPK(solo.Keys()[0], thor)
    gasPayer := txmanager.NewDelegator(solo.Keys()[1])
    txSender := txmanager.NewDelegatedManager(thor, origin, gasPayer)

    // Create a new account to receive the tokens
    recipient, _ := txmanager.GeneratePK(thor)
    recipientBalance := new(big.Int)

    // Call the balanceOf function
    err := vtho.Call("balanceOf", &recipientBalance, recipient.Address())
    slog.Info("recipient balance before", "balance", recipientBalance, "error", err)

    // Send 1000 tokens to the recipient
    tx, _ := vtho.Send(txSender, "transfer", recipient.Address(), big.NewInt(1000))
    receipt, _ := tx.Wait()
    slog.Info("receipt", "txID", receipt.Meta.TxID, "reverted", receipt.Reverted)

    // Call the balanceOf function again
    err = vtho.Call("balanceOf", &recipientBalance, recipient.Address())
    slog.Info("recipient balance after", "balance", recipientBalance, "error", err)
}

var (
    vthoABI = `[
        {
            "constant": true,
            "inputs": [{"internalType": "address", "name": "account", "type": "address"}],
            "name": "balanceOf",
            "outputs": [{"internalType": "uint256", "name": "", "type": "uint256"}],
            "stateMutability": "view",
            "type": "function"
        },
        {
            "constant": false,
            "inputs": [
                {"internalType": "address", "name": "recipient", "type": "address"},
                {"internalType": "uint256", "name": "amount", "type": "uint256"}
            ],
            "name": "transfer",
            "outputs": [{"internalType": "bool", "name": "", "type": "bool"}],
            "stateMutability": "nonpayable",
            "type": "function"
        }
    ]`
)
```


</details>
    


### 3: Multi Clause Transaction

<details>
  <summary>Expand</summary>

```golang
package main

import (
	"log/slog"
	"math/big"
	"strings"

	"github.com/darrenvechain/thorgo/crypto/tx"
	"github.com/darrenvechain/thorgo/solo"
	"github.com/darrenvechain/thorgo"
	"github.com/darrenvechain/thorgo/txmanager"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

func main() {
	thor, _ := thorgo.FromURL("http://localhost:8669")

	// Load a contract
	contractABI, _ := abi.JSON(strings.NewReader(vthoABI))
	vtho := thor.Account(common.HexToAddress("0x0000000000000000000000000000456e65726779")).Contract(&contractABI)

	origin := txmanager.FromPK(solo.Keys()[0], thor)

	// clause1
	clause1, _ := vtho.AsClause("transfer", common.HexToAddress("0x87AA2B76f29583E4A9095DBb6029A9C41994E25B"), big.NewInt(1000000))
	clause2, _ := vtho.AsClause("transfer", common.HexToAddress("0xdf1b32ec78c1f338F584a2a459f01fD70529dDBF"), big.NewInt(1000000))

	// Option 1 - Directly using the txmanager.Manager
	tx, _ := origin.SendClauses([]*transaction.Clause{clause1, clause2})
	receipt, _ := thor.Transaction(tx).Wait()
	slog.Info("transaction receipt 1", "id", receipt.Meta.TxID, "reverted", receipt.Reverted)

	// Option 2 - Using the transaction builder with txmanager.Signer
	tx2, _ := thor.Transactor([]*transaction.Clause{clause1, clause2}, origin.Address()).
		GasPriceCoef(255).
		Send(origin)
	receipt2, _ := tx2.Wait()
	slog.Info("transaction receipt 2", "id", receipt2.Meta.TxID, "reverted", receipt2.Reverted)
}
```

</details>

