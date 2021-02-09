package plugin

import (
	"strings"
	"testing"

	"github.com/grafana/grafana-edge-app/pkg/tags"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/experimental"
	"github.com/grafana/grafana-plugin-sdk-go/live"
	"github.com/stretchr/testify/assert"
)

func TestStream(t *testing.T) {
	t.Skip()

	client, err := live.InitGrafanaLiveClient(live.ConnectionInfo{
		URL: "http://localhost:3000",
	})
	assert.NoError(t, err, "error loading live server")

	// Initalize streams
	cfg, err := tags.LoadCaptureSetConfig("../../config/demo-streams.json")
	assert.NoError(t, err, "error loading cfg")

	s, err := NewSignalStreamer(cfg, client)
	assert.NoError(t, err, "error loading streamer")
	assert.Equal(t, 7, len(s.signal.Fields), "field length")

	fname := "../testdata/streamer-simple.golden.txt"
	frames, err := s.Frames()
	dr := &backend.DataResponse{
		Frames: frames,
		Error:  err,
	}
	if err := experimental.CheckGoldenDataResponse(fname, dr, true); err != nil {
		if !strings.Contains(err.Error(), "no such file or directory") {
			t.Fatal(err)
		}
	}

	//s.Start()
}
