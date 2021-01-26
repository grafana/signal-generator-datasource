package models

import (
	"encoding/json"
	"strings"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
)

type DatasurceSettings struct {
	// For live
	LiveURL string `json:"live,omitempty"`

	// Single line... split by newline to capture
	CaptureX string `json:"captureX,omitempty"`

	// Capture configs
	Capture []string `json:"capture,omitempty"`
}

func GetDatasurceSettings(dsInfo backend.DataSourceInstanceSettings) (*DatasurceSettings, error) {
	s := &DatasurceSettings{}
	if err := json.Unmarshal(dsInfo.JSONData, s); err != nil {
		return nil, err
	}

	if (len(s.Capture) == 0) && s.CaptureX != "" {
		for _, c := range strings.Split(s.CaptureX, "\n") {
			p := strings.TrimSpace(c)
			if strings.HasPrefix(p, "#") || len(p) < 1 {
				// skip
			} else {
				s.Capture = append(s.Capture, p)
			}
		}
	}

	// default streaming URL
	if len(s.LiveURL) < 1 {
		s.LiveURL = "http://localhost:3000/"
	}

	return s, nil
}
