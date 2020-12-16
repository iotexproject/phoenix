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

	"github.com/iotexproject/phoenix/log"
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
		Uri          string `yaml:"uri" json:"uri"`
		ApiKey       string `yaml:"apiKey" json:"apiKey"`
		SecretApiKey string `yaml:"secretApiKey" json:"secretApiKey"`
	}
	RateLimit struct {
		Enable       bool     `yaml:"enable" json:"enable"`
		LimitByKey   []string `yaml:"limitByKey" json:"limitByKey"` // support ["ip", "url", "user"]
		RequestLimit int      `yaml:"requestLimit" json:"requestLimit"`
		WindowLength int      `yaml:"windowLength" json:"windowLength"`
	}
	Cors struct {
		Enable         bool     `yaml:"enable" json:"enable"`
		AllowedOrigins []string `yaml:"allowedOrigins" json:"allowedOrigins"`
		AllowedMethods []string `yaml:"allowedMethods" json:"allowedMethods"`
		AllowedHeaders []string `yaml:"allowedHeaders" json:"allowedHeaders"`
	}
	Server struct {
		Port      string    `yaml:"port" json:"port"`
		RateLimit RateLimit `yaml:"rateLimit" json:"rateLimit"`
		Cors      Cors      `yaml:"cors" json:"cors"`
		DBPath    string    `yaml:"dbPath" json:"dbPath"`
	}
	Config struct {
		Pinata  Pinata                      `yaml:"pinata" json:"pinata"`
		Server  Server                      `yaml:"server" json:"server"`
		Log     log.GlobalConfig            `yaml:"log" yaml:"log"`
		SubLogs map[string]log.GlobalConfig `yaml:"subLogs" json:"subLogs"`
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
