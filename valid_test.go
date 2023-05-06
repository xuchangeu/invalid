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
	yamlKeyMissing(t)
	yamlTypeMismatch(t)
	testSwagger(t)
}

func testSwagger(t *testing.T) {
	file, err := os.Open(filepath.Join([]string{"test", "yaml-cases", "openapi.yaml"}...))
	assert.Nil(t, err)

	field, err := NewYAML(file)
	assert.Nil(t, err)
	assert.NotNil(t, field)

	file, err = os.OpenFile(filepath.Join("test", "exam", "openapi.yaml"), os.O_RDONLY, os.ModeSticky)
	assert.Nil(t, err)
	assert.NotNil(t, file)

	rule, err := NewRule(file)
	assert.Nil(t, err)
	assert.NotNil(t, rule)

	result := rule.Validate(field)
	assert.EqualValues(t, 2, len(result))
}

func yamlTypeMismatch(t *testing.T) {
	file, err := os.Open(filepath.Join([]string{"test", "yaml-cases", "type_mismatch.yaml"}...))
	assert.Nil(t, err)

	field, err := NewYAML(file)
	assert.Nil(t, err)
	assert.NotNil(t, field)

	file, err = os.OpenFile(filepath.Join("test", "exam", "type_mismatch.yaml"), os.O_RDONLY, os.ModeSticky)
	assert.Nil(t, err)
	assert.NotNil(t, file)

	rule, err := NewRule(file)
	assert.Nil(t, err)
	assert.NotNil(t, rule)

	errs := rule.Validate(field)
	assert.NotNil(t, errs)
	assert.EqualValues(t, 8, len(errs))
	for i := range errs {
		assert.EqualValues(t, TypeMismatch, errs[i].Type)
	}
	assert.EqualValues(t, NewTypeMismatchError("stringVal", string(RuleTypeStr)), errs[0].Error)
	assert.EqualValues(t, NewTypeMismatchError("intVal", string(RuleTypeInt)), errs[1].Error)
	assert.EqualValues(t, NewTypeMismatchError("booVal", string(RuleTypeBool)), errs[2].Error)
	assert.EqualValues(t, NewTypeMismatchError("floatVal", string(RuleTypeFloat)), errs[3].Error)
	assert.EqualValues(t, NewTypeMismatchError("nullVal", string(RuleTypeNil)), errs[4].Error)
	assert.EqualValues(t, NewTypeMismatchError("list.0", string(RuleTypeInt)), errs[5].Error)
	assert.EqualValues(t, NewTypeMismatchError("list.1", string(RuleTypeInt)), errs[6].Error)
	assert.EqualValues(t, NewTypeMismatchError("list.2", string(RuleTypeInt)), errs[7].Error)

}

func yamlKeyMissing(t *testing.T) {
	file, err := os.Open(filepath.Join([]string{"test", "yaml-cases", "key_missing.yaml"}...))
	assert.Nil(t, err)

	field, err := NewYAML(file)
	assert.Nil(t, err)
	assert.NotNil(t, field)

	file, err = os.OpenFile(filepath.Join("test", "exam", "key_missing.yaml"), os.O_RDONLY, os.ModeSticky)
	assert.Nil(t, err)
	assert.NotNil(t, file)

	rule, err := NewRule(file)
	assert.Nil(t, err)
	assert.NotNil(t, rule)

	errs := rule.Validate(field)
	assert.NotNil(t, errs)
	assert.EqualValues(t, 1, len(errs))
	assert.EqualValues(t, NewKeyMissingError("bar1"), errs[0].Error)
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

func BenchmarkValid(b *testing.B) {
	b.StopTimer()
	file, err := os.Open(filepath.Join([]string{"test", "yaml-cases", "valid.yaml"}...))

	if err != nil {
		log.Printf(err.Error())
		return
	}
	b.StartTimer()

	field, err := NewYAML(file)
	if err != nil {
		log.Printf(err.Error())
		return
	}

	b.StopTimer()
	file, err = os.OpenFile(filepath.Join("test", "exam", "valid.yaml"), os.O_RDONLY, os.ModeSticky)
	if err != nil {
		log.Printf(err.Error())
		return
	}

	b.StartTimer()
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
