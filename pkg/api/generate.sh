#!/usr/bin/env sh

# Copyright 2021 SpecializedGeneralist Authors
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

set -ex

# Generate api.proto
gnostic \
  --grpc-out=. \
  api.yaml

sed \
  -i "2i option go_package = \"github.com/SpecializedGeneralist/translator/pkg/api\";" \
  api.proto

# Generate api.pb.go
protoc \
  --go_out=. \
  --go_opt='paths=source_relative' \
  api.proto

# Generate api_grpc.pb.go
protoc \
  --go-grpc_out=. \
  --go-grpc_opt='paths=source_relative' \
  api.proto

# Generate api_descriptor.pb
protoc \
  --proto_path=. \
  --include_imports \
  --include_source_info \
  --descriptor_set_out=api_descriptor.pb \
  api.proto

# Generate api.pb.gw.go
protoc \
  --grpc-gateway_out=. \
  --grpc-gateway_opt='logtostderr=true' \
  --grpc-gateway_opt='paths=source_relative' \
  api.proto
