package newrelic

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/SennaSemakula/nrobservgen/pkg/http"
)

const api = "https://api.eu.newrelic.com/v2"

type Application struct {
	Apps []interface{} `json:"applications"`
}

func getAPMData(c http.HTTPClient, app, token string) error {
	// Move this out
	url := fmt.Sprintf("%s/applications.json", api)

	headers := map[string]string{
		"API-Key": token,
	}
	params := map[string]string{
		"filter[name]": app,
	}

	req, err := http.NewRequest(url, headers, params)
	if err != nil {
		return err
	}

	resp, err := c.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode > 299 || resp.StatusCode < 200 {
		return fmt.Errorf(req.Method+" %s got status code: %d", req.URL, resp.StatusCode)
	}
	defer resp.Body.Close()
	// Check if APM data exists
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	// Possibly seperate this out
	var application Application
	if err := json.Unmarshal(b, &application); err != nil {
		return err
	}

	if len(application.Apps) < 1 {
		return fmt.Errorf("no APM data found for app: %s", app)
	}

	return nil

}

func GetAPMData(app string) error {
	envVar := os.Getenv("NR_API_KEY")
	if len(os.Getenv(envVar)) < 1 {
		return errors.New("missing NR_API_KEY environment variable")
	}
	c := http.NewClient(envVar)
	if err := getAPMData(c, app, envVar); err != nil {
		return err
	}

	return nil
}
