package broker

import (
	"log"
	"net/http"

	// Import this library.
	"github.com/centrifugal/centrifuge"
)

// Authentication middleware. Centrifuge expects Credentials with current user ID.
// Without provided Credentials client connection won't be accepted.
func auth(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		setupResponse(&w, r)

		ctx := r.Context()
		// Put authentication Credentials into request Context. Since we don't
		// have any session backend here we simply set user ID as empty string.
		// Users with empty ID called anonymous users, in real app you should
		// decide whether anonymous users allowed to connect to your server
		// or not. There is also another way to set Credentials - returning them
		// from ConnectingHandler which is called after client sent first command
		// to server called Connect. See _examples folder in repo to find real-life
		// auth samples (OAuth2, Gin sessions, JWT etc).
		cred := &centrifuge.Credentials{
			UserID: "",
		}
		newCtx := centrifuge.SetCredentials(ctx, cred)
		r = r.WithContext(newCtx)
		h.ServeHTTP(w, r)
	})
}

// GrafanaBroker pretends to be the server
type GrafanaBroker struct {
	node *centrifuge.Node
}

// Publish sends the data to the channel
func (b *GrafanaBroker) Publish(channel string, data []byte) {
	b.node.Publish(channel, data)
}

// ListenAndServe starts a broker running at the given address (:3004)
func (b *GrafanaBroker) ListenAndServe(addr string) {
	// We use default config here as starting point. Default config contains
	// reasonable values for available options.
	cfg := centrifuge.DefaultConfig

	// Node is the core object in Centrifuge library responsible for many useful
	// things. For example Node allows to publish messages to channels from server
	// side with its Publish method, but in this example we will publish messages
	// only from client side.
	node, err := centrifuge.New(cfg)
	if err != nil {
		log.Fatal(err)
	}
	b.node = node

	// Set ConnectHandler called when client successfully connected to Node. Your code
	// inside handler must be synchronized since it will be called concurrently from
	// different goroutines (belonging to different client connections). This is also
	// true for other event handlers.
	node.OnConnect(func(c *centrifuge.Client) {
		// In our example transport will always be Websocket but it can also be SockJS.
		transportName := c.Transport().Name()
		// In our example clients connect with JSON protocol but it can also be Protobuf.
		transportEncoding := c.Transport().Encoding()
		log.Printf("client connected via %s (%s)", transportName, transportEncoding)
	})

	// Set SubscribeHandler to react on every channel subscription attempt
	// initiated by client. Here you can theoretically return an error or
	// disconnect client from server if needed. But now we just accept
	// all subscriptions to all channels. In real life you may use a more
	// complex permission check here.
	node.OnSubscribe(func(c *centrifuge.Client, e centrifuge.SubscribeEvent) (centrifuge.SubscribeReply, error) {
		log.Printf("client subscribes on channel %s", e.Channel)

		return centrifuge.SubscribeReply{}, nil
	})

	node.OnUnsubscribe(func(c *centrifuge.Client, e centrifuge.UnsubscribeEvent) {
		s, _ := node.PresenceStats(e.Channel)

		log.Printf("client unsubscribe from channel %s (clients:%d, users:%d)", e.Channel, s.NumClients, s.NumUsers)
	})

	// By default, clients can not publish messages into channels. By setting
	// PublishHandler we tell Centrifuge that publish from client side is possible.
	// Now each time client calls publish method this handler will be called and
	// you have a possibility to validate publication request before message will
	// be published into channel and reach active subscribers. In our simple chat
	// app we allow everyone to publish into any channel.
	node.OnPublish(func(c *centrifuge.Client, e centrifuge.PublishEvent) (centrifuge.PublishReply, error) {
		log.Printf("client publishes into channel %s: %s", e.Channel, string(e.Data))
		return centrifuge.PublishReply{}, nil
	})

	// Set Disconnect handler to react on client disconnect events.
	node.OnDisconnect(func(c *centrifuge.Client, e centrifuge.DisconnectEvent) {
		log.Printf("client disconnected")
	})

	// Run node. This method does not block.
	if err := node.Run(); err != nil {
		log.Fatal(err)
	}

	// Serve Websocket connections using WebsocketHandler.
	// wsHandler := centrifuge.NewWebsocketHandler(node, centrifuge.WebsocketConfig{})
	// http.Handle("/broker", auth(wsHandler))

	sockjsHandler := centrifuge.NewSockjsHandler(node, centrifuge.SockjsConfig{
		URL:                      "https://cdn.jsdelivr.net/npm/sockjs-client@1/dist/sockjs.min.js", //??
		HandlerPrefix:            "/broker/sockjs",
		WebsocketReadBufferSize:  1024,
		WebsocketWriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			log.Printf("check orgin? %s", r.RemoteAddr)
			return true
		},
		WebsocketCheckOrigin: func(r *http.Request) bool {
			log.Printf("check websocket orgin? %s", r.RemoteAddr)
			return true
		},
	})
	http.Handle("/broker/sockjs/", auth(sockjsHandler))

	log.Printf("Starting server, visit: %s/broker", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal(err)
	}
}

func setupResponse(w *http.ResponseWriter, req *http.Request) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "GET")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}
