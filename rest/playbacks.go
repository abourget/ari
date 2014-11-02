package rest

import (
	"fmt"

	"github.com/abourget/ari/models"
)

func (r *REST) PlaybacksGetById(playbackId string) (*models.Playback, error) {
	var out models.Playback

	if _, err := r.Get(fmt.Sprintf("/playbacks/%s", playbackId), nil, &out, nil); err != nil {
		return nil, err
	}
	return &out, nil
}

func (r *REST) PlaybacksDeleteByid(playbackId string) error {
	_, err := r.Delete(fmt.Sprintf("/playbacks/%s", playbackId), nil, nil)
	return err
}

func (r *REST) PlaybacksControlPostById(playbackId, operation string) (status int, err error) {
	query := map[string]string{"operation": operation}
	res, err := r.Post(fmt.Sprintf("/playbacks/%s/control", playbackId), query, nil, nil)
	if err != nil {
		if res != nil {
			return res.Status(), err
		} else {
			return 0, err
		}
	}
	return res.Status(), err
}
