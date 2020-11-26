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
		Uri          string `yaml:"uri" json:"uri"`
		ApiKey       string `yaml:"apiKey" json:"apiKey"`
		SecretApiKey string `yaml:"secretApiKey" json:"secretApiKey"`
	}
	RateLimit struct {
		RequestLimit int `yaml:"requestLimit" json:"requestLimit"`
		WindowLength int `yaml:"windowLength" json:"windowLength"`
	}
	Server struct {
		Port       string    `yaml:"port" json:"port"`
		AuthSecret string    `yaml:"authSecret" json:"authSecret"`
		RateLimit  RateLimit `yaml:"rateLimit" json:"rateLimit"`
	}
	Storage struct {
		Provider string `yaml:"provider" json:"provider"`
	}
	S3 struct {
		EndPoint  string `yaml:"endpoint" json:"endpoint"`
		AccessKey string `yaml:"accessKey" json:"accessKey"`
		SecretKey string `yaml:"secretKey" json:"secretKey"`
		Region    string `yaml:"region" json:"region"`
	}
	Config struct {
		Pinata  Pinata                      `yaml:"pinata" json:"pinata"`
		Server  Server                      `yaml:"server" json:"server"`
		Log     log.GlobalConfig            `yaml:"log" yaml:"log"`
		SubLogs map[string]log.GlobalConfig `yaml:"subLogs" json:"subLogs"`
		Storage Storage                     `yaml:"storage" json:"storage"`
		S3      S3                          `yaml:"s3" json:"s3"`
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
