package test

import (
	"github.com/stretchr/testify/assert"
	"github.com/xucheng/invalid"
	"os"
	"path/filepath"
	"testing"
)

func loadTestCase(files []string) (*invalid.YAMLRoot, error) {

	file, err := os.Open(filepath.Join(files...))
	if err != nil {
		return nil, err
	}

	r, err := invalid.NewYAML(file)
	if err != nil {
		return nil, err
	}

	return r, nil
}

func TestLinting(t *testing.T) {
	r, err := loadTestCase([]string{"yaml-cases", "test2.yaml"})
	r.Valid()
	assert.Nil(t, err)
	assert.NotNil(t, r)

}
