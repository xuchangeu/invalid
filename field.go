package invalid

import (
	"fmt"
	"gopkg.in/yaml.v3"
)

const (
	sequenceField FieldKind = 1 << iota
	mappingField
	scalarField
)

// Field
// There are several of kins of 'NodeKind' as following
// DocumentNode : the kind of the root of document (value 1)
// SequenceNode : slice (value 2)
// MappingNode :  like a map or an object contain a list of fields inside (value 4)
// ScalarNode : string, boolean,integer, float (value 8)
// AliasNode : using * to refer to an existing node, which we should avoid to use in swagger (value 16).
// style represent the literal writing format of field, style can be value as following
// DoubleQuotedStyle : using in template field usually, field in template quoted in double quote when template's using single quote
// SingleQuotedStyle : alike with "DoubleQuotedStyle"
// LiteralStyle : using in the string literal start with "|", it was recommended in using the string literal needs break line.
// FoldedStyle : using in the string literal start with ">". The difference from LiteralStyle is break lines in content texts
//
//	are ignored.It is also recommended in best practice to do non-wrapping text in this way
//
// FlowStyle : FlowStyle is another writing format for SequenceNode.eg,. example: [1, 2, 3] , while it was not recommended to using in best practice for the
//
//	difficulty of code reading when it's too long.

type yamlInt interface {
	convertSelf() *Field
}

type YAMLRoot struct {
	bytes    []byte
	children []*Field
	Node     *yaml.Node
}

func (r *YAMLRoot) Valid() {
	if len(r.Node.Content) == 1 && r.Node.Content[0].Kind == yaml.MappingNode {
		r.Node = r.Node.Content[0]
	}

	for i := 0; i < len(r.Node.Content)/2; i++ {
		keyNode := r.Node.Content[i*2]
		valueNode := r.Node.Content[i*2+1]
		var field yamlInt
		switch valueNode.Kind {
		case yaml.MappingNode:
			field = MappingField{Field{
				keyNode:   keyNode,
				valueNode: valueNode,
			}}
		case yaml.ScalarNode:
			field = ScalarField{Field{
				keyNode:   keyNode,
				valueNode: valueNode,
			}}
		case yaml.SequenceNode:
			field = SequenceField{Field{
				keyNode:   keyNode,
				valueNode: valueNode,
			}}
		}
		child := field.convertSelf()
		r.addContent(child)
	}

	//all fields set properly, start to find position
	for _, v := range r.children {
		v.findPosition()
	}
}

func (r *YAMLRoot) addContent(field *Field) {
	if r.children == nil {
		r.children = []*Field{}
	}
	r.children = append(r.children, field)
}

type FieldKind uint32

type Field struct {
	keyNode    *yaml.Node
	valueNode  *yaml.Node
	keyRange   []*FieldRange
	valueRange []*FieldRange
	key        string
	value      string
	kind       FieldKind
	style      yaml.Style
	content    []*Field
}

// make user function called after all field content filled
func (f *Field) findPosition() {
	f.findKeyPosition()
	f.findValuePosition()

	if f.content != nil && len(f.content) > 0 {
		for _, v := range f.content {
			v.findPosition()
		}
	}
}

func (f *Field) findKeyPosition() {
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

func (f *Field) findValuePosition() {
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

func (f *Field) setFieldName() {
	if f.keyNode != nil {
		f.key = f.keyNode.Value
	} else {
		f.key = ""
	}
}

func (f *Field) setFieldKind() {
	if f.valueNode == nil {
		return
	}
	var kind FieldKind
	switch f.valueNode.Kind {
	case yaml.MappingNode:
		kind = mappingField
	case yaml.ScalarNode:
		kind = scalarField
	case yaml.SequenceNode:
		kind = sequenceField
	default:
		kind = scalarField
	}
	f.kind = kind
}

func (f *Field) setFieldValue() {
	if f.valueNode == nil {
		return
	}
	f.value = f.valueNode.Value
}

func (f *Field) setFieldStyle() {
	if f.valueNode == nil {
		return
	}
	f.style = f.valueNode.Style
}

func (f *Field) addContent(content *Field) {
	if content == nil {
		return
	}
	if f.content == nil {
		f.content = []*Field{}
	}

	f.content = append(f.content, content)
}

type SequenceField struct {
	Field
}

func (f *Field) sameType() bool {
	//node kind in list are not equal
	for i := 0; i < len(f.valueNode.Content)-1; i++ {
		if f.valueNode.Content[i].Kind != f.valueNode.Content[i+1].Kind {
			return false
		}
	}

	return true
}

func (field SequenceField) findValuePosition() {
	fmt.Println("sequence field find value position called")
}

func (field SequenceField) convertSelf() *Field {
	field.setFieldName()
	field.setFieldKind()
	field.setFieldStyle()
	root := field.valueNode
	if field.sameType() {
		for i := 0; i < len(root.Content); i++ {
			var fieldInt yamlInt
			switch root.Content[i].Kind {
			case yaml.MappingNode:
				// content type of sequence is mapping
				fieldInt = MappingField{Field{
					keyNode:   nil,
					valueNode: root.Content[i],
				}}
				// content type of sequence is mapping
				// like list:
				//    	  - value1
				//        - value2
				//        - value3
				// or   list : ["value1","value2","value3"]

			case yaml.ScalarNode:
				fieldInt = ScalarField{Field{
					keyNode:   nil,
					valueNode: root.Content[i],
				}}
			case yaml.SequenceNode:
				fieldInt = SequenceField{Field{
					keyNode:   nil,
					valueNode: root.Content[i],
				}}
			}
			child := fieldInt.convertSelf()
			field.addContent(child)
		}
	}

	return &field.Field
}

type MappingField struct {
	Field
}

func (field MappingField) convertSelf() *Field {
	field.setFieldName()
	field.setFieldKind()
	field.setFieldStyle()
	root := field.valueNode
	for i := 0; i < len(root.Content)/2; i++ {
		//key 为奇数，value 为偶数
		keyNode := root.Content[i*2]
		valueNode := root.Content[i*2+1]
		var fieldInt yamlInt
		switch valueNode.Kind {
		case yaml.MappingNode:
			fieldInt = MappingField{Field{
				keyNode:   keyNode,
				valueNode: valueNode,
			}}
		case yaml.ScalarNode:
			fieldInt = ScalarField{Field{
				keyNode:   keyNode,
				valueNode: valueNode,
			}}
		case yaml.SequenceNode:
			fieldInt = SequenceField{Field{
				keyNode:   keyNode,
				valueNode: valueNode,
			}}
		}

		child := fieldInt.convertSelf()
		field.addContent(child)
	}
	return &field.Field
}

type ScalarField struct {
	Field
}

func (field ScalarField) convertSelf() *Field {
	field.setFieldName()
	field.setFieldKind()
	field.setFieldStyle()
	field.setFieldValue()
	return &field.Field
}
