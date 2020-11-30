// Copyright (c) 2020 IoTeX
// This is an alpha (internal) release and is not suitable for production. This source code is provided 'as is' and no
// warranties are given as to title or non-infringement, merchantability or fitness for purpose and, to the extent
// permitted by law, all liability for your use of the code is disclaimed. This source code is governed by Apache
// License 2.0 that can be found in the LICENSE file.

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
