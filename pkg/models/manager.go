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

package models

import (
	"fmt"
	"github.com/SpecializedGeneralist/translator/pkg/configuration"
	"github.com/rs/zerolog"
)

// Manager allows easy handling of multiple translation models.
type Manager struct {
	config *configuration.Config
	models modelsMap
	logger zerolog.Logger
}

// NewManager creates a new Manager.
func NewManager(config *configuration.Config, logger zerolog.Logger) *Manager {
	return &Manager{
		config: config,
		models: make(modelsMap, 1),
		logger: logger,
	}
}

// modelsMap maps [source][target] => *Model
type modelsMap map[string]map[string]*Model

// LoadModels loads all models according to the configuration.
// If a model path is not found, automatic download and conversion
// are performed using spaGO huggingface Downloader and Converter.
func (mng *Manager) LoadModels() error {
	mng.logger.Info().Msg("loading all models...")
	if len(mng.models) > 0 {
		return fmt.Errorf("models already loaded")
	}

	for _, lm := range mng.config.LanguageModels {
		err := mng.loadModel(lm)
		if err != nil {
			return err
		}
	}

	mng.logger.Info().Msg("all models loaded successfully")
	return nil
}

// GetModel returns a Model for translating texts from the given source
// language to the given target language. It also reports whether a model for
// that pair or languages is present (previously loaded).
func (mng *Manager) GetModel(source, target string) (*Model, bool) {
	sourceMap, sourceOk := mng.models[source]
	if !sourceOk {
		return nil, false
	}
	model, modelOk := sourceMap[target]
	if !modelOk {
		return nil, false
	}
	return model, true
}

// Translate is a convenience method to get a model and perform translation
// in a single step.
func (mng *Manager) Translate(source, target, text string) (string, error) {
	model, modelFound := mng.GetModel(source, target)
	if !modelFound {
		return "", fmt.Errorf("no model available for translation from %#v to %#v", source, target)
	}

	translatedText := model.Translate(text)
	return translatedText, nil
}

func (mng *Manager) loadModel(ln configuration.LanguageModel) error {
	if _, ok := mng.models[ln.Source]; !ok {
		mng.models[ln.Source] = make(map[string]*Model, 1)
	}
	if _, ok := mng.models[ln.Source][ln.Target]; ok {
		return fmt.Errorf("a model was already loaded for translation from %#v to %#v", ln.Source, ln.Target)
	}

	model := NewModel(mng.config, ln.Model, mng.logger)
	mng.models[ln.Source][ln.Target] = model

	return model.Load()
}
