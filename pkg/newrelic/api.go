package newrelic

import (
	"encoding/json"
	_ "encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/SennaSemakula/nrobservgen/internal/http"
)

const api = "https://api.eu.newrelic.com/v2"

type Application struct {
	Apps []interface{} `json:"applications"`
}

func getAPMData(c *http.Client, app string) error {
	// Move this out
	url := fmt.Sprintf("%s/applications.json", api)

	headers := map[string]string{
		"API-Key": c.Token,
	}
	params := map[string]string{
		"filter[name]": app,
	}

	req, err := c.NewRequest(url, headers, params)
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
	envVar := "NR_API_KEY"
	if len(os.Getenv(envVar)) < 1 {
		return fmt.Errorf("missing %s environment variable", envVar)
	}
	c := http.NewClient(os.Getenv(envVar))
	if err := getAPMData(c, app); err != nil {
		return err
	}

	return nil
}
