package models

import (
	"encoding/json"
	"fmt"

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
