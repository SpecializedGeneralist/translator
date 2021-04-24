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

import "github.com/SpecializedGeneralist/translator/pkg/api"

func (s *Server) makeErrors(req interface{}, err error) *api.ResponseErrors {
	s.logger.Trace().Err(err).Interface("request", req).Send()
	return &api.ResponseErrors{
		Value: []*api.ResponseError{
			{Message: err.Error()},
		},
	}
}

func (s *Server) makeFatalErrors(req interface{}, err error) *api.ResponseErrors {
	s.logger.Error().Err(err).Interface("request", req).Send()
	return &api.ResponseErrors{
		Value: []*api.ResponseError{
			{Message: err.Error()},
		},
	}
}
