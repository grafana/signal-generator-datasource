package awg

import (
	"context"
	"fmt"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/instancemgmt"
	"github.com/grafana/grafana-plugin-sdk-go/backend/resource"
	"github.com/grafana/waveform-datasource/pkg/models"
)

// WaveformDatasource is holds the instance
type WaveformDatasource struct {
	im instancemgmt.InstanceManager
}

func (ds *WaveformDatasource) getInstance(ctx backend.PluginContext) (*instanceSettings, error) {
	s, err := ds.im.Get(ctx)
	if err != nil {
		return nil, err
	}
	return s.(*instanceSettings), nil // ugly cast... but go ¯\_(ツ)_/¯
}

// QueryData handles multiple queries
func (ds *WaveformDatasource) QueryData(ctx context.Context, req *backend.QueryDataRequest) (*backend.QueryDataResponse, error) {
	s, err := ds.getInstance(req.PluginContext)
	if err != nil {
		return nil, err
	}

	// create response struct
	rsp := backend.NewQueryDataResponse()

	// loop over queries and execute them individually.
	for _, q := range req.Queries {
		query, err := models.GetQueryModel(q)
		if err != nil {
			rsp.Responses[q.RefID] = backend.DataResponse{
				Error: err,
			}
		} else {
			rsp.Responses[q.RefID] = backend.DataResponse{
				Error: fmt.Errorf("todo... actually run a query!!! %v/%v", s, query),
			}
		}
	}

	return rsp, nil
}

// CheckHealth handles health checks
func (ds *WaveformDatasource) CheckHealth(ctx context.Context, req *backend.CheckHealthRequest) (*backend.CheckHealthResult, error) {
	_, err := ds.getInstance(req.PluginContext)
	if err != nil {
		return &backend.CheckHealthResult{
			Status:  backend.HealthStatusError,
			Message: err.Error(),
		}, nil
	}

	return &backend.CheckHealthResult{
		Status:  backend.HealthStatusOk,
		Message: fmt.Sprintf("todo... actually check something"),
	}, nil
}

// CallResource HTTP style resource
func (ds *WaveformDatasource) CallResource(tx context.Context, req *backend.CallResourceRequest, sender backend.CallResourceResponseSender) error {
	if req.Path == "hello" {
		return resource.SendPlainText(sender, "world")
	}

	return fmt.Errorf("unknown resource")
}
