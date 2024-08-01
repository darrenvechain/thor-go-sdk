package internal

import (
	"github.com/darrenvechain/thor-go-sdk/client"
	"github.com/stretchr/testify/assert"
	"testing"
)

// TestTransfersByBlockRangeASC fetches transfers by block range in ascending order
func TestTransfersByBlockRangeASC(t *testing.T) {
	// Don't apply any criteria, just get all events
	events, err := thor.Transfers([]client.TransferCriteria{}).
		BlockRange(0, 1).
		Ascending().
		Apply(0, 100)
	assert.NoError(t, err)
	assert.NotNil(t, events)
}

// TestTransfersByTimeRangeDESC fetches transfers by time range in descending order
func TestTransfersByTimeRangeDESC(t *testing.T) {
	genesis, err := thor.Blocks.ByNumber(0)
	assert.NoError(t, err)
	best, err := thor.Blocks.Best()
	assert.NoError(t, err)

	events, err := thor.Transfers([]client.TransferCriteria{}).
		TimeRange(genesis.Timestamp, best.Timestamp).
		Descending().
		Apply(0, 100)

	assert.NoError(t, err)
	assert.NotNil(t, events)
}
