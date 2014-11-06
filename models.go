package ari

// Package models implements the Asterisk ARI Messages structures.  See https://wiki.asterisk.org/wiki/display/AST/Asterisk+12+REST+Data+Models

import "fmt"

type Message struct {
	Type string
}

type Event struct {
	Message
	Application string
	Timestamp   *Time
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

type BridgeCreated struct {
	Event
	Bridge *Bridge
}

type BridgeDestroyed struct {
	Event
	Bridge *Bridge
}

type BridgeMerged struct {
	Event
	Bridge     *Bridge
	BridgeFrom *Bridge `json:"bridge_from"`
}

type BridgeBlindTransfer struct {
	Event
	Bridge         *Bridge
	Channel        *Channel
	Context        string
	Exten          string
	IsExternal     bool     `json:"is_external"`
	ReplaceChannel *Channel `json:"replace_channel"`
	Result         string
	Transferee     *Channel
}

type BridgeAttendedTransfer struct {
	Event
	DestinationApplication     string   `json:"destination_application"`
	DestinationBridge          string   `json:"destination_bridge"`
	DestinationLinkFirstLeg    *Channel `json:"destination_link_first_leg"`
	DestinationLinkSecondLeg   *Channel `json:"destination_link_second_leg"`
	DestinationThreeWayBridge  *Bridge  `json:"destination_three_way_bridge"`
	DestinationThreeWayChannel *Channel `json:"destination_three_way_channel"`
	DestinationType            string   `json:"destination_type"`
	IsExternal                 bool     `json:"is_external"`
	ReplaceChannel             *Channel `json:"replace_channel"`
	Result                     string
	TransferTarget             *Channel `json:"transfer_target"`
	Transferee                 *Channel
	TransfererFirstLeg         *Channel `json:"transferer_first_leg"`
	TransfererFirstLegBridge   *Bridge  `json:"transferer_first_leg_bridge"`
	TransfererSecondLeg        *Channel `json:"transferer_second_leg"`
	TransfererSecondLegBridge  *Bridge  `json:"transferer_second_leg_bridge"`
}

type ChannelHangupRequest struct {
	Event
	Cause   int
	Channel *Channel
	Soft    bool
}

type ChannelDtmfReceived struct {
	Event
	Channel    *Channel
	Digit      string
	DurationMs int `json:"duration_ms"`
}

type ChannelTalkingStarted struct {
	Event
	Channel *Channel
}

type ChannelTalkingFinished struct {
	Event
	Channel  *Channel
	Duration int64
}

type ChannelStateChange struct {
	Event
	Channel *Channel
}

type ChannelEnteredBridge struct {
	Event
	Bridge  *Bridge
	Channel *Channel
}

type ChannelLeftBridge struct {
	Event
	Bridge  *Bridge
	Channel *Channel
}

type ChannelDialplan struct {
	Event
	Channel         *Channel
	DialplanApp     string `json:"dialplan_app"`
	DialplanAppData string `json:"dialplan_app_data"`
}

type ChannelCallerId struct {
	Event
	CallerPresentation    int64  `json:"caller_presentation"`
	CallerPresentationTxt string `json:"caller_presentation_txt"`
	Channel               *Channel
}

type ChannelCreated struct {
	Event
	Channel *Channel
}

type ChannelDestroyed struct {
	Event
	Channel  *Channel
	Cause    int64
	CauseTxt string `json:"cause_txt"`
}

type CallerID struct {
	Name   string
	Number string
}

func (c *CallerID) String() string {
	return fmt.Sprintf("%s <%s>", c.Name, c.Number)
}

type Sound struct {
	Formats []FormatLangPair
	Id      string
	Text    string
}

type DialplanCEP struct {
	Context  string
	Exten    string
	Priority int
}

type FormatLangPair struct {
	Format   string
	Language string
}

type PlaybackStarted struct {
	Playback *Playback
}

type PlaybackFinished struct {
	Playback *Playback
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
	LastReloadTime *Time `json:"last_reload_time"`
	StartupTime    *Time `json:"startup_time"`
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
