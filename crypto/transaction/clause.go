// Copyright (c) 2018 The VeChainThor developers

// Distributed under the GNU Lesser General Public License v3.0 software license, see the accompanying
// file LICENSE or <https://www.gnu.org/licenses/lgpl-3.0.html>

package transaction

import (
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
)

type clauseBody struct {
	To    *common.Address `rlp:"nil"`
	Value *big.Int
	Data  []byte
}

// Clause is the basic execution unit of a transaction.
type Clause struct {
	body clauseBody
}

// NewClause create a new clause instance.
func NewClause(to *common.Address) *Clause {
	if to != nil {
		// make a copy of 'to'
		cpy := *to
		to = &cpy
	}
	return &Clause{
		clauseBody{
			to,
			&big.Int{},
			nil,
		},
	}
}

// WithValue create a new clause copy with value changed.
func (c *Clause) WithValue(value *big.Int) *Clause {
	newClause := *c
	newClause.body.Value = new(big.Int).Set(value)
	return &newClause
}

// WithData create a new clause copy with data changed.
func (c *Clause) WithData(data []byte) *Clause {
	newClause := *c
	newClause.body.Data = append([]byte(nil), data...)
	return &newClause
}

// To returns 'To' address.
func (c *Clause) To() *common.Address {
	if c.body.To == nil {
		return nil
	}
	cpy := *c.body.To
	return &cpy
}

// Value returns 'Value'.
func (c *Clause) Value() *big.Int {
	return new(big.Int).Set(c.body.Value)
}

// Data returns 'Data'.
func (c *Clause) Data() []byte {
	return append([]byte(nil), c.body.Data...)
}

// IsCreatingContract return if this clause is going to create a contract.
func (c *Clause) IsCreatingContract() bool {
	return c.body.To == nil
}

// EncodeRLP implements rlp.Encoder
func (c *Clause) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, &c.body)
}

// DecodeRLP implements rlp.Decoder
func (c *Clause) DecodeRLP(s *rlp.Stream) error {
	var body clauseBody
	if err := s.Decode(&body); err != nil {
		return err
	}
	*c = Clause{body}
	return nil
}

func (c *Clause) String() string {
	var to string
	if c.body.To == nil {
		to = "nil"
	} else {
		to = c.body.To.String()
	}
	return fmt.Sprintf(`
		(To:	%v
		 Value:	%v
		 Data:	0x%x)`, to, c.body.Value, c.body.Data)
}

func (c *Clause) MarshalJSON() ([]byte, error) {
	body := make(map[string]interface{})
	if c.body.To != nil {
		body["to"] = c.body.To.String()
	} else {
		body["to"] = nil
	}
	body["value"] = "0x" + c.body.Value.Text(16)
	body["data"] = "0x" + common.Bytes2Hex(c.body.Data)
	return json.Marshal(body)
}

func (c *Clause) UnmarshalJSON(bytes []byte) error {
	var body map[string]*string
	if err := json.Unmarshal(bytes, &body); err != nil {
		return err
	}
	to, ok := body["to"]
	if !ok {
		return fmt.Errorf("missing 'to' field")
	}
	if to == nil {
		c.body.To = nil
	} else {
		addr := common.HexToAddress(*to)
		c.body.To = &addr
	}
	value, ok := body["value"]
	if !ok {
		return fmt.Errorf("missing 'value' field")
	}
	val, ok := new(big.Int).SetString(strings.TrimPrefix(*value, "0x"), 16)
	if !ok {
		return fmt.Errorf("invalid 'value' field")
	}
	c.body.Value = val
	data, ok := body["data"]
	if !ok {
		return fmt.Errorf("missing 'data' field")
	}
	c.body.Data = common.Hex2Bytes(*data)
	return nil
}
