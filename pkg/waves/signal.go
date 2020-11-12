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

type SignalGenX struct {
	args models.SignalConfig
	expr []*govaluate.EvaluableExpression
}

func NewSignalGenerator(args models.SignalConfig) (*SignalGenX, error) {
	count := len(args.Fields)
	expr := make([]*govaluate.EvaluableExpression, count)
	for i, field := range args.Fields {
		if len(field.Name) < 1 {
			return nil, fmt.Errorf("invalid name for field %d", i)
		}
		ex, err := govaluate.NewEvaluableExpressionWithFunctions(field.Expr, WaveformFunctions)
		if err != nil {
			return nil, err
		}
		expr[i] = ex
	}

	return &SignalGenX{
		args: args,
		expr: expr,
	}, nil
}

func DoSignalQuery(query *models.SignalQuery) (*data.Frame, error) {
	if len(query.Signal.Fields) < 1 {
		f0 := models.SignalField{
			Name: "Hello",
			Expr: "x",
		}
		f1 := models.SignalField{
			Name: "Hello",
			Expr: "Sine(x)/x",
		}

		query.Period = "10s"
		query.Signal = models.SignalConfig{
			Name:   "test",
			Fields: []models.SignalField{f0, f1},
		}
	}

	gen, err := NewSignalGenerator(query.Signal)
	if err != nil {
		return nil, err
	}

	input, err := MakeInputFields(query)
	if err != nil {
		return nil, err
	}

	fields, err := gen.Calculate(input)
	if err != nil {
		return nil, err
	}

	frame := data.NewFrame(query.Signal.Name, append([]*data.Field{input[0]}, fields...)...)
	return frame, nil
}

func (s *SignalGenX) Calculate(input []*data.Field) ([]*data.Field, error) {
	fieldCount := len(s.expr)
	rowCount := input[0].Len()
	fields := make([]*data.Field, fieldCount)

	// Setup the fields
	for i := 0; i < fieldCount; i++ {
		fields[i] = data.NewFieldFromFieldType(data.FieldTypeFloat64, rowCount)
		fields[i].Name = s.args.Fields[i].Name
		fields[i].Config = s.args.Fields[i].Config
		fields[i].Labels = s.args.Fields[i].Labels
	}

	parameters := make(map[string]interface{}, fieldCount+4)
	parameters["PI"] = math.Pi

	for row := 0; row < rowCount; row++ {
		for _, field := range input {
			parameters[field.Name] = field.At(row)
		}

		for i, ex := range s.expr {
			v, err := ex.Evaluate(parameters)
			if err != nil {
				v = nil
			}
			parameters[s.args.Fields[i].Name] = v
			fields[i].Set(row, v)
		}
	}

	return fields, nil
}

func MakeInputFields(query *models.SignalQuery) ([]*data.Field, error) {
	period := 0.0

	// Normalize the period args
	if strings.HasPrefix(query.Period, "range/") {
		f, err := strconv.ParseFloat(query.Period[6:], 64)
		if err != nil {
			return nil, fmt.Errorf("invalid range: %s", err.Error())
		}

		timeRange := query.TimeRange.From.Sub(query.TimeRange.To)
		period = timeRange.Seconds() / f
	} else if query.Period != "" {
		r, err := time.ParseDuration(query.Period)
		if err != nil {
			return nil, fmt.Errorf("invalid period: %s", err.Error())
		}
		period = r.Seconds()
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

	return data.Fields{time, percent, x}, nil
}
