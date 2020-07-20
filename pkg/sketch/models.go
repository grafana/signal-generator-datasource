package broker

import "time"

// The static stuff loaded from plugin.json
type pluginStaticInfo struct {
	PluginID     string   `json:"queryType,omitempty"`
	RPCFunctions []string // query, notifier, etc
	Verson       string
	Author       string

	// if exsts, will have a managed hashi-corp plugin
	Executable string
}

type pluginItem struct {
	Info          pluginStaticInfo
	SigStatus     string // checked at startup
	AllowUnsigned bool   // when unsigned, explicitly say it is OK
	Enabled       bool   // Support running this plugin (plugins can exist without )
	LogLevel      string // trace, etc, etc ???

	// Return startup errors to the browser
	ErrorInfo []string // if there were startup errors, return them to the browser

	// The runner that will execute
	Runner *pluginRunner
}

// Only exists *while* it is connected
type pluginRunner struct {
	plugins []pluginStaticInfo // The plugins this runner supports  1-N
	mode    string             // builtin, managed, connected (load balanced?)

	// The connected client
	clientID string    // centrefuge id
	started  time.Time // when the thing was loaded

	// All RPC methods
	checkHealth      interface{} // plugin, id, details?
	CallResource     interface{} // plugin/id/path
	queryData        interface{}
	sendNotification interface{} // FUTURE
	CollectMetrics   interface{} // return prometheus metrics for the runner (as string)
}

type pluginRegistry struct {
	pluginRunner
}

// Channels
//-----------
// runners:  private channel for processes that can launch plugins

// For plugins that support streams....
// ${plugin-id}/${id}/job/xyz
