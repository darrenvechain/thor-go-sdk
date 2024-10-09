package txmanager

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/darrenvechain/thor-go-sdk/crypto/transaction"
	"github.com/darrenvechain/thor-go-sdk/thorgo"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

// DelegatedManager is a transaction manager that delegates the payment of transaction fees to a Delegator
type DelegatedManager struct {
	thor     *thorgo.Thor
	gasPayer Delegator
	origin   Signer
}

func NewDelegatedManager(thor *thorgo.Thor, origin Signer, gasPayer Delegator) *DelegatedManager {
	return &DelegatedManager{
		thor:     thor,
		origin:   origin,
		gasPayer: gasPayer,
	}
}

func (d *DelegatedManager) SendClauses(clauses []*transaction.Clause) (common.Hash, error) {
	tx, err := d.thor.Transactor(clauses, d.Address()).Delegate().Build()
	if err != nil {
		return common.Hash{}, fmt.Errorf("failed to build transaction: %w", err)
	}
	delegatorSig, err := d.gasPayer.Delegate(tx, d.Address())
	if err != nil {
		return common.Hash{}, fmt.Errorf("failed to delegate: %w", err)
	}
	signature, err := d.origin.SignTransaction(tx)
	if err != nil {
		return common.Hash{}, fmt.Errorf("failed to sign transaction: %w", err)
	}

	signature = append(signature, delegatorSig...)
	tx = tx.WithSignature(signature)
	res, err := d.thor.Client.SendTransaction(tx)
	if err != nil {
		return common.Hash{}, fmt.Errorf("failed to send transaction: %w", err)
	}
	return res.ID, nil
}

// Address returns the address of the origin manager
func (d *DelegatedManager) Address() common.Address {
	return d.origin.Address()
}

// PKDelegator is a delegator that uses a private key to pay for transaction fees
type PKDelegator struct {
	key *ecdsa.PrivateKey
}

func NewDelegator(key *ecdsa.PrivateKey) *PKDelegator {
	return &PKDelegator{key: key}
}

func (p *PKDelegator) PublicKey() *ecdsa.PublicKey {
	return &p.key.PublicKey
}

func (p *PKDelegator) Address() (addr common.Address) {
	return crypto.PubkeyToAddress(p.key.PublicKey)
}

func (p *PKDelegator) Delegate(tx *transaction.Transaction, origin common.Address) ([]byte, error) {
	return crypto.Sign(tx.DelegatorSigningHash(origin).Bytes(), p.key)
}

// URLDelegator is a delegator that uses a remote URL to pay for transaction fees
type URLDelegator struct {
	url string
}

func NewUrlDelegator(url string) *URLDelegator {
	return &URLDelegator{url: url}
}

func (p *URLDelegator) Delegate(tx *transaction.Transaction, origin common.Address) ([]byte, error) {
	encoded, err := tx.Encoded()
	if err != nil {
		return nil, err
	}

	req := &DelegateRequest{
		Origin: origin.String(),
		Raw:    encoded,
	}

	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	res, err := http.Post(p.url, "application/json", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, errors.New("200 OK expected")
	}

	defer res.Body.Close()
	var response DelegateResponse
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, err
	}

	return common.FromHex(response.Signature), nil
}
