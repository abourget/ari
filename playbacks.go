package ari

import (
	"fmt"
)

//
// Service
//

type PlaybackService struct {
	client *Client
}

func (s *PlaybackService) Get(playbackID string) (*Playback, error) {
	var out Playback
	return &out, s.client.Get(fmt.Sprintf("/playbacks/%s", playbackID), nil, &out)
}

//
// Model
//

type Playback struct {
	ID        string
	Language  string
	MediaURI  string `json:"media_uri"`
	State     string
	TargetURI string `json:"target_uri"`

	// For further manipulations
	client *Client
}

func (p *Playback) setClient(client *Client) {
	p.client = client
}

func (p *Playback) Stop() error {
	return p.client.Delete(fmt.Sprintf("/playbacks/%s", p.ID), nil)
}

func (p *Playback) Control(operation string) (status int, err error) {
	query := map[string]string{"operation": operation}
	res, err := p.client.PostWithResponse(fmt.Sprintf("/playbacks/%s/control", p.ID), query, nil)
	if err != nil {
		if res == nil {
			return 0, err
		}
		return res.Status(), err
	}
	return res.Status(), err
}
