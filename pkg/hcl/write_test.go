package hcl

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewBody(t *testing.T) {
	version := "0.1.0"
	actual := NewBody(version)
	require.NotNil(t, actual.Body)
	require.NotNil(t, actual.file)
}
