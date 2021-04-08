package models

import (
	"encoding/json"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
)

type DatasurceSettings struct {
	// global settings
}

func GetDatasurceSettings(dsInfo backend.DataSourceInstanceSettings) (*DatasurceSettings, error) {
	s := &DatasurceSettings{}
	if err := json.Unmarshal(dsInfo.JSONData, s); err != nil {
		return nil, err
	}

	return s, nil
}
