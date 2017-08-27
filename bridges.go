package ari

import "fmt"

type BridgeService struct {
	client *Client
}

func (s *BridgeService) List() ([]*Bridge, error) {
	var out []*Bridge
	return out, s.client.Get("/bridges", nil, &out)
}

func (s *BridgeService) Create(params CreateBridgeParams) (*Bridge, error) {
	var out Bridge
	return &out, s.client.Post("/bridges", params, &out)
}

type CreateBridgeParams struct {
	Type     string `json:"type,omitempty"`
	BridgeID string `json:"bridgeId,omitempty"`
	Name     string `json:"name,omitempty"`
}

func (s *BridgeService) Get(bridgeID string) (*Bridge, error) {
	var out Bridge
	return &out, s.client.Get(fmt.Sprintf("/bridges/%s", bridgeID), nil, &out)
}

func (s *BridgeService) Destroy(bridgeID string) error {
	return s.client.Delete(fmt.Sprintf("/bridges/%s", bridgeID), nil)
}

type Bridge struct {
	ID          string
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
	b.client = client
}

func (b *Bridge) Destroy() error {
	return b.client.Delete(fmt.Sprintf("/bridges/%s", b.ID), nil)
}

// AddChannel adds a channel to a bridge. `role` can be `participant` or `announcer`
func (b *Bridge) AddChannel(channel string, role RoleType) error {
	params := map[string]string{
		"channel": channel,
		"role":    string(role),
	}
	return b.client.Post(fmt.Sprintf("/bridges/%s/addChannel", b.ID), params, nil)
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
	return b.client.Post(fmt.Sprintf("/bridges/%s/removeChannel", b.ID), params, nil)
}

// StartMOH starts Music on hold. If mohClass is "", it will not be sent as a param on the request.
func (b *Bridge) StartMOH(mohClass string) error {
	var payload interface{}
	if mohClass != "" {
		payload = map[string]string{"mohClass": mohClass}
	}
	return b.client.Post(fmt.Sprintf("/bridges/%s/moh", b.ID), payload, nil)

}

func (b *Bridge) StopMOH() error {
	return b.client.Delete(fmt.Sprintf("/bridges/%s/moh", b.ID), nil)
}

func (b *Bridge) Play(params PlayParams) (*Playback, error) {
	var out Playback
	return &out, b.client.Post(fmt.Sprintf("/bridges/%s/play", b.ID), &params, &out)
}

func (b *Bridge) Record(params RecordParams) (*LiveRecording, error) {
	var out LiveRecording

	return &out, b.client.Post(fmt.Sprintf("/bridges/%s/record", b.ID), &params, &out)
}
