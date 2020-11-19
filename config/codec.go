package config

import (
	"encoding/json"
	"errors"

	"gopkg.in/yaml.v2"
)

// Encode encode data based on format
func Encode(in interface{}, format string) ([]byte, error) {

	switch format {
	case "yaml", "yml":
		return yaml.Marshal(in)
	case "json":
		return json.MarshalIndent(in, "", "    ")
	default:
		return nil, errors.New("Unknown format " + format)
	}
}

// Decode decode data based on format
func Decode(data []byte, out interface{}, format string) error {

	switch format {
	case "yaml", "yml":
		return yaml.Unmarshal(data, out)
	case "json":
		return json.Unmarshal(data, out)
	default:
		return errors.New("Unknown format " + format)
	}
}
