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

	if _, err := s.client.Get("/channels", nil, &out); err != nil {
		return nil, err
	}

	s.client.setClientRecurse(out)
	return out, nil
}

func (s *ChannelService) Create(params OriginateParams) (*Channel, error) {
	var out Channel
	if _, err := s.client.Post("/channels", params, &out); err != nil {
		return nil, err
	}

	out.setClient(s.client)
	return &out, nil
}

func (s *ChannelService) Get(channelID string) (*Channel, error) {
	var out Channel

	if _, err := s.client.Get(fmt.Sprintf("/channels/%s", channelID), nil, &out); err != nil {
		return nil, err
	}

	out.setClient(s.client)
	return &out, nil
}

func (s *ChannelService) Hangup(channelID string) error {
	_, err := s.client.Delete(fmt.Sprintf("/channels/%s", channelID), nil)
	return err
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
	_, err := c.client.Delete(fmt.Sprintf("/channels/%s", c.ID), nil)
	return err
}

func (c *Channel) ContinueInDialplan(context, exten string, priority int, label string) error {
	if _, err := c.client.Post(fmt.Sprintf("/channels/%s/continue", c.ID), Dialplan{context, exten, priority, label}, nil); err != nil {
		return err
	}
	return nil
}

func (c *Channel) Answer() error {
	if _, err := c.client.Post(fmt.Sprintf("/channels/%s/answer", c.ID), nil, nil); err != nil {
		return err
	}
	return nil
}

func (c *Channel) Ring() error {
	if _, err := c.client.Post(fmt.Sprintf("/channels/%s/ring", c.ID), nil, nil); err != nil {
		return err
	}
	return nil
}

func (c *Channel) RingStop() error {
	_, err := c.client.Delete(fmt.Sprintf("/channels/%s/ring", c.ID), nil)
	return err
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

	if _, err := c.client.Post(fmt.Sprintf("/channels/%s/dtmf", c.ID), dtmfSend, nil); err != nil {
		return err
	}
	return nil
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
	if _, err := c.client.Post(fmt.Sprintf("/channels/%s/mute", c.ID), map[string]string{"direction": direction}, nil); err != nil {
		return err
	}
	return nil
}

// ChannelsMuteDeleteById unmutes a channel. Use `direction="both"` for default behavior.
func (c *Channel) Unmute(direction string) error {
	_, err := c.client.Delete(fmt.Sprintf("/channels/%s/mute?direction=%s", c.ID, direction), nil)
	return err
}

func (c *Channel) Hold() error {
	if _, err := c.client.Post(fmt.Sprintf("/channels/%s/hold", c.ID), nil, nil); err != nil {
		return err
	}
	return nil
}

func (c *Channel) StopHold() error {
	_, err := c.client.Delete(fmt.Sprintf("/channels/%s/hold", c.ID), nil)
	return err
}

// StartMOH starts Music on hold. If mohClass is "", it will not be sent as a param on the request.
func (c *Channel) StartMOH(mohClass string) error {
	var payload interface{}
	if mohClass != "" {
		payload = map[string]string{"mohClass": mohClass}
	}
	if _, err := c.client.Post(fmt.Sprintf("/channels/%s/moh", c.ID), payload, nil); err != nil {
		return err
	}
	return nil

}

func (c *Channel) StopMOH() error {
	_, err := c.client.Delete(fmt.Sprintf("/channels/%s/moh", c.ID), nil)
	return err
}

func (c *Channel) StartSilence() error {
	if _, err := c.client.Post(fmt.Sprintf("/channels/%s/silence", c.ID), nil, nil); err != nil {
		return err
	}
	return nil
}

func (c *Channel) StopSilence() error {
	_, err := c.client.Delete(fmt.Sprintf("/channels/%s/silence", c.ID), nil)
	return err
}

// Play plays media through channel. See: https://wiki.asterisk.org/wiki/display/AST/ARI+and+Channels%3A+Simple+Media+Manipulation
func (c *Channel) Play(params PlayParams) (*Playback, error) {
	var out Playback

	if _, err := c.client.Post(fmt.Sprintf("/channels/%s/play", c.ID), &params, &out); err != nil {
		return nil, err
	}

	out.setClient(c.client)
	return &out, nil
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

	if _, err := c.client.Post(fmt.Sprintf("/channels/%s/record", c.ID), &params, &out); err != nil {
		return nil, err
	}

	out.setClient(c.client)
	return &out, nil
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
	if _, err := c.client.Get(fmt.Sprintf("/channels/%s/variable", c.ID), &params, &out); err != nil {
		return "", err
	}
	return out.Value, nil
}

func (c *Channel) SetVar(variable, value string) error {
	payload := map[string]string{"variable": variable, "value": value}

	if _, err := c.client.Post(fmt.Sprintf("/channels/%s/variable", c.ID), payload, nil); err != nil {
		return err
	}
	return nil
}

func (c *Channel) Snoop(params SnoopParams) (*Channel, error) {
	var out Channel

	if _, err := c.client.Post(fmt.Sprintf("/channels/%s/snoop", c.ID), params, &out); err != nil {
		return nil, err
	}

	out.setClient(c.client)
	return &out, nil

}

type SnoopParams struct {
	App     string `json:"app"`
	AppArgs string `json:"appArgs,omitempty"`
	Spy     string `json:"spy,omitempty"`
	Whisper string `json:"whisper,omitempty"`
	SnoopID string `json:"snoopId,omitempty"`
}
