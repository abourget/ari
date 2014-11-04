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

func (s *PlaybackService) Get(playbackId string) (*Playback, error) {
	var out Playback

	if _, err := s.client.Get(fmt.Sprintf("/playbacks/%s", playbackId), nil, &out, nil); err != nil {
		return nil, err
	}

	out.setClient(s.client)
	return &out, nil
}

//
// Model
//

type Playback struct {
	Id        string
	Language  string
	MediaURI  string `json:"media_uri"`
	State     string
	TargetURI string `json:"target_uri"`

	// For further manipulations
	client *Client
}

func (p *Playback) setClient(client *Client) {
	if p != nil {
		p.client = client
	}
}

func (p *Playback) Stop() error {
	_, err := p.client.Delete(fmt.Sprintf("/playbacks/%s", p.Id), nil, nil)
	return err
}

func (p *Playback) Control(operation string) (status int, err error) {
	query := map[string]string{"operation": operation}
	res, err := p.client.Post(fmt.Sprintf("/playbacks/%s/control", p.Id), query, nil, nil)
	if err != nil {
		if res != nil {
			return res.Status(), err
		} else {
			return 0, err
		}
	}
	return res.Status(), err
}
