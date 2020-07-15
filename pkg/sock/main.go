package sock

import (
	"context"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
)

// RunChatServer runs a chat server
func RunChatServer() {
	err := run()
	if err != nil {
		log.DefaultLogger.Error(err.Error())
	}
}

// run initializes the chatServer and then
// starts a http.Server for the passed in address.
func run() error {
	l, err := net.Listen("tcp", "localhost:3003")
	if err != nil {
		return err
	}
	log.DefaultLogger.Info("listening on", "address", l.Addr())

	cs := newChatServer()
	s := &http.Server{
		Handler:     cs,
		ReadTimeout: time.Second * 10,
		//	WriteTimeout: time.Second * 10,
	}
	errc := make(chan error, 1)
	go func() {
		errc <- s.Serve(l)
	}()

	// Send signal stream to everyone
	go cs.streamSignalToSocket()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt)
	select {
	case err := <-errc:
		log.DefaultLogger.Info("failed to serve:", err)
	case sig := <-sigs:
		log.DefaultLogger.Info("terminating:", sig)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	return s.Shutdown(ctx)
}
