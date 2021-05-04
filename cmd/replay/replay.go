package main

import (
	"fmt"

	"github.com/grafana/signal-generator-datasource/pkg/replay"
)

func main() {
	fpath := "/home/ryan/Downloads/archer-sample-data.log"
	url := "ws://localhost:3000/api/live/push?gf_live_stream=telegraf"
	key := "eyJrIjoicExKYjlEN29yQmlrMEg4YmtodlRFSjN6R0FOUjRLMEQiLCJuIjoicHVibGlzaCIsImlkIjoxfQ=="

	ws := replay.NewWebSocket(url)
	ws.Headers = map[string]string{
		"Authorization": "Bearer " + key,
	}
	err := ws.Connect()
	if err != nil {
		panic(err)
	}

	fmt.Printf("hello! %s,%s,%s", fpath, url, key)
	replay.DoReplay(fpath, ws.Write)
}
