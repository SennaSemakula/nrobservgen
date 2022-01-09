package service

type Alert map[string]AlertConfig

// Do I need this?
type Defaults []string

type AlertConfig struct {
	Enabled       bool      `yaml:"enabled"`
	Query         string    `yaml:"query"`
	Runbook       string    `yaml:"runbook"`
	WarnThreshold Threshold `yaml:"warning_threshold,omitempty"`
	CritThreshold Threshold `yaml:"critical_threshold"`
}

type Threshold struct {
	Value     float32 `yaml:"threshold"`
	Duration  int32   `yaml:"duration_sec"`
	Condition string  `yaml:"condition"`
	// Enabled     bool    `yaml:"enabled"`
}

type Service struct {
	Name    string `yaml:"service_name"`
	Version string `yaml:"version"`
	Alerts  Alert  `yaml:"alerts"`
}
