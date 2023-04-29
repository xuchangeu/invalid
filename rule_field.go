package invalid

import (
	"errors"
	"fmt"
	"gopkg.in/yaml.v3"
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
	ConstraintKeyKReg       = `$key-reg`    //a regexp written in string to perform key validation.It can be used in scenario like checking extensible keys only prefix with ‘x’ in Swagger, key-reg exists under type $obj
	ConstraintKeyConstraint = `$constraint` //a type constraint for type $arr , valid for type $arr
)

var specKeyInObj = []string{ConstraintKeyType, ConstraintKeyRequired, ConstraintKeyOptional, ConstraintKeyKReg}

type Rule struct {
	required  bool //field's required
	keyNode   *yaml.Node
	valueNode *yaml.Node
	ruleType  RuleType //type field in validation file
	ruleMap   map[string]ruleInt
	ruleList  []ruleInt
}

func (rule *Rule) GetRuleMap() map[string]ruleInt {
	return rule.ruleMap
}

func (rule *Rule) GetType() RuleType {
	return rule.ruleType
}

func (rule *Rule) getContent() []*yaml.Node {
	if len(rule.valueNode.Content) > 0 {
		return rule.valueNode.Content
	} else {
		return []*yaml.Node{}
	}
}

func (rule *Rule) addRule(name string, r ruleInt) {
	if rule.ruleMap == nil {
		rule.ruleMap = map[string]ruleInt{}
	}
	rule.ruleMap[name] = r

	if rule.ruleList == nil {
		rule.ruleList = make([]ruleInt, 0)
	}
	rule.ruleList = append(rule.ruleList, r)
}

func (rule *Rule) addRules(ruleMap map[string]ruleInt) {
	for k, v := range ruleMap {
		rule.addRule(k, v)
	}
}

func (rule *Rule) restructure() error {
	panic("implement me")
}

type ObjFieldRule struct {
	Rule
	keyRegExp *regexp.Regexp
}

func (rule *ObjFieldRule) restructure() error {
	//get content except special key
	nodes := getContentExcept(rule.valueNode, specKeyInObj...)
	for i := 0; i < len(nodes)/2; i++ {
		k := nodes[i*2]
		v := nodes[i*2+1]
		r, err := resolveRule(k, v)
		if err != nil {
			return err
		}
		err = r.restructure()
		if err != nil {
			return err
		}
		rule.addRule(k.Value, r)
	}

	//handle key reg
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

	//handle required
	k, v, e = GetKVNodeByKeyName(ConstraintKeyRequired, rule.getContent())
	if k != nil && v != nil && e {
		if !validBoolNode(v) {
			return errors.New(fmt.Sprintf("value node must be boolean : [%s]", k.Value))
		} else if v.Value != "true" {
			return errors.New(fmt.Sprintf("value for required must be true"))
		}
		rule.required = true
	}

	return nil
}

// ArrFieldRule represent a rule of arr
type ArrFieldRule struct {
	Rule
	constraint    string
	constraintMap *ObjFieldRule
}

func (rule *ArrFieldRule) restructure() error {
	//check constraint
	k, v, exist := GetKVNodeByKeyName(ConstraintKeyConstraint, rule.getContent())
	if k != nil && v != nil && exist {
		if validMapNode(v) {

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

// StrFieldRule represent a rule field of string
type StrFieldRule struct {
	Rule
	max uint           //max length of field
	min uint           //min length of field
	reg *regexp.Regexp //regexp of field
}

func (rule *StrFieldRule) restructure() error {
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

	//check key reg
	k, v, e = GetKVNodeByKeyName(ConstraintKeyReg, rule.getContent())
	if k != nil && v != nil && e && validStrNode(v) {
		reg, err := regexp.Compile(v.Value)
		if err != nil {
			return errors.New(fmt.Sprintf("compile regexp error : [%s]", k.Value))
		}
		rule.reg = reg
	}
	return nil
}

// BoolFieldRule represent a rule of boolean
type BoolFieldRule struct {
	Rule
}

func (rule *BoolFieldRule) restructure() error {
	//I have nothing to do
	return nil
}

// FloatFieldRule represent a rule a float
type FloatFieldRule struct {
	Rule
}

func (rule *FloatFieldRule) restructure() error {
	//I have nothing to do
	return nil
}

// IntFieldRule represent a rule of int
type IntFieldRule struct {
	Rule
}

func (rule *IntFieldRule) restructure() error {
	rule.ruleType = RuleTypeInt
	return nil
}

// NilFieldRule represent a rule of nil
type NilFieldRule struct {
	Rule
}

func (rule *NilFieldRule) restructure() error {
	return nil
}
