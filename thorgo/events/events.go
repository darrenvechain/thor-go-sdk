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

type Events struct {
	client  *client.Client
	request *client.EventFilter
}

func New(c *client.Client, criteria []client.EventCriteria) *Events {
	return &Events{client: c, request: &client.EventFilter{
		Criteria: &criteria,
	}}
}

func (e *Events) Descending() *Events {
	e.request.Order = &descending
	return e
}

func (e *Events) Ascending() *Events {
	e.request.Order = &ascending
	return e
}

func (e *Events) BlockRange(from uint64, to uint64) *Events {
	e.request.Range = &client.FilterRange{
		From: &from,
		To:   &to,
		Unit: &block,
	}
	return e
}

func (e *Events) TimeRange(from uint64, to uint64) *Events {
	e.request.Range = &client.FilterRange{
		From: &from,
		To:   &to,
		Unit: &time,
	}
	return e
}

func (e *Events) Apply(offset uint64, limit uint64) (*[]client.EventLog, error) {
	if limit > 256 {
		return nil, errors.New("limit must be less than or equal to 256")
	}

	e.request.Options = &client.FilterOptions{
		Offset: &offset,
		Limit:  &limit,
	}

	return e.client.FilterEvents(e.request)
}
