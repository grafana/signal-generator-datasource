package replay_test

import (
	"strings"
	"testing"
	"time"

	"github.com/grafana/signal-generator-datasource/pkg/replay"
	"github.com/stretchr/testify/require"
)

func TestLocalFile(t *testing.T) {
	fpath := "../testdata/motor.log"
	count := 0
	player := func(msg []byte) error {
		lines := strings.Split(strings.TrimSpace(string(msg)), "\n")
		count += len(lines)
		return nil // NOOP
	}

	interval := 50 * time.Millisecond
	rows := replay.ReplayInfluxLog(fpath, interval, player)
	require.Equal(t, rows, count) // number of lines in the file
}
