package txmanager_test

import (
	"crypto/ecdsa"
	"encoding/json"
	"math/big"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/darrenvechain/thorgo"
	"github.com/darrenvechain/thorgo/accounts"
	"github.com/darrenvechain/thorgo/builtins"
	"github.com/darrenvechain/thorgo/crypto/transaction"
	"github.com/darrenvechain/thorgo/solo"
	"github.com/darrenvechain/thorgo/transactions"
	"github.com/darrenvechain/thorgo/txmanager"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
)

var (
	// PKDelegator should implement Delegator
	_ txmanager.Delegator = &txmanager.PKDelegator{}
	// URLDelegator should implement Delegator
	_ txmanager.Delegator = &txmanager.URLDelegator{}
	// DelegatedManager should implement transactions.Signer
	_ transactions.Signer = &txmanager.DelegatedManager{}
	// DelegatedManager should implement accounts.TxManager
	_ accounts.TxManager = &txmanager.DelegatedManager{}
)

func TestPKDelegator(t *testing.T) {
	origin := txmanager.FromPK(solo.Keys()[0], thor)
	delegator := txmanager.NewDelegator(solo.Keys()[1])

	clause := transaction.NewClause(&common.Address{}).WithValue(new(big.Int))
	tx, err := thor.Transactor([]*transaction.Clause{clause}).
		GasPayer(delegator.Address()).
		Delegate().
		Build(origin.Address())

	assert.NoError(t, err)

	delegatorSignature, err := delegator.Delegate(tx, origin.Address())
	assert.NoError(t, err)

	signature, err := origin.SignTransaction(tx)
	assert.NoError(t, err)
	signature = append(signature, delegatorSignature...)

	signedTx := tx.WithSignature(signature)
	originAddr, err := signedTx.Origin()
	assert.NoError(t, err)
	assert.Equal(t, origin.Address(), originAddr)
	delegatorAddr, err := signedTx.Delegator()
	assert.NoError(t, err)
	assert.Equal(t, delegator.Address(), *delegatorAddr)
}

func createDelegationServer(key *ecdsa.PrivateKey) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req txmanager.DelegateRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		tx, err := transaction.Decode(common.Hex2Bytes(req.Raw))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		origin := common.HexToAddress(req.Origin)
		signingHash := tx.DelegatorSigningHash(origin)
		signature, err := crypto.Sign(signingHash.Bytes(), key)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		resp := txmanager.DelegateResponse{Signature: common.Bytes2Hex(signature)}
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}))
}

func TestNewUrlDelegator(t *testing.T) {
	origin := txmanager.FromPK(solo.Keys()[0], thor)
	server := createDelegationServer(solo.Keys()[1])

	delegator := txmanager.NewUrlDelegator(server.URL)

	clause := transaction.NewClause(&common.Address{}).WithValue(new(big.Int))
	tx, err := thor.Transactor([]*transaction.Clause{clause}).
		Delegate().
		Build(origin.Address())
	assert.NoError(t, err)

	delegatorSignature, err := delegator.Delegate(tx, origin.Address())
	assert.NoError(t, err)

	signature, err := origin.SignTransaction(tx)
	assert.NoError(t, err)

	signature = append(signature, delegatorSignature...)

	//combine the 2 signatures to make 1 of 130 bytes
	signed := tx.WithSignature(signature)

	//verify the signature
	delegatorAddr, err := signed.Delegator()
	assert.NoError(t, err)
	assert.Equal(t, *delegatorAddr, crypto.PubkeyToAddress(solo.Keys()[1].PublicKey))
	originAddr, err := signed.Origin()
	assert.NoError(t, err)
	assert.Equal(t, originAddr, origin.Address())
}

func TestNewDelegatedManager(t *testing.T) {
	thor, err := thorgo.FromURL("http://localhost:8669")
	assert.NoError(t, err)

	origin := txmanager.FromPK(solo.Keys()[0], thor)
	gasPayer := txmanager.NewDelegator(solo.Keys()[1])
	manager := txmanager.NewDelegatedManager(thor, origin, gasPayer)

	contract := builtins.VTHO.Load(thor)

	tx, err := contract.Send(manager, "transfer", common.Address{100}, big.NewInt(1000))
	assert.NoError(t, err)

	receipt, err := tx.Wait()
	assert.NoError(t, err)
	assert.False(t, receipt.Reverted)

	assert.Equal(t, gasPayer.Address().Hex(), receipt.GasPayer.Hex())
	assert.Equal(t, origin.Address().Hex(), receipt.Meta.TxOrigin.Hex())
}
