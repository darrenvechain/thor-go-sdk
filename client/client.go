package client

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/darrenvechain/thor-go-sdk/transaction"
	"github.com/ethereum/go-ethereum/common"
)

type Client struct {
	client       http.Client
	url          string
	genesisBlock *Block
}

func New(url string, client http.Client) (*Client, error) {
	return newClient(url, client)
}

func FromURL(url string) (*Client, error) {
	return New(url, http.Client{})
}

func newClient(url string, client http.Client) (*Client, error) {
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
	url := "/accounts/" + addr.String()
	return Get(c, url, new(Account))
}

// AccountForRevision fetches the account information for the given address at the given revision.
func (c *Client) AccountForRevision(addr common.Address, revision common.Hash) (*Account, error) {
	url := "/accounts/" + addr.String() + "?revision=" + revision.String()
	return Get(c, url, new(Account))
}

// Inspect will send an array of clauses to the node to simulate the execution of the clauses.
// This can be used to:
// - Read contract(s) state
// - Simulate the execution of a transaction
func (c *Client) Inspect(body InspectRequest) ([]InspectResponse, error) {
	url := "/accounts/*"
	result, err := Post(c, url, body, new([]InspectResponse))
	if err != nil {
		return nil, err
	}
	return *result, nil
}

// AccountCode fetches the code for the account at the given address.
func (c *Client) AccountCode(addr common.Address) (*AccountCode, error) {
	url := "/accounts/" + addr.String() + "/code"
	return Get(c, url, new(AccountCode))
}

// AccountCodeForRevision fetches the code for the account at the given address at the given revision.
func (c *Client) AccountCodeForRevision(addr common.Address, revision common.Hash) (*AccountCode, error) {
	url := "/accounts/" + addr.String() + "/code?revision=" + revision.String()
	return Get(c, url, new(AccountCode))
}

// AccountStorage fetches the storage value for the account at the given address and key.
func (c *Client) AccountStorage(addr common.Address, key common.Hash) (*AccountStorage, error) {
	url := "/accounts/" + addr.String() + "/storage/" + key.Hex()
	return Get(c, url, new(AccountStorage))
}

// AccountStorageForRevision fetches the storage value for the account at the given address and key at the given revision.
func (c *Client) AccountStorageForRevision(
	addr common.Address,
	key common.Hash,
	revision common.Hash,
) (*AccountStorage, error) {
	url := "/accounts/" + addr.Hex() + "/storage/" + key.Hex() + "?revision=" + revision.String()
	return Get(c, url, new(AccountStorage))
}

// Block fetches the block for the given revision.
func (c *Client) Block(revision string) (*Block, error) {
	url := "/blocks/" + revision
	return Get(c, url, new(Block))
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
	return Get(c, url, new(ExpandedBlock))
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
	return Post(c, "/transactions", body, new(SendTransactionResponse))
}

// SendRawTransaction sends a raw transaction to the node.
func (c *Client) SendRawTransaction(raw string) (*SendTransactionResponse, error) {
	body := make(map[string]string)
	body["raw"] = raw
	return Post(c, "/transactions", body, new(SendTransactionResponse))
}

// Transaction fetches a transaction by its ID.
func (c *Client) Transaction(id common.Hash) (*Transaction, error) {
	url := "/transactions/" + id.Hex()
	return Get(c, url, new(Transaction))
}

// RawTransaction fetches a transaction by its ID and returns the raw transaction.
func (c *Client) RawTransaction(id common.Hash) (*RawTransaction, error) {
	url := "/transactions/" + id.Hex() + "?raw=true"
	return Get(c, url, new(RawTransaction))
}

// PendingTransaction includes the pending block when fetching a transaction.
func (c *Client) PendingTransaction(id common.Hash) (*Transaction, error) {
	url := "/transactions/" + id.Hex() + "?pending=true"
	return Get(c, url, new(Transaction))
}

func (c *Client) TransactionReceipt(id common.Hash) (*TransactionReceipt, error) {
	url := "/transactions/" + id.Hex() + "/receipt"
	return Get(c, url, new(TransactionReceipt))
}

func (c *Client) FilterEvents(filter *EventFilter) ([]EventLog, error) {
	path := "/logs/event"
	result, err := Post(c, path, filter, new([]EventLog))
	if err != nil {
		return nil, err
	}
	return *result, nil
}

func (c *Client) FilterTransfers(filter *TransferFilter) ([]TransferLog, error) {
	path := "/logs/transfer"
	result, err := Post(c, path, filter, new([]TransferLog))
	if err != nil {
		return nil, err
	}
	return *result, nil
}

func (c *Client) Peers() ([]Peer, error) {
	result, err := Get(c, "/node/network/peers", new([]Peer))
	if err != nil {
		return nil, err
	}
	return *result, nil
}

func Get[T any](c *Client, endpoint string, v *T) (*T, error) {
	req, err := http.NewRequest("GET", c.url+endpoint, nil)
	if err != nil {
		return nil, err
	}
	return Do(c, req, v)
}

func Post[T any](c *Client, path string, body interface{}, v *T) (*T, error) {
	reqBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	request, err := http.NewRequest("POST", c.url+path, strings.NewReader(string(reqBody)))
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json")
	return Do(c, request, v)
}

func Do[T any](c *Client, req *http.Request, v *T) (*T, error) {
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
	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	// Check if the body is "null"
	if strings.TrimSpace(string(bodyBytes)) == "null" {
		return nil, ErrNotFound
	}

	// Create a new reader from the buffered body
	bodyReader := bytes.NewReader(bodyBytes)

	// Decode the JSON response
	err = json.NewDecoder(bodyReader).Decode(v)
	if err != nil {
		return nil, err
	}

	return v, nil
}
