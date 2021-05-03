package models

import (
	"encoding/json"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
)

type DatasourceSettings struct {
	// global settings
}

func GetDatasourceSettings(dsInfo backend.DataSourceInstanceSettings) (*DatasourceSettings, error) {
	s := &DatasourceSettings{}
	if err := json.Unmarshal(dsInfo.JSONData, s); err != nil {
		return nil, err
	}
	return s, nil
}
