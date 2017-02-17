package ari

type EndpointService struct {
	client *Client
}

func (s *EndpointService) List() ([]*Endpoint, error) {
	var out []*Endpoint
	if _, err := s.client.Get("/endpoints", nil, &out); err != nil {
		return nil, err
	}

	s.client.setClientRecurse(out)
	return out, nil
}
