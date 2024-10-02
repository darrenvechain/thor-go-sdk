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

## Usage

### Example 1: Creating a New Client

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

### Example 2: Interacting with a contract + Delegated Transaction


<details>
  <summary>Expand</summary>

```golang
package main

import (
    "log/slog"
    "math/big"
    "strings"

    "github.com/darrenvechain/thor-go-sdk/solo"
    "github.com/darrenvechain/thor-go-sdk/thorgo"
    "github.com/darrenvechain/thor-go-sdk/txmanager"
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
    


### Example 3: Multi Clause Transaction

<details>
  <summary>Expand</summary>

```golang
package main

import (
    "github.com/darrenvechain/thor-go-sdk/solo"
    "github.com/darrenvechain/thor-go-sdk/thorgo"
    "github.com/darrenvechain/thor-go-sdk/transaction"
    "github.com/darrenvechain/thor-go-sdk/txmanager"
    "github.com/ethereum/go-ethereum/accounts/abi"
    "github.com/ethereum/go-ethereum/common"
    "log/slog"
    "math/big"
    "strings"
)

func main() {
    thor, _ := thorgo.FromURL("http://localhost:8669")

    // Load a contract
    contractABI, _ := abi.JSON(strings.NewReader(vthoABI))
    vtho := thor.Account(common.HexToAddress("0x0000000000000000000000000000456e65726779")).Contract(&contractABI)

    origin := txmanager.FromPK(solo.Keys()[0], thor)

    // clause1
    clause1, _ := vtho.AsClause("transfer", common.HexToAddress("0x87AA2B76f29583E4A9095DBb6029A9C41994E25B"), big.NewInt(1000000))
    clause2, _ := vtho.AsClause("transfer", common.HexToAddress("0x87AA2B76f29583E4A9095DBb6029A9C41994E25B"), big.NewInt(1000000))

    // Option 1 - Directly using the txmanager.Manager
    tx, _ := origin.SendClauses([]*transaction.Clause{clause1, clause2})
    receipt, _ := thor.Transaction(tx).Wait()
    slog.Info("transaction receipt 1", "id", receipt.Meta.TxID, "reverted", receipt.Reverted)

    // Option 2 - Using the transaction builder with txmanager.Signer
    tx2, _ := thor.TxBuilder([]*transaction.Clause{clause1, clause2}, origin.Address()).
        GasPriceCoef(255).
        Send(origin)
    receipt2, _ := tx2.Wait()
    slog.Info("transaction receipt 2", "id", receipt2.Meta.TxID, "reverted", receipt2.Reverted)
}
```

</details>

