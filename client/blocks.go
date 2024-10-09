package client

import (
	"github.com/darrenvechain/thor-go-sdk/crypto/transaction"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

type Block struct {
	Number       uint64         `json:"number"`
	ID           common.Hash    `json:"id"`
	Size         uint64         `json:"size"`
	ParentID     common.Hash    `json:"parentID"`
	Timestamp    uint64         `json:"timestamp"`
	GasLimit     uint64         `json:"gasLimit"`
	Beneficiary  common.Address `json:"beneficiary"`
	GasUsed      uint64         `json:"gasUsed"`
	TotalScore   uint64         `json:"totalScore"`
	TxsRoot      common.Hash    `json:"txsRoot"`
	TxsFeatures  uint64         `json:"txsFeatures"`
	StateRoot    common.Hash    `json:"stateRoot"`
	ReceiptsRoot common.Hash    `json:"receiptsRoot"`
	Com          bool           `json:"com"`
	Signer       common.Address `json:"signer"`
	IsTrunk      bool           `json:"isTrunk"`
	IsFinalized  bool           `json:"isFinalized"`
	Transactions []common.Hash  `json:"transactions"`
}

func (b *Block) ChainTag() byte {
	return b.ID[len(b.ID)-1]
}

func (b *Block) BlockRef() transaction.BlockRef {
	return transaction.NewBlockRefFromID(b.ID)
}

type BlockTransaction struct {
	ID           common.Hash          `json:"id"`
	ChainTag     byte                 `json:"chainTag"`
	BlockRef     transaction.BlockRef `json:"blockRef"`
	Expiration   uint64               `json:"expiration"`
	Clauses      []transaction.Clause `json:"clauses"`
	GasPriceCoef uint64               `json:"gasPriceCoef"`
	Gas          uint64               `json:"gas"`
	Origin       common.Address       `json:"origin"`
	Delegator    *common.Address      `json:"delegator,omitempty"`
	Nonce        hexutil.Big          `json:"nonce"`
	DependsOn    *common.Hash         `json:"dependsOn,omitempty"`
	Size         uint64               `json:"size"`
	GasUsed      uint64               `json:"gasUsed"`
	GasPayer     common.Address       `json:"gasPayer"`
	Paid         hexutil.Big          `json:"paid"`
	Reward       hexutil.Big          `json:"reward"`
	Reverted     bool                 `json:"reverted"`
	Outputs      []Output             `json:"outputs"`
}

type ExpandedBlock struct {
	Block
	Transactions []BlockTransaction `json:"transactions"`
}
