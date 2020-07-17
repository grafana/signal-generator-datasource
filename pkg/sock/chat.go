package sock

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	"github.com/grafana/grafana-plugin-sdk-go/backend/log"

	"math/rand"

	"github.com/grafana/waveform-datasource/pkg/models"
	"github.com/grafana/waveform-datasource/pkg/serializers"
	"golang.org/x/time/rate"

	"nhooyr.io/websocket"
)

// ChatServer enables broadcasting to a set of subscribers.
type ChatServer struct {
	// subscriberMessageBuffer controls the max number
	// of messages that can be queued for a subscriber
	// before it is kicked.
	//
	// Defaults to 16.
	subscriberMessageBuffer int

	// publishLimiter controls the rate limit applied to the publish endpoint.
	//
	// Defaults to one publish every 100ms with a burst of 8.
	publishLimiter *rate.Limiter

	// logf controls where logs are sent.
	// Defaults to log.Printf.
	logf func(f string, v ...interface{})

	// serveMux routes the various endpoints to the appropriate handler.
	serveMux http.ServeMux

	subscribersMu sync.Mutex
	subscribers   map[*subscriber]struct{}
}

// newChatServer constructs a ChatServer with the defaults.
func NewChatServer() *ChatServer {
	cs := &ChatServer{
		subscriberMessageBuffer: 16,
		logf: func(f string, v ...interface{}) {
			log.DefaultLogger.Info("LOGF", "msg", fmt.Sprintf(f, v))
		},
		subscribers:    make(map[*subscriber]struct{}),
		publishLimiter: rate.NewLimiter(rate.Every(time.Millisecond*100), 8),
	}
	cs.serveMux.Handle("/", http.FileServer(http.Dir("/Users/stephanie/src/plugins/waveform-datasource/pkg/sock")))
	cs.serveMux.HandleFunc("/subscribe", cs.subscribeHandler)
	cs.serveMux.HandleFunc("/publish", cs.publishHandler)
	cs.serveMux.HandleFunc("/stream", cs.streamSignal)

	return cs
}

// subscriber represents a subscriber.
// Messages are sent on the msgs channel and if the client
// cannot keep up with the messages, closeSlow is called.
type subscriber struct {
	msgs      chan []byte
	closeSlow func()
}

func (cs *ChatServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	cs.serveMux.ServeHTTP(w, r)
}

// subscribeHandler accepts the WebSocket connection and then subscribes
// it to all future messages.
func (cs *ChatServer) subscribeHandler(w http.ResponseWriter, r *http.Request) {
	c, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		InsecureSkipVerify: true, // alow cross orgin
	})
	if err != nil {
		log.DefaultLogger.Info("subscribe", "ACCEPT", err.Error())
		cs.logf("%v", err)
		return
	}
	defer c.Close(websocket.StatusInternalError, "")

	err = cs.subscribe(r.Context(), c)
	if errors.Is(err, context.Canceled) {
		log.DefaultLogger.Info("subscribe", "CANCEL", err.Error)
		return
	}
	if websocket.CloseStatus(err) == websocket.StatusNormalClosure ||
		websocket.CloseStatus(err) == websocket.StatusGoingAway {
		return
	}
	if err != nil {
		log.DefaultLogger.Info("subscribe", "ERROR", err.Error)
		cs.logf("%v", err)
		return
	}
}

// publishHandler reads the request body with a limit of 8192 bytes and then publishes
// the received message.
func (cs *ChatServer) publishHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	body := http.MaxBytesReader(w, r.Body, 8192)
	msg, err := ioutil.ReadAll(body)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusRequestEntityTooLarge), http.StatusRequestEntityTooLarge)
		return
	}

	txt := string(msg) + " // " + r.RemoteAddr

	cs.publish([]byte(txt))

	w.WriteHeader(http.StatusAccepted)
}

// subscribe subscribes the given WebSocket to all broadcast messages.
// It creates a subscriber with a buffered msgs chan to give some room to slower
// connections and then registers the subscriber. It then listens for all messages
// and writes them to the WebSocket. If the context is cancelled or
// an error occurs, it returns and deletes the subscription.
//
// It uses CloseRead to keep reading from the connection to process control
// messages and cancel the context if the connection drops.
func (cs *ChatServer) subscribe(ctx context.Context, c *websocket.Conn) error {
	ctx = c.CloseRead(ctx)

	s := &subscriber{
		msgs: make(chan []byte, cs.subscriberMessageBuffer),
		closeSlow: func() {
			c.Close(websocket.StatusPolicyViolation, "connection too slow to keep up with messages")
		},
	}
	cs.addSubscriber(s)
	defer cs.deleteSubscriber(s)

	for {
		select {
		case msg := <-s.msgs:
			err := writeTimeout(ctx, time.Second*5, c, msg)
			if err != nil {
				return err
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

// publish publishes the msg to all subscribers.
// It never blocks and so messages to slow subscribers
// are dropped.
func (cs *ChatServer) publish(msg []byte) {
	cs.subscribersMu.Lock()
	defer cs.subscribersMu.Unlock()

	cs.publishLimiter.Wait(context.Background())

	for s := range cs.subscribers {
		select {
		case s.msgs <- msg:
		default:
			go s.closeSlow()
		}
	}
}

func (cs *ChatServer) Publish(msg []byte) {
	cs.publish(msg)
}

// addSubscriber registers a subscriber.
func (cs *ChatServer) addSubscriber(s *subscriber) {
	cs.subscribersMu.Lock()
	cs.subscribers[s] = struct{}{}
	cs.subscribersMu.Unlock()
}

// deleteSubscriber deletes the given subscriber.
func (cs *ChatServer) deleteSubscriber(s *subscriber) {
	cs.subscribersMu.Lock()
	delete(cs.subscribers, s)
	cs.subscribersMu.Unlock()
}

func writeTimeout(ctx context.Context, timeout time.Duration, c *websocket.Conn, msg []byte) error {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	return c.Write(ctx, websocket.MessageText, msg)
}

///////

// write to a stream....
func (cs *ChatServer) streamSignalToSocket() {
	speed := 1000 // / 20 // 20 hz
	spread := 50.0

	walker := rand.Float64() * 100
	ticker := time.NewTicker(time.Duration(speed) * time.Millisecond)

	line := models.InfluxLine{
		Name:   "simple",
		Fields: make(map[string]interface{}),
		Tags:   make(map[string]string),
	}

	for t := range ticker.C {
		// if rand.Float64() > 0.1 {
		// 	continue
		// }

		delta := rand.Float64() - 0.5
		walker += delta

		//ms := t.UnixNano() / (int64(time.Millisecond) / int64(time.Nanosecond))

		line.Timestamp = t
		line.Fields["value"] = walker
		line.Fields["min"] = walker - ((rand.Float64() * spread) + 0.01)
		line.Fields["max"] = walker + ((rand.Float64() * spread) + 0.01)

		b, _ := json.Marshal(line)

		cs.publish(b)
	}
}

func (cs *ChatServer) streamMetricsToSocket(datac chan *models.InfluxLine, s serializers.Serializer) {
	for {
		metric := <-datac
		b, _ := s.Serialize(metric)
		cs.publish(b)
	}
}

// write to a stream....
func (cs *ChatServer) streamSignal(w http.ResponseWriter, r *http.Request) {
	setupResponse(&w, r)
	if (*r).Method == "OPTIONS" {
		return
	}

	flusher, ok := w.(http.Flusher)
	if !ok {
		panic("expected http.ResponseWriter to be an http.Flusher")
	}

	speed := 1000 / 20 // 20 hz
	spread := 50.0

	walker := rand.Float64() * 100
	ticker := time.NewTicker(time.Duration(speed) * time.Millisecond)

	w.Header().Set("Content-Type", "text/plain")

	line := models.InfluxLine{
		Name:   "simple",
		Fields: make(map[string]interface{}),
		Tags:   make(map[string]string),
	}

	newLine := []byte("\n")

	lastFlushed := time.Now()
	flushInterval := time.Second / 4

	for t := range ticker.C {
		delta := rand.Float64() - 0.5
		walker += delta

		line.Timestamp = t
		line.Fields["value"] = walker
		line.Fields["min"] = walker - ((rand.Float64() * spread) + 0.01)
		line.Fields["max"] = walker + ((rand.Float64() * spread) + 0.01)

		b, _ := json.Marshal(line)
		w.Write(b)
		w.Write(newLine)

		if time.Now().Sub(lastFlushed) > flushInterval {
			flusher.Flush()
			lastFlushed = time.Now()
		}
	}
}

func setupResponse(w *http.ResponseWriter, req *http.Request) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "GET")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}
