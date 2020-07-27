package broker

// Simple Go chat client for https://github.com/centrifugal/centrifuge/tree/master/examples/events example.

import (
	"log"

	"github.com/centrifugal/centrifuge-go"
	"github.com/dgrijalva/jwt-go"
)

// Actually in real life clients should never know secret key.
// This is only for example purposes to quickly generate JWT for
// connection.
const exampleTokenHmacSecret = "secret"

func connToken(user string, exp int64) string {
	// NOTE that JWT must be generated on backend side of your application!
	// Here we are generating it on client side only for example simplicity.
	claims := jwt.MapClaims{"sub": user}
	if exp > 0 {
		claims["exp"] = exp
	}
	t, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(exampleTokenHmacSecret))
	if err != nil {
		panic(err)
	}
	return t
}

// GrafanaLiveChannel lets you write data
type GrafanaLiveChannel struct {
	Channel string
	sub     *centrifuge.Subscription
}

// Publish sends the data to the channel
func (b *GrafanaLiveChannel) Publish(data []byte) {
	b.sub.Publish(data)
}

// ChatMessage is chat app specific message struct.
type ChatMessage struct {
	Input string `json:"input"`
}

type eventHandler struct{}

func (h *eventHandler) OnConnect(c *centrifuge.Client, e centrifuge.ConnectEvent) {
	log.Printf("Connected to chat with ID %s", e.ClientID)
	return
}

func (h *eventHandler) OnError(c *centrifuge.Client, e centrifuge.ErrorEvent) {
	log.Printf("Error: %s", e.Message)
	return
}

func (h *eventHandler) OnDisconnect(c *centrifuge.Client, e centrifuge.DisconnectEvent) {
	log.Printf("Disconnected from chat: %s", e.Reason)
	return
}

func (h *eventHandler) OnJoin(sub *centrifuge.Subscription, e centrifuge.JoinEvent) {
	log.Printf("Someone joined: user id %s, client id %s", e.User, e.Client)
}

func (h *eventHandler) OnLeave(sub *centrifuge.Subscription, e centrifuge.LeaveEvent) {
	log.Printf("Someone left: user id %s, client id %s", e.User, e.Client)
}

func (h *eventHandler) OnSubscribeSuccess(sub *centrifuge.Subscription, e centrifuge.SubscribeSuccessEvent) {
	log.Printf("Subscribed on channel %s, resubscribed: %v, recovered: %v", sub.Channel(), e.Resubscribed, e.Recovered)
}

func (h *eventHandler) OnSubscribeError(sub *centrifuge.Subscription, e centrifuge.SubscribeErrorEvent) {
	log.Printf("Subscribed on channel %s failed, error: %s", sub.Channel(), e.Error)
}

func (h *eventHandler) OnUnsubscribe(sub *centrifuge.Subscription, e centrifuge.UnsubscribeEvent) {
	log.Printf("Unsubscribed from channel %s", sub.Channel())
}

// InitGrafanaLiveChannel starts a chat server
func InitGrafanaLiveChannel(url string, channel string) (*GrafanaLiveChannel, error) {
	// url := "ws://localhost:3000/live/ws?format=protobuf"

	log.Printf("Connect to %s\n", url)

	b := &GrafanaLiveChannel{
		Channel: channel,
	}

	c := centrifuge.New(url, centrifuge.DefaultConfig())
	// Uncomment to make it work with Centrifugo and JWT auth.
	//c.SetToken(connToken("49", 0))
	//	defer c.Close()
	handler := &eventHandler{}
	c.OnConnect(handler)
	c.OnError(handler)
	c.OnDisconnect(handler)

	sub, err := c.NewSubscription(channel)
	if err != nil {
		return nil, err
	}
	b.sub = sub

	sub.OnJoin(handler)
	sub.OnLeave(handler)
	sub.OnSubscribeSuccess(handler)
	sub.OnSubscribeError(handler)
	sub.OnUnsubscribe(handler)

	err = sub.Subscribe()
	if err != nil {
		return nil, err
	}

	err = c.Connect()
	if err != nil {
		return nil, err
	}

	return b, nil
}
