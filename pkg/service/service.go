package service

import (
	"fmt"
	"os"

	"github.com/SennaSemakula/nrobservgen/internal/http"
	"github.com/SennaSemakula/nrobservgen/pkg/common"
	"github.com/SennaSemakula/nrobservgen/pkg/newrelic"
)

func NewService(name string) *Service {
	return &Service{name, "v0.1.0", nil}
}

func NewServiceDefaults(name string) *Service {
	return &Service{Name: name, Version: "v0.1.0", Alerts: Alert{
		"throughput": AlertConfig{
			Enabled:       true,
			Query:         fmt.Sprintf(`SELECT rate(count(apm.service.transaction.duration), 1 minute) FROM Metric WHERE appName = '%s'`, name),
			Runbook:       "https://fake-runbook.com",
			CritThreshold: Threshold{1, 3600, "below"},
			WarnThreshold: Threshold{12, 1800, "below"},
		},
		"response_time": AlertConfig{
			Enabled:       true,
			Query:         fmt.Sprintf(`SELECT average(apm.service.transaction.duration) * 1000 FROM Metric WHERE appName = '%s'`, name),
			Runbook:       "https://fake-runbook.com",
			CritThreshold: Threshold{1000, 600, "above"},
			WarnThreshold: Threshold{500, 600, "above"},
		},
		"error_rate": AlertConfig{
			Enabled:       true,
			Query:         fmt.Sprintf(`SELECT count(apm.service.transaction.error.count) / count(apm.service.transaction.duration) * 100 FROM Metric WHERE appName = '%s'`, name),
			Runbook:       "https://fake-runbook.com",
			CritThreshold: Threshold{5, 3600, "above"},
			WarnThreshold: Threshold{2, 1800, "above"},
		},
		"apdex": AlertConfig{
			Enabled:       true,
			Query:         fmt.Sprintf(`SELECT apdex(apm.service.apdex) WHERE appName = '%s' FROM Metric`, name),
			Runbook:       "https://fake-runbook.com",
			CritThreshold: Threshold{0.5, 3600, "below"},
			WarnThreshold: Threshold{0.7, 1800, "below"},
		},
		"cpu_utilization_time": AlertConfig{
			Enabled:       true,
			Query:         fmt.Sprintf(`SELECT average(newrelic.timeslice.value) * 1000 FROM Metric WHERE metricTimesliceName = 'CPU/User/Utilization' AND appName = '%s'`, name),
			Runbook:       "https://fake-runbook.com",
			CritThreshold: Threshold{80, 600, "above"},
			WarnThreshold: Threshold{70, 600, "above"},
		},
		"heap_memory_usage": AlertConfig{
			Enabled:       true,
			Query:         fmt.Sprintf(`SELECT average(newrelic.timeslice.value) * 100 FROM Metric WHERE metricTimesliceName = 'Memory/Heap/Utilization' AND appName = '%s'`, name),
			Runbook:       "https://fake-runbook.com",
			CritThreshold: Threshold{90, 600, "above"},
			WarnThreshold: Threshold{80, 600, "above"},
		},
		// "kafka_consumer_lag": AlertConfig{
		// 	Enabled:       true,
		// 	Runbook:       "https://fake-runbook.com",,
		// 	CritThreshold: Threshold{10, 3600, 0, "below"},
		// 	WarnThreshold: Threshold{5, 1800, 0, "below"},
		// },
		// "kafka_consumer_rate": AlertConfig{
		// 	Enabled:       true,
		// 	Runbook:       "https://fake-runbook.com",,
		// 	CritThreshold: Threshold{10, 3600, 0, "below"},
		// 	WarnThreshold: Threshold{5, 1800, 0, "below"},
		// },
		// "gc_cpu_time": AlertConfig{
		// 	Enabled:       true,
		// 	Runbook:       "https://fake-runbook.com",,
		// 	CritThreshold: Threshold{20, 0, 5, "below"},
		// },
		// "deadlocked_threads": AlertConfig{
		// 	Enabled:       true,
		// 	Runbook:       "https://fake-runbook.com",,
		// 	CritThreshold: Threshold{3, 0, 5, "below"},
		// },
	}}
}

func NewDefaultAlerts() Defaults {
	return []string{
		"error_rate", "cpu_utilization_time", "throughput",
		"response_time", "apdex", "heap_memory_usage",
	}
}

func Validate(cfg string) error {
	svc := NewService("")

	b, err := common.ReadYaml(cfg)
	if err != nil {
		return err
	}
	if err := LoadYaml(b, svc); err != nil {
		return err
	}
	if err := svc.validate(); err != nil {
		return err
	}

	return nil
}

func (s *Service) Load(config string, target *Service) error {
	b, err := common.ReadYaml(config)
	if err != nil {
		return err
	}

	if err := LoadYaml(b, target); err != nil {
		return err
	}

	return nil
}

func (s *Service) validate() error {
	if len(s.Name) <= 0 {
		return fmt.Errorf("missing required parameter `service_name`")
	}
	if len(s.Alerts) <= 0 {
		return fmt.Errorf("missing values for parameter `alerts`")
	}
	if err := validVersion(s.Version); err != nil {
		return err
	}
	if err := newrelic.GetAPMData(s.Name); err != nil {
		return err
	}
	if ok := s.Alerts.valid(); !ok {
		return fmt.Errorf("found invalid config")
	}

	valid, err := validateRunbooks(s.Alerts)
	if err != nil {
		return err
	}
	if !valid {
		return fmt.Errorf("found invalid runbooks")
	}

	return nil
}

func validVersion(vers string) error {
	c := http.NewClient("")
	user := os.Getenv("GITHUB_USER")
	pass := os.Getenv("GITHUB_PASS")

	if len(user) < 1 || len(pass) < 1 {
		return fmt.Errorf("missing GITHUB_USER or GITHUB_PASS environment variables")
	}
	// Move this out
	url := "ENTER_URL_HERE"
	path := "ENTER_PATH_HERE"

	contains := func(vers string, tags http.Tags) bool {
		for _, v := range tags.Tags {
			if vers == v.Version {
				return true
			}
		}
		return false
	}
	remoteTags, err := http.GetTags(c, user, pass, url+path)
	if err != nil {
		return fmt.Errorf("unable to find terraform module version %s: %v", vers, err)
	}
	if !contains(vers, *remoteTags) {
		return fmt.Errorf("unable to find terraform module version: %s", vers)
	}
	return nil
}

func validateRunbooks(alerts map[string]AlertConfig) (bool, error) {
	envVar := "CONFLUENCE_API_KEY"
	token := os.Getenv(envVar)
	if len(token) < 1 {
		return false, fmt.Errorf("missing %s", envVar)
	}
	c := http.NewInsecureClient(token)

	var invalid int
	for _, v := range alerts {
		// Alert is disabled so no need to perform validation checks on runbook
		if !v.Enabled {
			continue
		}
		if err := http.GetRunbook(c, v.Runbook); err != nil {
			fmt.Println(err)
			invalid++
		}
	}

	return invalid <= 0, nil
}

func (a Alert) valid() bool {
	defaults := NewDefaultAlerts()
	var invalid int
	for _, v := range defaults {
		if _, ok := a[v]; !ok {
			fmt.Printf("missing required alert `%s`\n", v)
			invalid++
		}
	}
	for _, config := range a {
		if err := config.valid(); err != nil {
			fmt.Println(err)
			invalid++
		}
	}
	return invalid <= 0
}

func (a *AlertConfig) valid() error {
	if err := a.CritThreshold.valid(); err != nil {
		return err
	}
	if err := a.WarnThreshold.valid(); err != nil {
		return err
	}

	return nil
}

func (t *Threshold) valid() error {
	validCondition := func(condition string) bool {
		for _, v := range []string{"below", "above"} {
			if condition == v {
				return true
			}
		}
		return false
	}
	// Critical threshold checks
	if !validCondition(t.Condition) {
		return fmt.Errorf("critical_threshold.condition only accepts 'above' or 'below': got %s", t.Condition)
	}
	if t.Value < 0 {
		return fmt.Errorf("critical_threshold.value cannot be negative: got %g", t.Value)
	}
	if t.Duration < 0 {
		return fmt.Errorf("critical_threshold.duration_secs cannot be negative: got %d", t.Duration)
	}

	// Warning threshold checks
	if !validCondition(t.Condition) {
		return fmt.Errorf("warning_threshold.condition only accepts 'above' or 'below': got %s", t.Condition)
	}
	if t.Value < 0 {
		return fmt.Errorf("warning_threshold.value cannot be negative: got %g", t.Value)
	}
	if t.Duration < 0 {
		return fmt.Errorf("warning_threshold.duration_secs cannot be negative: got %d", t.Duration)
	}
	return nil
}
