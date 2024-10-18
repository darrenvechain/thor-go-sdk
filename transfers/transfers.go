package transfers

import (
	"errors"

	"github.com/darrenvechain/thorgo/client"
)

var (
	descending = "desc"
	ascending  = "asc"
	time       = "time"
	block      = "block"
)

type Filter struct {
	client  *client.Client
	request *client.TransferFilter
}

func New(c *client.Client, criteria []client.TransferCriteria) *Filter {
	return &Filter{client: c, request: &client.TransferFilter{
		Criteria: &criteria,
	}}
}

// Desc sets the order of the transfers to descending. Default is ascending.
func (f *Filter) Desc() *Filter {
	f.request.Order = &descending
	return f
}

// Asc sets the order of the transfers to ascending. This is the default.
func (f *Filter) Asc() *Filter {
	f.request.Order = &ascending
	return f
}

// BlockRange sets the block range for the transfer filter.
func (f *Filter) BlockRange(from int64, to int64) *Filter {
	f.request.Range = &client.FilterRange{
		From: &from,
		To:   &to,
		Unit: &block,
	}
	return f
}

// TimeRange sets the time range for the transfer filter.
func (f *Filter) TimeRange(from int64, to int64) *Filter {
	f.request.Range = &client.FilterRange{
		From: &from,
		To:   &to,
		Unit: &time,
	}
	return f
}

// Apply sends the transfer filter to the node and returns the results.
func (f *Filter) Apply(offset int64, limit int64) ([]client.TransferLog, error) {
	if limit > 256 {
		return nil, errors.New("limit must be less than or equal to 256")
	}
	f.request.Options = &client.FilterOptions{
		Offset: &offset,
		Limit:  &limit,
	}

	return f.client.FilterTransfers(f.request)
}
