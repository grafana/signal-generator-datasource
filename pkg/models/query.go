package models

import (
	"encoding/json"
	"time"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/data"
)

const (
	QueryTypeAWG     = "AWG"
	QueryTypeEasings = "easing"
)

// x = 0-2π
// t = time in seconds
// PI = 3.1415

type ValueRange struct {
	Min   float64 `json:"min,omitempty"`   // The frame name
	Max   float64 `json:"max,omitempty"`   // time string
	Count int64   `json:"count,omitempty"` // time string
}

type SignalField struct {
	Name   string            `json:"name,omitempty"`
	Expr   string            `json:"expr,omitempty"`
	Config *data.FieldConfig `json:"config,omitempty"`
	Labels data.Labels       `json:"labels,omitempty"`
}

type SignalConfig struct {
	Name string `json:"name,omitempty"` // The frame name

	// The non-time fields
	Fields []SignalField `json:"fields,omitempty"`
}

type SignalQuery struct {
	Signal SignalConfig `json:"signal,omitempty"` // all components get added together

	// Parameters to generate the initial x value
	Period string      `json:"period,omitempty"` // time string
	XGen   *ValueRange `json:"xgen,omitempty"`   // OR xvalues
	Ease   string      `json:"ease,omitempty"`   // used in easing query

	// These are added from the base query
	Interval      time.Duration     `json:"-"`
	TimeRange     backend.TimeRange `json:"-"`
	MaxDataPoints int64             `json:"-"`
	QueryType     string            `json:"-"`
}

func GetSignalQuery(dq *backend.DataQuery) (*SignalQuery, error) {
	query := &SignalQuery{}
	if err := json.Unmarshal(dq.JSON, query); err != nil {
		return nil, err
	}

	// add on the DataQuery params
	query.TimeRange = dq.TimeRange
	query.Interval = dq.Interval
	query.MaxDataPoints = dq.MaxDataPoints
	query.QueryType = dq.QueryType

	return query, nil
}
