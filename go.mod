module github.com/grafana/signal-generator-datasource

go 1.14

// replace github.com/grafana/grafana-edge-app => ../grafana-edge-app

// replace github.com/grafana/grafana-plugin-sdk-go => ../../more/grafana-plugin-sdk-go

require (
	github.com/Knetic/govaluate v3.0.0+incompatible
	github.com/apache/arrow/go/arrow v0.0.0-20210209165900-d25660ed5f57
	github.com/grafana/grafana-edge-app v0.0.0-20210209173145-6692f438f8be
	github.com/grafana/grafana-plugin-sdk-go v0.86.0
	github.com/stretchr/testify v1.7.0
)
