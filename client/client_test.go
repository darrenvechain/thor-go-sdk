package client

import (
	"github.com/darrenvechain/thorgo/solo"
)

var client *Client

func init() {
	var err error
	client, err = FromURL(solo.URL)
	if err != nil {
		panic(err)
	}
}
