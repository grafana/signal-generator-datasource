package plugin

import (
	"context"
	"math"
	"time"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/grafana/signal-generator-datasource/pkg/waves"
)

// DatasourceHandler is the plugin entrypoint and implements all of the necessary handler functions for dataqueries, healthchecks, and resources.
type SignalStreamer struct {
	interval time.Duration
	signal   *waves.SignalGen
	frame    *data.Frame
	init     time.Time // TODO, periodically kill non
}

// func NewSignalStreamerFromConfig(extcfg *capture.CaptureSetConfig) (*SignalStreamer, error) {
// 	cfg := models.SignalConfig{
// 		Time: models.TimeFieldConfig{
// 			Period: "5s",
// 		},
// 		Fields: []models.ExpressionConfig{},
// 	}

// 	interval := time.Second * 2 // 2s
// 	if extcfg.Interval != "" {
// 		d, err := time.ParseDuration(extcfg.Interval)
// 		if err == nil {
// 			interval = d
// 		}
// 	}

// 	for idx := range extcfg.Input {
// 		tag := extcfg.Input[idx]
// 		if tag.Path == "time" {
// 			// TODO... configure the time period
// 			continue
// 		}
// 		name := tag.Name
// 		if len(name) < 1 {
// 			name = tag.Path
// 		}

// 		if len(name) < 1 {
// 			return nil, fmt.Errorf("invalid field name for tag: %v", tag)
// 		}

// 		if len(tag.Path) > 1 {
// 			tag.Config.Path = tag.Path
// 		}

// 		if tag.Value == nil || tag.Value == "" {
// 			return nil, fmt.Errorf("missing value for field: %s", tag.Path)
// 		}

// 		ft := tag.FieldType
// 		if ft == data.FieldTypeUnknown {
// 			ft = data.FieldTypeFloat64 // the default
// 		}

// 		cfg.Fields = append(cfg.Fields, models.ExpressionConfig{
// 			BaseSignalField: models.BaseSignalField{
// 				Name:   name,
// 				Config: &tag.Config,
// 				Labels: tag.Labels,
// 			},
// 			Expr:     fmt.Sprintf("%v", tag.Value),
// 			DataType: ft,
// 		})
// 	}

// 	gen, err := waves.NewSignalGenerator(cfg)
// 	if err != nil {
// 		return nil, err
// 	}

// 	rowCount := 1
// 	fields := make([]*data.Field, len(gen.Fields)+1)
// 	fields[0] = data.NewFieldFromFieldType(data.FieldTypeTime, rowCount)
// 	fields[0].Name = "Time"
// 	for i, f := range gen.Fields {
// 		cfg := f.GetConfig()
// 		fields[i+1] = data.NewFieldFromFieldType(cfg.DataType, rowCount)
// 		fields[i+1].Name = cfg.Name
// 		fields[i+1].Config = cfg.Config
// 		fields[i+1].Labels = cfg.Labels
// 	}

// 	frame := data.NewFrame(extcfg.Name, fields...)
// 	return &SignalStreamer{
// 		signal:   gen,
// 		frame:    frame,
// 		interval: interval,
// 		init:     time.Now(),
// 	}, nil
// }

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

		s.frame.Fields[idx+1].Set(0, v)
	}
	return nil
}

func (s *SignalStreamer) doStream(ctx context.Context, sender backend.StreamPacketSender) error {
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	paramCount := len(s.signal.Fields) + 4
	parameters := make(map[string]interface{}, paramCount)
	parameters["PI"] = math.Pi

	backend.Logger.Info("start streaming")

	// local copy
	fields := make([]*data.Field, len(s.frame.Fields))
	for idx, f := range s.frame.Fields {
		fields[idx] = data.NewFieldFromFieldType(f.Type(), 1)
	}
	frame := data.NewFrame("", fields...)

	for {
		select {
		case <-ctx.Done():
			backend.Logger.Info("stop streaming (context canceled)")
			return nil
		case t := <-ticker.C:
			frame.Fields[0].Set(0, t)

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

				frame.Fields[idx+1].Set(0, v)
			}

			bytes, err := data.FrameToJSON(frame, false, true)
			if err != nil {
				backend.Logger.Warn("error marshaling frame to JSON", "error", err)
				continue
			}
			packet := &backend.StreamPacket{
				Data: bytes,
			}

			err = sender.Send(packet)
			if err != nil {
				backend.Logger.Warn("Unable to send data", "error", err)
				// Broken stream? Return an error, stream should be re-established soon.
				return err
			}
		}
	}
}

func (s *SignalStreamer) Frames() (data.Frames, error) {
	return data.Frames{s.frame}, nil
}
