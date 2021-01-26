package plugin

import (
	"context"
	"fmt"
	"net/http"

	"github.com/grafana/grafana-edge-app/pkg/actions"
	"github.com/grafana/grafana-edge-app/pkg/tags"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/grafana/grafana-plugin-sdk-go/live"
	"github.com/grafana/signal-generator-datasource/pkg/models"
	"github.com/grafana/signal-generator-datasource/pkg/waves"
)

type Datasource struct {
	live     *live.GrafanaLiveClient
	settings *models.DatasurceSettings
	streams  map[string]*SignalStreamer
}

func NewDatasource(settings *models.DatasurceSettings) *Datasource {
	client, _ := live.InitGrafanaLiveClient(live.ConnectionInfo{
		URL: settings.LiveURL,
	})

	// Initalize streams
	streams := make(map[string]*SignalStreamer)
	for _, path := range settings.Capture {
		if client == nil {
			backend.Logger.Error("missing live server connection")
			continue
		}

		cfg, err := tags.LoadCaptureSetConfig(path)
		if err != nil {
			backend.Logger.Error("error loading config", "err", err, "path", path)
		} else if len(cfg.Flags) > 0 {
			backend.Logger.Error("flags not supported", "path", path)
		} else {
			stream, err := NewSignalStreamer(cfg, client)
			if err != nil {
				backend.Logger.Error("error initalizing stream", "err", err, "path", path)
			} else {
				streams[cfg.Name] = stream
			}
		}
	}

	return &Datasource{
		settings: settings,
		live:     client,
		streams:  streams,
	}
}

func (ds *Datasource) ExecuteAction(ctx context.Context, cmd actions.ActionCommand) actions.ActionResponse {
	s, ok := ds.streams[cmd.Path]
	if !ok {
		keys := make([]string, 0, len(ds.streams))
		for k := range ds.streams {
			keys = append(keys, k)
		}

		return actions.ActionResponse{
			Code:  http.StatusBadRequest,
			Error: fmt.Sprintf("stream not found in: %v", keys),
		}
	}

	vmap, ok := cmd.Value.(map[string]interface{})
	if !ok {
		return actions.ActionResponse{
			Code:  http.StatusBadRequest,
			Error: "value must be a map",
		}
	}

	err := s.UpdateValues(vmap)
	if err != nil {
		return actions.ActionResponse{
			Code:  http.StatusBadRequest,
			Error: err.Error(),
		}
	}

	return actions.ActionResponse{
		Code:  http.StatusOK,
		State: s.current,
	}
}

func (ds *Datasource) CallResource(ctx context.Context, req *backend.CallResourceRequest, sender backend.CallResourceResponseSender) error {

	// cmd := &actions.ActionCommand{}
	// if err := json.Unmarshal(req.Body, cmd); err != nil {
	// 	return err
	// }

	// for _, action := range cmd.Write {
	// 	if action.Path == "stream.start" {
	// 		backend.Logger.Info("START!!!")
	// 		ds.streamer.Start()
	// 	} else if action.Path == "stream.stop" {
	// 		backend.Logger.Info("STOP!!!")
	// 		ds.streamer.Stop()
	// 	} else {
	// 		backend.Logger.Info("???????????????")
	// 	}
	// }

	if req.Path == "action" {
		return actions.DoActionCommand(ctx, req, ds, sender)
	}

	return sender.Send(&backend.CallResourceResponse{
		Status: http.StatusOK,
		Body:   []byte("OK"),
	})
}

func (ds *Datasource) HealthCheck(ctx context.Context, req *backend.CheckHealthRequest) (*backend.CheckHealthResult, error) {
	streamCount := 0
	fieldCount := 0

	for _, s := range ds.streams {
		streamCount++
		fieldCount += len(s.frame.Fields)
	}

	return &backend.CheckHealthResult{
		Status:  backend.HealthStatusOk,
		Message: fmt.Sprintf("OK (%d streams, %d fields)", streamCount, fieldCount),
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
	case models.QueryStreams:
		return ds.doStream(ctx, query)
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

func (ds *Datasource) doStream(ctx context.Context, query *models.SignalQuery) (dr backend.DataResponse) {

	if len(query.Stream) > 1 {
		s, ok := ds.streams[query.Stream]
		if !ok {
			dr.Error = fmt.Errorf("unknown stream: %s", query.Stream)
			return
		}

		frames, err := s.Frames()
		dr.Frames = frames
		dr.Error = err
		return
	}

	for _, s := range ds.streams {
		f, _ := s.Frames()
		if f != nil {
			dr.Frames = append(dr.Frames, f...)
		}
	}
	return
}
