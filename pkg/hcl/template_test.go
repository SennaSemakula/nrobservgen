package hcl

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

const service = "fake-service"

func TestDefaultProviders(t *testing.T) {
	t.Parallel()
	expected := []*Provider{
		{"newrelic", "newrelic/newrelic"},
		{"aws", "hashicorp/aws"},
	}
	providers := DefaultProviders()
	require.Equal(t, 2, len(providers))
	require.Equal(t, expected, providers)
}

func TestDefaultVars(t *testing.T) {
	t.Parallel()
	expected := []*Variable{
		{"newrelic_account_id", "number", 0},
		{"newrelic_region", "string", "EU"},
		{"newrelic_environment", "string", ""},
		{"newrelic_api_key", "string", ""},
		{"aws_region", "string", "eu-west-1"},
		{"newrelic_local_creds", "bool", false},
		{"service_settings_file", "string", fmt.Sprintf("%s-alerts.yaml", service)},
	}
	actual := DefaultVars(service)
	require.Equal(t, len(expected), len(actual))
	require.Equal(t, expected, actual)
}

func TestDefaultLocalVars(t *testing.T) {
	t.Parallel()
	actual := DefaultLocalVars(service)
	expected := []*Variable{
		{"service_name", "string", "local.settings.service_name"},
		{"nrql_alerts", "string", "try(local.settings.alerts, {})"},
	}
	require.Equal(t, len(expected), len(actual))
	require.Equal(t, expected, actual)
}

func TestNewModule(t *testing.T) {
	t.Parallel()
	m := Module{
		Name:      "fake-module",
		Source:    "fake-source",
		Version:   "0.0.0",
		Variables: nil,
	}
	actual := NewModule(m.Name, m.Source, m.Version, m.Variables)
	require.Equal(t, m, actual)
}

func TestWriteMain(t *testing.T) {
	t.Parallel()
	tmp := t.TempDir()

	inputs := map[string]struct {
		input1 string
		input2 string
		input3 string
		input4 []Module
		want   error
	}{
		"valid variables.tf file with default directory source": {
			input1: service,
			input2: "2.1.0",
			input3: "",
			want:   nil,
		},
		"valid main.tf file with directory source": {
			input1: service,
			input2: "2.1.0",
			input3: tmp + "/",
			input4: []Module{
				{"fakemodule1", "source1", "1.0.0", nil, nil},
				{"fakemodule2", "source2", "5.4.0", nil, nil},
			},
			want: nil,
		},
		"non existent file path": {
			input1: service,
			input2: "2.3.0",
			input3: "dir_does_not_exist",
			input4: []Module{
				{"fakemodule1", "source1", "1.0.0", nil, nil},
				{"fakemodule2", "source2", "5.4.0", nil, nil},
			},
			want: fmt.Errorf("open %s/main.tf: no such file or directory", "dir_does_not_exist"),
		},
		"nil modules": {
			input1: service,
			input2: "2.3.0",
			input3: "dir_does_not_exist",
			want:   fmt.Errorf("open %s/main.tf: no such file or directory", "dir_does_not_exist"),
		},
	}

	for name, tc := range inputs {
		t.Run(name, func(t *testing.T) {
			actual := WriteMain(tc.input1, tc.input2, tc.input3, tc.input4)
			if tc.want != nil && actual != nil {
				require.Equal(t, tc.want.Error(), actual.Error())
				return
			}
			require.Equal(t, tc.want, actual)
		})
	}
}

func TestWriteVars(t *testing.T) {
	t.Parallel()
	tmp := t.TempDir()

	inputs := map[string]struct {
		input1 string
		input2 string
		input3 string
		want   error
	}{
		"valid variables.tf file with default directory source": {
			input1: service,
			input2: "0.1.0",
			input3: "",
			want:   nil,
		},
		"valid variables.tf file with directory source": {
			input1: service,
			input2: "1.0.0",
			input3: tmp + "/",
			want:   nil,
		},
		"non existent file path": {
			input1: service,
			input2: "0.1.0",
			input3: "dir_does_not_exist",
			want:   fmt.Errorf("open %s/variables.tf: no such file or directory", "dir_does_not_exist"),
		},
	}

	for name, tc := range inputs {
		t.Run(name, func(t *testing.T) {
			actual := WriteVars(tc.input1, tc.input2, tc.input3)
			if tc.want != nil && actual != nil {
				require.Equal(t, tc.want.Error(), actual.Error())
				return
			}
			require.Equal(t, tc.want, actual)
		})
	}
}

func TestWriteBackend(t *testing.T) {
	t.Parallel()
	tmp := t.TempDir()

	inputs := map[string]struct {
		// ugly - consider variadic functions args
		input1 string
		input2 string
		input3 string
		want   error
	}{
		"valid backend file with default directory source": {
			input1: service,
			input2: "0.1.0",
			input3: "",
			want:   nil,
		},
		"valid backend file with directory source": {
			input1: service,
			input2: "0.1.0",
			input3: tmp + "/",
			want:   nil,
		},
		"non existent file path": {
			input1: service,
			input2: "0.1.0",
			input3: "dir_does_not_exist",
			want:   fmt.Errorf("open %s/backend.tf: no such file or directory", "dir_does_not_exist"),
		},
	}

	for name, tc := range inputs {
		t.Run(name, func(t *testing.T) {
			actual := WriteBackend(tc.input1, tc.input2, tc.input3)
			if tc.want != nil && actual != nil {
				require.Equal(t, tc.want.Error(), actual.Error())
				return
			}
			require.NoError(t, actual)
			require.Equal(t, tc.want, actual)
			require.FileExists(t, tc.input3+"backend.tf")

		})
	}
}
