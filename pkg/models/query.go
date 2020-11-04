package models

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/signal-generator-datasource/pkg/waves"
)

const (
	QueryTypeAWG     = "AWG"
	QueryTypeEasings = "easing"
)

type SignalQuery struct {
	Wave []waves.WaveformArgs `json:"wave,omitempty"` // all components get added together
	Ease string               `json:"ease,omitempty"` // used in easing query

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

	// Convert the string period arguments to seconds
	if query.Wave != nil {
		for _, wave := range query.Wave {
			if strings.HasPrefix(wave.Period, "range/") {
				f, err := strconv.ParseFloat(wave.Period[6:], 64)
				if err != nil {
					return nil, fmt.Errorf("error reading wave period")
				}
				r := dq.TimeRange.To.Sub(dq.TimeRange.From).Seconds() / f
				wave.PeriodSec = r
			} else if wave.Period != "" {
				r, err := time.ParseDuration(wave.Period)
				if err != nil {
					return nil, fmt.Errorf("error reading wave period")
				}
				wave.PeriodSec = r.Seconds()
			}

			backend.Logger.Info("PARSE", "wave", wave.Period, "ppp", wave.PeriodSec)
		}
	}

	// add on the DataQuery params
	query.TimeRange = dq.TimeRange
	query.Interval = dq.Interval
	query.MaxDataPoints = dq.MaxDataPoints
	query.QueryType = dq.QueryType

	return query, nil
}
