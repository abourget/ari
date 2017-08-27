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
	return &ai, s.client.Get("/asterisk/info", nil, &ai)
}

func (s *AsteriskService) GetGlobalVar(variable string) (string, error) {
	var out Variable
	params := napping.Params{"variable": variable}.AsUrlValues()
	err := s.client.Get("/asterisk/variable", &params, &out)
	return out.Value, err
}

func (s *AsteriskService) SetGlobalVar(variable, value string) error {
	payload := map[string]string{
		"variable": variable,
		"value":    value,
	}
	return s.client.Post("/asterisk/variable", payload, nil)
}
