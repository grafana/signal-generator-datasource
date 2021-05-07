package main

import (
	"fmt"
	"time"

	"github.com/grafana/signal-generator-datasource/pkg/replay"
)

func main() {
	fpath := "/home/ryan/Downloads/influx-sample-data.log"
	url := "ws://localhost:3000/api/live/push/telegraf?gf_live_frame_format=labels_column"
	key := "eyJrIjoicExKYjlEN29yQmlrMEg4YmtodlRFSjN6R0FOUjRLMEQiLCJuIjoicHVibGlzaCIsImlkIjoxfQ=="

	ws := replay.NewWebSocket(url)
	ws.Headers = map[string]string{
		"Authorization": "Bearer " + key,
	}
	err := ws.Connect()
	if err != nil {
		panic(err)
	}

	interval := 50 * time.Millisecond
	count := replay.ReplayInfluxLog(fpath, interval, ws.Write)
	fmt.Printf("wrote: %d lines.\n", count)
}
