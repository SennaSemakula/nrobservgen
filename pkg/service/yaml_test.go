package service

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLoadYaml(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		input1 []byte
		input2 Service
		want   error
	}{
		"valid yaml config": {
			input1: []byte(`
            service_name: myservice
            version: v0.1.0
            alerts:
                apdex:
                    enabled: true
                    query: SELECT apdex(apm.service.apdex) WHERE appName = 'fakeservice' FROM Metric
                    runbook: https://fake-runbook.com
                    warning_threshold:
                        threshold: 0.7
                        duration_sec: 1800
                        condition: below
                    critical_threshold:
                        threshold: 0.5
                        duration_sec: 3600
                        condition: below
                cpu_utilization_time:
                    enabled: true
                    query: SELECT average(newrelic.timeslice.value) * 1000 FROM Metric WHERE metricTimesliceName = 'CPU/User/Utilization' AND appName = 'fakeservice'
                    runbook: https://fake-runbook.com
                    warning_threshold:
                        threshold: 70
                        duration_sec: 600
                        condition: above
                    critical_threshold:
                        threshold: 80
                        duration_sec: 600
                        condition: above
                error_rate:
                    enabled: true
                    query: SELECT count(apm.service.transaction.error.count) / count(apm.service.transaction.duration) * 100 FROM Metric WHERE appName = 'fakeservice'
                    runbook: https://fake-runbook.com
                    warning_threshold:
                        threshold: 2
                        duration_sec: 1800
                        condition: above
                    critical_threshold:
                        threshold: 5
                        duration_sec: 3600
                        condition: above
                heap_memory_usage:
                    enabled: true
                    query: SELECT average(newrelic.timeslice.value) * 100 FROM Metric WHERE metricTimesliceName = 'Memory/Heap/Utilization' AND appName = 'fakeservice'
                    runbook: https://fake-runbook.com
                    warning_threshold:
                        threshold: 80
                        duration_sec: 600
                        condition: above
                    critical_threshold:
                        threshold: 90
                        duration_sec: 600
                        condition: above
                response_time:
                    enabled: true
                    query: SELECT average(apm.service.transaction.duration) * 1000 FROM Metric WHERE appName = 'fakeservice'
                    runbook: https://fake-runbook.com
                    warning_threshold:
                        threshold: 500
                        duration_sec: 600
                        condition: above
                    critical_threshold:
                        threshold: 1000
                        duration_sec: 600
                        condition: above
                throughput:
                    enabled: true
                    query: SELECT rate(count(apm.service.transaction.duration), 1 minute) FROM Metric WHERE appName = 'fakeservice'
                    runbook: https://fake-runbook.com
                    warning_threshold:
                        threshold: 12
                        duration_sec: 1800
                        condition: below
                    critical_threshold:
                        threshold: 1
                        duration_sec: 3600
                        condition: below			
            `),
			input2: Service{},
			want:   nil,
		},
		"missing service_name attribute": {
			input1: []byte(`
            version: v0.1.0
            alerts:
                apdex:
                    enabled: false
                response_time:
                    enabled: false
                throughput:
                    enabled: false
                error_rate:
                    enabled: false
            `),
			input2: Service{},
			want:   fmt.Errorf("missing required parameter `service_name`"),
		},
		"wrong enabled attribute type": {
			input1: []byte(`
            service_name: myservice
            version: v0.1.0
            alerts:
                apdex:
                    enabled: 1232
            `),
			input2: Service{},
			want:   fmt.Errorf("cannot unmarshal !!int `1232` into bool. Is your config formatted correctly?"),
		},
		"wrong duration attribute type": {
			input1: []byte(`
            service_name: myservice
            version: v0.1.0
            alerts:
                apdex:
                    critical_threshold:
                        condition: below
                        duration: hello
            `),
			input2: Service{},
			want:   fmt.Errorf("Is your config formatted correctly?"),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			actual := LoadYaml(tc.input1, &tc.input2)
			t.Log(tc.input2)
			if tc.want != nil && actual != nil {
				require.Contains(t, actual.Error(), tc.want.Error())
				require.NotNil(t, actual)
				return
			}
			require.NoError(t, actual)
		})
	}
}
