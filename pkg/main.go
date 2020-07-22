package main

import (
	"math/rand"
	"os"
	"time"

	"github.com/grafana/grafana-plugin-sdk-go/backend/datasource"
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
	"github.com/grafana/waveform-datasource/pkg/awg"
	"github.com/grafana/waveform-datasource/pkg/broker"
	"github.com/grafana/waveform-datasource/pkg/models"
	"github.com/grafana/waveform-datasource/pkg/serializers"
)

func main() {
	// Background thread
	log.DefaultLogger.Info("starting http server")
	b := &broker.GrafanaBroker{}
	go b.ListenAndServe(":3007")
	go streamSignal(b, "example")

	log.DefaultLogger.Info("starting grpc server")
	err := datasource.Serve(awg.CreateDatasourcePlugin())

	// Log any error if we could start the plugin.
	if err != nil {
		log.DefaultLogger.Error(err.Error())
		os.Exit(1)
	}
}

// write to a stream....
func streamSignal(broker *broker.GrafanaBroker, channel string) {
	speed := 1000 / 20 // 20 hz
	spread := 50.0

	walker := rand.Float64() * 100
	ticker := time.NewTicker(time.Duration(speed) * time.Millisecond)

	line := models.InfluxLine{
		Name:   "simple",
		Fields: make(map[string]interface{}),
		Tags:   make(map[string]string),
	}

	s, _ := serializers.NewSerializer(&serializers.Config{
		DataFormat:     "json",
		TimestampUnits: time.Duration(1) * time.Millisecond,
	})

	for t := range ticker.C {
		delta := rand.Float64() - 0.5
		walker += delta

		line.Timestamp = t
		line.Fields["value"] = walker
		line.Fields["min"] = walker - ((rand.Float64() * spread) + 0.01)
		line.Fields["max"] = walker + ((rand.Float64() * spread) + 0.01)

		b, _ := s.SerializeBatch([]*models.InfluxLine{&line}) //json.Marshal(line)

		broker.Publish(channel, b)
	}
}
