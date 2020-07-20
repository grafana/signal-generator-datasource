package models

import (
	"time"
)

// InfluxLine looks like a line from telegraph json
type InfluxLine struct {
	Name      string                 `json:"name,omitempty"`
	Fields    map[string]interface{} `json:"fields,omitempty"`
	Tags      map[string]string      `json:"tags,omitempty"`
	Timestamp time.Time              `json:"timestamp,omitempty"`
}
