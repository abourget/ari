package ari

import (
	"fmt"
	"net/url"

	"github.com/jmcvetta/napping"
)

type REST struct {
	endpoint string
	session  *napping.Session
}

func NewRest(endpoint, username, password string) *REST {
	userinfo := url.UserPassword(username, password)

	return &REST{
		endpoint: endpoint,
		session: &napping.Session{
			Userinfo:        userinfo,
			UnsafeBasicAuth: true,
		},
	}
}

func (r *REST) Post(url string, payload, results, errMsg interface{}) (*napping.Response, error) {
	fullUrl := fmt.Sprintf("%s%s", r.endpoint, url)

	res, err := r.session.Post(fullUrl, payload, results, errMsg)
	if err == nil {
		if res.Status() != 200 {
			return nil, fmt.Errorf("Non-200 returned by server: %s", res.HttpResponse().Status)
		}
	}
	return res, err
}

func (r *REST) Get(url string, p *napping.Params, results, errMsg interface{}) (*napping.Response, error) {
	fullUrl := fmt.Sprintf("%s/ari%s", r.endpoint, url)

	res, err := r.session.Get(fullUrl, p, results, errMsg)
	if err == nil {
		if res.Status() != 200 {
			return nil, fmt.Errorf("Non-200 returned by server: %s", res.HttpResponse().Status)
		}
	}
	return res, err
}

func (r *REST) AsteriskInfoGet() (*AsteriskInfo, error) {
	ai := AsteriskInfo{}

	if _, err := r.Get("/asterisk/info", nil, &ai, nil); err != nil {
		return nil, err
	}
	return &ai, nil
}

func (r *REST) AsteriskVariableGet(variable string) (string, error) {
	var out Variable

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
