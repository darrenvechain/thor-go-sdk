package client

import (
	"github.com/darrenvechain/thorgo/crypto/tx"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

type Block struct {
	Number       int64          `json:"number"`
	ID           common.Hash    `json:"id"`
	Size         int64          `json:"size"`
	ParentID     common.Hash    `json:"parentID"`
	Timestamp    int64          `json:"timestamp"`
	GasLimit     int64          `json:"gasLimit"`
	Beneficiary  common.Address `json:"beneficiary"`
	GasUsed      int64          `json:"gasUsed"`
	TotalScore   int64          `json:"totalScore"`
	TxsRoot      common.Hash    `json:"txsRoot"`
	TxsFeatures  int64          `json:"txsFeatures"`
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

func (b *Block) BlockRef() tx.BlockRef {
	return tx.NewBlockRefFromID(b.ID)
}

type BlockTransaction struct {
	ID           common.Hash     `json:"id"`
	ChainTag     byte            `json:"chainTag"`
	BlockRef     tx.BlockRef     `json:"blockRef"`
	Expiration   int64           `json:"expiration"`
	Clauses      []tx.Clause     `json:"clauses"`
	GasPriceCoef int64           `json:"gasPriceCoef"`
	Gas          int64           `json:"gas"`
	Origin       common.Address  `json:"origin"`
	Delegator    *common.Address `json:"delegator,omitempty"`
	Nonce        hexutil.Big     `json:"nonce"`
	DependsOn    *common.Hash    `json:"dependsOn,omitempty"`
	Size         int64           `json:"size"`
	GasUsed      int64           `json:"gasUsed"`
	GasPayer     common.Address  `json:"gasPayer"`
	Paid         hexutil.Big     `json:"paid"`
	Reward       hexutil.Big     `json:"reward"`
	Reverted     bool            `json:"reverted"`
	Outputs      []Output        `json:"outputs"`
}

type ExpandedBlock struct {
	Block
	Transactions []BlockTransaction `json:"transactions"`
}
