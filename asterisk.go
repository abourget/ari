package ari

import "github.com/jmcvetta/napping"

//
// Asterisk endpoints wrappers
//

type AsteriskService struct {
	client *Client
}

func (s *AsteriskService) GetInfo() (*AsteriskInfo, error) {
	ai := AsteriskInfo{}

	if _, err := s.client.Get("/asterisk/info", nil, &ai); err != nil {
		return nil, err
	}
	return &ai, nil
}

func (s *AsteriskService) GetGlobalVar(variable string) (string, error) {
	var out Variable
	params := napping.Params{"variable": variable}.AsUrlValues()
	if _, err := s.client.Get("/asterisk/variable", &params, &out); err != nil {
		return "", err
	}
	return out.Value, nil
}

func (s *AsteriskService) SetGlobalVar(variable, value string) error {
	payload := map[string]string{
		"variable": variable,
		"value":    value,
	}
	if _, err := s.client.Post("/asterisk/variable", payload, nil); err != nil {
		return err
	}
	return nil
}
