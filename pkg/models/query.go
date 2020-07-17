package models

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
)

// QueryTypeEnum is the possible values for incoming queryType parameter
type QueryTypeEnum string

const (
	// QueryTypeUsers -- show a list of users
	QueryTypeUsers QueryTypeEnum = "users"

	// QueryTypeRepos -- list the repositories
	QueryTypeRepos = "repos"

	// QueryTypeBuilds -- the last
	QueryTypeBuilds = "builds"

	// QueryTypeLogs -- the last
	QueryTypeLogs = "logs"

	// QueryTypeIncomplete -- the running buidls
	QueryTypeIncomplete = "incomplete"

	// QueryTypeNodes -- the last
	QueryTypeNodes = "nodes"

	// QueryTypeServers -- the current servers
	QueryTypeServers = "servers"
)

// QueryModel represents
type QueryModel struct {
	QueryType QueryTypeEnum `json:"queryType,omitempty"`

	// Not from JSON
	TimeRange backend.TimeRange `json:"-"`
}

type InfluxLine struct {
	Name      string                 `json:"name,omitempty"`
	Fields    map[string]interface{} `json:"fields,omitempty"`
	Tags      map[string]string      `json:"tags,omitempty"`
	Timestamp time.Time              `json:"timestamp,omitempty"`
}

// GetQueryModel returns the well typed query model
func GetQueryModel(query backend.DataQuery) (*QueryModel, error) {
	model := &QueryModel{}

	err := json.Unmarshal(query.JSON, &model)
	if err != nil {
		return nil, fmt.Errorf("error reading query: %s", err.Error())
	}

	// Copy directly from the well typed query
	model.TimeRange = query.TimeRange
	return model, nil
}
