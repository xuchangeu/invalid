package invalid

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"gopkg.in/yaml.v3"
	"io"
	"sort"
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

// Field interface
type Field interface {
	restructure(sibling *yaml.Node) error
	getValueRange() *Range
	Key() string
	setKey(key string)
	Value() string
	ValueType() ValueType
	Kind() FieldKind
	Tag() string
	Fields() []Field
	Get(key string) (Field, bool)
	KeyRange() *Range
	ValueRange() *Range
	AddField(key string, field Field)
}

var lines []string

func readLines(by []byte) []string {
	lines = []string{}
	buffer := bytes.NewBuffer(by)
	rd := bufio.NewReader(buffer)
	i := 0
	for {
		b, _, err := rd.ReadLine()
		if err != nil && err == io.EOF {
			break
		} else {
			lines = append(lines, string(b))
			i++
		}
	}
	return lines
}

func NewYAML(r io.Reader) (Field, error) {
	by, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	node := &yaml.Node{}
	err = yaml.Unmarshal(by, node)
	if err != nil {
		return nil, err
	}

	lines = readLines(by)

	if len(node.Content) < 1 {
		return nil, errors.New("document must have at least one field")
	}
	docNode := node.Content[0]

	f, e := NewYamlField(nil, docNode)
	if e != nil {
		return nil, e
	}
	e = f.restructure(nil)
	if e != nil {
		return nil, e
	}
	return f, nil
}

func NewYamlField(keyNode, valueNode *yaml.Node) (Field, error) {
	var fieldInt Field
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
	keyNode     *yaml.Node
	valueNode   *yaml.Node
	siblingNode *yaml.Node
	keyRange    *Range
	valueRange  *Range
	key         string
	value       string
	valueType   ValueType
	kind        FieldKind
	tag         string
	style       yaml.Style
	children    map[string]Field
}

func (f *YAMLField) restructure(sibling *yaml.Node) error {
	//set properties
	f.siblingNode = sibling
	f.setKey("")
	f.setValue()
	f.setKind()
	f.setStyle()
	f.setValueType()
	f.setTag()

	//
	err := f.setKeyRange()
	if err != nil {
		return err
	}
	return nil
}

func (f *YAMLField) Key() string {
	return f.key
}

func (f *YAMLField) Value() string {
	return f.valueNode.Value
}

func (f *YAMLField) ValueType() ValueType {
	return f.valueType
}

func (f *YAMLField) Kind() FieldKind {
	return f.kind
}

func (f *YAMLField) Fields() []Field {
	result := make([]Field, 0)
	for _, v := range f.children {
		result = append(result, v)
	}
	sort.SliceStable(result, func(i, j int) bool {
		return result[i].Key() < result[j].Key()
	})
	return result
}

func (f *YAMLField) Get(key string) (Field, bool) {
	field, exist := f.children[key]
	return field, exist
}

func (f *YAMLField) KeyRange() *Range {
	return f.keyRange
}

func (f *YAMLField) ValueRange() *Range {
	return f.valueRange
}

func (f *YAMLField) AddField(key string, field Field) {
	f.addField(key, field)
}

func (f *YAMLField) setKey(key string) {
	if f.keyNode != nil {
		f.key = f.keyNode.Value
	} else {
		f.key = key
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

func (f *YAMLField) setTag() {
	if f.valueNode != nil {
		f.tag = f.valueNode.Tag
	}
}

func (f *YAMLField) getTag() string {
	return f.tag
}

func (f *YAMLField) Tag() string {
	return f.getTag()
}

func (f *YAMLField) addField(name string, child Field) {
	if f.children == nil {
		f.children = make(map[string]Field)
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

func (f *YAMLField) setKeyRange() error {
	if f.keyNode == nil {
		return nil
	}
	line, err := NewLineByYAMLNode(f.keyNode)
	if err == nil && line != nil {
		r := NewRange(line, line)
		f.keyRange = &r
	}
	return nil
}

func (f *YAMLField) setValueRange(r *Range) {
	f.valueRange = r
}

func (f *YAMLField) getValueRange() *Range {
	//calc self range
	if f.valueNode == nil {
		return nil
	}

	if f.valueRange != nil {
		return f.valueRange
	}

	line, err := NewLineByYAMLNode(f.valueNode)
	if err == nil && line != nil {
		r := NewRange(line, line)
		f.setValueRange(&r)
		return &r
	}

	return nil
}

// YAMLArrField represent an array field
type YAMLArrField struct {
	YAMLField
}

func (field *YAMLArrField) findValuePosition() {
	fmt.Println("sequence field find value position called")
}

func (field *YAMLArrField) restructure(sibling *yaml.Node) error {
	err := field.YAMLField.restructure(sibling)
	if err != nil {
		return err
	}

	//loop children
	content := field.valueNode.Content
	selfRange := field.getValueRange()
	for i := 0; i < len(content); i++ {
		// ruleList type of sequence is mapping
		fieldInt, err := NewYamlField(nil, content[i])
		if err != nil {
			return err
		}

		//replace sib with parent sibling if node has no sibling
		var sib *yaml.Node
		if len(content) > i+1 {
			sib = content[i+1]
		}
		if sib == nil {
			sib = sibling
		}

		//call restructure recursively
		err = fieldInt.restructure(sib)
		if err != nil {
			return err
		}

		//calc range
		r := fieldInt.getValueRange()
		selfRange = selfRange.expend(r)

		//use string format of index as key since node inside list have no key item.
		key := fmt.Sprintf("%d", i)
		fieldInt.setKey(key)
		field.addField(key, fieldInt)
	}
	field.setValueRange(selfRange)

	return nil
}

// YAMLMappingField mapping field with YAML
type YAMLMappingField struct {
	YAMLField
}

func (field *YAMLMappingField) restructure(sibling *yaml.Node) error {
	//call parent restructure
	err := field.YAMLField.restructure(sibling)
	if err != nil {
		return err
	}

	selfRange := field.getValueRange()
	content := field.valueNode.Content
	for i := 0; i < len(content)/2; i++ {
		//paired key value nodes
		keyNode := content[i*2]
		valueNode := content[i*2+1]

		//initialize field interface by nodes
		fieldInt, err := NewYamlField(keyNode, valueNode)
		if err != nil {
			return err
		}

		//resolve child sibling
		var sib *yaml.Node
		if len(content) > i*2+2 {
			sib = content[i*2+2]
		}
		//the child use parent's sibling if child's sibling is nil
		if sib == nil {
			sib = sibling
		}

		//recursion of restructure
		err = fieldInt.restructure(sib)
		if err != nil {
			return err
		}

		//resolve child's value range
		r := fieldInt.getValueRange()
		//log.Printf("solving -> %s , range -> strt : %d, %d %d, end : %d, %d, %d . value ==>%s", fieldInt.Key(), r.Start.Line,
		//	r.Start.ColumnStart, r.Start.ColumnEnd, r.End.Line, r.End.ColumnStart, r.End.ColumnEnd,
		//	lines[int(r.Start.Line)-1][int(r.Start.ColumnStart)-1:int(r.Start.ColumnEnd)-1])
		selfRange = selfRange.expend(r)

		//add field
		field.addField(keyNode.Value, fieldInt)
	}
	field.setValueRange(selfRange)
	return nil
}

// YAMLStrField string field for YAML
type YAMLStrField struct {
	YAMLField
}

func (field *YAMLStrField) restructure(sibling *yaml.Node) error {
	err := field.YAMLField.restructure(sibling)
	if err != nil {
		return err
	}
	return nil
}

// YAMLBoolField boolean field for YAML
type YAMLBoolField struct {
	YAMLField
}

func (field *YAMLBoolField) restructure(sibling *yaml.Node) error {
	err := field.YAMLField.restructure(sibling)
	if err != nil {
		return err
	}
	return nil
}

// YAMLFloatField float field form YAML
type YAMLFloatField struct {
	YAMLField
}

func (field *YAMLFloatField) restructure(sibling *yaml.Node) error {
	err := field.YAMLField.restructure(sibling)
	if err != nil {
		return err
	}
	return nil
}

// YAMLIntField integer field for YAML
type YAMLIntField struct {
	YAMLField
}

func (field *YAMLIntField) restructure(sibling *yaml.Node) error {
	err := field.YAMLField.restructure(sibling)
	if err != nil {
		return err
	}
	return nil
}

// YAMLNilField nil field for YAML
type YAMLNilField struct {
	YAMLField
}

func (field *YAMLNilField) restructure(sibling *yaml.Node) error {
	err := field.YAMLField.restructure(sibling)
	if err != nil {
		return err
	}
	return nil
}
