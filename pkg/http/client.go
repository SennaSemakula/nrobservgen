package http

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}
type Client struct {
	HTTPClient
	Token string
}
type RunbookStore interface {
	GetRunbook(HTTPClient) error
}
type Confluence struct {
	token string
}
type Repository interface {
	GetTags() error
}
type Bitbucket struct {
	url  string
	user string
	pass string
}

func NewBitbucket(url, user, pass string) *Bitbucket {
	return &Bitbucket{
		url,
		user,
		pass,
	}
}

func NewConfluence(token string) *Confluence {
	return &Confluence{
		token,
	}
}

func NewClient(token string) *Client {
	return &Client{
		&http.Client{},
		token,
	}
}

// Required due to zscaler and confluence cert errors on jenkins
func NewInsecureClient(token string) *Client {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	http := &http.Client{Transport: tr}

	return &Client{http, token}
}

func NewRequest(url string, headers, params map[string]string) (*http.Request, error) {
	if len(url) <= 0 {
		return nil, fmt.Errorf("url cannot be empty")
	}
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
func (con *Confluence) GetRunbook(c HTTPClient, url string) error {
	headers := map[string]string{
		"Authorization": "Bearer " + con.token,
	}
	req, err := NewRequest(url, headers, nil)
	if err != nil {
		return err
	}

	resp, err := c.Do(req)
	if err != nil {
		// TODO: this may return valid server errors e.g. confluence is down.
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
func (bb *Bitbucket) GetTags(c *Client) (*Tags, error) {
	req, err := NewRequest(bb.url, nil, nil)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(bb.user, bb.pass)

	resp, err := c.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetching tags from '%s': %v", bb.url, err)
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
