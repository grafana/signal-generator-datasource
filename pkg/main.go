package main

import (
	"context"
	"errors"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/grafana/grafana-plugin-sdk-go/backend/datasource"
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
	"github.com/grafana/waveform-datasource/pkg/awg"
	"github.com/grafana/waveform-datasource/pkg/sock"
)

func main() {
	err := run()
	if err != nil {
		log.DefaultLogger.Error(err.Error())
	}

	err = datasource.Serve(awg.CreateDatasourcePlugin())

	// Log any error if we could start the plugin.
	if err != nil {
		log.DefaultLogger.Error(err.Error())
		os.Exit(1)
	}
}

// run starts a http.Server for the passed in address
// with all requests handled by echoServer.
func run() error {
	if len(os.Args) < 2 {
		return errors.New("please provide an address to listen on as the first argument")
	}

	l, err := net.Listen("tcp", os.Args[1])
	if err != nil {
		return err
	}
	log.DefaultLogger.Debug("listening on http://%v", l.Addr())

	s := &http.Server{
		Handler: sock.EchoServer{
			Logf: log.DefaultLogger.Debug,
		},
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 10,
	}
	errc := make(chan error, 1)
	go func() {
		errc <- s.Serve(l)
	}()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt)
	select {
	case err := <-errc:
		log.DefaultLogger.Debug("failed to serve: %v", err)
	case sig := <-sigs:
		log.DefaultLogger.Debug("terminating: %v", sig)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	return s.Shutdown(ctx)
}
