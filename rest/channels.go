package rest

import (
	"fmt"

	"github.com/abourget/ari/models"
	"github.com/jmcvetta/napping"
)

//
// Channels, see https://wiki.asterisk.org/wiki/display/AST/Asterisk+12+Channels+REST+API
//

func (r *REST) ChannelsGet() ([]*models.Channel, error) {
	var out []*models.Channel

	if _, err := r.Get("/channels", nil, &out, nil); err != nil {
		return nil, err
	}
	return out, nil
}

func (r *REST) ChannelsPost(params ChannelsPostParams) (*models.Channel, error) {
	var out models.Channel
	if _, err := r.Post("/channels", params, &out, nil); err != nil {
		return nil, err
	}
	return &out, nil
}

type ChannelsPostParams struct {
	Endpoint       string            `json:"endpoint"`
	Extension      string            `json:"extension,omitempty"`
	Context        string            `json:"context,omitempty"`
	Priority       int64             `json:"priority,omitempty"`
	App            string            `json:"Aap,omitempty"`
	AppArgs        string            `json:"appArgs,omitempty"`
	CallerId       string            `json:"callerId,omitempty"`
	Timeout        int64             `json:"timeout,omitempty"`
	ChannelId      string            `json:"channelId,omitempty"`
	OtherChannelId string            `json:"otherChannelId,omitempty"`
	Variables      map[string]string `json:"variables,omitempty"`
}

func (r *REST) ChannelsGetById(channelId string) (*models.Channel, error) {
	var out *models.Channel

	if _, err := r.Get(fmt.Sprintf("/channels/%s", channelId), nil, &out, nil); err != nil {
		return nil, err
	}
	return out, nil
}

func (r *REST) ChannelsPostById(channelId string, params ChannelsPostParams) (*models.Channel, error) {
	params.ChannelId = ""
	var out models.Channel

	if _, err := r.Post(fmt.Sprintf("/channels/%s", channelId), params, &out, nil); err != nil {
		return nil, err
	}
	return &out, nil
}

func (r *REST) ChannelsDeleteById(channelId string) error {
	_, err := r.Delete(fmt.Sprintf("/channels/%s", channelId), nil, nil)
	return err
}

func (r *REST) ChannelsContinuePostById(channelId string, cep models.DialplanCEP) error {
	if _, err := r.Post(fmt.Sprintf("/channels/%s/continue", channelId), cep, nil, nil); err != nil {
		return err
	}
	return nil
}

func (r *REST) ChannelsAnswerPostById(channelId string) error {
	if _, err := r.Post(fmt.Sprintf("/channels/%s/answer", channelId), nil, nil, nil); err != nil {
		return err
	}
	return nil
}

func (r *REST) ChannelsRingPostById(channelId string) error {
	if _, err := r.Post(fmt.Sprintf("/channels/%s/ring", channelId), nil, nil, nil); err != nil {
		return err
	}
	return nil
}

func (r *REST) ChannelsRingDeleteById(channelId string) error {
	_, err := r.Delete(fmt.Sprintf("/channels/%s/ring", channelId), nil, nil)
	return err
}

// ChannelsDTMFPostById calls /channels/{channelId}/dtmf. It accepts either a string or a ChannelDTMFSend object.
func (r *REST) ChannelsDTMFPostById(channelId string, dtmf interface{}) error {
	var dtmfSend DTMFParams
	switch o := dtmf.(type) {
	case string:
		dtmfSend = DTMFParams{DTMF: o}
	case DTMFParams:
		dtmfSend = o
	default:
		panic("Invalid type for `dtmf` param in ChannelsDTMFPostById")
	}

	if _, err := r.Post(fmt.Sprintf("/channels/%s/dtmf", channelId), dtmfSend, nil, nil); err != nil {
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
func (r *REST) ChannelsMutePostById(channelId string, direction string) error {
	if _, err := r.Post(fmt.Sprintf("/channels/%s/mute", channelId), map[string]string{"direction": direction}, nil, nil); err != nil {
		return err
	}
	return nil
}

// ChannelsMuteDeleteById unmutes a channel. Use `direction="both"` for default behavior.
func (r *REST) ChannelsMuteDeleteById(channelId string, direction string) error {
	_, err := r.Delete(fmt.Sprintf("/channels/%s/mute?direction=%s", channelId, direction), nil, nil)
	return err
}

func (r *REST) ChannelsHoldPostById(channelId string) error {
	if _, err := r.Post(fmt.Sprintf("/channels/%s/hold", channelId), nil, nil, nil); err != nil {
		return err
	}
	return nil
}

func (r *REST) ChannelsHoldDeleteById(channelId string) error {
	_, err := r.Delete(fmt.Sprintf("/channels/%s/hold", channelId), nil, nil)
	return err
}

// ChannelsMOHPostById posts to /channels/{channelId}/moh. If mohClass is "", it will not be sent as a param on the request.
func (r *REST) ChannelsMOHPostById(channelId string, mohClass string) error {
	var payload interface{}
	if mohClass != "" {
		payload = map[string]string{"mohClass": mohClass}
	}
	if _, err := r.Post(fmt.Sprintf("/channels/%s/moh", channelId), payload, nil, nil); err != nil {
		return err
	}
	return nil

}

func (r *REST) ChannelsMOHDeleteById(channelId string) error {
	_, err := r.Delete(fmt.Sprintf("/channels/%s/moh", channelId), nil, nil)
	return err
}

func (r *REST) ChannelsSilencePostById(channelId string) error {
	if _, err := r.Post(fmt.Sprintf("/channels/%s/silence", channelId), nil, nil, nil); err != nil {
		return err
	}
	return nil
}

func (r *REST) ChannelsSilenceDeleteById(channelId string) error {
	_, err := r.Delete(fmt.Sprintf("/channels/%s/silence", channelId), nil, nil)
	return err
}

func (r *REST) ChannelsPlayPostById(channelId string, params PlayParams) (*models.Playback, error) {
	var out models.Playback

	if _, err := r.Post(fmt.Sprintf("/channels/%s/play", channelId), &params, &out, nil); err != nil {
		return nil, err
	}
	return &out, nil
}

type PlayParams struct {
	Media      string `json:"media"`
	Lang       string `json:"lang,omitempty"`
	OffsetMS   int64  `json:"offsetms,omitempty"`
	SkipMS     int64  `json:"skipms,omitempty"`
	PlaybackId string `"json:"playbackId,omitempty"`
}

func (r *REST) ChannelsRecordPostById(channelId string, params RecordParams) (*models.LiveRecording, error) {
	var out models.LiveRecording

	if _, err := r.Post(fmt.Sprintf("/channels/%s/record", channelId), &params, &out, nil); err != nil {
		return nil, err
	}
	return &out, nil
}

type RecordParams struct {
	Name               string `json:"name"`
	Format             string `json:"format"`
	MaxDurationSeconds int64  `json:"maxDurationSeconds"`
	MaxSilenceSeconds  int64  `json:"maxSilenceSeconds"`
	IfExists           string `json:"ifExists,omitempty"`
	Beep               bool   `json:"beep"`
	TerminateOn        string `json:"terminateOn,omitempty"`
}

func (r *REST) ChannelsVariableGetById(channelId, variable string) (string, error) {
	var out models.Variable

	if _, err := r.Get(fmt.Sprintf("/channels/%s/variable", channelId), &napping.Params{"variable": variable}, &out, nil); err != nil {
		return "", err
	}
	return out.Value, nil
}

func (r *REST) ChannelsVariablePostById(channelId, variable, value string) error {
	payload := map[string]string{"variable": variable, "value": value}

	if _, err := r.Post(fmt.Sprintf("/channels/%s/variable", channelId), payload, nil, nil); err != nil {
		return err
	}
	return nil
}

func (r *REST) ChannelsSnoopPostById(channelId string, params SnoopParams) (*models.Channel, error) {
	var out models.Channel

	if _, err := r.Post(fmt.Sprintf("/channels/%s/snoop", channelId), params, &out, nil); err != nil {
		return nil, err
	}
	return &out, nil

}

type SnoopParams struct {
	App     string `json:"app"`
	AppArgs string `json:"appArgs,omitempty"`
	Spy     string `json:"spy,omitempty"`
	Whisper string `json:"whisper,omitempty"`
	SnoopId string `json:"snoopId,omitempty"`
}
