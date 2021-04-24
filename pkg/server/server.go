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

package server

import (
	"context"
	"fmt"
	"github.com/SpecializedGeneralist/translator/pkg/api"
	"github.com/SpecializedGeneralist/translator/pkg/configuration"
	"github.com/SpecializedGeneralist/translator/pkg/models"
	"github.com/rs/zerolog"
	"runtime/debug"
	"time"
)

// Server is the main implementation of api.ApiServer.
type Server struct {
	api.UnimplementedApiServer
	config  *configuration.Config
	manager *models.Manager
	logger  zerolog.Logger
}

// New creates a new Server.
func New(config *configuration.Config, manager *models.Manager, logger zerolog.Logger) *Server {
	return &Server{
		config:  config,
		manager: manager,
		logger:  logger,
	}
}

// TranslateText translates a text.
func (s *Server) TranslateText(_ context.Context, req *api.TranslateTextRequest) (resp *api.TranslateTextResponse, _ error) {
	defer func() {
		if r := recover(); r != nil {
			st := string(debug.Stack())
			resp = &api.TranslateTextResponse{Errors: s.makeFatalErrors(req, fmt.Errorf("panic: %v\n%s", r, st))}
		}
	}()

	startTime := time.Now()

	in := req.GetTranslateTextInput()
	source := in.GetSourceLanguage()
	target := in.GetTargetLanguage()
	text := in.GetText()

	translatedText, err := s.manager.Translate(source, target, text)

	if err != nil {
		return &api.TranslateTextResponse{Errors: s.makeErrors(req, err)}, nil
	}

	elapsedTime := time.Since(startTime)
	resp = &api.TranslateTextResponse{
		Data: &api.TranslateTextData{
			TranslatedText: translatedText,
			Took:           float32(elapsedTime.Seconds()),
		},
	}
	return resp, nil
}
