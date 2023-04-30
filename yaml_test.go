package invalid

import (
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
)

func TestYAML(t *testing.T) {
	testSimpleCase(t)
	testK8SService(t)
	testVariousValue(t)
}

func testVariousValue(t *testing.T) {
	file, err := os.Open(filepath.Join([]string{"test", "yaml-cases", "various_value.yaml"}...))
	assert.Nil(t, err)

	field, err := NewYAML(file)
	assert.Nil(t, err)
	assert.NotNil(t, field)

	//negative
	negative := field.Get("negative")
	assert.NotNil(t, negative)
	assert.EqualValues(t, "-12", negative.Value())
	assert.EqualValues(t, FieldKindScalar, negative.Kind())
	assert.EqualValues(t, ValueTypeInt, negative.ValueType())

	//zero
	zero := field.Get("zero")
	assert.NotNil(t, zero)
	assert.EqualValues(t, "0", zero.Value())
	assert.EqualValues(t, FieldKindScalar, zero.Kind())
	assert.EqualValues(t, ValueTypeInt, zero.ValueType())

	//positive
	positive := field.Get("positive")
	assert.NotNil(t, positive)
	assert.EqualValues(t, "34", positive.Value())
	assert.EqualValues(t, FieldKindScalar, positive.Kind())
	assert.EqualValues(t, ValueTypeInt, positive.ValueType())

	//canonical
	canonical := field.Get("canonical")
	assert.NotNil(t, canonical)
	assert.EqualValues(t, "12345", canonical.Value())
	assert.EqualValues(t, FieldKindScalar, canonical.Kind())
	assert.EqualValues(t, ValueTypeInt, canonical.ValueType())

	//decimal
	decimal := field.Get("decimal")
	assert.NotNil(t, decimal)
	assert.EqualValues(t, "+12,345", decimal.Value())
	assert.EqualValues(t, FieldKindScalar, decimal.Kind())
	assert.EqualValues(t, ValueTypeStr, decimal.ValueType())

	//sexagesimal
	sexagesimal := field.Get("sexagesimal")
	assert.NotNil(t, sexagesimal)
	assert.EqualValues(t, "3:25:45", sexagesimal.Value())
	assert.EqualValues(t, FieldKindScalar, sexagesimal.Kind())
	assert.EqualValues(t, ValueTypeStr, sexagesimal.ValueType())

	//octal
	octal := field.Get("octal")
	assert.NotNil(t, octal)
	assert.EqualValues(t, "014", octal.Value())
	assert.EqualValues(t, FieldKindScalar, octal.Kind())
	assert.EqualValues(t, ValueTypeInt, octal.ValueType())

	//hexadecimal
	hexadecimal := field.Get("hexadecimal")
	assert.NotNil(t, hexadecimal)
	assert.EqualValues(t, "0xC", hexadecimal.Value())
	assert.EqualValues(t, FieldKindScalar, hexadecimal.Kind())
	assert.EqualValues(t, ValueTypeInt, hexadecimal.ValueType())

	//canonical2
	canonical2 := field.Get("canonical2")
	assert.NotNil(t, canonical2)
	assert.EqualValues(t, "1.23015e+3", canonical2.Value())
	assert.EqualValues(t, FieldKindScalar, canonical2.Kind())
	assert.EqualValues(t, ValueTypeFloat, canonical2.ValueType())

	//exponential
	exponential := field.Get("exponential")
	assert.NotNil(t, exponential)
	assert.EqualValues(t, "12.3015e+02", exponential.Value())
	assert.EqualValues(t, FieldKindScalar, exponential.Kind())
	assert.EqualValues(t, ValueTypeFloat, exponential.ValueType())

	//exponential
	sexagesimal2 := field.Get("sexagesimal2")
	assert.NotNil(t, sexagesimal2)
	assert.EqualValues(t, "20:30.15", sexagesimal2.Value())
	assert.EqualValues(t, FieldKindScalar, sexagesimal2.Kind())
	assert.EqualValues(t, ValueTypeStr, sexagesimal2.ValueType())

	//exponential
	fixed := field.Get("fixed")
	assert.NotNil(t, fixed)
	assert.EqualValues(t, "1,230.15", fixed.Value())
	assert.EqualValues(t, FieldKindScalar, fixed.Kind())
	assert.EqualValues(t, ValueTypeStr, fixed.ValueType())

	//negativeInfinity
	negativeInfinity := field.Get("negativeInfinity")
	assert.NotNil(t, negativeInfinity)
	assert.EqualValues(t, "-.inf", negativeInfinity.Value())
	assert.EqualValues(t, FieldKindScalar, negativeInfinity.Kind())
	assert.EqualValues(t, ValueTypeFloat, negativeInfinity.ValueType())

	//nan
	nan := field.Get("not a number")
	assert.NotNil(t, nan)
	assert.EqualValues(t, ".NaN", nan.Value())
	assert.EqualValues(t, FieldKindScalar, nan.Kind())
	assert.EqualValues(t, ValueTypeFloat, nan.ValueType())

	//null
	null := field.Get("null")
	assert.NotNil(t, nan)
	assert.EqualValues(t, "~", null.Value())
	assert.EqualValues(t, FieldKindScalar, null.Kind())
	assert.EqualValues(t, ValueTypeNil, null.ValueType())

	//true
	tr1 := field.Get("true")
	assert.NotNil(t, tr1)
	assert.EqualValues(t, "y", tr1.Value())
	assert.EqualValues(t, FieldKindScalar, tr1.Kind())
	assert.EqualValues(t, ValueTypeStr, tr1.ValueType())

	//false
	fa1 := field.Get("false")
	assert.NotNil(t, fa1)
	assert.EqualValues(t, "n", fa1.Value())
	assert.EqualValues(t, FieldKindScalar, fa1.Kind())
	assert.EqualValues(t, ValueTypeStr, fa1.ValueType())

	//string
	str := field.Get("string")
	assert.NotNil(t, str)
	assert.EqualValues(t, "12345", str.Value())
	assert.EqualValues(t, FieldKindScalar, str.Kind())
	assert.EqualValues(t, ValueTypeStr, str.ValueType())

	//true2
	tr2 := field.Get("true2")
	assert.NotNil(t, tr2)
	assert.EqualValues(t, "yes", tr2.Value())
	assert.EqualValues(t, FieldKindScalar, tr2.Kind())
	assert.EqualValues(t, ValueTypeBool, tr2.ValueType())

}

func testK8SService(t *testing.T) {
	file, err := os.Open(filepath.Join([]string{"test", "yaml-cases", "k8s-service.yaml"}...))
	assert.Nil(t, err)

	field, err := NewYAML(file)
	assert.Nil(t, err)
	assert.NotNil(t, field)

	//apiVersion
	apiVersion := field.Get("apiVersion")
	assert.NotNil(t, apiVersion)
	assert.EqualValues(t, "apiVersion", apiVersion.Key())
	assert.EqualValues(t, "v1", apiVersion.Value())
	assert.EqualValues(t, FieldKindScalar, apiVersion.Kind())
	assert.EqualValues(t, ValueTypeStr, apiVersion.ValueType())

	//kind
	kind := field.Get("kind")
	assert.NotNil(t, kind)
	assert.EqualValues(t, "kind", kind.Key())
	assert.EqualValues(t, "Service", kind.Value())
	assert.EqualValues(t, FieldKindScalar, kind.Kind())
	assert.EqualValues(t, ValueTypeStr, kind.ValueType())

	//metadata
	metadata := field.Get("metadata")
	assert.NotNil(t, kind)
	assert.EqualValues(t, "metadata", metadata.Key())
	assert.EqualValues(t, "", metadata.Value())
	assert.EqualValues(t, FieldKindMapping, metadata.Kind())
	assert.EqualValues(t, ValueTypeObj, metadata.ValueType())

	//metadata.name
	metadataName := metadata.Get("name")
	assert.NotNil(t, kind)
	assert.EqualValues(t, "name", metadataName.Key())
	assert.EqualValues(t, "my-service", metadataName.Value())
	assert.EqualValues(t, FieldKindScalar, metadataName.Kind())
	assert.EqualValues(t, ValueTypeStr, metadataName.ValueType())

	//spec.selector.app.kubernetes.io/name
	appName := field.Get("spec").Get("selector").Get("app.kubernetes.io/name")
	assert.NotNil(t, appName)
	assert.EqualValues(t, "MyApp", appName.Value())
	assert.EqualValues(t, FieldKindScalar, appName.Kind())
	assert.EqualValues(t, ValueTypeStr, appName.ValueType())

	//ports[0]
	ports := field.Get("spec").Get("ports")
	protocol := ports.Fields()[0].Get("protocol")
	assert.NotNil(t, protocol)
	assert.EqualValues(t, "TCP", protocol.Value())
	assert.EqualValues(t, FieldKindScalar, protocol.Kind())
	assert.EqualValues(t, ValueTypeStr, protocol.ValueType())

	port := ports.Fields()[0].Get("port")
	assert.NotNil(t, port)
	assert.EqualValues(t, "80", port.Value())
	assert.EqualValues(t, FieldKindScalar, port.Kind())
	assert.EqualValues(t, ValueTypeInt, port.ValueType())

	targetPort := ports.Fields()[0].Get("targetPort")
	assert.NotNil(t, targetPort)
	assert.EqualValues(t, "9376", targetPort.Value())
	assert.EqualValues(t, FieldKindScalar, targetPort.Kind())
	assert.EqualValues(t, ValueTypeInt, targetPort.ValueType())

}

func testSimpleCase(t *testing.T) {
	file, err := os.Open(filepath.Join([]string{"test", "yaml-cases", "simple.yaml"}...))
	assert.Nil(t, err)

	field, err := NewYAML(file)
	assert.Nil(t, err)
	assert.NotNil(t, field)

	m := field.Get("map")
	assert.NotNil(t, m)
	assert.EqualValues(t, "map", m.Key())
	assert.EqualValues(t, FieldKindMapping, m.Kind())
	assert.EqualValues(t, "", m.Value())
	assert.EqualValues(t, ValueTypeObj, m.ValueType())
	assert.EqualValues(t, nil, m.Get("notexist"))

	map2 := m.Get("map2")
	assert.NotNil(t, map2)
	assert.EqualValues(t, "map2", map2.Key())
	assert.EqualValues(t, FieldKindMapping, map2.Kind())
	assert.EqualValues(t, "", map2.Value())
	assert.EqualValues(t, ValueTypeObj, map2.ValueType())
	assert.EqualValues(t, nil, map2.Get("notexist"))

	strVal := map2.Get("strVal")
	assert.NotNil(t, strVal)
	assert.EqualValues(t, "strVal", strVal.Key())
	assert.EqualValues(t, FieldKindScalar, strVal.Kind())
	assert.EqualValues(t, "53minute", strVal.Value())
	assert.EqualValues(t, ValueTypeStr, strVal.ValueType())
	assert.EqualValues(t, nil, strVal.Get("notexist"))

	boolVal := map2.Get("boolVal")
	assert.NotNil(t, boolVal)
	assert.EqualValues(t, "boolVal", boolVal.Key())
	assert.EqualValues(t, FieldKindScalar, boolVal.Kind())
	assert.EqualValues(t, "true", boolVal.Value())
	assert.EqualValues(t, ValueTypeBool, boolVal.ValueType())
	assert.EqualValues(t, nil, boolVal.Get("notexist"))

	floatVal := map2.Get("floatVal")
	assert.NotNil(t, floatVal)
	assert.EqualValues(t, "floatVal", floatVal.Key())
	assert.EqualValues(t, FieldKindScalar, floatVal.Kind())
	assert.EqualValues(t, "1e2", floatVal.Value())
	assert.EqualValues(t, ValueTypeFloat, floatVal.ValueType())
	assert.EqualValues(t, nil, floatVal.Get("notexist"))

	nilVal := map2.Get("nilVal")
	assert.NotNil(t, nilVal)
	assert.EqualValues(t, "nilVal", nilVal.Key())
	assert.EqualValues(t, FieldKindScalar, nilVal.Kind())
	assert.EqualValues(t, "null", nilVal.Value())
	assert.EqualValues(t, ValueTypeNil, nilVal.ValueType())
	assert.EqualValues(t, nil, nilVal.Get("notexist"))

	//list

	list := m.Get("list")
	assert.NotNil(t, list)
	assert.EqualValues(t, "list", list.Key())
	assert.EqualValues(t, FieldKindSequence, list.Kind())
	assert.EqualValues(t, "", list.Value())
	assert.EqualValues(t, ValueTypeArr, list.ValueType())
	assert.EqualValues(t, nil, list.Get("notexist"))

	listResult := list.Fields()
	assert.EqualValues(t, 3, len(listResult))
	for _, v := range listResult {
		assert.EqualValues(t, FieldKindScalar, v.Kind())
		assert.EqualValues(t, ValueTypeStr, v.ValueType())
		assert.EqualValues(t, nil, v.Get("notexist"))
	}

	//list 2
	list2 := m.Get("list2")
	assert.NotNil(t, list)
	assert.EqualValues(t, "list2", list2.Key())
	assert.EqualValues(t, FieldKindSequence, list2.Kind())
	assert.EqualValues(t, "", list2.Value())
	assert.EqualValues(t, ValueTypeArr, list2.ValueType())
	assert.EqualValues(t, nil, list2.Get("notexist"))

	listResult = list2.Fields()
	assert.EqualValues(t, 2, len(listResult))
	for _, v := range listResult {
		foo := v.Get("foo")
		bar := v.Get("bar")
		assert.NotNil(t, foo)
		assert.NotNil(t, bar)
		assert.EqualValues(t, ValueTypeStr, foo.ValueType())
		assert.EqualValues(t, FieldKindScalar, foo.Kind())
		assert.EqualValues(t, ValueTypeStr, bar.ValueType())
		assert.EqualValues(t, FieldKindScalar, bar.Kind())
		assert.EqualValues(t, nil, foo.Get("notexist"))
		assert.EqualValues(t, nil, bar.Get("notexist"))
	}

	//list 3
	list3 := m.Get("list3")
	assert.NotNil(t, list)
	assert.EqualValues(t, "list3", list3.Key())
	assert.EqualValues(t, FieldKindSequence, list3.Kind())
	assert.EqualValues(t, "", list3.Value())
	assert.EqualValues(t, ValueTypeArr, list3.ValueType())
	assert.EqualValues(t, 3, len(list3.Fields()))
	assert.EqualValues(t, nil, list3.Get("notexist"))

	//map4
	map4 := m.Get("map4")
	assert.NotNil(t, map4)
	assert.EqualValues(t, "map4", map4.Key())
	assert.EqualValues(t, FieldKindScalar, map4.Kind())
	assert.EqualValues(t, "", map4.Value())
	assert.EqualValues(t, ValueTypeNil, map4.ValueType())
	assert.EqualValues(t, nil, map4.Get("notexist"))

	//map5
	map5 := m.Get("map5")
	assert.NotNil(t, map5)
	assert.EqualValues(t, "map5", map5.Key())
	assert.EqualValues(t, FieldKindMapping, map5.Kind())
	assert.EqualValues(t, "", map5.Value())
	assert.EqualValues(t, ValueTypeObj, map5.ValueType())
	assert.EqualValues(t, nil, map5.Get("notexist"))

	//strVal2
	strVal2 := map5.Get("strVal2")
	assert.NotNil(t, strVal2)
	assert.EqualValues(t, "strVal2", strVal2.Key())
	assert.EqualValues(t, FieldKindScalar, strVal2.Kind())
	assert.EqualValues(t, "123", strVal2.Value())
	assert.EqualValues(t, ValueTypeInt, strVal2.ValueType())
	assert.EqualValues(t, nil, strVal2.Get("notexist"))
}
