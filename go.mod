module github.com/grafana/signal-generator-datasource

go 1.16

// replace github.com/grafana/grafana-plugin-sdk-go => ../../more/grafana-plugin-sdk-go

require (
	github.com/Knetic/govaluate v3.0.1-0.20171022003610-9aa49832a739+incompatible
	github.com/apache/arrow/go/arrow v0.0.0-20210503194501-6c591d023b9e // indirect
	github.com/cortexproject/cortex v1.8.1
	github.com/gorilla/websocket v1.4.2
	github.com/grafana/grafana-edge-app v0.0.0-20210426175502-bd9fe7a6d797
	github.com/grafana/grafana-plugin-sdk-go v0.92.0
	github.com/stretchr/testify v1.7.0
)
