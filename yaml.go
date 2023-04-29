package invalid

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"io"
)

type FieldKind uint32

// an abstract of field kinds
const (
	FieldKindSequence FieldKind = 1 << iota
	FieldKindMapping
	FieldKindScalar
	FieldKindDocument
)

type ValueType string

const (
	ValueTypeNil   ValueType = "$nil"
	ValueTypeAny   ValueType = "$any"
	ValueTypeBool  ValueType = "$bool"
	ValueTypeInt   ValueType = "$int"
	ValueTypeFloat ValueType = "$float"
	ValueTypeStr   ValueType = "$str"
	ValueTypeObj   ValueType = "$obj" //an object value contains sub-ruleMap inside, mostly it's a map
	ValueTypeSeq   ValueType = "$seq" //a list with value in any type
	ValueTypeArr   ValueType = "$arr"
)

// FieldInt Field interface
type FieldInt interface {
	Restructure() error
	GetKey() string
	GetValue() string
	GetValueType() ValueType
	GetKind() FieldKind
	GetFields() []FieldInt
	GetField(key string) FieldInt
	GetKeyRange() []*FieldRange
	GetValueRange() []*FieldRange
	position() error
	AddField(key string, fieldInt FieldInt)
}

func NewYAML(r io.Reader) (FieldInt, error) {
	by, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	node := &yaml.Node{}
	err = yaml.Unmarshal(by, node)
	if err != nil {
		return nil, err
	}

	f, e := NewYamlField(nil, node.Content[0])
	if e != nil {
		return nil, e
	}
	e = f.Restructure()
	if e != nil {
		return nil, e
	}
	return f, nil
}

func NewYamlField(keyNode, valueNode *yaml.Node) (FieldInt, error) {
	var fieldInt FieldInt
	if validMapNode(valueNode) {
		fieldInt = &YAMLMappingField{YAMLField{
			keyNode:   keyNode,
			valueNode: valueNode,
		}}
	} else if validArrNode(valueNode) {
		fieldInt = &YAMLArrField{YAMLField{
			keyNode:   keyNode,
			valueNode: valueNode,
		}}
	} else if validStrNode(valueNode) {
		fieldInt = &YAMLStrField{YAMLField{
			keyNode:   keyNode,
			valueNode: valueNode,
		}}
	} else if validBoolNode(valueNode) {
		fieldInt = &YAMLBoolField{YAMLField{
			keyNode:   keyNode,
			valueNode: valueNode,
		}}
	} else if validIntNode(valueNode) {
		fieldInt = &YAMLIntField{YAMLField{
			keyNode:   keyNode,
			valueNode: valueNode,
		}}
	} else if validFloatNode(valueNode) {
		fieldInt = &YAMLFloatField{YAMLField{
			keyNode:   keyNode,
			valueNode: valueNode,
		}}
	} else if validNilNode(valueNode) {
		fieldInt = &YAMLNilField{YAMLField{
			keyNode:   keyNode,
			valueNode: valueNode,
		}}
	}

	return fieldInt, nil
}

type YAMLField struct {
	keyNode    *yaml.Node
	valueNode  *yaml.Node
	keyRange   []*FieldRange
	valueRange []*FieldRange
	key        string
	value      string
	valueType  ValueType
	kind       FieldKind
	style      yaml.Style
	children   map[string]FieldInt
}

func (f *YAMLField) Restructure() error {
	f.setKey()
	f.setValue()
	f.setKind()
	f.setStyle()
	f.setValueType()
	return nil
}

func (f *YAMLField) GetKey() string {
	return f.key
}

func (f *YAMLField) GetValue() string {
	return f.valueNode.Value
}

func (f *YAMLField) GetValueType() ValueType {
	return f.valueType
}

func (f *YAMLField) GetKind() FieldKind {
	return f.kind
}

func (f *YAMLField) GetFields() []FieldInt {
	result := make([]FieldInt, 0)
	for _, v := range f.children {
		result = append(result, v)
	}
	return result
}

func (f *YAMLField) GetField(key string) FieldInt {
	field, exist := f.children[key]
	if exist {
		return field
	}
	return nil
}

func (f *YAMLField) GetKeyRange() []*FieldRange {
	return f.keyRange
}

func (f *YAMLField) GetValueRange() []*FieldRange {
	return f.valueRange
}

func (f *YAMLField) AddField(key string, field FieldInt) {
	f.addField(key, field)
}

func (f *YAMLField) setKey() {
	if f.keyNode != nil {
		f.key = f.keyNode.Value
	}
}

func (f *YAMLField) setKind() {
	if f.valueNode == nil {
		return
	}
	var kind FieldKind
	switch f.valueNode.Kind {
	case yaml.DocumentNode:
		kind = FieldKindDocument
	case yaml.MappingNode:
		kind = FieldKindMapping
	case yaml.ScalarNode:
		kind = FieldKindScalar
	case yaml.SequenceNode:
		kind = FieldKindSequence
	default:
		kind = FieldKindScalar
	}
	f.kind = kind
}

func (f *YAMLField) setValue() {
	if f.valueNode != nil {
		f.value = f.valueNode.Value
	}
}

func (f *YAMLField) setValueType() {
	if f.valueNode == nil {
		return
	}

	if validMapNode(f.valueNode) {
		f.valueType = ValueTypeObj
	} else if validArrNode(f.valueNode) {
		f.valueType = ValueTypeArr
	} else if validStrNode(f.valueNode) {
		f.valueType = ValueTypeStr
	} else if validIntNode(f.valueNode) {
		f.valueType = ValueTypeInt
	} else if validFloatNode(f.valueNode) {
		f.valueType = ValueTypeFloat
	} else if validBoolNode(f.valueNode) {
		f.valueType = ValueTypeBool
	} else if validNilNode(f.valueNode) {
		f.valueType = ValueTypeNil
	}
}

func (f *YAMLField) setStyle() {
	if f.valueNode != nil {
		f.style = f.valueNode.Style
	}
}

// make user function called after all field fieldList filled
func (f *YAMLField) position() error {
	f.findKeyPosition()
	f.findValuePosition()

	if f.children != nil && len(f.children) > 0 {
		for _, v := range f.children {
			//v.(FieldInt).FindPosition()
			err := v.position()
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (f *YAMLField) findKeyPosition() {
	if f.keyNode != nil {
		var l, cs, ce int
		l = f.keyNode.Line
		cs = f.keyNode.Column
		ce = f.keyNode.Column + len(f.key)
		f.keyRange = []*FieldRange{
			{
				Line:        l,
				ColumnStart: cs,
				ColumnEnd:   ce,
			},
		}
	}
}

func (f *YAMLField) findValuePosition() {
	// TODO: find the value position for folded-kind node,in which case node should find the sibling or next topper level node position to determinate itself.
	if f.valueNode != nil {
		var l, cs, ce int
		l = f.valueNode.Column
		cs = f.valueNode.Column
		ce = f.valueNode.Column + len(f.value)

		r1 := &FieldRange{
			Line:        l,
			ColumnStart: cs,
			ColumnEnd:   ce,
		}
		f.keyRange = []*FieldRange{
			r1,
		}
	}
}

func (f *YAMLField) addField(name string, child FieldInt) {
	if f.children == nil {
		f.children = make(map[string]FieldInt)
	}
	f.children[name] = child
}

func (f *YAMLField) sameKind() bool {
	//node kind in list are not equal
	for i := 0; i < len(f.valueNode.Content)-1; i++ {
		if f.valueNode.Content[i].Kind != f.valueNode.Content[i+1].Kind {
			return false
		}
	}

	return true
}

// YAMLArrField represent an array field
type YAMLArrField struct {
	YAMLField
}

func (field *YAMLArrField) findValuePosition() {
	fmt.Println("sequence field find value position called")
}

func (field *YAMLArrField) Restructure() error {
	err := field.YAMLField.Restructure()
	if err != nil {
		return err
	}
	//if field.sameKind() {
	for i := 0; i < len(field.valueNode.Content); i++ {
		// ruleList type of sequence is mapping
		fieldInt, err := NewYamlField(nil, field.valueNode.Content[i])
		if err != nil {
			return err
		}
		err = fieldInt.Restructure()
		if err != nil {
			return err
		}
		key := fmt.Sprintf("%d", i)
		field.addField(key, fieldInt)
	}
	//} else {
	//	return errors.New(fmt.Sprintf("array should have some type inside : [%s]", field.GetKey()))
	//}

	return nil
}

type YAMLMappingField struct {
	YAMLField
}

func (field *YAMLMappingField) Restructure() error {
	err := field.YAMLField.Restructure()
	if err != nil {
		return err
	}
	root := field.valueNode
	for i := 0; i < len(root.Content)/2; i++ {
		keyNode := root.Content[i*2]
		valueNode := root.Content[i*2+1]
		fieldInt, err := NewYamlField(keyNode, valueNode)
		if err != nil {
			return err
		}

		err = fieldInt.Restructure()
		if err != nil {
			return err
		}
		field.addField(keyNode.Value, fieldInt)
	}
	return nil
}

type YAMLStrField struct {
	YAMLField
}

func (field *YAMLStrField) Restructure() error {
	err := field.YAMLField.Restructure()
	if err != nil {
		return err
	}
	return nil
}

type YAMLBoolField struct {
	YAMLField
}

func (field *YAMLBoolField) Restructure() error {
	err := field.YAMLField.Restructure()
	if err != nil {
		return err
	}
	return nil
}

type YAMLFloatField struct {
	YAMLField
}

func (field *YAMLFloatField) Restructure() error {
	err := field.YAMLField.Restructure()
	if err != nil {
		return err
	}
	return nil
}

type YAMLIntField struct {
	YAMLField
}

func (field *YAMLIntField) Restructure() error {
	err := field.YAMLField.Restructure()
	if err != nil {
		return err
	}
	return nil
}

type YAMLNilField struct {
	YAMLField
}

func (field *YAMLNilField) Restructure() error {
	err := field.YAMLField.Restructure()
	if err != nil {
		return err
	}
	return nil
}
