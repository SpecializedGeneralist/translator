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
	"github.com/SpecializedGeneralist/translator/pkg/osutils"
	"github.com/nlpodyssey/spago/pkg/ml/ag"
	"github.com/nlpodyssey/spago/pkg/ml/nn"
	"github.com/nlpodyssey/spago/pkg/nlp/tokenizers/sentencepiece"
	bartconfig "github.com/nlpodyssey/spago/pkg/nlp/transformers/bart/config"
	"github.com/nlpodyssey/spago/pkg/nlp/transformers/bart/head/conditionalgeneration"
	"github.com/nlpodyssey/spago/pkg/nlp/transformers/bart/loader"
	"github.com/nlpodyssey/spago/pkg/nlp/transformers/huggingface"
	"github.com/rs/zerolog"
	"path"
)

// Model provides high-level functionalities for loading spaGO models and
// performing automatic translation.
type Model struct {
	config    *configuration.Config
	name      string
	model     nn.Model
	tokenizer *sentencepiece.Tokenizer
	logger    zerolog.Logger
}

// NewModel creates a new Model.
func NewModel(config *configuration.Config, name string, logger zerolog.Logger) *Model {
	return &Model{
		config:    config,
		name:      name,
		model:     nil,
		tokenizer: nil,
		logger:    logger.With().Str("model", name).Logger(),
	}
}

// Load loads the underlying spaGO model.
// If the model path is not found, automatic download and conversion
// are performed using spaGO huggingface Downloader and Converter.
func (m *Model) Load() error {
	if m.model != nil {
		return fmt.Errorf("model already loaded")
	}

	err := m.downloadSpagoModelIfMissing()
	if err != nil {
		return err
	}

	err = m.convertModelIfNecessary()
	if err != nil {
		return err
	}

	m.logger.Info().Msg("loading model...")

	modelPath := path.Join(m.config.ModelsPath, m.name)
	m.model, err = loader.Load(modelPath)
	if err != nil {
		return err
	}

	m.logger.Info().Msg("loading tokenizer...")

	m.tokenizer, err = sentencepiece.NewFromModelFolder(modelPath, false)
	if err != nil {
		return err
	}

	m.logger.Info().Msg("model and tokenizer loaded successfully")
	return nil
}

// Translate performs automatic translation of the given text.
func (m *Model) Translate(text string) string {
	g := ag.NewGraph(ag.IncrementalForward(false))
	defer g.Clear()

	proc := nn.Reify(nn.Context{Graph: g, Mode: nn.Inference}, m.model).(*conditionalgeneration.Model)
	bartConfig := proc.BART.Config

	tokens := m.tokenizer.Tokenize(text)
	tokenIDs := m.tokenizer.TokensToIDs(tokens)

	tokenIDs = append(tokenIDs, bartConfig.EosTokenID)

	rawGeneratedIDs := proc.Generate(tokenIDs)
	generatedIDs := stripBadTokens(rawGeneratedIDs, bartConfig)

	generatedTokens := m.tokenizer.IDsToTokens(generatedIDs)
	return m.tokenizer.Detokenize(generatedTokens)
}

func stripBadTokens(ids []int, bcfg bartconfig.Config) []int {
	result := make([]int, 0, len(ids))
	for _, id := range ids {
		if id == bcfg.EosTokenID || id == bcfg.PadTokenID ||
			id == bcfg.BosTokenID || id == bcfg.DecoderStartTokenID {
			continue
		}
		result = append(result, id)
	}
	return result
}

func (m *Model) downloadSpagoModelIfMissing() error {
	modelPath := path.Join(m.config.ModelsPath, m.name)

	modelPathExists, err := osutils.DirExists(modelPath)
	if err != nil {
		return err
	}
	if modelPathExists {
		return nil
	}

	m.logger.Info().Msgf("%#v does not exist", modelPath)
	return m.downloadSpagoModel()
}

func (m *Model) downloadSpagoModel() error {
	m.logger.Info().Msg("downloading model from Hugging Face models hub...")

	modelsPath := m.config.ModelsPath
	modelsPathExists, err := osutils.DirExists(modelsPath)
	if err != nil {
		return err
	}
	if !modelsPathExists {
		return fmt.Errorf("models path %#v does not exist", modelsPath)
	}

	downloader := huggingface.NewDownloader(modelsPath, m.name, false)
	err = downloader.Download()
	if err != nil {
		return err
	}

	m.logger.Info().Msg("model downloaded successfully")

	return m.convertModel()
}

const defaultSpagoModelFilename = "spago_model.bin"

func (m *Model) convertModelIfNecessary() error {
	modelPath := path.Join(m.config.ModelsPath, m.name)
	spagoModelFilename := path.Join(modelPath, defaultSpagoModelFilename)

	spagoModelFileExists, err := osutils.FileExists(spagoModelFilename)
	if err != nil {
		return err
	}
	if spagoModelFileExists {
		return nil
	}

	m.logger.Info().Msgf("%#v does not exist", spagoModelFilename)
	return m.convertModel()
}

func (m *Model) convertModel() error {
	m.logger.Info().Msg("converting model...")
	converter := huggingface.NewConverter(m.config.ModelsPath, m.name)
	err := converter.Convert()
	if err != nil {
		return err
	}
	m.logger.Info().Msg("model converted successfully")
	return nil
}
