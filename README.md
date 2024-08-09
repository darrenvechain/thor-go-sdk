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
