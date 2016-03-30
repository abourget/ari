package ari

import (
	"fmt"

	"github.com/jmcvetta/napping"
)

type SoundService struct {
	client *Client
}

// SoundsGet retrieves sounds. Both `lang` and `format` can be empty strings.
func (s *SoundService) List(lang, format string) ([]*Sound, error) {
	var out []*Sound
	p := napping.Params{}
	if lang != "" {
		p["lang"] = lang
	}
	if format != "" {
		p["format"] = format
	}

	params := p.AsUrlValues()
	if _, err := s.client.Get("/sounds", &params, &out); err != nil {
		return nil, err
	}
	return out, nil
}

func (s *SoundService) Get(soundId string) (*Sound, error) {
	var out *Sound

	if _, err := s.client.Get(fmt.Sprintf("/sounds/%s", soundId), nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}
