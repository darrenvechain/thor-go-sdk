package events

import (
	"testing"

	"github.com/darrenvechain/thor-go-sdk/client"
	"github.com/darrenvechain/thor-go-sdk/solo"
	"github.com/darrenvechain/thor-go-sdk/thorgo/blocks"
	"github.com/stretchr/testify/assert"
)

var (
	thorClient, _ = client.FromURL(solo.URL)
)

// TestEventByBlockRangeASC fetches events by block range in ascending order
func TestEventByBlockRangeASC(t *testing.T) {
	// Don't apply any criteria, just get all events
	events, err := New(thorClient, []client.EventCriteria{}).
		BlockRange(0, 1).
		Ascending().
		Apply(0, 100)
	assert.NoError(t, err)
	assert.NotNil(t, events)
}

// TestEventsByTimeRangeDESC fetches events by time range in descending order
func TestEventsByTimeRangeDESC(t *testing.T) {
	genesis, err := blocks.New(thorClient).ByNumber(0)
	assert.NoError(t, err)
	best, err := blocks.New(thorClient).Best()
	assert.NoError(t, err)

	events, err := New(thorClient, []client.EventCriteria{}).
		TimeRange(genesis.Timestamp, best.Timestamp).
		Descending().
		Apply(0, 100)

	assert.NoError(t, err)
	assert.NotNil(t, events)
}
