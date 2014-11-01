package ari

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

type DialplanCEP struct {
	Context  string
	Exten    string
	Priority int
}
