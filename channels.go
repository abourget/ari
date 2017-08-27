package ari

import (
	"fmt"

	"github.com/jmcvetta/napping"
)

//
// Channels, see https://wiki.asterisk.org/wiki/display/AST/Asterisk+12+Channels+REST+API
//

type ChannelService struct {
	client *Client
}

func (s *ChannelService) List() ([]*Channel, error) {
	var out []*Channel
	return out, s.client.Get("/channels", nil, &out)
}

func (s *ChannelService) Create(params OriginateParams) (*Channel, error) {
	var out Channel
	return &out, s.client.Post("/channels", params, &out)
}

func (s *ChannelService) Get(channelID string) (*Channel, error) {
	var out Channel
	return &out, s.client.Get(fmt.Sprintf("/channels/%s", channelID), nil, &out)
}

func (s *ChannelService) Hangup(channelID string) error {
	return s.client.Delete(fmt.Sprintf("/channels/%s", channelID), nil)
}

type OriginateParams struct {
	Endpoint       string            `json:"endpoint"`
	Extension      string            `json:"extension,omitempty"`
	Context        string            `json:"context,omitempty"`
	Priority       int64             `json:"priority,omitempty"`
	App            string            `json:"app,omitempty"`
	AppArgs        string            `json:"appArgs,omitempty"`
	CallerID       string            `json:"callerId,omitempty"`
	Timeout        int64             `json:"timeout,omitempty"`
	ChannelID      string            `json:"channelId,omitempty"`
	OtherChannelID string            `json:"otherChannelId,omitempty"`
	Variables      map[string]string `json:"variables,omitempty"`
}

//
// Model
//

type Channel struct {
	ID           string
	AccountCode  string
	Caller       *CallerID
	Connected    *CallerID
	CreationTime *Time
	Dialplan     *DialplanCEP
	Name         string
	State        string

	// For further manipulations
	client *Client
}

func (c *Channel) setClient(client *Client) {
	c.client = client
}

func (c *Channel) String() string {
	s := fmt.Sprintf("id=%s", c.ID)
	if c.Caller != nil {
		s = fmt.Sprintf("%s,caller=%s", s, c.Caller)
	}
	if c.Connected != nil {
		s = fmt.Sprintf("%s,with=%s", s, c.Connected)
	}

	s = fmt.Sprintf("%s,state=%s", s, c.State)
	return s
}

func (c *Channel) Hangup() error {
	return c.client.Delete(fmt.Sprintf("/channels/%s", c.ID), nil)
}

func (c *Channel) ContinueInDialplan(context, exten string, priority int, label string) error {
	return c.client.Post(fmt.Sprintf("/channels/%s/continue", c.ID), Dialplan{context, exten, priority, label}, nil)
}

func (c *Channel) Answer() error {
	return c.client.Post(fmt.Sprintf("/channels/%s/answer", c.ID), nil, nil)
}

func (c *Channel) Ring() error {
	return c.client.Post(fmt.Sprintf("/channels/%s/ring", c.ID), nil, nil)
}

func (c *Channel) RingStop() error {
	return c.client.Delete(fmt.Sprintf("/channels/%s/ring", c.ID), nil)
}

// SendDTMF sends DTMF signals to the channel. It accepts either a string or a ChannelDTMFSend object.
func (c *Channel) SendDTMF(dtmf interface{}) error {
	var dtmfSend DTMFParams
	switch o := dtmf.(type) {
	case string:
		dtmfSend = DTMFParams{DTMF: o}
	case DTMFParams:
		dtmfSend = o
	default:
		panic("Invalid type for `dtmf` param in ChannelsDTMFPostById")
	}

	return c.client.Post(fmt.Sprintf("/channels/%s/dtmf", c.ID), dtmfSend, nil)
}

type DTMFParams struct {
	DTMF     string `json:"dtmf"`
	Before   int64  `json:"before,omitempty"`
	Between  int64  `json:"between,omitempty"`
	Duration int64  `json:"duration,omitempty"`
	After    int64  `json:"after,omitempty"`
}

// ChannelsMutePostById mutes a channel. Use `direction="both"` for default behavior.
func (c *Channel) Mute(direction string) error {
	return c.client.Post(fmt.Sprintf("/channels/%s/mute", c.ID), map[string]string{"direction": direction}, nil)
}

// ChannelsMuteDeleteById unmutes a channel. Use `direction="both"` for default behavior.
func (c *Channel) Unmute(direction string) error {
	return c.client.Delete(fmt.Sprintf("/channels/%s/mute?direction=%s", c.ID, direction), nil)
}

func (c *Channel) Hold() error {
	return c.client.Post(fmt.Sprintf("/channels/%s/hold", c.ID), nil, nil)
}

func (c *Channel) StopHold() error {
	return c.client.Delete(fmt.Sprintf("/channels/%s/hold", c.ID), nil)
}

// StartMOH starts Music on hold. If mohClass is "", it will not be sent as a param on the request.
func (c *Channel) StartMOH(mohClass string) error {
	var payload interface{}
	if mohClass != "" {
		payload = map[string]string{"mohClass": mohClass}
	}
	return c.client.Post(fmt.Sprintf("/channels/%s/moh", c.ID), payload, nil)

}

func (c *Channel) StopMOH() error {
	return c.client.Delete(fmt.Sprintf("/channels/%s/moh", c.ID), nil)
}

func (c *Channel) StartSilence() error {
	return c.client.Post(fmt.Sprintf("/channels/%s/silence", c.ID), nil, nil)
}

func (c *Channel) StopSilence() error {
	return c.client.Delete(fmt.Sprintf("/channels/%s/silence", c.ID), nil)
}

// Play plays media through channel. See: https://wiki.asterisk.org/wiki/display/AST/ARI+and+Channels%3A+Simple+Media+Manipulation
func (c *Channel) Play(params PlayParams) (*Playback, error) {
	var out Playback
	return &out, c.client.Post(fmt.Sprintf("/channels/%s/play", c.ID), &params, &out)
}

type PlayParams struct {
	Media      string `json:"media"`
	Lang       string `json:"lang,omitempty"`
	OffsetMS   int64  `json:"offsetms,omitempty"`
	SkipMS     int64  `json:"skipms,omitempty"`
	PlaybackID string `json:"playbackId,omitempty"`
}

func (c *Channel) Record(params RecordParams) (*LiveRecording, error) {
	var out LiveRecording
	return &out, c.client.Post(fmt.Sprintf("/channels/%s/record", c.ID), &params, &out)
}

type RecordParams struct {
	Name               string `json:"name"`
	Format             string `json:"format,omitempty"`
	MaxDurationSeconds int64  `json:"maxDurationSeconds"`
	MaxSilenceSeconds  int64  `json:"maxSilenceSeconds"`
	IfExists           string `json:"ifExists,omitempty"`
	Beep               bool   `json:"beep"`
	TerminateOn        string `json:"terminateOn,omitempty"`
}

func (c *Channel) GetVar(variable string) (string, error) {
	var out Variable
	params := napping.Params{"variable": variable}.AsUrlValues()
	err := c.client.Get(fmt.Sprintf("/channels/%s/variable", c.ID), &params, &out)
	return out.Value, err
}

func (c *Channel) SetVar(variable, value string) error {
	payload := map[string]string{"variable": variable, "value": value}

	return c.client.Post(fmt.Sprintf("/channels/%s/variable", c.ID), payload, nil)
}

func (c *Channel) Snoop(params SnoopParams) (*Channel, error) {
	var out Channel
	return &out, c.client.Post(fmt.Sprintf("/channels/%s/snoop", c.ID), params, &out)

}

type SnoopParams struct {
	App     string `json:"app"`
	AppArgs string `json:"appArgs,omitempty"`
	Spy     string `json:"spy,omitempty"`
	Whisper string `json:"whisper,omitempty"`
	SnoopID string `json:"snoopId,omitempty"`
}
