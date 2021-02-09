package main

import (
	"flag"
	"os"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/datasource"
	"github.com/grafana/signal-generator-datasource/pkg/plugin"
)

type args struct {
	server *bool
}

func main() {
	args := args{
		server: flag.Bool("server", false, "Run server"),
	}
	flag.Parse()

	if *args.server {
		// server.RunServer()
		os.Exit(0)
	}

	backend.SetupPluginEnvironment("signal-generator-datasource")
	err := datasource.Serve(plugin.GetDatasourceServeOpts())

	// Log any error if we could start the plugin.
	if err != nil {
		backend.Logger.Error(err.Error())
		os.Exit(1)
	}
}
