package client

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/darrenvechain/thor-go-sdk/crypto/transaction"
	"github.com/ethereum/go-ethereum/common"
)

type Client struct {
	client       *http.Client
	url          string
	genesisBlock *Block
}

func New(url string, client *http.Client) (*Client, error) {
	return newClient(url, client)
}

func FromURL(url string) (*Client, error) {
	return New(url, &http.Client{})
}

func newClient(url string, client *http.Client) (*Client, error) {
	url = strings.TrimSuffix(url, "/")

	c := &Client{
		client: client,
		url:    url,
	}

	block, err := c.Block("0")
	if err != nil {
		return nil, err
	}
	c.genesisBlock = block

	return c, nil
}

// Account fetches the account information for the given address.
func (c *Client) Account(addr common.Address) (*Account, error) {
	url := "/accounts/" + addr.Hex()
	return httpGet(c, url, &Account{})
}

// AccountAt fetches the account information for an address at the given revision.
func (c *Client) AccountAt(addr common.Address, revision common.Hash) (*Account, error) {
	url := "/accounts/" + addr.Hex() + "?revision=" + revision.Hex()
	return httpGet(c, url, &Account{})
}

// Inspect will send an array of clauses to the node to simulate the execution of the clauses.
// This can be used to:
// - Read contract(s) state
// - Simulate the execution of a transaction
func (c *Client) Inspect(body InspectRequest) ([]InspectResponse, error) {
	url := "/accounts/*"
	response := make([]InspectResponse, 0)
	_, err := httpPost(c, url, body, &response)
	if err != nil {
		return nil, err
	}
	return response, nil
}

// InspectAt will send an array of clauses to the node to simulate the execution of the clauses at the given revision.
func (c *Client) InspectAt(body InspectRequest, revision common.Hash) ([]InspectResponse, error) {
	url := "/accounts/*?revision=" + revision.Hex()
	response := make([]InspectResponse, 0)
	_, err := httpPost(c, url, body, &response)
	if err != nil {
		return nil, err
	}
	return response, nil
}

// AccountCode fetches the code for the account at the given address.
func (c *Client) AccountCode(addr common.Address) (*AccountCode, error) {
	url := "/accounts/" + addr.Hex() + "/code"
	return httpGet(c, url, &AccountCode{})
}

// AccountCodeAt fetches the code for the account at the given address and revision.
func (c *Client) AccountCodeAt(addr common.Address, revision common.Hash) (*AccountCode, error) {
	url := "/accounts/" + addr.Hex() + "/code?revision=" + revision.Hex()
	return httpGet(c, url, &AccountCode{})
}

// AccountStorage fetches the storage value for the account at the given address and key.
func (c *Client) AccountStorage(addr common.Address, key common.Hash) (*AccountStorage, error) {
	url := "/accounts/" + addr.Hex() + "/storage/" + key.Hex()
	return httpGet(c, url, &AccountStorage{})
}

// AccountStorageAt fetches the storage value for the account at the given address and key at the given revision.
func (c *Client) AccountStorageAt(
	addr common.Address,
	key common.Hash,
	revision common.Hash,
) (*AccountStorage, error) {
	url := "/accounts/" + addr.Hex() + "/storage/" + key.Hex() + "?revision=" + revision.Hex()
	return httpGet(c, url, &AccountStorage{})
}

// Block fetches the block for the given revision.
func (c *Client) Block(revision string) (*Block, error) {
	url := "/blocks/" + revision
	return httpGet(c, url, &Block{})
}

// BestBlock returns the best block.
func (c *Client) BestBlock() (*Block, error) {
	return c.Block("best")
}

// GenesisBlock returns the genesis block.
func (c *Client) GenesisBlock() *Block {
	return c.genesisBlock
}

// ExpandedBlock fetches the block at the given revision with all the transactions expanded.
func (c *Client) ExpandedBlock(revision string) (*ExpandedBlock, error) {
	url := "/blocks/" + revision + "?expanded=true"
	return httpGet(c, url, &ExpandedBlock{})
}

// ChainTag returns the chain tag of the genesis block.
func (c *Client) ChainTag() byte {
	return c.genesisBlock.ChainTag()
}

// SendTransaction sends a transaction to the node.
func (c *Client) SendTransaction(tx *transaction.Transaction) (*SendTransactionResponse, error) {
	body := make(map[string]string)
	encoded, err := tx.Encoded()
	if err != nil {
		return nil, err
	}
	body["raw"] = "0x" + encoded
	return httpPost(c, "/transactions", body, &SendTransactionResponse{})
}

// SendRawTransaction sends a raw transaction to the node.
func (c *Client) SendRawTransaction(raw string) (*SendTransactionResponse, error) {
	body := make(map[string]string)
	body["raw"] = raw
	return httpPost(c, "/transactions", body, &SendTransactionResponse{})
}

// Transaction fetches a transaction by its ID.
func (c *Client) Transaction(id common.Hash) (*Transaction, error) {
	url := "/transactions/" + id.Hex()
	return httpGet(c, url, &Transaction{})
}

// TransactionAt fetches a transaction by its ID for the given head block ID.
func (c *Client) TransactionAt(id common.Hash, head common.Hash) (*Transaction, error) {
	url := "/transactions/" + id.Hex() + "?head=" + head.Hex()
	return httpGet(c, url, &Transaction{})
}

// RawTransaction fetches a transaction by its ID and returns the raw transaction.
func (c *Client) RawTransaction(id common.Hash) (*RawTransaction, error) {
	url := "/transactions/" + id.Hex() + "?raw=true"
	return httpGet(c, url, &RawTransaction{})
}

// RawTransactionAt fetches a transaction by its ID for the given head block ID and returns the raw transaction.
func (c *Client) RawTransactionAt(id common.Hash, head common.Hash) (*RawTransaction, error) {
	url := "/transactions/" + id.Hex() + "?head=" + head.Hex() + "&raw=true"
	return httpGet(c, url, &RawTransaction{})
}

// PendingTransaction includes the pending block when fetching a transaction.
func (c *Client) PendingTransaction(id common.Hash) (*Transaction, error) {
	url := "/transactions/" + id.Hex() + "?pending=true"
	return httpGet(c, url, &Transaction{})
}

// TransactionReceipt fetches a transaction receipt by its ID.
func (c *Client) TransactionReceipt(id common.Hash) (*TransactionReceipt, error) {
	url := "/transactions/" + id.Hex() + "/receipt"
	return httpGet(c, url, &TransactionReceipt{})
}

// TransactionReceiptAt fetches a transaction receipt by its ID for the given head block ID.
func (c *Client) TransactionReceiptAt(id common.Hash, head common.Hash) (*TransactionReceipt, error) {
	url := "/transactions/" + id.Hex() + "/receipt?revision=" + head.Hex()
	return httpGet(c, url, &TransactionReceipt{})
}

func (c *Client) FilterEvents(filter *EventFilter) ([]EventLog, error) {
	path := "/logs/event"
	events := make([]EventLog, 0)
	_, err := httpPost(c, path, filter, &events)
	if err != nil {
		return nil, err
	}
	return events, nil
}

func (c *Client) FilterTransfers(filter *TransferFilter) ([]TransferLog, error) {
	path := "/logs/transfer"
	transfers := make([]TransferLog, 0)
	_, err := httpPost(c, path, filter, &transfers)
	if err != nil {
		return nil, err
	}
	return transfers, nil
}

func (c *Client) Peers() ([]Peer, error) {
	path := "/node/network/peers"
	peers := make([]Peer, 0)
	_, err := httpGet(c, path, &peers)
	if err != nil {
		return nil, err
	}
	return peers, nil
}

func httpGet[T any](c *Client, endpoint string, v *T) (*T, error) {
	req, err := http.NewRequest(http.MethodGet, c.url+endpoint, nil)
	if err != nil {
		return nil, err
	}
	return httpDo(c, req, v)
}

func httpPost[T any](c *Client, path string, body interface{}, v *T) (*T, error) {
	reqBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	request, err := http.NewRequest(http.MethodPost, c.url+path, strings.NewReader(string(reqBody)))
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json")
	return httpDo(c, request, v)
}

func httpDo[T any](c *Client, req *http.Request, v *T) (*T, error) {
	response, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	statusOK := response.StatusCode >= 200 && response.StatusCode < 300
	if !statusOK {
		return nil, newHttpError(response)
	}
	defer response.Body.Close()

	// Read the entire body into a buffer
	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	// Check if the body is "null"
	if strings.TrimSpace(string(responseBody)) == "null" {
		return nil, ErrNotFound
	}

	// Decode the JSON response
	err = json.NewDecoder(bytes.NewReader(responseBody)).Decode(v)
	if err != nil {
		return nil, err
	}

	return v, nil
}
