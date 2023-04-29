package invalid

import (
	"errors"
	"fmt"
	"gopkg.in/yaml.v3"
)

type ruleInt interface {
	restructure() error
	GetType() RuleType
}

//func NewRule(r io.Reader) (*Rule, error) {
//	by, err := io.ReadAll(r)
//	if err != nil {
//		return nil, err
//	}
//
//	node := &yaml.Node{}
//	err = yaml.Unmarshal(by, node)
//	if err != nil {
//		return nil, err
//	}
//
//	rule := Rul
//
//	return &rule, nil
//}

//type RuleRoot struct {
//	by         []byte
//	node       *yaml.Node
//	ruleIntMap map[string]ruleInt
//	ruleInList []ruleInt
//}

//func (r *RuleRoot) addRule(key string, rule ruleInt) {
//	if r.ruleIntMap == nil {
//		r.ruleIntMap = make(map[string]ruleInt, 0)
//	}
//	r.ruleIntMap[key] = rule
//
//	if r.ruleInList == nil {
//		r.ruleInList = make([]ruleInt, 0)
//	}
//	r.ruleInList = append(r.ruleInList, rule)
//}
//
//func (r *RuleRoot) addRules(ruleMap map[string]ruleInt) {
//	for k, v := range ruleMap {
//		r.addRule(k, v)
//	}
//}
//
//func (r *RuleRoot) ValidYAML(root *YAMLRoot) {
//
//}
//
//func (r *RuleRoot) GetFieldRules() map[string]ruleInt {
//	return r.ruleIntMap
//}
//
//func (r *RuleRoot) Restructure() error {
//	if r.node == nil {
//		return errors.New("none yaml nodes available")
//	}
//
//	if len(r.node.Content) == 1 && r.node.Kind == yaml.DocumentNode {
//		r.node = r.node.Content[0]
//	}
//
//	for i := 0; i < len(r.node.Content)/2; i++ {
//		k := r.node.Content[i*2]
//		v := r.node.Content[i*2+1]
//		rule, err := resolveRule(k, v)
//		if err != nil {
//			return err
//		}
//		err = rule.restructure()
//		if err != nil {
//			return err
//		}
//		r.addRule(k.Value, rule)
//	}
//
//	return nil
//
//}

func resolveRule(keyNode, valueNode *yaml.Node) (ruleInt, error) {

	if !validMapNode(valueNode) {
		return nil, errors.New(fmt.Sprintf("value node must be map : [%s]", keyNode.Value))
	}

	k, v, e := GetKVNodeByKeyName(ConstraintKeyType, valueNode.Content)
	if !(k != nil && v != nil && e) {
		return nil, errors.New(fmt.Sprintf("type not found : [%s]", keyNode.Value))
	}

	switch RuleType(v.Value) {
	case RuleTypeArr:
		return &ArrFieldRule{
			Rule: Rule{
				ruleType:  RuleTypeArr,
				keyNode:   keyNode,
				valueNode: valueNode,
			}}, nil
	case RuleTypeSeq, RuleTypeAny:
		//TODO : tbc
	case RuleTypeObj:
		return &ObjFieldRule{
			Rule: Rule{
				ruleType:  RuleTypeObj,
				keyNode:   keyNode,
				valueNode: valueNode,
			}}, nil
	case RuleTypeInt:
		return &IntFieldRule{Rule{
			ruleType:  RuleTypeInt,
			keyNode:   keyNode,
			valueNode: valueNode,
		}}, nil
	case RuleTypeStr:
		return &StrFieldRule{
			Rule: Rule{
				ruleType:  RuleTypeStr,
				keyNode:   keyNode,
				valueNode: valueNode,
			}}, nil
	case RuleTypeBool:
		return &BoolFieldRule{Rule{
			ruleType:  RuleTypeBool,
			keyNode:   keyNode,
			valueNode: valueNode,
		}}, nil
	case RuleTypeFloat:
		return &FloatFieldRule{Rule{
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
