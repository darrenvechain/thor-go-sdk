package transfers

import (
	"testing"

	"github.com/darrenvechain/thor-go-sdk/client"
	"github.com/darrenvechain/thor-go-sdk/solo"
	blocks2 "github.com/darrenvechain/thor-go-sdk/thorgo/blocks"
	"github.com/stretchr/testify/assert"
)

var (
	thorClient, _ = client.FromURL(solo.URL)
	blocks        = blocks2.New(thorClient)
)

// TestTransfersByBlockRangeASC fetches transfers by block range in ascending order
func TestTransfersByBlockRangeASC(t *testing.T) {
	// Don't apply any criteria, just get all events
	events, err := New(thorClient, []client.TransferCriteria{}).
		BlockRange(0, 1).
		Ascending().
		Apply(0, 100)
	assert.NoError(t, err)
	assert.NotNil(t, events)
}

// TestTransfersByTimeRangeDESC fetches transfers by time range in descending order
func TestTransfersByTimeRangeDESC(t *testing.T) {
	genesis, err := blocks.ByNumber(0)
	assert.NoError(t, err)
	best, err := blocks.Best()
	assert.NoError(t, err)

	events, err := New(thorClient, []client.TransferCriteria{}).
		TimeRange(genesis.Timestamp, best.Timestamp).
		Descending().
		Apply(0, 100)

	assert.NoError(t, err)
	assert.NotNil(t, events)
}
