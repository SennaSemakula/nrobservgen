package common

import (
	"io/ioutil"
	"os"
)

func ReadYaml(name string) ([]byte, error) {
	b, err := ioutil.ReadFile(name)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func WriteFile(svc, dest string, b []byte) (string, error) {
	// TODO: stop passing in service. This should be a general function. it currently depends on service
	f, err := os.Create(dest)
	if err != nil {
		return "", err
	}
	if _, err := f.Write(b); err != nil {
		return "", err
	}
	pwd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	config := svc + "-alerts.yaml"
	if dest != config {
		return dest, nil
	}
	return pwd + "/" + config, nil
}

func Contains(a string, items []string) bool {
	for i := range items {
		if items[i] == a {
			return true
		}
	}
	return false
}
