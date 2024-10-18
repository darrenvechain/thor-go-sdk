package client

import (
	"github.com/darrenvechain/thorgo/crypto/tx"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

type SendTransactionResponse struct {
	ID common.Hash `json:"id"`
}

type RawTransaction struct {
	Raw  string `json:"raw"`
	Meta TxMeta `json:"meta"`
}

type TransactionReceipt struct {
	GasUsed  int64          `json:"gasUsed"`
	GasPayer common.Address `json:"gasPayer"`
	Paid     *hexutil.Big   `json:"paid"`
	Reward   *hexutil.Big   `json:"reward"`
	Reverted bool           `json:"reverted"`
	Meta     ReceiptMeta    `json:"meta"`
	Outputs  []Output       `json:"outputs"`
}

type Transaction struct {
	ID           common.Hash     `json:"id"`
	ChainTag     int64           `json:"chainTag"`
	BlockRef     tx.BlockRef     `json:"blockRef"`
	Expiration   int64           `json:"expiration"`
	Clauses      []tx.Clause     `json:"clauses"`
	GasPriceCoef int64           `json:"gasPriceCoef"`
	Gas          int64           `json:"gas"`
	Origin       common.Address  `json:"origin"`
	Delegator    *common.Address `json:"delegator"`
	Nonce        hexutil.Big     `json:"nonce"`
	DependsOn    *common.Hash    `json:"dependsOn"`
	Size         int64           `json:"size"`
	Meta         TxMeta          `json:"meta"`
}

type Transfer struct {
	Sender    common.Address `json:"sender"`
	Recipient common.Address `json:"recipient"`
	Amount    *hexutil.Big   `json:"amount"`
}

type Output struct {
	ContractAddress string     `json:"contractAddress"`
	Events          []Event    `json:"events"`
	Transfers       []Transfer `json:"transfers"`
}

type Event struct {
	Address common.Address `json:"address"`
	Topics  []common.Hash  `json:"topics"`
	Data    string         `json:"data"`
}

type ReceiptMeta struct {
	BlockID        common.Hash    `json:"blockID"`
	BlockNumber    int64          `json:"blockNumber"`
	BlockTimestamp int64          `json:"blockTimestamp"`
	TxID           common.Hash    `json:"txID"`
	TxOrigin       common.Address `json:"txOrigin"`
}

type TxMeta struct {
	BlockID        common.Hash `json:"blockID"`
	BlockNumber    int64       `json:"blockNumber"`
	BlockTimestamp int64       `json:"blockTimestamp"`
}
