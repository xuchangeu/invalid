package invalid

import (
	"context"
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
	RuleTypeNil   RuleType = "$null"
	RuleTypeAny   RuleType = "$any"
	RuleTypeBool  RuleType = "$bool"
	RuleTypeInt   RuleType = "$int"
	RuleTypeFloat RuleType = "$float"
	RuleTypeStr   RuleType = "$str"

	//collection types
	RuleTypeObj RuleType = "$obj" //an object value contains sub-ruleMap inside, mostly it's a map
	RuleTypeSeq RuleType = "$seq" //a list with value in any type
	RuleTypeArr RuleType = "$arr"

	YAMLTypeNull  = "!!null"
	YAMLTypeBool  = "!!bool"
	YAMLTypeInt   = "!!int"
	YAMLTypeFloat = "!!float"
	YAMLTypeStr   = "!!str"
	YAMLTypeMap   = "!!map"
	YAMLTypeSeq   = "!!seq"
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
	ConstraintKeyKReg       = `$key-reg`    //a regexp written in string to perform key-name validation.It can be used in scenario like checking extensible keys only prefix with ‘x’ in Swagger, key-regexp exists under type $obj
	ConstraintKeyConstraint = `$constraint` //a type constraint for type $arr , valid for type $arr
)

var specKeyInObj = []string{ConstraintKeyType, ConstraintKeyRequired, ConstraintKeyOptional, ConstraintKeyKReg}

type Ruler interface {
	Get(key string) (Ruler, bool)
	GetRules() []Ruler
	RuleType() RuleType
	Key() string
	Required() bool
	Validate(f Field) []*Result
	restructure() error
}

func NewRule(r io.Reader) (Ruler, error) {
	byte, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	node := &yaml.Node{}
	err = yaml.Unmarshal(byte, node)
	if err != nil {
		return nil, err
	}

	if len(node.Content) < 1 {
		return nil, errors.New("document must have at least one field")
	}
	node = node.Content[0]

	ruler, err := newRuler(nil, node, true)
	if err != nil {
		return nil, err
	}
	err = ruler.restructure()
	if err != nil {
		return nil, err
	}

	return ruler, nil
}

type Rule struct {
	required  bool //field's required
	keyNode   *yaml.Node
	valueNode *yaml.Node
	ruleType  RuleType //type field in validation file
	ruleMap   map[string]Ruler
	ruleList  []Ruler
}

func (rule *Rule) Validate(f Field) []*Result {

	ctx, cancel := context.WithCancel(context.Background())
	result := doValidate(ctx, cancel, rule, f, nil)
	if *result == nil {
		x := make([]*Result, 0)
		return x
	}
	return *result

}

func doValidate(ctx context.Context, cancel context.CancelFunc, rule Ruler, field Field, result *[]*Result) *[]*Result {
	if result == nil {
		result = new([]*Result)
	}

	if ctx.Err() == context.Canceled {
		return result
	}

	for i := 0; i < len(rule.GetRules()); i++ {
		if ctx.Err() == context.Canceled {
			return result
		}

		r := rule.GetRules()[i]
		f, e := field.Get(r.Key())
		//check key required missing
		if !e && r.Required() {
			err := NewResult(KeyMissing, NewKeyMissingError(r.Key()), field.Fields()[i].getValueRange())
			x := *result
			v := append(x, &err)
			cancel()
			return &v
		}
		switch v := r.(type) {
		case *ObjRule:
			result = doValidate(ctx, cancel, r, f, result)
		case *ArrRule:
			switch v.constraint.(type) {
			//scalar constraint
			case string:
				for i := 0; i < len(f.Fields()); i++ {
					if string(f.Fields()[i].ValueType()) != v.constraint {
						e := NewResult(TypeMismatch, NewTypeMismatchError(fmt.Sprintf("%s.%s", f.Key(),
							f.Fields()[i].Key()), v.constraint.(string)), f.getValueRange())
						x := *result
						y := append(x, &e)
						result = &y
					}
				}
			//for ruler object
			case Ruler:
				for i := 0; i < len(f.Fields()); i++ {
					if ctx.Err() == context.Canceled {
						return result
					}
					result = doValidate(ctx, cancel, v.constraint.(Ruler), f.Fields()[i], result)
				}
			}

		case *StrRule:
			if f.Tag() != YAMLTypeStr {
				e := NewResult(TypeMismatch, NewTypeMismatchError(f.Key(), string(RuleTypeStr)), f.getValueRange())
				x := *result
				y := append(x, &e)
				result = &y
			}

			//check min or max
			if v.max != 0 || v.min != 0 {
				if v.min != 0 && len(f.Value()) < int(v.min) {
					e := NewResult(StrLengthMismatch, NewStrLengthError1(r.Key(), int(v.min)), f.getValueRange())
					x := *result
					y := append(x, &e)
					result = &y
				} else if v.max != 0 && len(f.Value()) > int(v.max) {
					e := NewResult(StrLengthMismatch, NewStrLengthError2(r.Key(), int(v.max)), f.getValueRange())
					x := *result
					y := append(x, &e)
					result = &y
				}
			}

			//check regexp
			if v.GetReg() != nil {
				m := v.regexp.Match([]byte(f.Value()))
				if !m {
					e := NewResult(RegxMismatch, NewRegxError(r.Key(), v.GetReg().String()), f.getValueRange())
					x := *result
					y := append(x, &e)
					result = &y
				}
			}

		case *IntRule:
			if f.Tag() != YAMLTypeInt {
				e := NewResult(TypeMismatch, NewTypeMismatchError(f.Key(), string(RuleTypeInt)), f.getValueRange())
				x := *result
				y := append(x, &e)
				result = &y
			}

		case *FloatRule:
			if f.Tag() != YAMLTypeFloat {
				e := NewResult(TypeMismatch, NewTypeMismatchError(f.Key(), string(RuleTypeFloat)), f.getValueRange())
				x := *result
				y := append(x, &e)
				result = &y
			}

		case *BoolRule:
			if f.Tag() != YAMLTypeBool {
				e := NewResult(TypeMismatch, NewTypeMismatchError(f.Key(), string(RuleTypeBool)), f.getValueRange())
				x := *result
				y := append(x, &e)
				result = &y
			}

		case *NullFieldRule:
			if f.Tag() != YAMLTypeNull {
				e := NewResult(TypeMismatch, NewTypeMismatchError(f.Key(), string(RuleTypeNil)), f.getValueRange())
				x := *result
				y := append(x, &e)
				result = &y
			}
		}
	}

	return result
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

func (rule *Rule) Key() string {
	if rule.keyNode == nil {
		return ""
	}
	return rule.keyNode.Value
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

func (rule *Rule) GetRules() []Ruler {
	return rule.getRules()
}

func (rule *Rule) getRules() []Ruler {
	return rule.ruleList
}

func (rule *Rule) restructure() error {
	//panic("implement me")
	//handle required
	key, value, exist := GetKVNodeByKeyName(ConstraintKeyOptional, rule.getContent())
	if key != nil && value != nil && exist {
		if !validBoolNode(value) {
			return errors.New(fmt.Sprintf("value node must be boolean : [%s]", key.Value))
		} else if value.Value != "true" {
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
		r, e := newRuler(k, v, false)
		if e != nil {
			return e
		}
		e = r.restructure()
		if e != nil {
			return e
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

type Constraint interface{}

// ArrRule represent a rule of arr
type ArrRule struct {
	Rule
	constraint Constraint
}

func (rule *ArrRule) GetConstraint() interface{} {
	return rule.constraint
}

func (rule *ArrRule) restructure() error {
	err := rule.Rule.restructure()
	if err != nil {
		return err
	}

	//check constraint
	key, value, exist := GetKVNodeByKeyName(ConstraintKeyConstraint, rule.getContent())
	if key != nil && value != nil && exist {
		//constraint is node
		if validMapNode(value) {
			ruler, err := newRuler(key, value, true)
			if err != nil {
				return err
			}

			err = ruler.restructure()
			if err != nil {
				return err
			}
			rule.constraint = ruler

		} else if validStrNode(value) {
			if !contains(scalarTypes, value.Value) {
				return errors.New(fmt.Sprintf("constraint should be one of %v", value.Value))
			} else {
				rule.constraint = value.Value
			}
		} else {
			return errors.New("constraint format should be a string value of scalar type or obj")
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
	key, value, exist := GetKVNodeByKeyName(ConstraintKeyLength, rule.getContent())
	if key != nil && value != nil && exist {
		min, err := GetIntValue(ConstraintKeyMin, value.Content)
		if err != nil {
			return err
		}
		rule.min = uint(min)

		//check max
		max, err := GetIntValue(ConstraintKeyMax, value.Content)
		if err != nil {
			return err
		}
		rule.max = uint(max)
	}

	//check key regexp
	key, value, exist = GetKVNodeByKeyName(ConstraintKeyReg, rule.getContent())
	if key != nil && value != nil && exist && validStrNode(value) {
		reg, err := regexp.Compile(value.Value)
		if err != nil {
			return errors.New(fmt.Sprintf("compile regexp error : [%s]", key.Value))
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
	return rule.Rule.restructure()
}

// FloatRule represent a rule a float
type FloatRule struct {
	Rule
}

func (rule *FloatRule) restructure() error {
	return rule.Rule.restructure()
}

// IntRule represent a rule of int
type IntRule struct {
	Rule
}

func (rule *IntRule) restructure() error {
	return rule.Rule.restructure()
}

// NullFieldRule represent a rule of nil
type NullFieldRule struct {
	Rule
}

func (rule *NullFieldRule) restructure() error {
	return rule.Rule.restructure()
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
		return &NullFieldRule{Rule{
			ruleType:  RuleTypeNil,
			keyNode:   keyNode,
			valueNode: valueNode,
		}}, nil

	}
	return nil, errors.New(fmt.Sprintf("type not match : [%s]", keyNode.Value))
}
