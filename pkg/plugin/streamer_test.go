package plugin

import (
	"testing"

	"github.com/grafana/grafana-edge-app/pkg/tags"
	"github.com/grafana/grafana-plugin-sdk-go/live"
	"github.com/stretchr/testify/assert"
)

func TestStream(t *testing.T) {
	// t.Skip()

	client, err := live.InitGrafanaLiveClient(live.ConnectionInfo{
		URL: "http://localhost:3000",
	})
	assert.NoError(t, err, "error loading live server")

	// Initalize streams
	cfg, err := tags.LoadCaptureSetConfig("../../config/demo-streams.json")
	assert.NoError(t, err, "error loading cfg")

	s, err := NewSignalStreamer(cfg, client)
	assert.NoError(t, err, "error loading streamer")
	assert.Equal(t, 3, len(s.signal.Fields), "field length")

	s.Start()
}
