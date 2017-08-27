package ari

import (
	"fmt"
	"net/url"

	"github.com/jmcvetta/napping"
)

//
// napping Post/Get/Delete wrappers
//

// Post does POST request
func (c *Client) Post(path string, payload, results interface{}) error {
	_, err := c.PostWithResponse(path, payload, results)
	return err
}

// PostWithResponse does POST request and returns the response
func (c *Client) PostWithResponse(path string, payload, results interface{}) (*napping.Response, error) {
	url := c.makeURL(path)
	var errMsg errorResponse
	c.Log("Sending POST request to %s", url)
	res, err := c.session.Post(url, payload, results, &errMsg)
	if results != nil {
		c.setClientRecurse(results)
	}
	return res, c.checkNappingError(res, err, errMsg)
}

// Get does GET request
func (c *Client) Get(path string, p *url.Values, results interface{}) error {
	url := c.makeURL(path)
	var errMsg errorResponse
	c.Log("Sending GET request to %s", url)
	res, err := c.session.Get(url, p, results, &errMsg)
	if results != nil {
		c.setClientRecurse(results)
	}
	return c.checkNappingError(res, err, errMsg)
}

// Delete does DELETE request
func (c *Client) Delete(path string, results interface{}) error {
	url := c.makeURL(path)
	var errMsg errorResponse
	c.Log("Sending DELETE request to %s", url)
	res, err := c.session.Delete(url, nil, results, &errMsg)
	if results != nil {
		c.setClientRecurse(results)
	}
	return c.checkNappingError(res, err, errMsg)
}

type errorResponse struct {
	Message string
}

func (c *Client) makeURL(path string) string {
	return fmt.Sprintf("%s/ari%s", c.endpoint, path)
}

func (c *Client) checkNappingError(res *napping.Response, err error, errMsg errorResponse) error {
	if err != nil {
		return err
	}
	status := res.Status()
	if status > 299 {
		err := fmt.Errorf("Non-2XX returned by server (%s)", res.HttpResponse().Status)
		if errMsg.Message != "" {
			err = fmt.Errorf("%s: %s", err.Error(), errMsg.Message)
		}
		c.Log(fmt.Sprintf(" - %s", err.Error()))
		return err
	}
	c.Log(" - Success")
	return err
}
