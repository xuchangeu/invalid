package invalid

import (
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
)

func TestRule(t *testing.T) {
	testRule1(t)
	testConstraintOfInvalid(t)
	testConstraintOfInvalid2(t)
	testConstraintOfValid(t)
}

func testConstraintOfInvalid(t *testing.T) {
	file, err := os.OpenFile(filepath.Join("test", "yaml-cases", "constraint_of_invalid.yaml"), os.O_RDONLY, os.ModeSticky)
	assert.Nil(t, err)
	assert.NotNil(t, file)

	ruler, err := NewRule(file)
	assert.NotNil(t, err)
	assert.Nil(t, ruler)
	assert.Equal(t, err, OfTypeError("strVal.2", string(RuleTypeStr)))

}

func testConstraintOfInvalid2(t *testing.T) {
	file, err := os.OpenFile(filepath.Join("test", "yaml-cases", "constraint_of_invalid2.yaml"), os.O_RDONLY, os.ModeSticky)
	assert.Nil(t, err)
	assert.NotNil(t, file)

	ruler, err := NewRule(file)
	assert.NotNil(t, err)
	assert.Nil(t, ruler)
	assert.Equal(t, ConstraintTypeError("strVal", yamlNodeTypeSeq), err)

}

func testConstraintOfValid(t *testing.T) {
	file, err := os.OpenFile(filepath.Join("test", "yaml-cases", "constraint_of_valid.yaml"), os.O_RDONLY, os.ModeSticky)
	assert.Nil(t, err)
	assert.NotNil(t, file)

	ruler, err := NewRule(file)
	assert.Nil(t, err)
	assert.NotNil(t, ruler)

}

func testRule1(t *testing.T) {
	file, err := os.OpenFile(filepath.Join("test", "exam", "simple.yaml"), os.O_RDONLY, os.ModeSticky)
	assert.Nil(t, err)
	assert.NotNil(t, file)

	rule, err := NewRule(file)
	assert.Nil(t, err)
	assert.NotNil(t, rule)

	//map
	m, _ := rule.Get("map")
	assert.EqualValues(t, RuleTypeObj, m.RuleType())
	assert.EqualValues(t, true, m.Required())
	mx := m.(*ObjRule)
	assert.EqualValues(t, ".*", mx.GetKeyReg().String())

	//map2
	m2, _ := m.Get("map2")
	assert.NotNil(t, m2)
	assert.EqualValues(t, RuleTypeObj, m2.RuleType())
	assert.EqualValues(t, true, m2.Required())

	//map5
	m5, _ := m.Get("map5")
	assert.NotNil(t, m5)
	assert.EqualValues(t, RuleTypeObj, m5.RuleType())
	assert.EqualValues(t, false, m5.Required())

	//strVal
	strVal, _ := m2.Get("strVal")
	assert.NotNil(t, strVal)
	assert.EqualValues(t, RuleTypeStr, strVal.RuleType())
	assert.True(t, strVal.Required())

	strValX, valid := strVal.(*StrRule)
	assert.EqualValues(t, true, valid)
	assert.EqualValues(t, ".*", strValX.regexp.String())
	assert.EqualValues(t, 20, strValX.max)
	assert.EqualValues(t, 10, strValX.min)

	//list
	list, _ := m2.Get("list")
	assert.NotNil(t, list)
	assert.EqualValues(t, RuleTypeArr, list.RuleType())
	assert.True(t, list.Required())

	listX, valid := list.(*ArrRule)
	assert.True(t, valid)
	assert.EqualValues(t, RuleTypeStr, listX.constraint)

	//list2
	list2, _ := m2.Get("list2")
	assert.NotNil(t, list2)
	assert.EqualValues(t, RuleTypeArr, list2.RuleType())
	assert.True(t, list2.Required())

	list2X, valid := list2.(*ArrRule)
	assert.True(t, valid)
	assert.NotNil(t, list2X.constraint)

	//list2.constraint.name
	list2Cx := list2X.constraint.(Ruler)
	list2Name, _ := list2Cx.Get("name")
	assert.NotNil(t, list2Name)
	assert.EqualValues(t, RuleTypeStr, list2Name.RuleType())
	assert.True(t, list2Name.Required())

	//list2.constraint.description
	list2Desc, _ := list2Cx.Get("description")
	assert.NotNil(t, list2Desc)
	assert.EqualValues(t, RuleTypeStr, list2Desc.RuleType())
	assert.True(t, list2Desc.Required())

	list2DescX, valid := list2Desc.(*StrRule)
	assert.True(t, valid)
	assert.EqualValues(t, ".*", list2DescX.GetReg().String())

	//map4
	m4, _ := m2.Get("map4")
	assert.NotNil(t, m4)
	assert.EqualValues(t, RuleTypeObj, m4.RuleType())
	assert.True(t, m4.Required())

	//m5.strVal
	strVal2, _ := m5.Get("strVal2")
	assert.EqualValues(t, RuleTypeStr, strVal2.RuleType())
	assert.True(t, strVal2.Required())
}
