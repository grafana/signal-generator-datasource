package models

import "time"

// Managed by each plugin -- exposed via resource
type TopicMeta struct {
	Topic        string      `json:"topic,omitempty"` // the path typically: `path`
	Description  string      `json:"queryType,omitempty"`
	Opened       time.Time   `json:"queryType,omitempty"`
	MessageTTLs  int32       // Time the topic will stay alive even without data (0 forever)
	LonesomeTTLs int32       // Time the topic will stay alive even without subscribers (0 forever)
	Config       interface{} // json (depends on the implementation)
}

// Managed by grafana core - time relative to
type Client struct {
	SessionID string    // uniq for this session
	UserId    int64     // userinfo
	Connected time.Time `json:"queryType,omitempty"`

	Subscriptions map[string]string // plugin[topic]
}

// PROTOBUF
type Message struct {
	Plugin  string `json:"plugin,omitempty"` // the plugin that sends the request
	Topic   string `json:"topic,omitempty"`  // the topic id
	Session string // optionaly send it to a *single* subscriber session
	Payload []byte // depends on the topic what to do with the bytes
	Action  int32  // DataFrame | AppendRows | text | JSON | TopicShutdown
}

// Interval: Cyclic | OnChange
// Layout: Wide/Long

// Subscribe( topic, user, key ) >> Message | error >>> message is the first one sent on the channel

// Grafana core will maintain the list of subscribers
