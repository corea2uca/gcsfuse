// Copyright 2021 Google Inc. All Rights Reserved.
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

// Provides integration tests for read flows with implicit_dir flag set.
package implicitdir_test

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/googlecloudplatform/gcsfuse/tools/integration_tests/setup"
)

func TestReadAfterWrite(t *testing.T) {
	tmpDir, err := ioutil.TempDir(setup.MntDir(), "tmpDir")
	if err != nil {
		t.Errorf("Mkdir at %q: %v", setup.MntDir(), err)
		return
	}

	for i := 0; i < 10; i++ {
		tmpFile, err := ioutil.TempFile(tmpDir, "tmpFile")
		if err != nil {
			t.Errorf("Create file at %q: %v", tmpDir, err)
			return
		}

		fileName := tmpFile.Name()
		if _, err := tmpFile.WriteString("line 1\n"); err != nil {
			t.Errorf("WriteString: %v", err)
		}
		if err := tmpFile.Close(); err != nil {
			t.Errorf("Close: %v", err)
		}

		// After write, data will be cached by kernel. So subsequent read will be
		// served using cached data by kernel instead of calling gcsfuse.
		// Clearing kernel cache to ensure that gcsfuse is invoked during read operation.
		err = setup.ClearKernelCache()
		if err != nil {
			t.Errorf("Clear Kernel Cache: %v", err)
		}

		tmpFile, err = os.Open(fileName)
		if err != nil {
			t.Errorf("Open %q: %v", fileName, err)
			return
		}

		content, err := ioutil.ReadAll(tmpFile)
		if err != nil {
			t.Errorf("ReadAll: %v", err)
		}
		if got, want := string(content), "line 1\n"; got != want {
			t.Errorf("File content %q not match %q", got, want)
		}
	}
}
