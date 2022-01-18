package http

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

type MockClient struct {
	token string
}

func (c *MockClient) Do(req *http.Request) (*http.Response, error) {
	json := `{"runbook_contents": "The server is on fire!"}`
	r := ioutil.NopCloser(bytes.NewReader([]byte(json)))
	switch req.URL.Host + req.URL.Path {
	case "<your valid link>":
		return &http.Response{StatusCode: 200, Body: r}, nil
	case "<your invalid link>":
		return &http.Response{StatusCode: 404, Body: r}, nil
	case "dead-server.com":
		return &http.Response{StatusCode: 500, Body: r}, nil
	default:
		return &http.Response{StatusCode: 404, Body: r}, nil
	}
}

func TestNewClient(t *testing.T) {
	t.Parallel()
	actual := NewClient("CLIENT_TOKEN")
	require.NotNil(t, actual)
	require.NotNil(t, actual.HTTPClient)
	require.Equal(t, "CLIENT_TOKEN", actual.Token)
}

func TestNewInsecureClient(t *testing.T) {
	t.Parallel()
	actual := NewInsecureClient("FAKE_TOKEN")
	require.NotNil(t, actual)
	require.NotNil(t, actual.HTTPClient)
	require.Equal(t, "FAKE_TOKEN", actual.Token)
}

func TestGetConfluenceRunbook(t *testing.T) {
	t.Parallel()
	client := &MockClient{token: "MOCK_TOKEN"}
	inputs := map[string]struct {
		url      string
		expected error
	}{
		"Valid runbook": {
			url:      "https://confluence.co.uk",
			expected: nil,
		},
		"Dead runbook link": {
			url: "https://dead-runbook.com",
			// How does one self reference url field in Golang? :) If you find out, let me know!
			expected: fmt.Errorf("invalid runbook: https://dead-runbook.com. got status code 404"),
		},
		"Unresponsive server": {
			url:      "https://dead-server.com",
			expected: fmt.Errorf("invalid runbook: https://dead-server.com. got status code 500"),
		},
	}

	con := NewConfluence("FAKE_TOKEN")
	for name, tc := range inputs {
		t.Run(name, func(t *testing.T) {
			actual := con.GetRunbook(client, tc.url)
			if tc.expected != nil {
				require.Equal(t, tc.expected.Error(), actual.Error())
				return
			}
			require.Equal(t, tc.expected, actual)
		})
	}
}
func TestNewRequest(t *testing.T) {
	t.Parallel()
	input := map[string]struct {
		url      string
		headers  map[string]string
		params   map[string]string
		expected error
	}{
		"Valid request": {
			url:      "www.google.com",
			headers:  nil,
			params:   nil,
			expected: nil,
		},
		"Empty url": {
			url: "",
			headers: map[string]string{
				"API_TOKEN":  "FAKE_API_TOKEN",
				"USER_AGENT": "MOZILLA MAC",
			},
			params:   nil,
			expected: fmt.Errorf("url cannot be empty"),
		},
		"Request with headers": {
			url: "www.test-url.com",
			headers: map[string]string{
				"API_TOKEN":  "FAKE_API_TOKEN",
				"USER_AGENT": "MOZILLA MAC",
			},
		},
		"Request with url parameters": {
			url: "www.test-url.com",
			params: map[string]string{
				"PAGE": "12",
			},
		},
	}

	for name, tc := range input {
		t.Run(name, func(t *testing.T) {
			actual, err := NewRequest(tc.url, tc.headers, tc.params)
			if tc.expected != nil && err != nil {
				require.Equal(t, tc.expected.Error(), err.Error())
				return
			}
			require.NoError(t, err)
			require.NotNil(t, actual)
			// Check headers
			for header := range tc.headers {
				require.NotEmpty(t, actual.Header.Get(header), "cannot find %s in request headers")
			}
			// Check url parameters
			for key := range tc.params {
				require.NotEmptyf(t, actual.URL.Query().Get(key), "cannot find %s parameter in url", key)
			}
		})
	}

}
