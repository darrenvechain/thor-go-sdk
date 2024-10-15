package client

import (
	"github.com/ethereum/go-ethereum/common"
)

type FilterRange struct {
	Unit *string `json:"unit,omitempty"`
	From *int64  `json:"from,omitempty"`
	To   *int64  `json:"to,omitempty"`
}

type FilterOptions struct {
	Offset *int64 `json:"offset,omitempty"`
	Limit  *int64 `json:"limit,omitempty"`
}

type LogMeta struct {
	BlockID     common.Hash    `json:"blockID"`
	BlockNumber int64          `json:"blockNumber"`
	BlockTime   int64          `json:"blockTimestamp"`
	TxID        common.Hash    `json:"txID"`
	TxOrigin    common.Address `json:"txOrigin"`
	ClauseIndex int64          `json:"clauseIndex"`
}
