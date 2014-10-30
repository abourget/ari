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
	Channel        *Channel `json:"channel"`
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

type Channel struct {
	Id           string
	AccountCode  string
	Caller       *CallerID
	Connected    *CallerID
	CreationTime *AriTime
	Dialplan     *DialplanCEP
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
