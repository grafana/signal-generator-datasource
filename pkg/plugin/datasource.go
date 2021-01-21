package plugin

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/grafana/grafana-edge-app/pkg/actions"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/grafana/signal-generator-datasource/pkg/models"
	"github.com/grafana/signal-generator-datasource/pkg/waves"
)

type Datasource struct {
	settings *models.DatasurceSettings
	streamer *SignalStreamer
}

func NewDatasource(settings *models.DatasurceSettings) *Datasource {
	return &Datasource{
		settings: settings,
		streamer: &SignalStreamer{
			speedMillis: 50, // 20hz
		},
	}
}

func (ds *Datasource) CallResource(ctx context.Context, req *backend.CallResourceRequest, sender backend.CallResourceResponseSender) error {

	cmd := &actions.ActionCommand{}
	if err := json.Unmarshal(req.Body, cmd); err != nil {
		return err
	}

	for _, action := range cmd.Write {
		if action.Path == "stream.start" {
			backend.Logger.Info("START!!!")
			ds.streamer.Start()
		} else if action.Path == "stream.stop" {
			backend.Logger.Info("STOP!!!")
			ds.streamer.Stop()
		} else {
			backend.Logger.Info("???????????????")
		}
	}
	return sender.Send(&backend.CallResourceResponse{
		Status: http.StatusOK,
		Body:   []byte("OK"),
	})
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
		// case models.QueryTypeEasings:
		// 	return ds.doEasing(ctx, query)
	}
	return backend.DataResponse{
		Error: fmt.Errorf("unsupported query: %s", query.QueryType),
	}
}

// func (ds *Datasource) doEasing(ctx context.Context, query *models.SignalQuery) (dr backend.DataResponse) {
// 	if query.Ease == "" {
// 		query.Ease = "*"
// 	}

// 	g, err := glob.Compile(query.Ease)
// 	if err != nil {
// 		dr.Error = err
// 		return
// 	}

// 	input, err := waves.MakeInputFields(query)
// 	if err != nil {
// 		dr.Error = err
// 		return
// 	}
// 	time := input[0]
// 	percent := input[1]

// 	frame := data.NewFrame("", time)
// 	count := time.Len()

// 	ease := make([]waves.EaseFunc, 0)
// 	for key, f := range waves.EaseFunctions {
// 		if g.Match(key) {
// 			ease = append(ease, f)

// 			val := data.NewFieldFromFieldType(data.FieldTypeFloat64, count)
// 			val.Name = key
// 			frame.Fields = append(frame.Fields, val)
// 		}
// 	}

// 	for i := 0; i < count; i++ {
// 		p, _ := percent.FloatAt(i)
// 		for idx, f := range ease {
// 			v := f(p)
// 			frame.Fields[idx+1].Set(i, v)
// 		}
// 	}

// 	dr.Frames = data.Frames{frame}
// 	return
// }

func (ds *Datasource) doAWG(ctx context.Context, query *models.SignalQuery) (dr backend.DataResponse) {
	frame, err := waves.DoSignalQuery(query)
	dr.Frames = data.Frames{frame}
	dr.Error = err
	return
}
