package waves

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/grafana/grafana-plugin-sdk-go/data"
)

type SignalArgs struct {
	Name      string            `json:"name,omitempty"`
	Component []WaveformArgs    `json:"component,omitempty"`
	Config    *data.FieldConfig `json:"config,omitempty"`
	Labels    data.Labels       `json:"labels,omitempty"`
}

type SignalGen struct {
	wave SignalArgs
	comp []WaveformFunc
	args []*WaveformArgs
}

func NewSignalGen(wave SignalArgs, timeRange time.Duration) (*SignalGen, error) {
	count := len(wave.Component)
	comp := make([]WaveformFunc, count)
	args := make([]*WaveformArgs, count)
	for i := 0; i < count; i++ {
		c := wave.Component[i]

		f, ok := WaveformFunctions[c.Type]
		if !ok {
			return nil, fmt.Errorf("invalid waveform type: %s", c.Type)
		}
		comp[i] = f
		args[i] = &c

		// Normalize the period args
		if strings.HasPrefix(c.Period, "range/") {
			f, err := strconv.ParseFloat(c.Period[6:], 64)
			if err != nil {
				return nil, fmt.Errorf("error reading wave period")
			}
			r := timeRange.Seconds() / f
			c.PeriodSec = r
		} else if c.Period != "" {
			r, err := time.ParseDuration(c.Period)
			if err != nil {
				return nil, fmt.Errorf("error reading wave period")
			}
			c.PeriodSec = r.Seconds()
		}
	}

	return &SignalGen{
		args: args,
		comp: comp,
	}, nil
}

// Get the value for a
func (s *SignalGen) GetValue(t time.Time) float64 {
	v := float64(0)
	for i, f := range s.comp {
		v += f(t, s.args[i])
	}
	return v
}

func (s *SignalGen) GetField(timeField *data.Field) *data.Field {
	count := timeField.Len()
	val := data.NewFieldFromFieldType(data.FieldTypeFloat64, count)
	val.Name = s.wave.Name
	val.Config = s.wave.Config
	val.Labels = s.wave.Labels

	for i := 0; i < count; i++ {
		t := timeField.At(i).(time.Time)
		val.Set(i, s.GetValue(t))
	}
	return val
}
