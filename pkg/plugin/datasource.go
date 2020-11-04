package plugin

import (
	"context"
	"fmt"
	"time"

	"github.com/gobwas/glob"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/grafana/signal-generator-datasource/pkg/models"
	"github.com/grafana/signal-generator-datasource/pkg/waves"
)

type Datasource struct {
	settings *models.DatasurceSettings
}

func NewDatasource(settings *models.DatasurceSettings) *Datasource {
	return &Datasource{
		settings: settings,
	}
}

func (ds *Datasource) HealthCheck(ctx context.Context, req *backend.CheckHealthRequest) (*backend.CheckHealthResult, error) {
	return &backend.CheckHealthResult{
		Status:  backend.HealthStatusOk,
		Message: "OK",
	}, nil
}

func (ds *Datasource) QueryData(ctx context.Context, req *backend.QueryDataRequest) (*backend.QueryDataResponse, error) {
	res := backend.NewQueryDataResponse()
	for idx := range req.Queries {
		v := &req.Queries[idx]
		q, err := models.GetSignalQuery(v)
		if err != nil {
			res.Responses[v.RefID] = backend.DataResponse{
				Error: err,
			}
		} else {
			res.Responses[v.RefID] = ds.doQuery(ctx, q)
		}
	}
	return res, nil
}

func (ds *Datasource) doQuery(ctx context.Context, query *models.SignalQuery) backend.DataResponse {
	switch query.QueryType {
	case models.QueryTypeAWG:
		return ds.doAWG(ctx, query)
	case models.QueryTypeEasings:
		return ds.doEasing(ctx, query)
	}
	return backend.DataResponse{
		Error: fmt.Errorf("unsupported query: %s", query.QueryType),
	}
}

func makeTimeAndPercent(query *models.SignalQuery) (*data.Field, *data.Field) {
	total := query.TimeRange.To.Sub(query.TimeRange.From)
	count := int(query.MaxDataPoints - 1)
	if count < 1 {
		count = 1
	}
	interval := total / time.Duration(count)

	time := data.NewFieldFromFieldType(data.FieldTypeTime, count+1)
	time.Name = "Time"

	percent := data.NewFieldFromFieldType(data.FieldTypeFloat64, count+1)
	percent.Name = "Percent"

	t := query.TimeRange.From
	for i := 0; i <= count; i++ {
		p := float64(i) / float64(count)
		percent.Set(i, p)
		time.Set(i, t)
		t = t.Add(interval)
	}

	return time, percent
}

func (ds *Datasource) doEasing(ctx context.Context, query *models.SignalQuery) (dr backend.DataResponse) {
	if query.Ease == "" {
		query.Ease = "*"
	}

	g, err := glob.Compile(query.Ease)
	if err != nil {
		dr.Error = err
		return
	}

	time, percent := makeTimeAndPercent(query)
	frame := data.NewFrame("", time)
	count := time.Len()

	ease := make([]waves.EaseFunc, 0)
	for key, f := range waves.EaseFunctions {
		if g.Match(key) {
			ease = append(ease, f)

			val := data.NewFieldFromFieldType(data.FieldTypeFloat64, count)
			val.Name = key
			frame.Fields = append(frame.Fields, val)
		}
	}

	for i := 0; i < count; i++ {
		p, _ := percent.FloatAt(i)
		for idx, f := range ease {
			v := f(p)
			frame.Fields[idx+1].Set(i, v)
		}
	}

	dr.Frames = data.Frames{frame}
	return
}

func (ds *Datasource) doAWG(ctx context.Context, query *models.SignalQuery) (dr backend.DataResponse) {
	if len(query.Wave) < 1 {
		query.Wave = make([]waves.WaveformArgs, 1)
		query.Wave[0] = waves.WaveformArgs{
			PeriodSec: 30,
			Amplitude: 1,
			Type:      "Sin",
		}
		backend.Logger.Info("adding default wave", "wave", query.Wave)
	}
	backend.Logger.Info("AWG", "wave", query.Wave)

	wave := make([]waves.WaveformFunc, len(query.Wave))
	for i, w := range query.Wave {
		f, ok := waves.WaveformFunctions[w.Type]
		if !ok {
			dr.Error = fmt.Errorf("unknown waveform: %s", w.Type)
			return
		}
		wave[i] = f

		backend.Logger.Info("RUN", "wave", w.PeriodSec, "www", w.Amplitude)
	}

	timef, val := makeTimeAndPercent(query)
	frame := data.NewFrame("", timef, val)
	count := timef.Len()
	val.Name = "Value"

	for i := 0; i < count; i++ {
		t := timef.At(i).(time.Time)
		v := float64(0)
		for j, w := range wave {
			args := query.Wave[j]
			v += w(t, &args)
		}
		val.Set(i, v) // the calculated value
	}

	dr.Frames = data.Frames{frame}
	return
}
