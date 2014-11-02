package rest

import (
	"fmt"

	"github.com/abourget/ari/models"
	"github.com/jmcvetta/napping"
)

// SoundsGet retrieves sounds. Both `lang` and `format` can be empty strings.
func (r *REST) SoundsGet(lang, format string) ([]*models.Sound, error) {
	var out []*models.Sound
	p := napping.Params{}
	if lang != "" {
		p["lang"] = lang
	}
	if format != "" {
		p["format"] = format
	}

	if _, err := r.Get("/sounds", &p, &out, nil); err != nil {
		return nil, err
	}
	return out, nil
}

func (r *REST) SoundsGetById(soundId string) (*models.Sound, error) {
	var out *models.Sound

	if _, err := r.Get(fmt.Sprintf("/sounds/%s", soundId), nil, &out, nil); err != nil {
		return nil, err
	}
	return out, nil
}
