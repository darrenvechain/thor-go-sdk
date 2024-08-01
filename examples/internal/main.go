package internal

import (
	"github.com/darrenvechain/thor-go-sdk/solo"
	"github.com/darrenvechain/thor-go-sdk/thorgo"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

var (
	thor         *thorgo.Thor
	account1, _  = solo.Key(0)
	account1Addr = crypto.PubkeyToAddress(account1.PublicKey)
	account2, _  = solo.Key(1)
	account2Addr = crypto.PubkeyToAddress(account2.PublicKey)
	vthoAddr     = common.HexToAddress("0x0000000000000000000000000000456e65726779")
)

func init() {
	var err error
	thor, err = thorgo.FromURL(solo.URL)
	if err != nil {
		panic(err)
	}
}
