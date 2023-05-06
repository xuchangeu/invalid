package invalid

import (
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
)

func TestYAML(t *testing.T) {
	testSimpleCase(t)
	//log.Println("===================")
	testK8SService(t)
	//log.Println("===================")
	testVariousValue(t)
}

func BenchmarkYAML(b *testing.B) {
	file, _ := os.Open(filepath.Join([]string{"test", "yaml-cases", "various_value.yaml"}...))
	b.ResetTimer() //reset benchmark timer
	for i := 0; i < b.N; i++ {
		NewYAML(file)
	}
}

func testVariousValue(t *testing.T) {
	file, err := os.Open(filepath.Join([]string{"test", "yaml-cases", "various_value.yaml"}...))
	assert.Nil(t, err)

	field, err := NewYAML(file)
	assert.Nil(t, err)
	assert.NotNil(t, field)

	//negative
	negative, _ := field.Get("negative")
	assert.NotNil(t, negative)
	assert.EqualValues(t, "-12", negative.Value())
	assert.EqualValues(t, FieldKindScalar, negative.Kind())
	assert.EqualValues(t, ValueTypeInt, negative.ValueType())

	//zero
	zero, _ := field.Get("zero")
	assert.NotNil(t, zero)
	assert.EqualValues(t, "0", zero.Value())
	assert.EqualValues(t, FieldKindScalar, zero.Kind())
	assert.EqualValues(t, ValueTypeInt, zero.ValueType())

	//positive
	positive, _ := field.Get("positive")
	assert.NotNil(t, positive)
	assert.EqualValues(t, "34", positive.Value())
	assert.EqualValues(t, FieldKindScalar, positive.Kind())
	assert.EqualValues(t, ValueTypeInt, positive.ValueType())

	//canonical
	canonical, _ := field.Get("canonical")
	assert.NotNil(t, canonical)
	assert.EqualValues(t, "12345", canonical.Value())
	assert.EqualValues(t, FieldKindScalar, canonical.Kind())
	assert.EqualValues(t, ValueTypeInt, canonical.ValueType())

	//decimal
	decimal, _ := field.Get("decimal")
	assert.NotNil(t, decimal)
	assert.EqualValues(t, "+12,345", decimal.Value())
	assert.EqualValues(t, FieldKindScalar, decimal.Kind())
	assert.EqualValues(t, ValueTypeStr, decimal.ValueType())

	//sexagesimal
	sexagesimal, _ := field.Get("sexagesimal")
	assert.NotNil(t, sexagesimal)
	assert.EqualValues(t, "3:25:45", sexagesimal.Value())
	assert.EqualValues(t, FieldKindScalar, sexagesimal.Kind())
	assert.EqualValues(t, ValueTypeStr, sexagesimal.ValueType())

	//octal
	octal, _ := field.Get("octal")
	assert.NotNil(t, octal)
	assert.EqualValues(t, "014", octal.Value())
	assert.EqualValues(t, FieldKindScalar, octal.Kind())
	assert.EqualValues(t, ValueTypeInt, octal.ValueType())

	//hexadecimal
	hexadecimal, _ := field.Get("hexadecimal")
	assert.NotNil(t, hexadecimal)
	assert.EqualValues(t, "0xC", hexadecimal.Value())
	assert.EqualValues(t, FieldKindScalar, hexadecimal.Kind())
	assert.EqualValues(t, ValueTypeInt, hexadecimal.ValueType())

	//canonical2
	canonical2, _ := field.Get("canonical2")
	assert.NotNil(t, canonical2)
	assert.EqualValues(t, "1.23015e+3", canonical2.Value())
	assert.EqualValues(t, FieldKindScalar, canonical2.Kind())
	assert.EqualValues(t, ValueTypeFloat, canonical2.ValueType())

	//exponential
	exponential, _ := field.Get("exponential")
	assert.NotNil(t, exponential)
	assert.EqualValues(t, "12.3015e+02", exponential.Value())
	assert.EqualValues(t, FieldKindScalar, exponential.Kind())
	assert.EqualValues(t, ValueTypeFloat, exponential.ValueType())

	//exponential
	sexagesimal2, _ := field.Get("sexagesimal2")
	assert.NotNil(t, sexagesimal2)
	assert.EqualValues(t, "20:30.15", sexagesimal2.Value())
	assert.EqualValues(t, FieldKindScalar, sexagesimal2.Kind())
	assert.EqualValues(t, ValueTypeStr, sexagesimal2.ValueType())

	//exponential
	fixed, _ := field.Get("fixed")
	assert.NotNil(t, fixed)
	assert.EqualValues(t, "1,230.15", fixed.Value())
	assert.EqualValues(t, FieldKindScalar, fixed.Kind())
	assert.EqualValues(t, ValueTypeStr, fixed.ValueType())

	//negativeInfinity
	negativeInfinity, _ := field.Get("negativeInfinity")
	assert.NotNil(t, negativeInfinity)
	assert.EqualValues(t, "-.inf", negativeInfinity.Value())
	assert.EqualValues(t, FieldKindScalar, negativeInfinity.Kind())
	assert.EqualValues(t, ValueTypeFloat, negativeInfinity.ValueType())

	//nan
	nan, _ := field.Get("not a number")
	assert.NotNil(t, nan)
	assert.EqualValues(t, ".NaN", nan.Value())
	assert.EqualValues(t, FieldKindScalar, nan.Kind())
	assert.EqualValues(t, ValueTypeFloat, nan.ValueType())

	//null
	null, _ := field.Get("null")
	assert.NotNil(t, nan)
	assert.EqualValues(t, "~", null.Value())
	assert.EqualValues(t, FieldKindScalar, null.Kind())
	assert.EqualValues(t, ValueTypeNil, null.ValueType())

	//true
	tr1, _ := field.Get("true")
	assert.NotNil(t, tr1)
	assert.EqualValues(t, "y", tr1.Value())
	assert.EqualValues(t, FieldKindScalar, tr1.Kind())
	assert.EqualValues(t, ValueTypeStr, tr1.ValueType())

	//false
	fa1, _ := field.Get("false")
	assert.NotNil(t, fa1)
	assert.EqualValues(t, "n", fa1.Value())
	assert.EqualValues(t, FieldKindScalar, fa1.Kind())
	assert.EqualValues(t, ValueTypeStr, fa1.ValueType())

	//string
	str, _ := field.Get("string")
	assert.NotNil(t, str)
	assert.EqualValues(t, "12345", str.Value())
	assert.EqualValues(t, FieldKindScalar, str.Kind())
	assert.EqualValues(t, ValueTypeStr, str.ValueType())

	//true2
	tr2, _ := field.Get("true2")
	assert.NotNil(t, tr2)
	assert.EqualValues(t, "yes", tr2.Value())
	assert.EqualValues(t, FieldKindScalar, tr2.Kind())
	assert.EqualValues(t, ValueTypeStr, tr2.ValueType())

	//true3
	tr3, _ := field.Get("true3")
	assert.NotNil(t, tr3)
	assert.EqualValues(t, "true", tr3.Value())
	assert.EqualValues(t, FieldKindScalar, tr3.Kind())
	assert.EqualValues(t, ValueTypeBool, tr3.ValueType())

	//true4
	tr4, _ := field.Get("true4")
	assert.NotNil(t, tr4)
	assert.EqualValues(t, "false", tr4.Value())
	assert.EqualValues(t, FieldKindScalar, tr4.Kind())
	assert.EqualValues(t, ValueTypeBool, tr4.ValueType())

}

func testK8SService(t *testing.T) {
	file, err := os.Open(filepath.Join([]string{"test", "yaml-cases", "k8s-service.yaml"}...))
	assert.Nil(t, err)

	field, err := NewYAML(file)
	assert.Nil(t, err)
	assert.NotNil(t, field)

	//apiVersion
	apiVersion, _ := field.Get("apiVersion")
	assert.NotNil(t, apiVersion)
	assert.EqualValues(t, "apiVersion", apiVersion.Key())
	assert.EqualValues(t, "v1", apiVersion.Value())
	assert.EqualValues(t, FieldKindScalar, apiVersion.Kind())
	assert.EqualValues(t, ValueTypeStr, apiVersion.ValueType())

	//kind
	kind, _ := field.Get("kind")
	assert.NotNil(t, kind)
	assert.EqualValues(t, "kind", kind.Key())
	assert.EqualValues(t, "Service", kind.Value())
	assert.EqualValues(t, FieldKindScalar, kind.Kind())
	assert.EqualValues(t, ValueTypeStr, kind.ValueType())

	//metadata
	metadata, _ := field.Get("metadata")
	assert.NotNil(t, kind)
	assert.EqualValues(t, "metadata", metadata.Key())
	assert.EqualValues(t, "", metadata.Value())
	assert.EqualValues(t, FieldKindMapping, metadata.Kind())
	assert.EqualValues(t, ValueTypeObj, metadata.ValueType())

	//metadata.name
	metadataName, _ := metadata.Get("name")
	assert.NotNil(t, kind)
	assert.EqualValues(t, "name", metadataName.Key())
	assert.EqualValues(t, "my-service", metadataName.Value())
	assert.EqualValues(t, FieldKindScalar, metadataName.Kind())
	assert.EqualValues(t, ValueTypeStr, metadataName.ValueType())

	spec, _ := field.Get("spec")
	selector, _ := spec.Get("selector")
	//spec.selector.app.kubernetes.io/name
	appName, _ := selector.Get("app.kubernetes.io/name")
	assert.NotNil(t, appName)
	assert.EqualValues(t, "MyApp", appName.Value())
	assert.EqualValues(t, FieldKindScalar, appName.Kind())
	assert.EqualValues(t, ValueTypeStr, appName.ValueType())

	//ports[0]

	ports, _ := spec.Get("ports")
	protocol, _ := ports.Fields()[0].Get("protocol")
	assert.NotNil(t, protocol)
	assert.EqualValues(t, "TCP", protocol.Value())
	assert.EqualValues(t, FieldKindScalar, protocol.Kind())
	assert.EqualValues(t, ValueTypeStr, protocol.ValueType())

	port, _ := ports.Fields()[0].Get("port")
	assert.NotNil(t, port)
	assert.EqualValues(t, "80", port.Value())
	assert.EqualValues(t, FieldKindScalar, port.Kind())
	assert.EqualValues(t, ValueTypeInt, port.ValueType())

	targetPort, _ := ports.Fields()[0].Get("targetPort")
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

	m, _ := field.Get("map")
	ne, _ := m.Get("notexist")
	assert.NotNil(t, m)
	assert.EqualValues(t, "map", m.Key())
	assert.EqualValues(t, FieldKindMapping, m.Kind())
	assert.EqualValues(t, "", m.Value())
	assert.EqualValues(t, ValueTypeObj, m.ValueType())
	assert.EqualValues(t, nil, ne)

	map2, _ := m.Get("map2")
	ne, _ = map2.Get("notexist")
	assert.NotNil(t, map2)
	assert.EqualValues(t, "map2", map2.Key())
	assert.EqualValues(t, FieldKindMapping, map2.Kind())
	assert.EqualValues(t, "", map2.Value())
	assert.EqualValues(t, ValueTypeObj, map2.ValueType())
	assert.EqualValues(t, nil, ne)

	strVal, _ := map2.Get("strVal")
	ne, _ = strVal.Get("notexist")
	assert.NotNil(t, strVal)
	assert.EqualValues(t, "strVal", strVal.Key())
	assert.EqualValues(t, FieldKindScalar, strVal.Kind())
	assert.EqualValues(t, "53minute", strVal.Value())
	assert.EqualValues(t, ValueTypeStr, strVal.ValueType())
	assert.EqualValues(t, nil, ne)

	boolVal, _ := map2.Get("boolVal")
	ne, _ = strVal.Get("notexist")
	assert.NotNil(t, boolVal)
	assert.EqualValues(t, "boolVal", boolVal.Key())
	assert.EqualValues(t, FieldKindScalar, boolVal.Kind())
	assert.EqualValues(t, "true", boolVal.Value())
	assert.EqualValues(t, ValueTypeBool, boolVal.ValueType())
	assert.EqualValues(t, nil, ne)

	floatVal, _ := map2.Get("floatVal")
	ne, _ = floatVal.Get("notexist")
	assert.NotNil(t, floatVal)
	assert.EqualValues(t, "floatVal", floatVal.Key())
	assert.EqualValues(t, FieldKindScalar, floatVal.Kind())
	assert.EqualValues(t, "1e2", floatVal.Value())
	assert.EqualValues(t, ValueTypeFloat, floatVal.ValueType())
	assert.EqualValues(t, nil, ne)

	nilVal, _ := map2.Get("nilVal")
	ne, _ = nilVal.Get("notexist")
	assert.NotNil(t, nilVal)
	assert.EqualValues(t, "nilVal", nilVal.Key())
	assert.EqualValues(t, FieldKindScalar, nilVal.Kind())
	assert.EqualValues(t, "null", nilVal.Value())
	assert.EqualValues(t, ValueTypeNil, nilVal.ValueType())
	assert.EqualValues(t, nil, ne)

	//list

	list, _ := m.Get("list")
	ne, _ = list.Get("notexist")
	assert.NotNil(t, list)
	assert.EqualValues(t, "list", list.Key())
	assert.EqualValues(t, FieldKindSequence, list.Kind())
	assert.EqualValues(t, "", list.Value())
	assert.EqualValues(t, ValueTypeArr, list.ValueType())
	assert.EqualValues(t, nil, ne)

	listResult := list.Fields()
	assert.EqualValues(t, 3, len(listResult))
	for _, v := range listResult {
		ne, _ = v.Get("notexist")
		assert.EqualValues(t, FieldKindScalar, v.Kind())
		assert.EqualValues(t, ValueTypeStr, v.ValueType())
		assert.EqualValues(t, nil, ne)
	}

	//list 2
	list2, _ := m.Get("list2")
	ne, _ = list2.Get("notexist")
	assert.NotNil(t, list)
	assert.EqualValues(t, "list2", list2.Key())
	assert.EqualValues(t, FieldKindSequence, list2.Kind())
	assert.EqualValues(t, "", list2.Value())
	assert.EqualValues(t, ValueTypeArr, list2.ValueType())
	assert.EqualValues(t, nil, ne)

	listResult = list2.Fields()
	assert.EqualValues(t, 2, len(listResult))
	for _, v := range listResult {
		foo, _ := v.Get("foo")
		bar, _ := v.Get("bar")
		ne1, _ := foo.Get("notexist")
		ne2, _ := bar.Get("notexist")
		assert.NotNil(t, foo)
		assert.NotNil(t, bar)
		assert.EqualValues(t, ValueTypeStr, foo.ValueType())
		assert.EqualValues(t, FieldKindScalar, foo.Kind())
		assert.EqualValues(t, ValueTypeStr, bar.ValueType())
		assert.EqualValues(t, FieldKindScalar, bar.Kind())
		assert.EqualValues(t, nil, ne1)
		assert.EqualValues(t, nil, ne2)
	}

	//list 3
	list3, _ := m.Get("list3")
	ne, _ = list3.Get("notexist")
	assert.NotNil(t, list)
	assert.EqualValues(t, "list3", list3.Key())
	assert.EqualValues(t, FieldKindSequence, list3.Kind())
	assert.EqualValues(t, "", list3.Value())
	assert.EqualValues(t, ValueTypeArr, list3.ValueType())
	assert.EqualValues(t, 3, len(list3.Fields()))
	assert.EqualValues(t, nil, ne)

	//map4
	map4, _ := m.Get("map4")
	ne, _ = map4.Get("notexist")
	assert.NotNil(t, map4)
	assert.EqualValues(t, "map4", map4.Key())
	assert.EqualValues(t, FieldKindScalar, map4.Kind())
	assert.EqualValues(t, "", map4.Value())
	assert.EqualValues(t, ValueTypeNil, map4.ValueType())
	assert.EqualValues(t, nil, ne)

	//map5
	map5, _ := m.Get("map5")
	ne, _ = map5.Get("notexist")
	assert.NotNil(t, map5)
	assert.EqualValues(t, "map5", map5.Key())
	assert.EqualValues(t, FieldKindMapping, map5.Kind())
	assert.EqualValues(t, "", map5.Value())
	assert.EqualValues(t, ValueTypeObj, map5.ValueType())
	assert.EqualValues(t, nil, ne)

	//strVal2
	strVal2, _ := map5.Get("strVal2")
	ne, _ = strVal2.Get("notexist")
	assert.NotNil(t, strVal2)
	assert.EqualValues(t, "strVal2", strVal2.Key())
	assert.EqualValues(t, FieldKindScalar, strVal2.Kind())
	assert.EqualValues(t, "123", strVal2.Value())
	assert.EqualValues(t, ValueTypeInt, strVal2.ValueType())
	assert.EqualValues(t, nil, ne)
}
