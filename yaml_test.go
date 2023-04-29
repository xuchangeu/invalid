package invalid

import (
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
)

func TestYAML(t *testing.T) {
	testSimpleCase(t)

}

func testSimpleCase(t *testing.T) {
	file, err := os.Open(filepath.Join([]string{"test", "yaml-cases", "simple.yaml"}...))
	assert.Nil(t, err)

	field, err := NewYAML(file)
	assert.Nil(t, err)
	assert.NotNil(t, field)

	m := field.GetField("map")
	assert.NotNil(t, m)
	assert.EqualValues(t, "map", m.GetKey())
	assert.EqualValues(t, FieldKindMapping, m.GetKind())
	assert.EqualValues(t, "", m.GetValue())
	assert.EqualValues(t, ValueTypeObj, m.GetValueType())
	assert.EqualValues(t, nil, m.GetField("notexist"))

	map2 := m.GetField("map2")
	assert.NotNil(t, map2)
	assert.EqualValues(t, "map2", map2.GetKey())
	assert.EqualValues(t, FieldKindMapping, map2.GetKind())
	assert.EqualValues(t, "", map2.GetValue())
	assert.EqualValues(t, ValueTypeObj, map2.GetValueType())
	assert.EqualValues(t, nil, map2.GetField("notexist"))

	strVal := map2.GetField("strVal")
	assert.NotNil(t, strVal)
	assert.EqualValues(t, "strVal", strVal.GetKey())
	assert.EqualValues(t, FieldKindScalar, strVal.GetKind())
	assert.EqualValues(t, "53minute", strVal.GetValue())
	assert.EqualValues(t, ValueTypeStr, strVal.GetValueType())
	assert.EqualValues(t, nil, strVal.GetField("notexist"))

	boolVal := map2.GetField("boolVal")
	assert.NotNil(t, boolVal)
	assert.EqualValues(t, "boolVal", boolVal.GetKey())
	assert.EqualValues(t, FieldKindScalar, boolVal.GetKind())
	assert.EqualValues(t, "true", boolVal.GetValue())
	assert.EqualValues(t, ValueTypeBool, boolVal.GetValueType())
	assert.EqualValues(t, nil, boolVal.GetField("notexist"))

	floatVal := map2.GetField("floatVal")
	assert.NotNil(t, floatVal)
	assert.EqualValues(t, "floatVal", floatVal.GetKey())
	assert.EqualValues(t, FieldKindScalar, floatVal.GetKind())
	assert.EqualValues(t, "1e2", floatVal.GetValue())
	assert.EqualValues(t, ValueTypeFloat, floatVal.GetValueType())
	assert.EqualValues(t, nil, floatVal.GetField("notexist"))

	nilVal := map2.GetField("nilVal")
	assert.NotNil(t, nilVal)
	assert.EqualValues(t, "nilVal", nilVal.GetKey())
	assert.EqualValues(t, FieldKindScalar, nilVal.GetKind())
	assert.EqualValues(t, "null", nilVal.GetValue())
	assert.EqualValues(t, ValueTypeNil, nilVal.GetValueType())
	assert.EqualValues(t, nil, nilVal.GetField("notexist"))

	//list

	list := m.GetField("list")
	assert.NotNil(t, list)
	assert.EqualValues(t, "list", list.GetKey())
	assert.EqualValues(t, FieldKindSequence, list.GetKind())
	assert.EqualValues(t, "", list.GetValue())
	assert.EqualValues(t, ValueTypeArr, list.GetValueType())
	assert.EqualValues(t, nil, list.GetField("notexist"))

	listResult := list.GetFields()
	assert.EqualValues(t, 3, len(listResult))
	for _, v := range listResult {
		assert.EqualValues(t, FieldKindScalar, v.GetKind())
		assert.EqualValues(t, ValueTypeStr, v.GetValueType())
		assert.EqualValues(t, nil, v.GetField("notexist"))
	}

	//list 2
	list2 := m.GetField("list2")
	assert.NotNil(t, list)
	assert.EqualValues(t, "list2", list2.GetKey())
	assert.EqualValues(t, FieldKindSequence, list2.GetKind())
	assert.EqualValues(t, "", list2.GetValue())
	assert.EqualValues(t, ValueTypeArr, list2.GetValueType())
	assert.EqualValues(t, nil, list2.GetField("notexist"))

	listResult = list2.GetFields()
	assert.EqualValues(t, 2, len(listResult))
	for _, v := range listResult {
		foo := v.GetField("foo")
		bar := v.GetField("bar")
		assert.NotNil(t, foo)
		assert.NotNil(t, bar)
		assert.EqualValues(t, ValueTypeStr, foo.GetValueType())
		assert.EqualValues(t, FieldKindScalar, foo.GetKind())
		assert.EqualValues(t, ValueTypeStr, bar.GetValueType())
		assert.EqualValues(t, FieldKindScalar, bar.GetKind())
		assert.EqualValues(t, nil, foo.GetField("notexist"))
		assert.EqualValues(t, nil, bar.GetField("notexist"))
	}

	//list 3
	list3 := m.GetField("list3")
	assert.NotNil(t, list)
	assert.EqualValues(t, "list3", list3.GetKey())
	assert.EqualValues(t, FieldKindSequence, list3.GetKind())
	assert.EqualValues(t, "", list3.GetValue())
	assert.EqualValues(t, ValueTypeArr, list3.GetValueType())
	assert.EqualValues(t, 3, len(list3.GetFields()))
	assert.EqualValues(t, nil, list3.GetField("notexist"))

	//map4
	map4 := m.GetField("map4")
	assert.NotNil(t, map4)
	assert.EqualValues(t, "map4", map4.GetKey())
	assert.EqualValues(t, FieldKindScalar, map4.GetKind())
	assert.EqualValues(t, "", map4.GetValue())
	assert.EqualValues(t, ValueTypeNil, map4.GetValueType())
	assert.EqualValues(t, nil, map4.GetField("notexist"))

	//map5
	map5 := m.GetField("map5")
	assert.NotNil(t, map5)
	assert.EqualValues(t, "map5", map5.GetKey())
	assert.EqualValues(t, FieldKindMapping, map5.GetKind())
	assert.EqualValues(t, "", map5.GetValue())
	assert.EqualValues(t, ValueTypeObj, map5.GetValueType())
	assert.EqualValues(t, nil, map5.GetField("notexist"))

	//strVal2
	strVal2 := map5.GetField("strVal2")
	assert.NotNil(t, strVal2)
	assert.EqualValues(t, "strVal2", strVal2.GetKey())
	assert.EqualValues(t, FieldKindScalar, strVal2.GetKind())
	assert.EqualValues(t, "123", strVal2.GetValue())
	assert.EqualValues(t, ValueTypeInt, strVal2.GetValueType())
	assert.EqualValues(t, nil, strVal2.GetField("notexist"))
}
