package transfers

import (
	"errors"

	"github.com/darrenvechain/thor-go-sdk/client"
)

var (
	descending = "desc"
	ascending  = "asc"
	time       = "time"
	block      = "block"
)

type Transfers struct {
	client  *client.Client
	request *client.TransferFilter
}

func New(c *client.Client, criteria []client.TransferCriteria) *Transfers {
	return &Transfers{client: c, request: &client.TransferFilter{
		Criteria: &criteria,
	}}
}

func (t *Transfers) Descending() *Transfers {
	t.request.Order = &descending
	return t
}

func (t *Transfers) Ascending() *Transfers {
	t.request.Order = &ascending
	return t
}

// BlockRange sets the block range for the transfer filter.
func (t *Transfers) BlockRange(from uint64, to uint64) *Transfers {
	t.request.Range = &client.FilterRange{
		From: &from,
		To:   &to,
		Unit: &block,
	}
	return t
}

// TimeRange sets the time range for the transfer filter.
func (t *Transfers) TimeRange(from uint64, to uint64) *Transfers {
	t.request.Range = &client.FilterRange{
		From: &from,
		To:   &to,
		Unit: &time,
	}
	return t
}

// Apply sends the transfer filter to the node and returns the results.
func (t *Transfers) Apply(offset uint64, limit uint64) (*[]client.TransferLog, error) {
	if limit > 256 {
		return nil, errors.New("limit must be less than or equal to 256")
	}
	t.request.Options = &client.FilterOptions{
		Offset: &offset,
		Limit:  &limit,
	}

	return t.client.FilterTransfers(t.request)
}
