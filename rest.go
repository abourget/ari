package ari

import (
	"fmt"
	"net/url"

	"github.com/jmcvetta/napping"
)

//
// napping Post/Get/Delete wrappers
//

func (c *Client) Post(url string, payload, results interface{}) (*napping.Response, error) {
	fullUrl := c.makeFullUrl(url)
	var errMsg errorResponse
	c.Log("Sending POST request to %s", fullUrl)
	res, err := c.session.Post(fullUrl, payload, results, &errMsg)
	return c.checkNappingError(res, err, errMsg)
}

func (c *Client) Get(url string, p *url.Values, results interface{}) (*napping.Response, error) {
	fullUrl := c.makeFullUrl(url)
	var errMsg errorResponse
	c.Log("Sending GET request to %s", fullUrl)
	res, err := c.session.Get(fullUrl, p, results, &errMsg)
	return c.checkNappingError(res, err, errMsg)
}

func (c *Client) Delete(url string, results interface{}) (*napping.Response, error) {
	fullUrl := c.makeFullUrl(url)
	var errMsg errorResponse
	c.Log("Sending DELETE request to %s", fullUrl)
	res, err := c.session.Delete(fullUrl, nil, results, &errMsg)
	return c.checkNappingError(res, err, errMsg)
}

type errorResponse struct {
	Message string
}

func (c *Client) makeFullUrl(url string) string {
	return fmt.Sprintf("%s/ari%s", c.endpoint, url)
}

func (c *Client) checkNappingError(res *napping.Response, err error, errMsg errorResponse) (*napping.Response, error) {
	if err == nil {
		status := res.Status()
		if status > 299 {
			err := fmt.Errorf("Non-2XX returned by server (%s)", res.HttpResponse().Status)
			if errMsg.Message != "" {
				err = fmt.Errorf("%s: %s", err.Error(), errMsg.Message)
			}
			c.Log(fmt.Sprintf(" - %s", err.Error()))
			return res, err
		}
	}
	c.Log(" - Success")
	return res, err
}
