package invalid

import (
	"errors"
	"fmt"
	"gopkg.in/yaml.v3"
	"strconv"
)

//func deepFieldWithDot(keys []string) string {
//	if len(keys) == 0 {
//		return ""
//	}
//	return strings.Join(keys, ConstraintKeyDepthIndicator)
//}

// get the content under the specific *yaml.Node by key only if it's Tag's !!map
// function return en empty slice if the tag is scalar node
// the value node should be a map
//func getChildValueNodes(key string, valueNode *yaml.Node) []*yaml.Node {
//	if valueNode.Tag != "!!map" || len(valueNode.Content) == 0 {
//		return []*yaml.Node{}
//	}
//	//find target
//	var targetNode *yaml.Node
//	content := valueNode.Content
//	for i := 0; i < len(content)/2; i++ {
//		k := content[i*2]
//		v := content[i*2+1]
//		if k.Value == key {
//			targetNode = v
//			break
//		}
//	}
//
//	//target exist
//	if targetNode != nil && targetNode.Tag == "!!map" && len(targetNode.Content) > 0 {
//		return targetNode.Content
//	}
//
//	return []*yaml.Node{}
//}

// getNodesExcept return yaml node's content except the key mentioned, also the value match to key
func getContentExcept(node *yaml.Node, keys ...string) []*yaml.Node {
	if !validMapNode(node) {
		return []*yaml.Node{}
	}

	result := make([]*yaml.Node, 0)
	for i := 0; i < len(node.Content)/2; i++ {
		k := node.Content[i*2]
		v := node.Content[i*2+1]
		if !contains(keys, k.Value) {
			result = append(result, k, v)
		}
	}

	return result
}

// GetFloatValue get float value of content by key name
// return error when tag mismatch
func GetFloatValue(key string, nodes []*yaml.Node) (float64, error) {
	k, v, exist := GetKVNodeByKeyName(key, nodes)
	if k != nil && v != nil && exist {
		if !validFloatNode(v) {
			return 0, errors.New(fmt.Sprintf("node tag is not float for key : [%s]", key))
		}

		val, err := strconv.ParseFloat(v.Value, 64)
		if err != nil {
			return 0, err
		}

		return val, nil
	}

	//return 0, errors.New(fmt.Sprintf("value not found for key : [%s]", key))
	return 0, nil
}

// GetIntValue get int value of content by key name
// return error when tag mismatch
func GetIntValue(key string, nodes []*yaml.Node) (int, error) {
	k, v, exist := GetKVNodeByKeyName(key, nodes)
	if k != nil && v != nil && exist {
		if !validIntNode(v) {
			return 0, errors.New(fmt.Sprintf("node tag is not int for key : [%s]", key))
		}

		val, err := strconv.Atoi(v.Value)
		if err != nil {
			return 0, err
		}

		return val, nil
	}

	//return 0, errors.New(fmt.Sprintf("value not found for key : [%s]", key))
	return 0, nil
}

// GetStringValue get string value of content by key name
// return error when tag mismatch
//func GetStringValue(key string, nodes []*yaml.Node) (string, error) {
//	k, v, exist := GetKVNodeByKeyName(key, nodes)
//	if k != nil && v != nil && exist {
//		if !validStrNode(v) {
//			return "", errors.New(fmt.Sprintf("node tag is not int for key : [%s]", key))
//		}
//		return v.Value, nil
//	}
//	return "", errors.New(fmt.Sprintf("value not found for key : [%s]", key))
//}

// GetKVNodeByKeyName function return keyNode,valueNode,exist by key name
func GetKVNodeByKeyName(key string, nodes []*yaml.Node) (*yaml.Node, *yaml.Node, bool) {
	for k, v := range nodes {
		if v.Kind == yaml.ScalarNode && v.Value == key && len(nodes) > k {
			return nodes[k], nodes[k+1], true
		}
	}
	return nil, nil, false
}

// weather tag of node is !!str
func validStrNode(node *yaml.Node) bool {
	return node.Tag == yamlNodeTypeStr
}

// weather tag of node is !!bool
func validBoolNode(node *yaml.Node) bool {
	return node.Tag == yamlNodeTypeBool
}

// weather tag of node is !!float
func validFloatNode(node *yaml.Node) bool {
	return node.Tag == yamlNodeTypeFloat
}

// weather tag of node is !!arr
func validArrNode(node *yaml.Node) bool {
	return node.Tag == yamlNodeTypeSeq
}

// weather tag of node is !!int
func validIntNode(node *yaml.Node) bool {
	return node.Tag == yamlNodeTypeInt
}

// weather tag of node is !!null
func validNullNode(node *yaml.Node) bool {
	return node.Tag == yamlNodeTypeNull
}

// weather tag of node is !!map
func validMapNode(node *yaml.Node) bool {
	return node.Tag == yamlNodeTypeMap
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

//const (
//	ConstraintKeyDepthIndicator = "."
//)
