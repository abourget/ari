package ari

import "fmt"

type BridgeService struct {
	client *Client
}

func (s *BridgeService) List() ([]*Bridge, error) {
	var out []*Bridge

	if _, err := s.client.Get("/bridges", nil, &out); err != nil {
		return nil, err
	}

	s.client.setClientRecurse(out)
	return out, nil
}

func (s *BridgeService) Create(params CreateBridgeParams) (*Bridge, error) {
	var out Bridge

	if _, err := s.client.Post("/bridges", params, &out); err != nil {
		return nil, err
	}

	out.setClient(s.client)
	return &out, nil
}

type CreateBridgeParams struct {
	Type     string `json:"type,omitempty"`
	BridgeId string `json:"bridgeId,omitempty"`
	Name     string `json:"name,omitempty"`
}

func (s *BridgeService) Get(bridgeId string) (*Bridge, error) {
	var out Bridge

	if _, err := s.client.Get(fmt.Sprintf("/bridges/%s", bridgeId), nil, &out); err != nil {
		return nil, err
	}

	out.setClient(s.client)
	return &out, nil
}

func (s *BridgeService) Destroy(bridgeId string) error {
	_, err := s.client.Delete(fmt.Sprintf("/bridges/%s", bridgeId), nil)
	return err
}

type Bridge struct {
	Id          string
	Name        string
	Technology  string
	Creator     string
	Channels    []string
	BridgeType  string `json:"bridge_type"`
	BridgeClass string `json:"bridge_class"`

	// For further manipulations
	client *Client
}

func (b *Bridge) setClient(client *Client) {
	if b != nil {
		b.client = client
	}
}

func (b *Bridge) Destroy() error {
	_, err := b.client.Delete(fmt.Sprintf("/bridges/%s", b.Id), nil)
	return err
}

// AddChannel adds a channel to a bridge. `role` can be `participant` or `announcer`
func (b *Bridge) AddChannel(channel string, role RoleType) error {
	params := map[string]string{
		"channel": channel,
		"role":    string(role),
	}
	if _, err := b.client.Post(fmt.Sprintf("/bridges/%s/addChannel", b.Id), params, nil); err != nil {
		return err
	}
	return nil
}

type RoleType string

const (
	Participant RoleType = "participant"
	Announcer   RoleType = "announcer"
)

func (b *Bridge) RemoveChannel(channel string) error {
	params := map[string]string{
		"channel": channel,
	}
	if _, err := b.client.Post(fmt.Sprintf("/bridges/%s/removeChannel", b.Id), params, nil); err != nil {
		return err
	}
	return nil
}

// StartMOH starts Music on hold. If mohClass is "", it will not be sent as a param on the request.
func (b *Bridge) StartMOH(mohClass string) error {
	var payload interface{}
	if mohClass != "" {
		payload = map[string]string{"mohClass": mohClass}
	}
	if _, err := b.client.Post(fmt.Sprintf("/bridges/%s/moh", b.Id), payload, nil); err != nil {
		return err
	}
	return nil

}

func (b *Bridge) StopMOH() error {
	_, err := b.client.Delete(fmt.Sprintf("/bridges/%s/moh", b.Id), nil)
	return err
}

func (b *Bridge) Play(params PlayParams) (*Playback, error) {
	var out Playback

	if _, err := b.client.Post(fmt.Sprintf("/bridges/%s/play", b.Id), &params, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (b *Bridge) Record(params RecordParams) (*LiveRecording, error) {
	var out LiveRecording

	if _, err := b.client.Post(fmt.Sprintf("/bridges/%s/record", b.Id), &params, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
