package service

import (
	"fmt"
	"os"

	"github.com/SennaSemakula/nrobservgen/pkg/common"
	"gopkg.in/yaml.v3"
)

// Required to populate line numbers
// func (alert AlertConfig) UnmarshalYAML(value *yaml.Node) error {
// 	type config struct {
// 		LineNum       int
// 		Enabled       bool      `yaml:"enabled"`
// 		Runbook       string    `yaml:"runbook"`
// 		WarnThreshold Threshold `yaml:"warning_threshold"`
// 		CritThreshold Threshold `yaml:"critical_threshold"`
// 	}
// 	var params struct {
// 		config
// 	}
// 	err := value.Decode(&params.config)
// 	if err != nil {
// 		return err
// 	}

// 	// Save the line number
// 	params.config.LineNum = value.Line
// 	log.Println(params)
// 	// alerts = params
// 	return nil
// }

// This is now coupled with Service struct. Not common anymore
func LoadYaml(b []byte, svc *Service) error {
	if err := yaml.Unmarshal(b, &svc); err != nil {
		return fmt.Errorf("%v. Is your config formatted correctly?", err)
	}
	if len(svc.Name) == 0 {
		return fmt.Errorf("missing required parameter `service_name`")
	}

	return nil
}

func WriteYaml(svc *Service, dest string) (string, error) {
	// TODO: cleanup. should not be in this function
	// so ugly
	if len(dest) > 0 {
		if _, err := os.Stat(dest); err != nil {
			return "", fmt.Errorf("%s: no such file or directory", dest)
		}
	}
	out, err := yaml.Marshal(svc)
	if err != nil {
		return "", err
	}
	config := svc.Name + "-alerts.yaml"
	manifest := dest + "/" + config
	if len(dest) == 0 {
		manifest = config
	}

	res, err := common.WriteFile(svc.Name, manifest, out)
	if err != nil {
		return "", err
	}

	return res, nil
}

func (s *Service) UnmarshalYAML(value *yaml.Node) error {
	var params struct {
		Service string                 `yaml:"service_name"`
		Version string                 `yaml:"version"`
		Alerts  map[string]AlertConfig `yaml:"alerts"`
	}
	if err := value.Decode(&params); err != nil {
		if _, ok := err.(*yaml.TypeError); !ok {
			return err
		}
		return err
		s.Name = params.Service
		s.Version = params.Version
		s.Alerts = params.Alerts
	}
	s.Name = params.Service
	s.Alerts = params.Alerts
	s.Version = params.Version
	if len(params.Version) == 0 {
		s.Version = "v0.1.0"
	}

	return nil
}
