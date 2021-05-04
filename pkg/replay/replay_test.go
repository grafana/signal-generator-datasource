package replay_test

import (
	"testing"
	"time"

	"github.com/grafana/signal-generator-datasource/pkg/replay"
)

func TestLocalFile(t *testing.T) {
	fpath := "/home/ryan/Downloads/archer-sample-data.log"
	player := func(msg []byte) error {
		return nil // NOOP
	}

	interval := 50 * time.Millisecond
	replay.DoReplay(fpath, interval, player)
}
