package plugin

// Simple Go chat client for https://github.com/centrifugal/centrifuge/tree/master/examples/events example.

import (
	"log"
	"time"

	"github.com/centrifugal/centrifuge-go"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
)

// GrafanaLiveChannel lets you write data
type GrafanaLiveChannel struct {
	Channel   string
	sub       *centrifuge.Subscription
	lastWarn  time.Time
	connected bool
}

// Publish sends the data to the channel
func (b *GrafanaLiveChannel) Publish(data []byte) {
	if b.connected {
		_, err := b.sub.Publish(data)
		if err != nil {
			backend.Logger.Error("err", "eee", err)
		}
	} else if time.Since(b.lastWarn) > time.Second*5 {
		b.lastWarn = time.Now()
		log.Printf("grafana live channel not connected: %s\n", b.Channel)
	}
}

func (b *GrafanaLiveChannel) OnConnect(c *centrifuge.Client, e centrifuge.ConnectEvent) {
	log.Printf("Connected to chat with ID %s", e.ClientID)
	b.connected = true
}

func (b *GrafanaLiveChannel) OnError(c *centrifuge.Client, e centrifuge.ErrorEvent) {
	log.Printf("Error: %s", e.Message)
}

func (b *GrafanaLiveChannel) OnDisconnect(c *centrifuge.Client, e centrifuge.DisconnectEvent) {
	log.Printf("Disconnected from chat: %s", e.Reason)
	b.connected = false
}

func (b *GrafanaLiveChannel) OnSubscribeSuccess(sub *centrifuge.Subscription, e centrifuge.SubscribeSuccessEvent) {
	log.Printf("Subscribed on channel %s, resubscribed: %v, recovered: %v", sub.Channel(), e.Resubscribed, e.Recovered)
}

func (b *GrafanaLiveChannel) OnSubscribeError(sub *centrifuge.Subscription, e centrifuge.SubscribeErrorEvent) {
	log.Printf("Subscribed on channel %s failed, error: %s", sub.Channel(), e.Error)
}

func (b *GrafanaLiveChannel) OnUnsubscribe(sub *centrifuge.Subscription, e centrifuge.UnsubscribeEvent) {
	log.Printf("Unsubscribed from channel %s", sub.Channel())
}

// InitGrafanaLiveChannel starts a chat server
func InitGrafanaLiveChannel(url string, channel string) (*GrafanaLiveChannel, error) {
	log.Printf("Connect to %s\n", url)

	b := &GrafanaLiveChannel{
		Channel: channel,
	}

	c := centrifuge.New(url, centrifuge.DefaultConfig())
	c.OnConnect(b)
	c.OnError(b)
	c.OnDisconnect(b)

	sub, err := c.NewSubscription(channel)
	if err != nil {
		return nil, err
	}
	b.sub = sub

	sub.OnSubscribeSuccess(b)
	sub.OnSubscribeError(b)
	sub.OnUnsubscribe(b)

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
