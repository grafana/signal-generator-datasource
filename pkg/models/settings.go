package models

import (
	"encoding/json"
	"fmt"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
)

// DatasourceSettings contains Settings needed for datasource
type DatasourceSettings struct {
	URL string

	// Secure JSON
	Token string
}

// LoadSettings gets the relevant settings from the plugin context
func LoadSettings(setting backend.DataSourceInstanceSettings) (*DatasourceSettings, error) {
	model := &DatasourceSettings{}

	err := json.Unmarshal(setting.JSONData, &model)
	if err != nil {
		return nil, fmt.Errorf("error reading settings: %s", err.Error())
	}

	return model, nil
}
