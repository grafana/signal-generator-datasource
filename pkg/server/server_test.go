package server

import (
	"io/ioutil"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBuildFrames(t *testing.T) {
	req := httptest.NewRequest("GET", "/flight", nil)
	w := httptest.NewRecorder()
	handlerFlight(w, req)

	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)

	require.NotEmpty(t, body)
}
