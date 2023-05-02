package invalid

import (
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"path/filepath"
	"testing"
)

func TestValid(t *testing.T) {
	yamlValid(t)
}

func BenchmarkValid(b *testing.B) {
	file, err := os.Open(filepath.Join([]string{"test", "yaml-cases", "valid.yaml"}...))
	if err != nil {
		log.Printf(err.Error())
		return
	}
	field, err := NewYAML(file)
	if err != nil {
		log.Printf(err.Error())
		return
	}

	file, err = os.OpenFile(filepath.Join("test", "exam", "valid.yaml"), os.O_RDONLY, os.ModeSticky)
	if err != nil {
		log.Printf(err.Error())
		return
	}

	rule, err := NewRule(file)
	if err != nil {
		log.Printf(err.Error())
		return
	}

	for i := 0; i < b.N; i++ {
		rule.Validate(field)
		//log.Printf("%v", result)
	}
}

func yamlValid(t *testing.T) {
	file, err := os.Open(filepath.Join([]string{"test", "yaml-cases", "valid.yaml"}...))
	assert.Nil(t, err)

	field, err := NewYAML(file)
	assert.Nil(t, err)
	assert.NotNil(t, field)

	file, err = os.OpenFile(filepath.Join("test", "exam", "valid.yaml"), os.O_RDONLY, os.ModeSticky)
	assert.Nil(t, err)
	assert.NotNil(t, file)

	rule, err := NewRule(file)
	assert.Nil(t, err)
	assert.NotNil(t, rule)

	result := rule.Validate(field)
	assert.NotNil(t, result)
	assert.EqualValues(t, 0, len(result))
}
