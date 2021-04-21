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

package osutils

import (
	"fmt"
	"os"
)

// DirExists reports whether a directory exists.
// It returns an error if the underlying os.Stat call fails, or if the given
// path does not correspond to a directory.
func DirExists(name string) (bool, error) {
	info, err := os.Stat(name)
	if os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("error checking if %#v exists: %w", name, err)
	}
	if !info.IsDir() {
		return false, fmt.Errorf("%#v exists but is not a directory", name)
	}
	return true, nil
}

// FileExists reports whether a file exists.
// It returns an error if the underlying os.Stat call fails, or if the given
// path does not correspond is a directory instead of a file.
func FileExists(name string) (bool, error) {
	info, err := os.Stat(name)
	if os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("error checking if %#v exists: %w", name, err)
	}
	if info.IsDir() {
		return false, fmt.Errorf("%#v exists but is not a file", name)
	}
	return true, nil
}
