package ari

import (
	"fmt"
	"net/url"
)

type ApplicationService struct {
	client *Client
}

func (s *ApplicationService) List() ([]*Application, error) {
	var out []*Application
	return out, s.client.Get("/applications", nil, &out)
}

func (s *ApplicationService) Get(applicationName string) (*Application, error) {
	var out Application

	return &out, s.client.Get(fmt.Sprintf("/applications/%s", applicationName), nil, &out)
}

type Application struct {
	BridgeIds   []string `json:"bridge_ids"`
	ChannelIds  []string `json:"channel_ids"`
	DeviceNames []string `json:"device_names"`
	EndpointIds []string `json:"endpoint_ids"`
	Name        string

	// For further mutations
	client *Client
}

func (a *Application) setClient(client *Client) {
	a.client = client
}

func (a *Application) Subscribe(eventSource string) (*Application, error) {
	var out Application
	params := map[string]string{
		"eventSource": eventSource,
	}
	return &out, a.client.Post(fmt.Sprintf("/applications/%s/subscription", a.Name), params, &out)
}

func (a *Application) Unsubscribe(eventSource string) (*Application, error) {
	var out Application
	return &out, a.client.Delete(fmt.Sprintf("/applications/%s/subscription?eventSource=%s", a.Name, url.QueryEscape(eventSource)), &out)
}
