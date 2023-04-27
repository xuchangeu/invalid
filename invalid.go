package invalid

import (
	"gopkg.in/yaml.v3"
	"io"
)

func NewYAML(r io.Reader) (*YAMLRoot, error) {
	by, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	node := &yaml.Node{}
	err = yaml.Unmarshal(by, node)
	if err != nil {
		return nil, err
	}

	yamlRoot := YAMLRoot{
		bytes: by,
		Node:  node,
	}

	return &yamlRoot, nil
}
