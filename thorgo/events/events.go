package events

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

type Filter struct {
	client  *client.Client
	request *client.EventFilter
}

func New(c *client.Client, criteria []client.EventCriteria) *Filter {
	return &Filter{client: c, request: &client.EventFilter{
		Criteria: &criteria,
	}}
}

func (f *Filter) Descending() *Filter {
	f.request.Order = &descending
	return f
}

func (f *Filter) Ascending() *Filter {
	f.request.Order = &ascending
	return f
}

func (f *Filter) BlockRange(from uint64, to uint64) *Filter {
	f.request.Range = &client.FilterRange{
		From: &from,
		To:   &to,
		Unit: &block,
	}
	return f
}

func (f *Filter) TimeRange(from uint64, to uint64) *Filter {
	f.request.Range = &client.FilterRange{
		From: &from,
		To:   &to,
		Unit: &time,
	}
	return f
}

func (f *Filter) Apply(offset uint64, limit uint64) ([]client.EventLog, error) {
	if limit > 256 {
		return nil, errors.New("limit must be less than or equal to 256")
	}

	f.request.Options = &client.FilterOptions{
		Offset: &offset,
		Limit:  &limit,
	}

	return f.client.FilterEvents(f.request)
}
