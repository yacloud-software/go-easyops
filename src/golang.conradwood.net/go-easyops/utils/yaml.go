package utils

import (
	"bytes"
	"gopkg.in/yaml.v3"
)

// reads a raml file and parses it (strict) into interface. error if unknown tags are encountered in yaml
func ReadYaml(filename string, target interface{}) error {
	b, err := ReadFile(filename)
	if err != nil {
		return err
	}

	decoder := yaml.NewDecoder(bytes.NewReader(b))
	decoder.KnownFields(true)
	err = decoder.Decode(target)
	//	err = yaml.Unmarshal(b, target)
	if err != nil {
		return err
	}
	return nil
}

// interprets the bytes as yaml and decodes it (strict) into interface. error if unknown tags are encountered in yaml
func UnmarshalYaml(buf []byte, target interface{}) error {
	decoder := yaml.NewDecoder(bytes.NewReader(buf))
	decoder.KnownFields(true)
	err := decoder.Decode(target)
	//	err = yaml.Unmarshal(b, target)
	if err != nil {
		return err
	}
	return nil
}

func MarshalYaml(src interface{}) ([]byte, error) {
	return yaml.Marshal(src)
}
