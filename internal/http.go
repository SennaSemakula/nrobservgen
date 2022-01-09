package http

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Client struct {
	http.Client
	Token string
}

func NewClient(token string) *Client {
	return &Client{
		Token: token,
	}
}

func NewInsecureClient(token string) *Client {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	http := &http.Client{Transport: tr}

	return &Client{*http, token}
}

func (c *Client) NewRequest(url string, headers, params map[string]string) (*http.Request, error) {
	var endpoint = url
	if len(params) > 0 {
		for k, v := range params {
			endpoint += "?" + k + "=" + v
		}
	}
	req, err := http.NewRequest("GET", endpoint, nil)

	if err != nil {
		return nil, err
	}

	for k, v := range headers {
		req.Header.Add(k, v)
	}

	return req, nil
}

// Should be moved out of this file. Make this a receiver
func GetRunbook(c *Client, url string) error {
	headers := map[string]string{
		"Authorization": "Bearer " + c.Token,
	}
	req, err := c.NewRequest(url, headers, nil)
	if err != nil {
		return err
	}

	resp, err := c.Do(req)
	if err != nil {
		// TODO: this may return valid server errors e.g. server is down.
		// find a better way to identify broken links
		return fmt.Errorf("invalid runbook: %s", url)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return fmt.Errorf("invalid runbook: %s. got status code %d", url, resp.StatusCode)
	}
	return nil
}

type Tags struct {
	Tags []struct {
		Version string `json:"displayId"`
	} `json:"values"`
}

// Should be moved out of this file. Make this a receiver
func GetTags(c *Client, user, pass, url string) (*Tags, error) {
	req, err := c.NewRequest(url, nil, nil)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(user, pass)

	resp, err := c.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetching tags from '%s': %v", url, err)
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	// Move this logic out
	var tags Tags
	if err := json.Unmarshal(b, &tags); err != nil {
		return nil, err
	}

	return &tags, nil
}
