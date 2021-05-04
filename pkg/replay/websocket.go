package replay

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"

	ws "github.com/gorilla/websocket"
)

const (
	defaultConnectTimeout = 5 * time.Second
	defaultWriteTimeout   = 5 * time.Second
	defaultReadTimeout    = 30 * time.Second
)

// WebSocket can output to WebSocket endpoint.
// This is a direct copy of the telegraf websocket output plugin
type WebSocket struct {
	URL            string            `toml:"url"`
	ConnectTimeout time.Duration     `toml:"connect_timeout"`
	WriteTimeout   time.Duration     `toml:"write_timeout"`
	ReadTimeout    time.Duration     `toml:"read_timeout"`
	Headers        map[string]string `toml:"headers"`
	UseTextFrames  bool              `toml:"use_text_frames"`

	conn *ws.Conn
}

var errInvalidURL = errors.New("invalid websocket URL")

// Connect to the output endpoint.
func (w *WebSocket) Connect() error {
	if parsedURL, err := url.Parse(w.URL); err != nil || (parsedURL.Scheme != "ws" && parsedURL.Scheme != "wss") {
		return fmt.Errorf("%w: \"%s\"", errInvalidURL, w.URL)
	}

	dialer := &ws.Dialer{
		Proxy:            http.ProxyFromEnvironment,
		HandshakeTimeout: time.Duration(w.ConnectTimeout),
		TLSClientConfig:  nil, // tlsCfg,
	}

	headers := http.Header{}
	for k, v := range w.Headers {
		headers.Set(k, v)
	}

	conn, resp, err := dialer.Dial(w.URL, headers)
	if err != nil {
		return fmt.Errorf("error dial: %v", err)
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusSwitchingProtocols {
		return fmt.Errorf("wrong status code while connecting to server: %d", resp.StatusCode)
	}

	w.conn = conn
	go w.read(conn)

	return nil
}

func (w *WebSocket) read(conn *ws.Conn) {
	defer func() { _ = conn.Close() }()
	if w.ReadTimeout > 0 {
		if err := conn.SetReadDeadline(time.Now().Add(time.Duration(w.ReadTimeout))); err != nil {
			return
		}
		conn.SetPingHandler(func(string) error {
			err := conn.SetReadDeadline(time.Now().Add(time.Duration(w.ReadTimeout)))
			if err != nil {
				return err
			}
			return conn.WriteControl(ws.PongMessage, nil, time.Now().Add(time.Duration(w.WriteTimeout)))
		})
	}
	for {
		// Need to read a connection (to properly process pings from a server).
		_, _, err := conn.ReadMessage()
		if err != nil {
			return
		}
		if w.ReadTimeout > 0 {
			if err := conn.SetReadDeadline(time.Now().Add(time.Duration(w.ReadTimeout))); err != nil {
				return
			}
		}
	}
}

// Write writes the given metrics to the destination. Not thread-safe.
func (w *WebSocket) Write(messageData []byte) error {
	if w.conn == nil {
		// Previous write failed with error and ws conn was closed.
		if err := w.Connect(); err != nil {
			return err
		}
	}

	messageType := ws.BinaryMessage
	if w.UseTextFrames {
		messageType = ws.TextMessage
	}

	if w.WriteTimeout > 0 {
		if err := w.conn.SetWriteDeadline(time.Now().Add(time.Duration(w.WriteTimeout))); err != nil {
			return fmt.Errorf("error setting write deadline: %v", err)
		}
	}
	err := w.conn.WriteMessage(messageType, messageData)
	if err != nil {
		_ = w.conn.Close()
		w.conn = nil
		return fmt.Errorf("error writing to connection: %v", err)
	}
	return nil
}

// Close closes the connection. Noop if already closed.
func (w *WebSocket) Close() error {
	if w.conn == nil {
		return nil
	}
	err := w.conn.Close()
	w.conn = nil
	return err
}

func NewWebSocket(url string) *WebSocket {
	return &WebSocket{
		URL:            url,
		ConnectTimeout: time.Duration(defaultConnectTimeout),
		WriteTimeout:   time.Duration(defaultWriteTimeout),
		ReadTimeout:    time.Duration(defaultReadTimeout),
		UseTextFrames:  true,
	}
}
