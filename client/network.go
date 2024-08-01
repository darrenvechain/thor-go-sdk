package client

import "github.com/ethereum/go-ethereum/common"

type Peer struct {
	Name        string      `json:"name"`
	BestBlockID common.Hash `json:"bestBlockID"`
	TotalScore  uint64      `json:"totalScore"`
	PeerID      string      `json:"peerID"`
	NetAddr     string      `json:"netAddr"`
	Inbound     bool        `json:"inbound"`
	Duration    uint64      `json:"duration"`
}
