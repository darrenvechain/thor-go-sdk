package certificate

import (
	"encoding/json"
	"fmt"

	"github.com/darrenvechain/thor-go-sdk/crypto/hash"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

type Certificate struct {
	Domain    string  `json:"domain"`
	Payload   Payload `json:"payload"`
	Purpose   string  `json:"purpose"`
	Signer    string  `json:"signer"`
	Timestamp uint64  `json:"timestamp"`
}

type Payload struct {
	Content string `json:"content"`
	Type    string `json:"type"`
}

func (c *Certificate) Encode() ([]byte, error) {
	return json.Marshal(c)
}

func (c *Certificate) SigningHash() (common.Hash, error) {
	encoded, err := c.Encode()
	if err != nil {
		return common.Hash{}, fmt.Errorf("failed to encode certificate: %w", err)
	}

	return hash.Blake2b(encoded), nil
}

func (c *Certificate) Verify(signature []byte) bool {
	signingHash, err := c.SigningHash()
	if err != nil {
		return false
	}
	pubkey, err := crypto.SigToPub(signingHash.Bytes(), signature)
	if err != nil {
		return false
	}
	return crypto.PubkeyToAddress(*pubkey) == common.HexToAddress(c.Signer)
}
