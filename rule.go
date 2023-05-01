package invalid

import (
	"errors"
	"fmt"
	"gopkg.in/yaml.v3"
	"io"
	"regexp"
)

type RuleType string

// basic type of schema
const (
	//single types
	RuleTypeNil   RuleType = "$nil"
	RuleTypeAny   RuleType = "$any"
	RuleTypeBool  RuleType = "$bool"
	RuleTypeInt   RuleType = "$int"
	RuleTypeFloat RuleType = "$float"
	RuleTypeStr   RuleType = "$str"

	//collection types
	RuleTypeObj RuleType = "$obj" //an object value contains sub-ruleMap inside, mostly it's a map
	RuleTypeSeq RuleType = "$seq" //a list with value in any type
	RuleTypeArr RuleType = "$arr"
)

// yaml scalar nodes, include bool, integer, float, string and null, but null was not included here.
var scalarTypes = []string{string(RuleTypeBool), string(RuleTypeInt),
	string(RuleTypeFloat), string(RuleTypeStr)}

const (
	//constraint keys
	ConstraintKeyType       = `$type`       //type definition
	ConstraintKeyRequired   = `$required`   //ruleMap must exist,  exists under type $obj
	ConstraintKeyOptional   = `$optional`   //ruleMap which are optional,  exists under type $obj alike required
	ConstraintKeyLength     = `$length`     //length of character, valid under type $obj
	ConstraintKeyReg        = `$reg`        //regexp pattern written in string, valid in type $str
	ConstraintKeyMin        = `$min`        //minimum length of string, valid in type $str
	ConstraintKeyMax        = `$max`        //maximum length of string, valid in type $str
	ConstraintKeyKReg       = `$key-reg`    //a regexp written in string to perform key validation.It can be used in scenario like checking extensible keys only prefix with ‘x’ in Swagger, key-regexp exists under type $obj
	ConstraintKeyConstraint = `$constraint` //a type constraint for type $arr , valid for type $arr
)

var specKeyInObj = []string{ConstraintKeyType, ConstraintKeyRequired, ConstraintKeyOptional, ConstraintKeyKReg}

type Ruler interface {
	restructure() error
	RuleType() RuleType
	Get(key string) (Ruler, bool)
	Required() bool
}

func NewRule(r io.Reader) (Ruler, error) {
	by, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	node := &yaml.Node{}
	err = yaml.Unmarshal(by, node)
	if err != nil {
		return nil, err
	}

	if len(node.Content) < 1 {
		return nil, errors.New("document must have at least one field")
	}
	docNode := node.Content[0]

	f, e := newRuler(nil, docNode, true)
	if e != nil {
		return nil, e
	}
	e = f.restructure()
	if e != nil {
		return nil, e
	}

	return f, nil
}

type Rule struct {
	required  bool //field's required
	keyNode   *yaml.Node
	valueNode *yaml.Node
	ruleType  RuleType //type field in validation file
	ruleMap   map[string]Ruler
	ruleList  []Ruler
}

func (rule *Rule) GetRuleMap() map[string]Ruler {
	return rule.ruleMap
}

func (rule *Rule) RuleType() RuleType {
	return rule.ruleType
}

func (rule *Rule) Required() bool {
	return rule.required
}

func (rule *Rule) Get(key string) (Ruler, bool) {
	if rule.ruleMap == nil {
		return nil, false
	}
	r, e := rule.ruleMap[key]
	return r, e
}

func (rule *Rule) getContent() []*yaml.Node {
	if len(rule.valueNode.Content) > 0 {
		return rule.valueNode.Content
	} else {
		return []*yaml.Node{}
	}
}

func (rule *Rule) addRule(name string, r Ruler) {
	if rule.ruleMap == nil {
		rule.ruleMap = map[string]Ruler{}
	}
	rule.ruleMap[name] = r

	if rule.ruleList == nil {
		rule.ruleList = make([]Ruler, 0)
	}
	rule.ruleList = append(rule.ruleList, r)
}

func (rule *Rule) addRules(ruleMap map[string]Ruler) {
	for k, v := range ruleMap {
		rule.addRule(k, v)
	}
}

func (rule *Rule) restructure() error {
	//panic("implement me")
	//handle required
	k, v, e := GetKVNodeByKeyName(ConstraintKeyOptional, rule.getContent())
	if k != nil && v != nil && e {
		if !validBoolNode(v) {
			return errors.New(fmt.Sprintf("value node must be boolean : [%s]", k.Value))
		} else if v.Value != "true" {
			return errors.New(fmt.Sprintf("value for required must be true"))
		}
		rule.required = false
	} else {
		rule.required = true
	}

	return nil
}

type ObjRule struct {
	Rule
	keyRegExp *regexp.Regexp
}

func (rule *ObjRule) GetKeyReg() *regexp.Regexp {
	return rule.keyRegExp
}

func (rule *ObjRule) restructure() error {
	err := rule.Rule.restructure()
	if err != nil {
		return err
	}

	//get content except special key
	nodes := getContentExcept(rule.valueNode, specKeyInObj...)
	for i := 0; i < len(nodes)/2; i++ {
		k := nodes[i*2]
		v := nodes[i*2+1]
		r, err := newRuler(k, v, false)
		if err != nil {
			return err
		}
		err = r.restructure()
		if err != nil {
			return err
		}
		rule.addRule(k.Value, r)
	}

	//handle key regexp
	k, v, e := GetKVNodeByKeyName(ConstraintKeyKReg, rule.getContent())
	if k != nil && v != nil && e {
		if !validStrNode(v) {
			return errors.New(fmt.Sprintf("value node must be string : [%s]", k.Value))
		}
		reg, err := regexp.Compile(v.Value)
		if err != nil {
			return errors.New(fmt.Sprintf("regexp compile error : [%s]", k.Value))
		}
		rule.keyRegExp = reg
	}

	return nil
}

type Constraint Ruler

// ArrRule represent a rule of arr
type ArrRule struct {
	Rule
	constraint    string
	constraintObj Constraint
}

func (rule *ArrRule) GetConstraint() string {
	return rule.constraint
}

func (rule *ArrRule) GetConstraintObj() *ObjRule {
	return rule.GetConstraintObj()
}

func (rule *ArrRule) restructure() error {
	err := rule.Rule.restructure()
	if err != nil {
		return err
	}

	//check constraint
	k, v, exist := GetKVNodeByKeyName(ConstraintKeyConstraint, rule.getContent())
	if k != nil && v != nil && exist {
		if validMapNode(v) {
			ruleInt, err := newRuler(k, v, true)
			if err != nil {
				return err
			}

			err = ruleInt.restructure()
			if err != nil {
				return err
			}
			rule.constraintObj = ruleInt

		} else if validStrNode(v) {
			if !contains(scalarTypes, v.Value) {
				return errors.New(fmt.Sprintf("constraint should be one of %v", v.Value))
			} else {
				rule.constraint = v.Value
			}
		}
	}
	return nil
}

// StrRule represent a rule field of string
type StrRule struct {
	Rule
	max    uint           //max length of field
	min    uint           //min length of field
	regexp *regexp.Regexp //regexp of field
}

func (rule *StrRule) GetMax() uint {
	return rule.min
}

func (rule *StrRule) GetMin() uint {
	return rule.min
}

func (rule *StrRule) GetReg() *regexp.Regexp {
	return rule.regexp
}

func (rule *StrRule) restructure() error {
	err := rule.Rule.restructure()
	if err != nil {
		return err
	}

	//check min & max
	k, v, e := GetKVNodeByKeyName(ConstraintKeyLength, rule.getContent())
	if k != nil && v != nil && e {
		min, err := GetIntValue(ConstraintKeyMin, v.Content)
		if err != nil {
			return err
		}
		rule.min = uint(min)

		//check max
		max, err := GetIntValue(ConstraintKeyMax, v.Content)
		if err != nil {
			return err
		}
		rule.max = uint(max)

	}

	//check key regexp
	k, v, e = GetKVNodeByKeyName(ConstraintKeyReg, rule.getContent())
	if k != nil && v != nil && e && validStrNode(v) {
		reg, err := regexp.Compile(v.Value)
		if err != nil {
			return errors.New(fmt.Sprintf("compile regexp error : [%s]", k.Value))
		}
		rule.regexp = reg
	}
	return nil
}

// BoolRule represent a rule of boolean
type BoolRule struct {
	Rule
}

func (rule *BoolRule) restructure() error {
	err := rule.Rule.restructure()
	if err != nil {
		return err
	}

	return nil
}

// FloatRule represent a rule a float
type FloatRule struct {
	Rule
}

func (rule *FloatRule) restructure() error {
	err := rule.Rule.restructure()
	if err != nil {
		return err
	}

	return nil
}

// IntRule represent a rule of int
type IntRule struct {
	Rule
}

func (rule *IntRule) restructure() error {
	err := rule.Rule.restructure()
	if err != nil {
		return err
	}

	return nil
}

// NilFieldRule represent a rule of nil
type NilFieldRule struct {
	Rule
}

func (rule *NilFieldRule) restructure() error {
	err := rule.Rule.restructure()
	if err != nil {
		return err
	}

	return nil
}

func newRuler(keyNode, valueNode *yaml.Node, document bool) (Ruler, error) {
	if !validMapNode(valueNode) {
		return nil, errors.New(fmt.Sprintf("value node must be map : [%s]", keyNode.Value))
	}

	if document {
		return &ObjRule{
			Rule: Rule{
				ruleType:  RuleTypeObj,
				keyNode:   keyNode,
				valueNode: valueNode,
			}}, nil
	}

	k, v, e := GetKVNodeByKeyName(ConstraintKeyType, valueNode.Content)
	if !(k != nil && v != nil && e) {
		return nil, errors.New(fmt.Sprintf("type not found : [%s]", keyNode.Value))
	}

	switch RuleType(v.Value) {
	case RuleTypeArr:
		return &ArrRule{
			Rule: Rule{
				ruleType:  RuleTypeArr,
				keyNode:   keyNode,
				valueNode: valueNode,
			}}, nil
	case RuleTypeSeq, RuleTypeAny:
		//TODO : tbc
	case RuleTypeObj:
		return &ObjRule{
			Rule: Rule{
				ruleType:  RuleTypeObj,
				keyNode:   keyNode,
				valueNode: valueNode,
			}}, nil
	case RuleTypeInt:
		return &IntRule{Rule{
			ruleType:  RuleTypeInt,
			keyNode:   keyNode,
			valueNode: valueNode,
		}}, nil
	case RuleTypeStr:
		return &StrRule{
			Rule: Rule{
				ruleType:  RuleTypeStr,
				keyNode:   keyNode,
				valueNode: valueNode,
			}}, nil
	case RuleTypeBool:
		return &BoolRule{Rule{
			ruleType:  RuleTypeBool,
			keyNode:   keyNode,
			valueNode: valueNode,
		}}, nil
	case RuleTypeFloat:
		return &FloatRule{Rule{
			ruleType:  RuleTypeFloat,
			keyNode:   keyNode,
			valueNode: valueNode,
		}}, nil
	case RuleTypeNil:
		return &NilFieldRule{Rule{
			ruleType:  RuleTypeNil,
			keyNode:   keyNode,
			valueNode: valueNode,
		}}, nil

	}
	return nil, errors.New(fmt.Sprintf("type not match : [%s]", keyNode.Value))
}
