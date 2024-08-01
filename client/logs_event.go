package client

import (
	"github.com/ethereum/go-ethereum/common"
)

type EventLog struct {
	Address *common.Address `json:"address,omitempty"`
	Topics  []common.Hash   `json:"topics"`
	Data    string          `json:"data"`
	Meta    LogMeta         `json:"meta"`
}

type EventFilter struct {
	Range    *FilterRange     `json:"range,omitempty"`
	Options  *FilterOptions   `json:"options,omitempty"`
	Criteria *[]EventCriteria `json:"criteriaSet,omitempty"`
	Order    *string          `json:"order,omitempty"`
}

type EventCriteria struct {
	Address *common.Address `json:"address,omitempty"`
	Topics  *[]common.Hash  `json:"topics,omitempty"`
}
