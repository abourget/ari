package rest

import (
	"github.com/abourget/ari/models"
	"github.com/jmcvetta/napping"
)

//
// Asterisk endpoints wrappers
//

func (r *REST) AsteriskInfoGet() (*models.AsteriskInfo, error) {
	ai := models.AsteriskInfo{}

	if _, err := r.Get("/asterisk/info", nil, &ai, nil); err != nil {
		return nil, err
	}
	return &ai, nil
}

func (r *REST) AsteriskVariableGet(variable string) (string, error) {
	var out models.Variable

	if _, err := r.Get("/asterisk/variable", &napping.Params{"variable": variable}, &out, nil); err != nil {
		return "", err
	}
	return out.Value, nil
}

func (r *REST) AsteriskVariablePost(variable, value string) error {
	payload := map[string]string{
		"variable": variable,
		"value":    value,
	}
	if _, err := r.Post("/asterisk/variable", payload, nil, nil); err != nil {
		return err
	}
	return nil
}
