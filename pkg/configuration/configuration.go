// Copyright 2021 SpecializedGeneralist Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package configuration

import (
	"bytes"
	"fmt"
	"github.com/rs/zerolog"
	"gopkg.in/yaml.v3"
	"os"
)

// Config is the main translator server configuration.
type Config struct {
	// LogLevel is the minimum severity level for log messages.
	LogLevel LogLevel `yaml:"log_level"`
	// Host is the server binding address.
	Host string `yaml:"host"`
	// Port is the server listening port.
	Port int `yaml:"port"`
	// MaxConcurrentComputations is the maximum amount of concurrent
	// computations allowed.
	MaxConcurrentComputations int `yaml:"max_concurrent_computations"`
	// TLSEnabled reports whether to enable TLS.
	TLSEnabled bool `yaml:"tls_enabled"`
	// TLSCert is the TLS cert file. It is ignored if TLSEnabled is false.
	TLSCert string `yaml:"tls_cert"`
	// TLSKey is the TLS key file. It is ignored if TLSEnabled is false.
	TLSKey string `yaml:"tls_key"`
	// ModelsPath is the local path for all spaGO-compatible models.
	ModelsPath string `yaml:"models_path"`
	// LanguageModels provides the configuration for translation models
	// to be loaded and handled by the service.
	LanguageModels []LanguageModel `yaml:"language_models"`
}

// LanguageModel identifies a single model and the identifiers for the source
// and target languages it supports.
type LanguageModel struct {
	// Source is an identifier for the source language of translation.
	Source string `yaml:"source"`
	// Target is an identifier for the target language of translation.
	Target string `yaml:"target"`
	// Model is the name of a spaGO-compatible model.
	Model string `yaml:"model"`
}

// LogLevel is a redefinition of zerolog.Level which satisfies
// encoding.TextUnmarshaler.
type LogLevel zerolog.Level

// UnmarshalText unmarshals the text to a LogLevel.
func (l *LogLevel) UnmarshalText(text []byte) (err error) {
	zl, err := zerolog.ParseLevel(string(text))
	*l = LogLevel(zl)
	return err
}

// FromYAMLFile reads a Config object from a YAML file.
func FromYAMLFile(filename string) (*Config, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("error reading configuration file %#v: %w", filename, err)
	}
	content = []byte(os.ExpandEnv(string(content)))
	config := new(Config)
	err = yaml.NewDecoder(bytes.NewReader(content)).Decode(&config)
	if err != nil {
		return nil, fmt.Errorf("error decoding configuration YAML file %#v: %w", filename, err)
	}
	return config, nil
}
