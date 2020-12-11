package plugin

import (
	"testing"
)

func TestStream(t *testing.T) {
	t.Skip()

	s := &SignalStreamer{
		speedMillis: 5000,
	}

	s.Start()
}
