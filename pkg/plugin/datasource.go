package plugin

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/instancemgmt"
	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/grafana/grafana-plugin-sdk-go/experimental/actions"
	"github.com/grafana/signal-generator-datasource/pkg/models"
	"github.com/grafana/signal-generator-datasource/pkg/waves"
)

type Datasource struct {
	mu            sync.RWMutex
	settings      *models.DatasourceSettings
	streams       map[string]*SignalStreamer
	channelPrefix string
	closeCh       chan struct{}
}

// Make sure SampleDatasource implements required interfaces.
// This is important to do since otherwise we will only get a
// not implemented error response from plugin in runtime.
var (
	_ backend.QueryDataHandler      = (*Datasource)(nil)
	_ backend.CheckHealthHandler    = (*Datasource)(nil)
	_ backend.StreamHandler         = (*Datasource)(nil)
	_ backend.CallResourceHandler   = (*Datasource)(nil)
	_ instancemgmt.InstanceDisposer = (*Datasource)(nil)
)

// NewMQTTDatasource creates a new datasource instance.
func NewDatasource(s backend.DataSourceInstanceSettings) (instancemgmt.Instance, error) {
	settings, err := models.GetDatasourceSettings(s)
	if err != nil {
		return nil, err
	}
	return &Datasource{
		settings:      settings,
		streams:       make(map[string]*SignalStreamer),
		channelPrefix: fmt.Sprintf("ds/%s/", s.UID),
		closeCh:       make(chan struct{}),
	}, nil
}

// Dispose here tells plugin SDK that plugin wants to clean up resources
// when a new instance created. As soon as datasource settings change detected
// by SDK old datasource instance will be disposed and a new one will be created
// using NewSampleDatasource factory function.
func (ds *Datasource) Dispose() {
	close(ds.closeCh)
}

func (ds *Datasource) ExecuteAction(ctx context.Context, cmd actions.ActionCommand) actions.ActionResponse {
	ds.mu.RLock()
	defer ds.mu.RUnlock()
	s, ok := ds.streams[cmd.Path]
	if !ok {
		keys := make([]string, 0, len(ds.streams))
		for k := range ds.streams {
			keys = append(keys, k)
		}

		return actions.ActionResponse{
			Code:  http.StatusBadRequest,
			Error: fmt.Sprintf("'%s' not found in: %v", cmd.Path, keys),
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
		State: s.frame,
	}
}

func (ds *Datasource) CallResource(ctx context.Context, req *backend.CallResourceRequest, sender backend.CallResourceResponseSender) error {
	if req.Path == "action" {
		//return actions.DoActionCommand(ctx, req, ds, sender)
		return sender.Send(&backend.CallResourceResponse{
			Status: http.StatusNotFound,
			Body:   []byte("not implemented yet"),
		})
	}
	return sender.Send(&backend.CallResourceResponse{
		Status: http.StatusOK,
		Body:   []byte("OK"),
	})
}

func (ds *Datasource) CheckHealth(ctx context.Context, req *backend.CheckHealthRequest) (*backend.CheckHealthResult, error) {
	ds.mu.RLock()
	defer ds.mu.RUnlock()

	streamCount := 0
	fieldCount := 0

	for _, s := range ds.streams {
		streamCount++
		fieldCount += len(s.frame.Fields)
	}

	backend.Logger.Error("datasource ID", "id", req.PluginContext.DataSourceInstanceSettings.ID)

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

func (ds *Datasource) initStream(query *models.SignalQuery, gen *waves.SignalGen, frame *data.Frame) string {
	key := models.GetStreamKey(query)

	ds.mu.Lock()
	defer ds.mu.Unlock()
	_, ok := ds.streams[key]
	if !ok {
		ds.streams[key] = &SignalStreamer{
			interval: query.Interval,
			signal:   gen,
			frame:    frame.EmptyCopy(),
			init:     time.Now(),
		}
	}
	return key
}

func (ds *Datasource) doAWG(_ context.Context, query *models.SignalQuery) (dr backend.DataResponse) {
	frame, gen, err := waves.DoSignalQuery(query)
	if err != nil {
		dr.Error = err
		return
	}
	if query.Stream {
		key := ds.initStream(query, gen, frame)
		frame.SetMeta(&data.FrameMeta{
			Channel: ds.channelPrefix + key, // ds/${id}/${key}
		})
	}
	dr.Frames = data.Frames{frame}
	dr.Error = err
	return
}

func (ds *Datasource) SubscribeStream(_ context.Context, req *backend.SubscribeStreamRequest) (*backend.SubscribeStreamResponse, error) {
	ds.mu.RLock()
	defer ds.mu.RUnlock()

	s, ok := ds.streams[req.Path]
	if s == nil || !ok {
		return &backend.SubscribeStreamResponse{
			Status: backend.SubscribeStreamStatusNotFound,
		}, nil
	}

	bytes, err := data.FrameToJSON(s.frame, true, false) // only schema
	if err != nil {
		return nil, err
	}

	return &backend.SubscribeStreamResponse{
		Status: backend.SubscribeStreamStatusOK,
		Data:   bytes, // just the schema
	}, nil
}

func (ds *Datasource) RunStream(ctx context.Context, req *backend.RunStreamRequest, sender backend.StreamPacketSender) error {
	ds.mu.RLock()
	s, ok := ds.streams[req.Path]
	if s == nil || !ok {
		ds.mu.RUnlock()
		// Return nil to stop RunStream till next subscriber. Any error here
		// will result into RunStream re-establishment.
		return nil
	}
	ds.mu.RUnlock()

	// When the stream is done, remove it.
	defer func() {
		ds.mu.Lock()
		defer ds.mu.Unlock()
		delete(ds.streams, req.Path)
	}()

	return s.doStream(ctx, sender)
}

func (ds *Datasource) PublishStream(ctx context.Context, req *backend.PublishStreamRequest) (*backend.PublishStreamResponse, error) {
	return &backend.PublishStreamResponse{
		Status: backend.PublishStreamStatusPermissionDenied, // ?? Unsupported
	}, nil
}
