package models

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/data"
)

const (
	QueryTypeAWG     = "AWG"
	QueryTypeEasings = "easing"
	QueryStreams     = "streams"
)

// x = 0-2π
// t = time in seconds
// PI = 3.1415

type BaseSignalField struct {
	Name   string            `json:"name,omitempty"`
	Config *data.FieldConfig `json:"config,omitempty"`
	Labels data.Labels       `json:"labels,omitempty"`
}

func (b *BaseSignalField) SetValue(v string) error {
	return fmt.Errorf("unimplemented")
}

type ExpressionConfig struct {
	BaseSignalField

	Expr     string         `json:"expr,omitempty"`
	DataType data.FieldType `json:"type,omitempty"`
}

type TimeFieldConfig struct {
	BaseSignalField

	Period string `json:"period,omitempty"` // time string
}

type RangeFieldConfig struct {
	BaseSignalField

	Min   float64 `json:"min,omitempty"`   // The frame name
	Max   float64 `json:"max,omitempty"`   // time string
	Count int64   `json:"count,omitempty"` // 0 will use maxDataPoints
	Ease  string  `json:"ease,omitempty"`  // easing funciton 0-1
}

type SignalConfig struct {
	Name string `json:"name,omitempty"` // The frame name

	Time  TimeFieldConfig  `json:"time,omitempty"`
	Range RangeFieldConfig `json:"range,omitempty"`

	// The non-time fields
	Fields []ExpressionConfig `json:"fields,omitempty"`
}

type SignalQuery struct {
	Signal SignalConfig `json:"signal,omitempty"` // all components get added together
	Stream string       `json:"stream,omitempty"` // used for streams

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
