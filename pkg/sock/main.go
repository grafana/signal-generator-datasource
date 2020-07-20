package sock

import (
	"context"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
	"github.com/grafana/waveform-datasource/pkg/models"
	"github.com/grafana/waveform-datasource/pkg/parsers"
	"github.com/grafana/waveform-datasource/pkg/replay"
	"github.com/grafana/waveform-datasource/pkg/serializers"
)

// RunChatServer runs a chat server
func RunChatServer(address string) {
	err := run(address)
	if err != nil {
		log.DefaultLogger.Error(err.Error())
	}
}

// run initializes the chatServer and then
// starts a http.Server for the passed in address.
func run(address string) error {
	log.DefaultLogger.Debug("running chat server")
	datac := make(chan *models.InfluxLine, 100)

	replay := &replay.Replay{
		Files:      []string{"/Users/stephanie/src/enterprise-plugins/telegraf/plugins/inputs/replay/dev/testfiles/avionics_hvbms_HvBmsData.csv"},
		Iterations: -1,
	}

	parser, err := parsers.NewParser(&parsers.Config{
		DataFormat:         "csv",
		CSVHeaderRowCount:  1,
		CSVTimestampColumn: "time",
		CSVTimestampFormat: "unix_ns",
		CSVTrimSpace:       true,
	})

	serializer, err := serializers.NewSerializer(&serializers.Config{
		DataFormat:     "json",
		TimestampUnits: time.Duration(1) * time.Millisecond,
	})

	replay.SetParser(parser)
	err = replay.Start(datac)

	l, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}
	log.DefaultLogger.Info("listening on", "address", l.Addr())

	cs := NewChatServer()
	s := &http.Server{
		Handler:     cs,
		ReadTimeout: time.Second * 10,
		//	WriteTimeout: time.Second * 10,
	}
	errc := make(chan error, 1)
	go func() {
		errc <- s.Serve(l)
	}()

	if err != nil {
		log.DefaultLogger.Debug("Could not start", err)
		return err
	}

	// Send signal stream to everyone
	//go cs.streamSignalToSocket()
	go cs.streamMetricsToSocket(datac, serializer)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt)
	select {
	case err := <-errc:
		log.DefaultLogger.Info("failed to serve:", err)
	case sig := <-sigs:
		log.DefaultLogger.Info("terminating:", sig)
		replay.Stop()
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	return s.Shutdown(ctx)
}
