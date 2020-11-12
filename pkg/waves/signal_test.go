package waves

import (
	"strings"
	"testing"
	"time"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/grafana/grafana-plugin-sdk-go/experimental"
	"github.com/grafana/signal-generator-datasource/pkg/models"
)

func TestEval(t *testing.T) {
	f0 := models.SignalField{
		Name: "Hello",
		Expr: "x",
	}

	query := &models.SignalQuery{
		Period:        "10s",
		MaxDataPoints: 10,
		Signal: models.SignalConfig{
			Name:   "test",
			Fields: []models.SignalField{f0},
		},
		TimeRange: backend.TimeRange{
			From: time.Unix(0, 0),
			To:   time.Unix(20, 0),
		},
	}

	frame, err := DoSignalQuery(query)
	dr := backend.DataResponse{
		Frames: data.Frames{frame},
		Error:  err,
	}

	fname := "../testdata/simple.golden.txt"
	if err := experimental.CheckGoldenDataResponse(fname, &dr, true); err != nil {
		if !strings.Contains(err.Error(), "no such file or directory") {
			t.Fatal(err)
		}
	}
}
