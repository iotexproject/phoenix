// Copyright (c) 2020 IoTeX Foundation
// This is an alpha (internal) release and is not suitable for production. This source code is provided 'as is' and no
// warranties are given as to title or non-infringement, merchantability or fitness for purpose and, to the extent
// permitted by law, all liability for your use of the code is disclaimed. This source code is governed by Apache
// License 2.0 that can be found in the LICENSE file.

package config

import (
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/iotexproject/phoenix-gem/log"
	"github.com/pkg/errors"
)

var (
	// Default is the default config
	Default = Config{
		SubLogs: make(map[string]log.GlobalConfig),
	}
)

type (
	Pinata struct {
		Uri          string `yaml:"uri"`
		ApiKey       string `yaml:"apiKey"`
		SecretApiKey string `yaml:"secretApiKey"`
	}
	Server struct {
		Port string `yaml:"port"`
	}
	Config struct {
		Pinata  Pinata                      `pinata`
		Server  Server                      `yaml:"server"`
		Log     log.GlobalConfig            `yaml:"log"`
		SubLogs map[string]log.GlobalConfig `yaml:"subLogs"`
	}
)

func New(path string) (cfg *Config, err error) {
	body, err := ioutil.ReadFile(path)
	if err != nil {
		return cfg, errors.Wrap(err, "failed to read config content")
	}
	fileExt := "yaml"
	extWithDot := filepath.Ext(path)
	if strings.HasPrefix(extWithDot, ".") {
		fileExt = extWithDot[1:]
	}
	cfg = &Default
	if err = Decode(body, cfg, fileExt); err != nil {
		return cfg, errors.Wrap(err, "failed to unmarshal config to struct")
	}
	return
}
