package data

import (
	"io/fs"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCrds(t *testing.T) {
	data, err := Crds()
	assert.NoError(t, err)
	files := []string{
		"envoy.kyverno.io_authorizationpolicies.yaml",
	}
	for _, file := range files {
		file, err := fs.Stat(data, file)
		assert.NoError(t, err)
		assert.NotNil(t, file)
		assert.False(t, file.IsDir())
	}
}
