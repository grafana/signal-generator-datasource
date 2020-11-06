package models

import (
	"encoding/json"
	"time"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/signal-generator-datasource/pkg/waves"
)

const (
	QueryTypeAWG     = "AWG"
	QueryTypeEasings = "easing"
)

type SignalQuery struct {
	Signal []waves.SignalArgs `json:"signals,omitempty"` // all components get added together
	Ease   string             `json:"ease,omitempty"`    // used in easing query

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
