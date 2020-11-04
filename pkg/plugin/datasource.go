package plugin

import (
	"context"
	"fmt"
	"time"

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
	case models.QueryTypeEasings:
		return ds.doEasing(ctx, query)
	}
	return backend.DataResponse{
		Error: fmt.Errorf("unsupported query: %s", query.QueryType),
	}
}

func (ds *Datasource) doEasing(ctx context.Context, query *models.SignalQuery) backend.DataResponse {

	total := query.TimeRange.To.Sub(query.TimeRange.From)
	count := int(query.MaxDataPoints - 1)
	if count < 1 {
		count = 1
	}
	interval := total / time.Duration(count)

	time := data.NewFieldFromFieldType(data.FieldTypeTime, count+1)
	time.Name = "Time"
	frame := data.NewFrame("", time)

	ease := make([]waves.EaseFunc, 0)
	for key, f := range waves.EaseFunctions {
		ease = append(ease, f)

		val := data.NewFieldFromFieldType(data.FieldTypeFloat64, count+1)
		val.Name = key
		frame.Fields = append(frame.Fields, val)
	}

	t := query.TimeRange.From
	for i := 0; i <= count; i++ {
		p := float64(i) / float64(count)

		for idx, f := range ease {
			v := f(p)
			frame.Fields[idx+1].Set(i, v)
		}

		time.Set(i, t)
		t = t.Add(interval)
	}

	return backend.DataResponse{
		Frames: data.Frames{frame},
	}
}
