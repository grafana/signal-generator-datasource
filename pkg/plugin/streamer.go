package plugin

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"time"

	"github.com/grafana/grafana-edge-app/pkg/capture"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/grafana/grafana-plugin-sdk-go/live"
	"github.com/grafana/grafana-plugin-sdk-go/measurement"
	"github.com/grafana/signal-generator-datasource/pkg/models"
	"github.com/grafana/signal-generator-datasource/pkg/waves"
)

// DatasourceHandler is the plugin entrypoint and implements all of the necessary handler functions for dataqueries, healthchecks, and resources.
type SignalStreamer struct {
	signal      *waves.SignalGen
	client      *live.GrafanaLiveClient
	channel     *live.GrafanaLiveChannel
	running     bool
	speedMillis int64
	current     measurement.Measurement
	frame       *data.Frame
}

func NewSignalStreamer(extcfg *capture.CaptureSetConfig, client *live.GrafanaLiveClient) (*SignalStreamer, error) {
	cfg := models.SignalConfig{
		Time: models.TimeFieldConfig{
			Period: "5s",
		},
		Fields: []models.ExpressionConfig{},
	}

	speedMillis := int64(1500)
	if extcfg.Interval != "" {
		d, err := time.ParseDuration(extcfg.Interval)
		if err == nil {
			speedMillis = d.Milliseconds()
		}
	}

	for idx := range extcfg.Input {
		tag := extcfg.Input[idx]
		if tag.Path == "time" {
			// TODO... configure the time period
			continue
		}
		name := tag.Name
		if len(name) < 1 {
			name = tag.Path
		}

		if len(name) < 1 {
			return nil, fmt.Errorf("invalid field name for tag: %v", tag)
		}

		if len(tag.Path) > 1 {
			tag.Config.Path = tag.Path
		}

		if tag.Value == nil || tag.Value == "" {
			return nil, fmt.Errorf("missing value for field: %s", tag.Path)
		}

		cfg.Fields = append(cfg.Fields, models.ExpressionConfig{
			BaseSignalField: models.BaseSignalField{
				Name:   name,
				Config: &tag.Config,
				Labels: tag.Labels,
			},
			Expr: fmt.Sprintf("%v", tag.Value),
		})
	}

	gen, err := waves.NewSignalGenerator(cfg)
	if err != nil {
		return nil, err
	}

	rowCount := 1
	fields := make([]*data.Field, len(gen.Fields)+1)
	fields[0] = data.NewFieldFromFieldType(data.FieldTypeTime, rowCount)
	fields[0].Name = "Time"
	for i, f := range gen.Fields {
		cfg := f.GetConfig()
		fields[i+1] = data.NewFieldFromFieldType(cfg.DataType, rowCount)
		fields[i+1].Name = cfg.Name
		fields[i+1].Config = cfg.Config
		fields[i+1].Labels = cfg.Labels
	}

	frame := data.NewFrame(extcfg.Name, fields...)
	frame.SetMeta(&data.FrameMeta{
		Custom: &models.CustomFrameMeta{
			StreamKey: "signal/" + extcfg.Name,
		},
	})

	m := measurement.Measurement{
		Name:   extcfg.Name,
		Time:   0,
		Values: make(map[string]interface{}, 5),
	}

	return &SignalStreamer{
		signal:      gen,
		client:      client,
		current:     m,
		frame:       frame,
		speedMillis: speedMillis,
	}, nil
}

func (s *SignalStreamer) Stop() {}

func (s *SignalStreamer) Start() {}

func (s *SignalStreamer) UpdateValues(props map[string]interface{}) error {
	err := s.signal.UpdateValues(props)
	if err != nil {
		return err
	}

	paramCount := len(s.signal.Fields) + 4
	parameters := make(map[string]interface{}, paramCount)
	parameters["PI"] = math.Pi

	t := time.Now()
	s.frame.Fields[0].Set(0, t)
	s.current.Time = t.UnixNano() / int64(time.Millisecond)

	// Set the time
	for _, i := range s.signal.Inputs {
		err := i.UpdateEnv(&t, parameters)
		if err != nil {
			backend.Logger.Warn("ERROR updating time", "error", err)
		}
	}

	// Calculate each value
	for idx, f := range s.signal.Fields {
		v, err := f.GetValue(parameters)
		if err != nil {
			v = float64(0)
		}
		name := f.GetConfig().Name
		parameters[name] = v
		s.current.Values[name] = v

		s.frame.Fields[idx+1].Set(0, v)
	}
	return nil
}

func (s *SignalStreamer) doStream(ctx context.Context, sender backend.StreamPacketSender) {
	ticker := time.NewTicker(time.Duration(s.speedMillis) * time.Millisecond)
	defer ticker.Stop()

	msg := measurement.Batch{
		Measurements: make([]measurement.Measurement, 1), // always a single measurement
	}

	paramCount := len(s.signal.Fields) + 4
	parameters := make(map[string]interface{}, paramCount)
	parameters["PI"] = math.Pi

	backend.Logger.Info("start streaming")

	for {
		select {
		case <-ctx.Done():
			backend.Logger.Info("stop streaming (context canceled)")
			return
		case t := <-ticker.C:
			s.frame.Fields[0].Set(0, t)
			s.current.Time = t.UnixNano() / int64(time.Millisecond)

			// Set the time
			for _, i := range s.signal.Inputs {
				err := i.UpdateEnv(&t, parameters)
				if err != nil {
					backend.Logger.Warn("ERROR updating time", "error", err)
				}
			}

			// Calculate each value
			for idx, f := range s.signal.Fields {
				v, err := f.GetValue(parameters)
				if err != nil {
					v = float64(0) // TODO!!!! better error support!!!
				}
				name := f.GetConfig().Name
				parameters[name] = v
				s.current.Values[name] = v

				s.frame.Fields[idx+1].Set(0, v)
			}

			msg.Measurements[0] = s.current

			bytes, err := json.Marshal(&msg)
			if err != nil {
				backend.Logger.Warn("unable to marshal line", "error", err)
				continue
			}
			err = sender.Send(&backend.StreamPacket{
				Payload: bytes,
			})
			if err != nil {
				backend.Logger.Warn("unable to send data", "error", err)
				continue
			}
		}
	}
}

func (s *SignalStreamer) Frames() (data.Frames, error) {
	return data.Frames{s.frame}, nil
}
