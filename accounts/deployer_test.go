package accounts_test

import (
	"strings"
	"testing"

	"github.com/darrenvechain/thorgo/accounts"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
)

func TestDeployer_Deploy(t *testing.T) {
	abi, err := abi.JSON(strings.NewReader(contractABI))
	assert.NoError(t, err)

	deployedName := "MyERC20"
	deployer := accounts.NewDeployer(thorClient, common.Hex2Bytes(erc20Bytecode), &abi)
	erc20, txID, err := deployer.Deploy(account1, deployedName, "ERC20")
	assert.NoError(t, err)
	assert.NotEqual(t, common.Hash{}, txID)

	var name string
	err = erc20.Call("name", &name)
	assert.NoError(t, err)
	assert.Equal(t, deployedName, name)
}

var erc20Bytecode = "60806040523480156200001157600080fd5b50604051620014c9380380620014c98339818101604052810190620000379190620001fa565b818181600390816200004a9190620004ca565b5080600490816200005c9190620004ca565b5050505050620005b1565b6000604051905090565b600080fd5b600080fd5b600080fd5b600080fd5b6000601f19601f8301169050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b620000d08262000085565b810181811067ffffffffffffffff82111715620000f257620000f162000096565b5b80604052505050565b60006200010762000067565b9050620001158282620000c5565b919050565b600067ffffffffffffffff82111562000138576200013762000096565b5b620001438262000085565b9050602081019050919050565b60005b838110156200017057808201518184015260208101905062000153565b60008484015250505050565b6000620001936200018d846200011a565b620000fb565b905082815260208101848484011115620001b257620001b162000080565b5b620001bf84828562000150565b509392505050565b600082601f830112620001df57620001de6200007b565b5b8151620001f18482602086016200017c565b91505092915050565b6000806040838503121562000214576200021362000071565b5b600083015167ffffffffffffffff81111562000235576200023462000076565b5b6200024385828601620001c7565b925050602083015167ffffffffffffffff81111562000267576200026662000076565b5b6200027585828601620001c7565b9150509250929050565b600081519050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b60006002820490506001821680620002d257607f821691505b602082108103620002e857620002e76200028a565b5b50919050565b60008190508160005260206000209050919050565b60006020601f8301049050919050565b600082821b905092915050565b600060088302620003527fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8262000313565b6200035e868362000313565b95508019841693508086168417925050509392505050565b6000819050919050565b6000819050919050565b6000620003ab620003a56200039f8462000376565b62000380565b62000376565b9050919050565b6000819050919050565b620003c7836200038a565b620003df620003d682620003b2565b84845462000320565b825550505050565b600090565b620003f6620003e7565b62000403818484620003bc565b505050565b5b818110156200042b576200041f600082620003ec565b60018101905062000409565b5050565b601f8211156200047a576200044481620002ee565b6200044f8462000303565b810160208510156200045f578190505b620004776200046e8562000303565b83018262000408565b50505b505050565b600082821c905092915050565b60006200049f600019846008026200047f565b1980831691505092915050565b6000620004ba83836200048c565b9150826002028217905092915050565b620004d5826200027f565b67ffffffffffffffff811115620004f157620004f062000096565b5b620004fd8254620002b9565b6200050a8282856200042f565b600060209050601f8311600181146200054257600084156200052d578287015190505b620005398582620004ac565b865550620005a9565b601f1984166200055286620002ee565b60005b828110156200057c5784890151825560018201915060208501945060208101905062000555565b868310156200059c578489015162000598601f8916826200048c565b8355505b6001600288020188555050505b505050505050565b610f0880620005c16000396000f3fe608060405234801561001057600080fd5b506004361061009e5760003560e01c806340c10f191161006657806340c10f191461015d57806370a082311461017957806395d89b41146101a9578063a9059cbb146101c7578063dd62ed3e146101f75761009e565b806306fdde03146100a3578063095ea7b3146100c157806318160ddd146100f157806323b872dd1461010f578063313ce5671461013f575b600080fd5b6100ab610227565b6040516100b89190610b5c565b60405180910390f35b6100db60048036038101906100d69190610c17565b6102b9565b6040516100e89190610c72565b60405180910390f35b6100f96102dc565b6040516101069190610c9c565b60405180910390f35b61012960048036038101906101249190610cb7565b6102e6565b6040516101369190610c72565b60405180910390f35b610147610315565b6040516101549190610d26565b60405180910390f35b61017760048036038101906101729190610c17565b61031a565b005b610193600480360381019061018e9190610d41565b610328565b6040516101a09190610c9c565b60405180910390f35b6101b1610370565b6040516101be9190610b5c565b60405180910390f35b6101e160048036038101906101dc9190610c17565b610402565b6040516101ee9190610c72565b60405180910390f35b610211600480360381019061020c9190610d6e565b610425565b60405161021e9190610c9c565b60405180910390f35b60606003805461023690610ddd565b80601f016020809104026020016040519081016040528092919081815260200182805461026290610ddd565b80156102af5780601f10610284576101008083540402835291602001916102af565b820191906000526020600020905b81548152906001019060200180831161029257829003601f168201915b5050505050905090565b6000806102c46104ac565b90506102d18185856104b4565b600191505092915050565b6000600254905090565b6000806102f16104ac565b90506102fe8582856104c6565b61030985858561055a565b60019150509392505050565b600090565b610324828261064e565b5050565b60008060008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020549050919050565b60606004805461037f90610ddd565b80601f01602080910402602001604051908101604052809291908181526020018280546103ab90610ddd565b80156103f85780601f106103cd576101008083540402835291602001916103f8565b820191906000526020600020905b8154815290600101906020018083116103db57829003601f168201915b5050505050905090565b60008061040d6104ac565b905061041a81858561055a565b600191505092915050565b6000600160008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054905092915050565b600033905090565b6104c183838360016106d0565b505050565b60006104d28484610425565b90507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff81146105545781811015610544578281836040517ffb8f41b200000000000000000000000000000000000000000000000000000000815260040161053b93929190610e1d565b60405180910390fd5b610553848484840360006106d0565b5b50505050565b600073ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff16036105cc5760006040517f96c6fd1e0000000000000000000000000000000000000000000000000000000081526004016105c39190610e54565b60405180910390fd5b600073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff160361063e5760006040517fec442f050000000000000000000000000000000000000000000000000000000081526004016106359190610e54565b60405180910390fd5b6106498383836108a7565b505050565b600073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff16036106c05760006040517fec442f050000000000000000000000000000000000000000000000000000000081526004016106b79190610e54565b60405180910390fd5b6106cc600083836108a7565b5050565b600073ffffffffffffffffffffffffffffffffffffffff168473ffffffffffffffffffffffffffffffffffffffff16036107425760006040517fe602df050000000000000000000000000000000000000000000000000000000081526004016107399190610e54565b60405180910390fd5b600073ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff16036107b45760006040517f94280d620000000000000000000000000000000000000000000000000000000081526004016107ab9190610e54565b60405180910390fd5b81600160008673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000208190555080156108a1578273ffffffffffffffffffffffffffffffffffffffff168473ffffffffffffffffffffffffffffffffffffffff167f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925846040516108989190610c9c565b60405180910390a35b50505050565b600073ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff16036108f95780600260008282546108ed9190610e9e565b925050819055506109cc565b60008060008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054905081811015610985578381836040517fe450d38c00000000000000000000000000000000000000000000000000000000815260040161097c93929190610e1d565b60405180910390fd5b8181036000808673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002081905550505b600073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff1603610a155780600260008282540392505081905550610a62565b806000808473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020600082825401925050819055505b8173ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef83604051610abf9190610c9c565b60405180910390a3505050565b600081519050919050565b600082825260208201905092915050565b60005b83811015610b06578082015181840152602081019050610aeb565b60008484015250505050565b6000601f19601f8301169050919050565b6000610b2e82610acc565b610b388185610ad7565b9350610b48818560208601610ae8565b610b5181610b12565b840191505092915050565b60006020820190508181036000830152610b768184610b23565b905092915050565b600080fd5b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b6000610bae82610b83565b9050919050565b610bbe81610ba3565b8114610bc957600080fd5b50565b600081359050610bdb81610bb5565b92915050565b6000819050919050565b610bf481610be1565b8114610bff57600080fd5b50565b600081359050610c1181610beb565b92915050565b60008060408385031215610c2e57610c2d610b7e565b5b6000610c3c85828601610bcc565b9250506020610c4d85828601610c02565b9150509250929050565b60008115159050919050565b610c6c81610c57565b82525050565b6000602082019050610c876000830184610c63565b92915050565b610c9681610be1565b82525050565b6000602082019050610cb16000830184610c8d565b92915050565b600080600060608486031215610cd057610ccf610b7e565b5b6000610cde86828701610bcc565b9350506020610cef86828701610bcc565b9250506040610d0086828701610c02565b9150509250925092565b600060ff82169050919050565b610d2081610d0a565b82525050565b6000602082019050610d3b6000830184610d17565b92915050565b600060208284031215610d5757610d56610b7e565b5b6000610d6584828501610bcc565b91505092915050565b60008060408385031215610d8557610d84610b7e565b5b6000610d9385828601610bcc565b9250506020610da485828601610bcc565b9150509250929050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b60006002820490506001821680610df557607f821691505b602082108103610e0857610e07610dae565b5b50919050565b610e1781610ba3565b82525050565b6000606082019050610e326000830186610e0e565b610e3f6020830185610c8d565b610e4c6040830184610c8d565b949350505050565b6000602082019050610e696000830184610e0e565b92915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b6000610ea982610be1565b9150610eb483610be1565b9250828201905080821115610ecc57610ecb610e6f565b5b9291505056fea2646970667358221220e38c2ea7a55d79f2695d7b57320f013a28b9dc41e8b492ba111ddb3eeefc626064736f6c63430008140033"

var contractABI = ` [
    {
      "inputs": [
        {
          "internalType": "string",
          "name": "name_",
          "type": "string"
        },
        {
          "internalType": "string",
          "name": "symbol_",
          "type": "string"
        }
      ],
      "stateMutability": "nonpayable",
      "type": "constructor"
    },
    {
      "inputs": [
        {
          "internalType": "address",
          "name": "spender",
          "type": "address"
        },
        {
          "internalType": "uint256",
          "name": "allowance",
          "type": "uint256"
        },
        {
          "internalType": "uint256",
          "name": "needed",
          "type": "uint256"
        }
      ],
      "name": "ERC20InsufficientAllowance",
      "type": "error"
    },
    {
      "inputs": [
        {
          "internalType": "address",
          "name": "sender",
          "type": "address"
        },
        {
          "internalType": "uint256",
          "name": "balance",
          "type": "uint256"
        },
        {
          "internalType": "uint256",
          "name": "needed",
          "type": "uint256"
        }
      ],
      "name": "ERC20InsufficientBalance",
      "type": "error"
    },
    {
      "inputs": [
        {
          "internalType": "address",
          "name": "approver",
          "type": "address"
        }
      ],
      "name": "ERC20InvalidApprover",
      "type": "error"
    },
    {
      "inputs": [
        {
          "internalType": "address",
          "name": "receiver",
          "type": "address"
        }
      ],
      "name": "ERC20InvalidReceiver",
      "type": "error"
    },
    {
      "inputs": [
        {
          "internalType": "address",
          "name": "sender",
          "type": "address"
        }
      ],
      "name": "ERC20InvalidSender",
      "type": "error"
    },
    {
      "inputs": [
        {
          "internalType": "address",
          "name": "spender",
          "type": "address"
        }
      ],
      "name": "ERC20InvalidSpender",
      "type": "error"
    },
    {
      "anonymous": false,
      "inputs": [
        {
          "indexed": true,
          "internalType": "address",
          "name": "owner",
          "type": "address"
        },
        {
          "indexed": true,
          "internalType": "address",
          "name": "spender",
          "type": "address"
        },
        {
          "indexed": false,
          "internalType": "uint256",
          "name": "value",
          "type": "uint256"
        }
      ],
      "name": "Approval",
      "type": "event"
    },
    {
      "anonymous": false,
      "inputs": [
        {
          "indexed": true,
          "internalType": "address",
          "name": "from",
          "type": "address"
        },
        {
          "indexed": true,
          "internalType": "address",
          "name": "to",
          "type": "address"
        },
        {
          "indexed": false,
          "internalType": "uint256",
          "name": "value",
          "type": "uint256"
        }
      ],
      "name": "Transfer",
      "type": "event"
    },
    {
      "inputs": [
        {
          "internalType": "address",
          "name": "owner",
          "type": "address"
        },
        {
          "internalType": "address",
          "name": "spender",
          "type": "address"
        }
      ],
      "name": "allowance",
      "outputs": [
        {
          "internalType": "uint256",
          "name": "",
          "type": "uint256"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "address",
          "name": "spender",
          "type": "address"
        },
        {
          "internalType": "uint256",
          "name": "value",
          "type": "uint256"
        }
      ],
      "name": "approve",
      "outputs": [
        {
          "internalType": "bool",
          "name": "",
          "type": "bool"
        }
      ],
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "address",
          "name": "account",
          "type": "address"
        }
      ],
      "name": "balanceOf",
      "outputs": [
        {
          "internalType": "uint256",
          "name": "",
          "type": "uint256"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [],
      "name": "decimals",
      "outputs": [
        {
          "internalType": "uint8",
          "name": "",
          "type": "uint8"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "address",
          "name": "to",
          "type": "address"
        },
        {
          "internalType": "uint256",
          "name": "amount",
          "type": "uint256"
        }
      ],
      "name": "mint",
      "outputs": [],
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "inputs": [],
      "name": "name",
      "outputs": [
        {
          "internalType": "string",
          "name": "",
          "type": "string"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [],
      "name": "symbol",
      "outputs": [
        {
          "internalType": "string",
          "name": "",
          "type": "string"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [],
      "name": "totalSupply",
      "outputs": [
        {
          "internalType": "uint256",
          "name": "",
          "type": "uint256"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "address",
          "name": "to",
          "type": "address"
        },
        {
          "internalType": "uint256",
          "name": "value",
          "type": "uint256"
        }
      ],
      "name": "transfer",
      "outputs": [
        {
          "internalType": "bool",
          "name": "",
          "type": "bool"
        }
      ],
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "address",
          "name": "from",
          "type": "address"
        },
        {
          "internalType": "address",
          "name": "to",
          "type": "address"
        },
        {
          "internalType": "uint256",
          "name": "value",
          "type": "uint256"
        }
      ],
      "name": "transferFrom",
      "outputs": [
        {
          "internalType": "bool",
          "name": "",
          "type": "bool"
        }
      ],
      "stateMutability": "nonpayable",
      "type": "function"
    }
  ]`
