package invalid

import (
	"errors"
	"gopkg.in/yaml.v3"
)

type ruleInt interface {
	mustConstraint() []string
	optionalConstraint() []string
}

type Rule struct {
	Node       *yaml.Node
	fieldRules []*FieldRule
}

func (r *Rule) Decode() error {
	if r.Node == nil {
		return errors.New("none yaml nodes available")
	}

	return nil
}

type FieldRule struct {
}

func (rule *FieldRange) mustConstraint() []string {
	//implement me
	return nil
}

func (rule *FieldRange) optionalConstraint() []string {
	//implement me
	return nil
}

type ObjFieldRule struct {
	FieldRule
}

func (rule *ObjFieldRule) mustConstraint() []string {
	return []string{ConstraintKeyType}
}

func (rule *ObjFieldRule) optionalConstraint() []string {
	return []string{ConstraintKeyKReg}
}

type StrFieldRule struct {
	FieldRule
}

func (rule *StrFieldRule) mustConstraint() []string {
	return []string{ConstraintKeyType}
}

func (rule *StrFieldRule) optionalConstraint() []string {
	return []string{
		deepFieldWithDot([]string{ConstraintKeyLength, ConstraintKeyMin}),
		deepFieldWithDot([]string{ConstraintKeyLength, ConstraintKeyMax}),
		ConstraintKeyReg,
	}
}

type BoolFieldRule struct {
	FieldRule
}

func (rule *BoolFieldRule) mustConstraint() []string {
	return []string{ConstraintKeyType}
}

func (rule *BoolFieldRule) optionalConstraint() []string {
	return nil
}

type ArrFieldRule struct {
	FieldRule
}

func (rule *ArrFieldRule) mustConstraint() []string {
	return []string{ConstraintKeyType}
}

func (rule *ArrFieldRule) optionalConstraint() []string {
	return []string{ConstraintKeyConstraint}
}

type NumFieldRule struct {
	FieldRule
}

func (rule *NumFieldRule) mustConstraint() []string {
	return []string{ConstraintKeyType}
}

func (rule *NumFieldRule) optionalConstraint() []string {
	return nil
}

type IntFieldRule struct {
	FieldRule
}

func (rule *IntFieldRule) mustConstraint() []string {
	return []string{ConstraintKeyType}
}

func (rule *IntFieldRule) optionalConstraint() []string {
	return nil
}

type NilFieldRule struct {
	FieldRule
}

func (rule *NilFieldRule) mustConstraint() []string {
	return []string{ConstraintKeyType}
}

func (rule *NilFieldRule) optionalConstraint() []string {
	return nil
}

type RuleFieldType string

// basic type of schema
const (
	//single types
	RuleTypeNil     = "$nil"
	RuleTypeAny     = "$any"
	RuleTypeBoolean = "$bool"
	RuleTypeInt     = "$int"
	RuleTypeFloat   = "$float"
	RuleTypeString  = "$str"

	//collection types
	TypeObj = "$obj" //an object value contains sub-fields inside, mostly it's a map
	TypeSeq = "$seq" //a list with value in any type
	TypeArr = "$arr"
)

const (
	//constraint keys
	ConstraintKeyType       = "type"       //type definition
	ConstraintKeyRequired   = "required"   //fields must exist,  exists under type $obj
	ConstraintKeyOptional   = "optional"   //fields which are optional,  exists under type $obj alike required
	ConstraintKeyLength     = "length"     //length of character, valid under type $obj
	ConstraintKeyReg        = "reg"        //regexp pattern written in string, valid in type $str
	ConstraintKeyMin        = "min"        //minimum length of string, valid in type $str
	ConstraintKeyMax        = "max"        //maximum length of string, valid in type $str
	ConstraintKeyKReg       = "key-reg"    //a regexp written in string to perform key key validation.It can be used in scenario like checking extensible keys only prefix with ‘x’ in Swagger, key-reg exists under type $obj
	ConstraintKeyConstraint = "constraint" //a type constraint for type $arr , valid for type $arr

)
