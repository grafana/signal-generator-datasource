package plugin

import (
	"encoding/json"
	"time"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/signal-generator-datasource/pkg/models"
	"github.com/grafana/signal-generator-datasource/pkg/waves"
)

// DatasourceHandler is the plugin entrypoint and implements all of the necessary handler functions for dataqueries, healthchecks, and resources.
type SignalStreamer struct {
	signal      *waves.SignalGen
	channel     *GrafanaLiveChannel
	running     bool
	speedMillis int64
}

func (s *SignalStreamer) Stop() {
	s.running = false
}

func (s *SignalStreamer) Start() {
	if s.running {
		backend.Logger.Info("already running")
		return
	}

	if s.signal == nil {
		cfg := models.SignalConfig{
			Time: models.TimeFieldConfig{
				Period: "5s",
			},
			Fields: []models.ExpressionConfig{},
		}
		cfg.Fields = append(cfg.Fields, models.ExpressionConfig{
			BaseSignalField: models.BaseSignalField{
				Name: "A",
			},
			Expr: "Sine(x)",
		})
		cfg.Fields = append(cfg.Fields, models.ExpressionConfig{
			BaseSignalField: models.BaseSignalField{
				Name: "B",
			},
			Expr: "Sine(x+1)",
		})
		cfg.Fields = append(cfg.Fields, models.ExpressionConfig{
			BaseSignalField: models.BaseSignalField{
				Name: "C",
			},
			Expr: "Sine(x+1.5)*2",
		})

		gen, _ := waves.NewSignalGenerator(cfg)
		if gen != nil {
			s.signal = gen
		}
	}

	if s.channel == nil {
		c, err := InitGrafanaLiveChannel("ws://localhost:3000/live/ws", "grafana/measurements/signal")
		if err != nil {
			backend.Logger.Error("error starting live")
			return
		}
		s.channel = c
	}

	if s.speedMillis < 10 {
		s.speedMillis = 1000
	}

	go s.doStream()
}

func (s *SignalStreamer) doStream() {
	s.running = true
	ticker := time.NewTicker(time.Duration(s.speedMillis) * time.Millisecond)

	measurement := models.Measurement{
		Name:   "Example",
		Time:   0,
		Values: make(map[string]interface{}, 5),
	}
	msg := models.MeasurementBatch{
		Measurements: []models.Measurement{measurement}, // always a single measurement
	}

	for t := range ticker.C {
		if !s.running {
			backend.Logger.Info("stoppint!!!")
			return
		}

		v := 7 //s.signal.GetValue(t)

		measurement.Time = t.UnixNano() / int64(time.Millisecond)
		measurement.Values["value"] = v

		bytes, err := json.Marshal(&msg)
		if err != nil {
			backend.Logger.Warn("unable to marshal line", "error", err)
			continue
		}
		s.channel.Publish(bytes)
	}
}
