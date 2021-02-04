module github.com/grafana/signal-generator-datasource

go 1.14

replace github.com/grafana/grafana-edge-app => ../grafana-edge-app

replace github.com/grafana/grafana-plugin-sdk-go => ../../more/grafana-plugin-sdk-go

require (
	github.com/Knetic/govaluate v3.0.0+incompatible
	github.com/apache/arrow/go/arrow v0.0.0-20210203230224-5e5c2b4890fd
	github.com/grafana/grafana-edge-app v0.0.0-20210121054608-7902bfcb2ff6
	github.com/grafana/grafana-plugin-sdk-go v0.86.1-0.20210204064750-13cca6dae85b
	github.com/mattetti/filebuffer v1.0.1
	github.com/stretchr/testify v1.7.0
)
