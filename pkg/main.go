package main

import (
	"os"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/experimental"
	"github.com/grafana/signal-generator-datasource/pkg/plugin"
)

func main() {
	err := experimental.DoGRPC("signal-generator-datasource", plugin.GetDatasourceServeOpts())

	// Log any error if we could start the plugin.
	if err != nil {
		backend.Logger.Error(err.Error())
		os.Exit(1)
	}
}
