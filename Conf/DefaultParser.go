package Conf

import (
	"gopkg.in/yaml.v2"
)

type yamlParser struct {
}

func NewYamlParser() *yamlParser {
	return &yamlParser{}
}

func (yp *yamlParser) Unmarshal(data []byte, v interface{}) error {
	return yaml.Unmarshal(data, v)
}
