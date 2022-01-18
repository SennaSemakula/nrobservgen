package newrelic

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

func (m *MockClient) Do(req *http.Request) (*http.Response, error) {
	json := `{"applications": ["service1", "service2"]}`
	resp := ioutil.NopCloser(bytes.NewReader([]byte(json)))
	switch string(req.URL.Query().Get("filter[name]")) {
	case "service1":
		return &http.Response{StatusCode: 200, Body: resp}, nil
	case "service2":
		return &http.Response{StatusCode: 200, Body: resp}, nil
	case "noservice":
		return &http.Response{StatusCode: 404, Body: resp}, fmt.Errorf("got status code: 404")
	case "fakeservice":
		return &http.Response{StatusCode: 404, Body: resp}, fmt.Errorf("got status code: 404")
	default:
		return &http.Response{StatusCode: 404, Body: resp}, fmt.Errorf("got status code: 404")
	}
}

func TestGetAPMData(t *testing.T) {
	t.Parallel()
	client := MockClient{token: ""}
	inputs := map[string]struct {
		app      string
		expected error
	}{
		"fetch service1 apm data": {
			app:      "service1",
			expected: nil,
		},
		"fetch creditmanager apm data": {
			app:      "service2",
			expected: nil,
		},
		"fetch noservice apm data": {
			app:      "noservice",
			expected: fmt.Errorf("got status code: 404"),
		},
		"fetch fakeservice apm data": {
			app:      "noservice",
			expected: fmt.Errorf("got status code: 404"),
		},
	}

	for name, tc := range inputs {
		t.Run(name, func(t *testing.T) {
			actual := getAPMData(&client, tc.app, client.token)
			if tc.expected != nil && actual != nil {
				require.Equal(t, tc.expected.Error(), actual.Error())
				return
			}
			require.NoError(t, actual)
		})
	}

}
