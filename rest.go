package ari

import (
	"fmt"

	"github.com/jmcvetta/napping"
)

//
// napping Post/Get/Delete wrappers
//

func (c *Client) Post(url string, payload, results, errMsg interface{}) (*napping.Response, error) {
	fullUrl := c.makeFullUrl(url)
	c.Log("Sending POST request to %s", fullUrl)
	res, err := c.session.Post(fullUrl, payload, results, errMsg)
	return c.checkNappingError(res, err)
}

func (c *Client) Get(url string, p *napping.Params, results, errMsg interface{}) (*napping.Response, error) {
	fullUrl := c.makeFullUrl(url)
	c.Log("Sending GET request to %s", fullUrl)
	res, err := c.session.Get(c.makeFullUrl(url), p, results, errMsg)
	return c.checkNappingError(res, err)
}

func (c *Client) Delete(url string, results, errMsg interface{}) (*napping.Response, error) {
	fullUrl := c.makeFullUrl(url)
	c.Log("Sending DELETE request to %s", fullUrl)
	res, err := c.session.Delete(fullUrl, results, errMsg)
	return c.checkNappingError(res, err)
}

func (c *Client) makeFullUrl(url string) string {
	return fmt.Sprintf("%s/ari%s", c.endpoint, url)
}

func (c *Client) checkNappingError(res *napping.Response, err error) (*napping.Response, error) {
	if err == nil {
		status := res.Status()
		if status > 299 {
			c.Log(" - Non-2XX returned by server: %s", res.HttpResponse().Status)
			return res, fmt.Errorf("Non-2XX returned by server: %s", res.HttpResponse().Status)
		}
	}
	c.Log(" - Success")
	return res, err
}
