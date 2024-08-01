// Copyright (c) 2018 The VeChainThor developers

// Distributed under the GNU Lesser General Public License v3.0 software license, see the accompanying
// file LICENSE or <https://www.gnu.org/licenses/lgpl-3.0.html>

package transaction

import (
	"encoding/binary"

	"github.com/ethereum/go-ethereum/common"
)

// BlockRef is block reference.
type BlockRef [8]byte

// Number extracts block number.
func (b BlockRef) Number() uint32 {
	return binary.BigEndian.Uint32(b[:])
}

// NewBlockRef create block reference with block number.
func NewBlockRef(blockNum uint32) (br BlockRef) {
	binary.BigEndian.PutUint32(br[:], blockNum)
	return
}

// NewBlockRefFromID create block reference from block id.
func NewBlockRefFromID(blockID common.Hash) (br BlockRef) {
	copy(br[:], blockID[:])
	return
}

func (b BlockRef) UnmarshalJSON(data []byte) error {
	// block ref is returned as a hex string from the API
	decoded := common.Hex2Bytes(string(data[1 : len(data)-1]))
	copy(b[:], decoded)
	return nil
}
