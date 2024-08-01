package client

import (
	"github.com/ethereum/go-ethereum/common"
)

type FilterRange struct {
	Unit *string `json:"unit,omitempty"`
	From *uint64 `json:"from,omitempty"`
	To   *uint64 `json:"to,omitempty"`
}

type FilterOptions struct {
	Offset *uint64 `json:"offset,omitempty"`
	Limit  *uint64 `json:"limit,omitempty"`
}

type LogMeta struct {
	BlockID     common.Hash    `json:"blockID"`
	BlockNumber uint64         `json:"blockNumber"`
	BlockTime   uint64         `json:"blockTimestamp"`
	TxID        common.Hash    `json:"txID"`
	TxOrigin    common.Address `json:"txOrigin"`
	ClauseIndex uint64         `json:"clauseIndex"`
}
