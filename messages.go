package ari

import "fmt"

type Message struct {
	Type string
}

type Event struct {
	Message
	Application string
	Timestamp   *AriTime
}

type StasisStart struct {
	Event
	Args           []string
	Channel        *Channel
	ReplaceChannel *Channel `json:"replace_channel"`
}

type StasisEnd struct {
	Event
	Channel *Channel
}

type ChannelVarset struct {
	Event
	Channel  *Channel // optionnal
	Value    string
	Variable string
}

type Variable struct {
	Value string
}

type ChannelHangupRequest struct {
	Event
	Cause   int
	Channel *Channel
	Soft    bool
}

type Channel struct {
	Id           string
	AccountCode  string
	Caller       *CallerID
	Connected    *CallerID
	CreationTime *AriTime
	Dialplan     *DialplanCEP
}

func (c Channel) String() string {
	s := fmt.Sprintf("Channel %s", c.Id)
	if c.Caller != nil {
		s = fmt.Sprintf(", caller=%s", c.Caller)
	}
	if c.Connected != nil {
		s = fmt.Sprintf(", with=%s", c.Connected)
	}
	return s
}

type ChannelDtmfReceived struct {
	Event
	Channel    *Channel
	Digit      string
	DurationMs int `json:"duration_ms"`
}

type CallerID struct {
	Name   string
	Number string
}

func (c *CallerID) String() string {
	return fmt.Sprintf("%s <%s>", c.Name, c.Number)
}

type DialplanCEP struct {
	Context  string
	Exten    string
	Priority int
}

//
// AsteriskInfo-related
//
type AsteriskInfo struct {
	Build  *BuildInfo
	Config *ConfigInfo
	Status *StatusInfo
	System *SystemInfo
}

type BuildInfo struct {
	Date    string
	Kernel  string
	Machine string
	Options string
	Os      string
	User    string
}

type ConfigInfo struct {
	DefaultLanguage string  `json:"default_language"`
	MaxChannels     int64   `json:"max_channels"`
	MaxLoad         float64 `json:"max_load"`
	MaxOpenFiles    int64   `json:"max_open_files"`
	Name            string
	SetId           SetId
}

type SetId struct {
	Group string
	User  string
}

type StatusInfo struct {
	LastReloadTime *AriTime `json:"last_reload_time"`
	StartupTime    *AriTime `json:"startup_time"`
}

type SystemInfo struct {
	EntityId string `json:"entity_id"`
	Version  string
}

// AriConnected is an Go library specific message, indicating a successful WebSocket connection.
type AriConnected struct {
	Reconnections int
}

// AriDisonnected is an Go library specific message, indicating an error or a disconnection of the WebSocket connection.
type AriDisconnected struct {
}
