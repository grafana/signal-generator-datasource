package plugin

import (
	"encoding/json"
	"math"
	"time"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/measurement"
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
			Expr: "Sine(x+1.5) * 2 + Noise() * 0.4", // + Noise()*.5",
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
		s.speedMillis = 2500
	}

	go s.doStream()
	//s.doStream()
}

func (s *SignalStreamer) doStream() {
	s.running = true
	ticker := time.NewTicker(time.Duration(s.speedMillis) * time.Millisecond)

	m := measurement.Measurement{
		Name:   "Example",
		Time:   0,
		Values: make(map[string]interface{}, 5),
	}
	msg := measurement.Batch{
		Measurements: []measurement.Measurement{m}, // always a single measurement
	}

	paramCount := len(s.signal.Fields) + 4
	parameters := make(map[string]interface{}, paramCount)
	parameters["PI"] = math.Pi

	backend.Logger.Info("START STREAMING", "sig", s.signal)

	for t := range ticker.C {
		if !s.running {
			backend.Logger.Info("stoppint!!!")
			return
		}

		m.Time = t.UnixNano() / int64(time.Millisecond)

		// Set the time
		for _, i := range s.signal.Inputs {
			err := i.UpdateEnv(&t, parameters)
			if err != nil {
				backend.Logger.Warn("ERROR updating time", "error", err)
			}
		}

		// Calculate each value
		for _, f := range s.signal.Fields {
			v, err := f.GetValue(parameters)
			if err != nil {
				v = nil
			}
			name := f.GetConfig().Name
			parameters[name] = v
			m.Values[name] = v
		}

		bytes, err := json.Marshal(&msg)
		if err != nil {
			backend.Logger.Warn("unable to marshal line", "error", err)
			continue
		}
		s.channel.Publish(bytes)
	}
}
