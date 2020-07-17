package main

import (
	"os"

	"github.com/grafana/grafana-plugin-sdk-go/backend/datasource"
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
	"github.com/grafana/waveform-datasource/pkg/awg"
	"github.com/grafana/waveform-datasource/pkg/sock"
)

func main() {
	// Background thread
	log.DefaultLogger.Info("starting http server")
	go sock.RunChatServer("localhost:3003")

	log.DefaultLogger.Info("starting grpc server")
	err := datasource.Serve(awg.CreateDatasourcePlugin())

	// Log any error if we could start the plugin.
	if err != nil {
		log.DefaultLogger.Error(err.Error())
		os.Exit(1)
	}
}
