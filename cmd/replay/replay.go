package main

import (
	"fmt"
	"os"
	"time"

	"github.com/grafana/signal-generator-datasource/pkg/replay"
)

// go run cmd/replay/replay.go dev/influx-sample-data.log "ws://localhost:3000/api/live/push/telegraf" "eyJrIjoicExKYjlEN29yQmlrMEg4YmtodlRFSjN6R0FOUjRLMEQiLCJuIjoicHVibGlzaCIsImlkIjoxfQ=="

// to loop forever:
// go run cmd/replay/replay.go dev/influx-sample-data.log "ws://localhost:3000/api/live/push/telegraf" "eyJrIjoicExKYjlEN29yQmlrMEg4YmtodlRFSjN6R0FOUjRLMEQiLCJuIjoicHVibGlzaCIsImlkIjoxfQ==" loop

func main() {
	if len(os.Args) < 3 {
		fmt.Printf("Expected:\n")
		fmt.Printf("%s {path} {url} {key}\n", os.Args[0])
		return
	}

	fpath := os.Args[1]
	url := os.Args[2]
	loop := ""
	if len(os.Args) > 4 {
		loop = os.Args[4]
	}

	// fpath := "/home/ryan/Downloads/influx-sample-data.log"
	// url := "ws://localhost:3000/api/live/push/telegraf?gf_live_frame_format=labels_column"
	// key := "eyJrIjoicExKYjlEN29yQmlrMEg4YmtodlRFSjN6R0FOUjRLMEQiLCJuIjoicHVibGlzaCIsImlkIjoxfQ=="

	ws := replay.NewWebSocket(url)

	if len(os.Args) > 3 {
		ws.Headers = map[string]string{
			"Authorization": "Bearer " + os.Args[3],
		}
	}
	err := ws.Connect()
	if err != nil {
		panic(err)
	}

	interval := 50 * time.Millisecond
	for {
		count := replay.ReplayInfluxLog(fpath, interval, ws.Write)
		fmt.Printf("wrote: %d lines.\n", count)
		if loop == "loop" {
			continue
		} else {
			break
		}
	}
}
