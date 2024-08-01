package client

import (
	"github.com/darrenvechain/thor-go-sdk/hex"
	"github.com/darrenvechain/thor-go-sdk/transaction"
	"github.com/ethereum/go-ethereum/common"
)

type Account struct {
	Balance hex.Int `json:"balance"`
	Energy  hex.Int `json:"energy"`
	HasCode bool    `json:"hasCode"`
}

type AccountCode struct {
	Code string `json:"code"`
}

type AccountStorage struct {
	Value string `json:"value"`
}

type InspectRequest struct {
	Gas        *uint64               `json:"gas,omitempty"`
	GasPrice   *uint64               `json:"gasPrice,omitempty"`
	Caller     *common.Address       `json:"caller,omitempty"`
	ProvedWork *string               `json:"provedWork,omitempty"`
	GasPayer   *common.Address       `json:"gasPayer,omitempty"`
	Expiration *uint64               `json:"expiration,omitempty"`
	BlockRef   *string               `json:"blockRef,omitempty"`
	Clauses    []*transaction.Clause `json:"clauses"`
}

type InspectResponse struct {
	Data      string     `json:"data"`
	Events    []Event    `json:"events"`
	Transfers []Transfer `json:"transfers"`
	GasUsed   uint64     `json:"gasUsed"`
	Reverted  bool       `json:"reverted"`
	VmError   string     `json:"vmError"`
}
