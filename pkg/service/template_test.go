package service

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewYamlTemplate(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	inputs := []struct {
		name string
		dest string
		err  error
	}{
		{"fake_service", dir, nil},
		{"eel_service", "hello", fmt.Errorf("hello: no such file or directory")},
	}
	for _, v := range inputs {
		actual, err := NewYamlTemplate(v.name, v.dest)
		if v.err != nil {
			require.Empty(t, actual)
			require.NotNil(t, err)
			continue
		}
		expected := dir + "/" + v.name + "-alerts.yaml"
		require.NoError(t, err)
		require.NotEmpty(t, actual)
		require.Equal(t, v.err, err)
		require.Equal(t, expected, actual)
	}
	// Check the contents of the file created
}
func TestNewTFTemplate(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	tests := map[string]struct {
		service Service
		input   string
		want    string
		err     bool
	}{
		"valid with default version": {
			service: Service{Name: "test_service", Version: "v0.1.0", Alerts: nil},
			input:   dir,
			want:    dir,
			err:     false,
		},
		"valid with version v0.4.0": {
			service: Service{Name: "fake_service", Version: "v0.4.0", Alerts: nil},
			input:   dir,
			want:    dir,
			err:     false,
		},
		"valid with version v0.4.0 and empty alerts": {
			service: Service{Name: "fake_service", Version: "v0.4.0", Alerts: Alert{}},
			input:   dir,
			want:    dir,
			err:     false,
		},
		"directory that doesn't exist": {
			service: Service{Name: "fake_service", Version: "v0.4.0", Alerts: Alert{}},
			input:   "doesnt_exist_dir",
			want:    "",
			err:     true,
		},
	}

	for name, tc := range tests {
		actual, err := NewTFTemplate(tc.service, tc.input)
		if tc.err {
			require.NotNil(t, err)
			require.Empty(t, actual)
			continue
		}
		require.NoErrorf(t, err, "%d expected no error but got %v", name, err)
		require.NotEmpty(t, actual)
		require.Equal(t, tc.want, actual)
	}
	// As of Go 1.16, os.ReadDir is a more efficient and correct choice:
	// it returns a list of fs.DirEntry instead of fs.FileInfo,
	// and it returns partial results in the case of an error midway through reading a directory
	files, err := ioutil.ReadDir(dir)
	require.NoError(t, err)
	require.Equal(t, 4, len(files))

	contains := func(file string) bool {
		for _, v := range []string{"main.tf", "variables.tf", "backend.tf", "provider.tf"} {
			if file == v {
				return true
			}
		}
		return false
	}
	for _, v := range files {
		if !contains(v.Name()) {
			require.Error(t, fmt.Errorf("terraform template did not generate file: %s", v.Name()))
		}
	}
}
