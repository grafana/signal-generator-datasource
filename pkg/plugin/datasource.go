package plugin

import (
	"context"
	"fmt"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/signal-generator-datasource/pkg/models"
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
	// switch query.QueryType {
	// case models.QueryTypeListTags:
	// 	return ds.handleListTags(ctx, query)
	// case models.QueryTypeListTagGroups:
	// 	return ds.handleListTagGroups(ctx, query)
	// case models.QueryTypeGetTagValue:
	// 	return ds.handleGetTagValue(ctx, query)
	// case models.QueryTypeGetTagConfig:
	// 	return ds.handleGetTagConfig(ctx, query)
	// }
	return backend.DataResponse{
		Error: fmt.Errorf("unsupported query: %s", query.QueryType),
	}
}
