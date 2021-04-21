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

package osutils_test

import (
	"github.com/SpecializedGeneralist/translator/pkg/osutils"
	"github.com/stretchr/testify/assert"
	"os"
	"path"
	"testing"
)

func TestDirExists(t *testing.T) {
	t.Parallel()

	t.Run("existing directory", func(t *testing.T) {
		t.Parallel()
		dir := createTempDir(t)
		defer deleteDir(t, dir)
		exists, err := osutils.DirExists(dir)
		assert.True(t, exists)
		assert.NoError(t, err)
	})

	t.Run("nonexisting directory", func(t *testing.T) {
		t.Parallel()
		dir := createTempDir(t)
		defer deleteDir(t, dir)
		exists, err := osutils.DirExists(path.Join(dir, "foo"))
		assert.False(t, exists)
		assert.NoError(t, err)
	})

	t.Run("file", func(t *testing.T) {
		t.Parallel()
		filename := createTempFile(t)
		defer deleteFile(t, filename)
		exists, err := osutils.DirExists(filename)
		assert.False(t, exists)
		assert.Error(t, err)
	})
}

func TestFileExists(t *testing.T) {
	t.Parallel()

	t.Run("existing file", func(t *testing.T) {
		t.Parallel()
		filename := createTempFile(t)
		defer deleteFile(t, filename)
		exists, err := osutils.FileExists(filename)
		assert.True(t, exists)
		assert.NoError(t, err)
	})

	t.Run("nonexisting file", func(t *testing.T) {
		t.Parallel()
		dir := createTempDir(t)
		defer deleteDir(t, dir)
		exists, err := osutils.FileExists(path.Join(dir, "foo"))
		assert.False(t, exists)
		assert.NoError(t, err)
	})

	t.Run("directory", func(t *testing.T) {
		t.Parallel()
		dir := createTempDir(t)
		defer deleteDir(t, dir)
		exists, err := osutils.FileExists(dir)
		assert.False(t, exists)
		assert.Error(t, err)
	})
}

func createTempDir(t *testing.T) string {
	dir, err := os.MkdirTemp("", "osutils_test")
	assert.NoError(t, err)
	return dir
}

func deleteDir(t *testing.T, dir string) {
	err := os.RemoveAll(dir)
	assert.NoError(t, err)
}

func createTempFile(t *testing.T) string {
	file, err := os.CreateTemp("", "osutils_test")
	assert.NoError(t, err)
	filename := file.Name()
	err = file.Close()
	assert.NoError(t, err)
	return filename
}

func deleteFile(t *testing.T, name string) {
	err := os.Remove(name)
	assert.NoError(t, err)
}
