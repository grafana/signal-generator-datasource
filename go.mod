module github.com/grafana/signal-generator-datasource

go 1.16

// replace github.com/grafana/grafana-edge-app => ../grafana-edge-app

// replace github.com/grafana/grafana-plugin-sdk-go => ../../more/grafana-plugin-sdk-go

require (
	github.com/Knetic/govaluate v3.0.1-0.20171022003610-9aa49832a739+incompatible
	github.com/apache/arrow/go/arrow v0.0.0-20210304100011-f55f657914af
	github.com/grafana/grafana-edge-app v0.0.0-20210309010455-fc89dec2e783
	github.com/grafana/grafana-plugin-sdk-go v0.91.1-0.20210406164418-35999a631d39
	github.com/stretchr/testify v1.7.0
)
