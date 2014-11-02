package rest

// Package rest implements the Asterisk ARI REST interface. See: https://wiki.asterisk.org/wiki/display/AST/Asterisk+12+ARI

import (
	"fmt"
	"log"
	"net/url"

	"github.com/jmcvetta/napping"
)

type REST struct {
	Debug    bool
	endpoint string
	session  *napping.Session
}

func New(endpoint, username, password string) *REST {
	userinfo := url.UserPassword(username, password)

	return &REST{
		endpoint: endpoint,
		session: &napping.Session{
			Userinfo:        userinfo,
			UnsafeBasicAuth: true,
		},
	}
}

func (r *REST) Log(format string, v ...interface{}) {
	if r.Debug {
		log.Printf(fmt.Sprintf("%s\n", format), v...)
	}
}

//
// napping Post/Get/Delete wrappers
//

func (r *REST) Post(url string, payload, results, errMsg interface{}) (*napping.Response, error) {
	fullUrl := r.makeFullUrl(url)
	r.Log("Sending POST request to %s", fullUrl)
	res, err := r.session.Post(fullUrl, payload, results, errMsg)
	return r.checkNappingError(res, err)
}

func (r *REST) Get(url string, p *napping.Params, results, errMsg interface{}) (*napping.Response, error) {
	fullUrl := r.makeFullUrl(url)
	r.Log("Sending GET request to %s", fullUrl)
	res, err := r.session.Get(r.makeFullUrl(url), p, results, errMsg)
	return r.checkNappingError(res, err)
}

func (r *REST) Delete(url string, results, errMsg interface{}) (*napping.Response, error) {
	fullUrl := r.makeFullUrl(url)
	r.Log("Sending DELETE request to %s", fullUrl)
	res, err := r.session.Delete(fullUrl, results, errMsg)
	return r.checkNappingError(res, err)
}

func (r *REST) makeFullUrl(url string) string {
	return fmt.Sprintf("%s/ari%s", r.endpoint, url)
}

func (r *REST) checkNappingError(res *napping.Response, err error) (*napping.Response, error) {
	if err == nil {
		status := res.Status()
		if status > 299 {
			r.Log(" - Non-2XX returned by server: %s", res.HttpResponse().Status)
			return res, fmt.Errorf("Non-2XX returned by server: %s", res.HttpResponse().Status)
		}
	}
	r.Log(" - Success")
	return res, err
}
