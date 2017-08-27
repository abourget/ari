package ari

type EndpointService struct {
	client *Client
}

func (s *EndpointService) List() ([]*Endpoint, error) {
	var out []*Endpoint
	return out, s.client.Get("/endpoints", nil, &out)
}
