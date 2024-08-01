package internal

import (
	"github.com/darrenvechain/thor-go-sdk/client"
	"github.com/stretchr/testify/assert"
	"testing"
)

// TestEventByBlockRangeASC fetches events by block range in ascending order
func TestEventByBlockRangeASC(t *testing.T) {
	// Don't apply any criteria, just get all events
	events, err := thor.Events([]client.EventCriteria{}).
		BlockRange(0, 1).
		Ascending().
		Apply(0, 100)
	assert.NoError(t, err)
	assert.NotNil(t, events)
}

// TestEventsByTimeRangeDESC fetches events by time range in descending order
func TestEventsByTimeRangeDESC(t *testing.T) {
	genesis, err := thor.Blocks.ByNumber(0)
	assert.NoError(t, err)
	best, err := thor.Blocks.Best()
	assert.NoError(t, err)

	events, err := thor.Events([]client.EventCriteria{}).
		TimeRange(genesis.Timestamp, best.Timestamp).
		Descending().
		Apply(0, 100)

	assert.NoError(t, err)
	assert.NotNil(t, events)
}
