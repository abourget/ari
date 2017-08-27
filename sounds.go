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
	return out, s.client.Get("/sounds", &params, &out)
}

func (s *SoundService) Get(soundID string) (*Sound, error) {
	var out Sound
	return &out, s.client.Get(fmt.Sprintf("/sounds/%s", soundID), nil, &out)
}
