package waves

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/Knetic/govaluate"
	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/grafana/signal-generator-datasource/pkg/models"
)

type SignalEnv = map[string]interface{}

type SignalField interface {
	GetConfig() *models.ExpressionConfig
	GetValue(env map[string]interface{}) (interface{}, error)
	SetValue(val interface{}) error
}

type InputField interface {
	GetConfig() *models.TimeFieldConfig
	GetValues(query *models.SignalQuery) (rsp []*data.Field, env []*data.Field, err error)
	UpdateEnv(time *time.Time, env map[string]interface{}) error
}

type signalFieldWithEval struct {
	config *models.ExpressionConfig
	expr   *govaluate.EvaluableExpression
}

func (f *signalFieldWithEval) GetConfig() *models.ExpressionConfig {
	return f.config
}

func (f *signalFieldWithEval) GetValue(env map[string]interface{}) (interface{}, error) {
	return f.expr.Evaluate(env)
}

func (f *signalFieldWithEval) SetValue(val interface{}) error {
	expr := fmt.Sprintf("%v", val)
	ex, err := govaluate.NewEvaluableExpressionWithFunctions(expr, WaveformFunctions)
	if err != nil {
		return err
	}

	f.config.Expr = expr
	f.expr = ex
	return nil
}

func NewEvalSignalField(config *models.ExpressionConfig) (SignalField, error) {
	if len(config.Name) < 1 {
		return nil, fmt.Errorf("invalid name for field %v", config)
	}
	ex, err := govaluate.NewEvaluableExpressionWithFunctions(config.Expr, WaveformFunctions)
	if err != nil {
		return nil, err
	}

	// For now all expressions are numbers
	config.DataType = data.FieldTypeFloat64
	return &signalFieldWithEval{
		config: config,
		expr:   ex,
	}, nil
}

type timeInputField struct {
	config   *models.TimeFieldConfig
	period   float64 // seconds
	rangeDiv float64 // if the period depends on the range
}

func (f *timeInputField) GetConfig() *models.TimeFieldConfig {
	return f.config
}

func (f *timeInputField) GetValues(query *models.SignalQuery) ([]*data.Field, []*data.Field, error) {
	period := f.period
	if f.rangeDiv > 0 {
		timeRange := query.TimeRange.From.Sub(query.TimeRange.To)
		period = timeRange.Seconds() / f.rangeDiv
	}

	total := query.TimeRange.To.Sub(query.TimeRange.From)
	count := int(query.MaxDataPoints - 1)
	if count < 1 {
		count = 1
	}
	interval := total / time.Duration(count)

	time := data.NewFieldFromFieldType(data.FieldTypeTime, count+1)
	time.Name = "time"

	percent := data.NewFieldFromFieldType(data.FieldTypeFloat64, count+1)
	percent.Name = "p"

	x := data.NewFieldFromFieldType(data.FieldTypeFloat64, count+1)
	x.Name = "x"

	rad := 0.0
	t := query.TimeRange.From
	for i := 0; i <= count; i++ {
		p := float64(i) / float64(count)

		if period > 0 {
			ms := t.UnixNano() % int64(period*1000000000)
			rad = ((float64(ms) / (period * 1000000000)) * 2 * math.Pi) // 0 >> 2Pi
		}

		percent.Set(i, p)
		time.Set(i, t)
		x.Set(i, rad)
		t = t.Add(interval)
	}

	return []*data.Field{time}, []*data.Field{x, percent}, nil
}

func (f *timeInputField) UpdateEnv(t *time.Time, env map[string]interface{}) error {
	ms := int64(0)
	rad := 0.0

	if f.period > 0 {
		ms = t.UnixNano() % int64(f.period*1000000000)
		rad = ((float64(ms) / (f.period * 1000000000)) * 2 * math.Pi) // 0 >> 2Pi
	}
	env["x"] = rad
	return nil
}

func NewTimeInputField(config *models.TimeFieldConfig) (InputField, error) {
	period := 0.0
	rangeDiv := 0.0

	if strings.HasPrefix(config.Period, "range/") {
		f, err := strconv.ParseFloat(config.Period[6:], 64)
		if err != nil {
			return nil, fmt.Errorf("invalid range: %s", err.Error())
		}
		rangeDiv = f
	} else if config.Period != "" {
		r, err := time.ParseDuration(config.Period)
		if err != nil {
			return nil, fmt.Errorf("invalid period: %s", err.Error())
		}
		period = r.Seconds()
	}

	return &timeInputField{
		config:   config,
		period:   period,
		rangeDiv: rangeDiv,
	}, nil
}

type SignalGen struct {
	Inputs []InputField
	Fields []SignalField
}

func (s *SignalGen) UpdateValues(props map[string]interface{}) error {
	byName := make(map[string]int)
	for idx, f := range s.Fields {
		cfg := f.GetConfig()
		byName[cfg.Name] = idx
		if len(cfg.Config.Path) > 0 {
			byName[cfg.Config.Path] = idx
		}
	}

	// First update values
	for k := range props {
		_, ok := byName[k]
		if !ok {
			return fmt.Errorf("can not find field: '%s' // %v", k, byName)
		}
	}

	// Now update
	for k, v := range props {
		idx := byName[k]
		err := s.Fields[idx].SetValue(v)
		if err != nil {
			return err
		}
	}

	return nil
}

func NewSignalGenerator(args models.SignalConfig) (*SignalGen, error) {
	gen := &SignalGen{
		Inputs: make([]InputField, 1),
		Fields: make([]SignalField, len(args.Fields)),
	}

	t, err := NewTimeInputField(&args.Time)
	if err != nil {
		return nil, err
	}
	gen.Inputs[0] = t

	for i := range args.Fields {
		f, err := NewEvalSignalField(&args.Fields[i])
		if err != nil {
			return nil, err
		}
		gen.Fields[i] = f
	}

	return gen, nil
}

func DoSignalQuery(query *models.SignalQuery) (*data.Frame, *SignalGen, error) {
	gen, err := NewSignalGenerator(query.Signal)
	if err != nil {
		return nil, nil, err
	}

	// Setup the initial fields
	outfields := make([]*data.Field, 0)
	envFields := make([]*data.Field, 0)
	for _, i := range gen.Inputs {
		rsp, env, err := i.GetValues(query)
		if err != nil {
			return nil, nil, err
		}

		if rsp != nil {
			outfields = append(outfields, rsp...)
		}
		if env != nil {
			envFields = append(envFields, env...)
		}
	}

	rowCount := outfields[0].Len()
	fields := make([]*data.Field, len(gen.Fields))
	for i, f := range gen.Fields {
		cfg := f.GetConfig()
		fields[i] = data.NewFieldFromFieldType(data.FieldTypeFloat64, rowCount)
		fields[i].Name = cfg.Name
		fields[i].Config = cfg.Config
		fields[i].Labels = cfg.Labels
	}

	paramCount := len(gen.Fields) + len(envFields) + 4
	parameters := make(map[string]interface{}, paramCount)
	parameters["PI"] = math.Pi

	for row := 0; row < rowCount; row++ {
		for _, field := range envFields {
			parameters[field.Name] = field.At(row)
		}

		for i, ex := range gen.Fields {
			v, err := ex.GetValue(parameters)
			if err != nil {
				v = nil
			}
			parameters[fields[i].Name] = v
			fields[i].Set(row, v)
		}
	}

	outfields = append(outfields, fields...)
	frame := data.NewFrame(query.Signal.Name, outfields...)
	return frame, gen, nil
}
