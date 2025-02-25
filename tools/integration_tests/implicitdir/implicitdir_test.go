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

// Provides integration tests when implicit_dir flag is set.
package implicitdir_test

import (
	"flag"
	"log"
	"os"
	"testing"

	"github.com/googlecloudplatform/gcsfuse/tools/integration_tests/setup"
)

func TestMain(m *testing.M) {
	flag.Parse()

	if setup.TestBucket() == "" && setup.MountedDirectory() == "" {
		log.Printf("--testbucket or --mountedDirectory must be specified")
		os.Exit(0)
	} else if setup.TestBucket() != "" && setup.MountedDirectory() != "" {
		log.Printf("Both --testbucket and --mountedDirectory can't be specified at the same time.")
		os.Exit(0)
	}

	if setup.MountedDirectory() != "" {
		setup.SetMntDir(setup.MountedDirectory())
		successCode := setup.ExecuteTest(m)
		os.Exit(successCode)
	}

	if err := setup.SetUpTestDir(); err != nil {
		log.Printf("setUpTestDir: %v\n", err)
		os.Exit(1)
	}

	flags := [][]string{{"--enable-storage-client-library=true", "--implicit-dirs=true"},
		{"--enable-storage-client-library=false"},
		{"--implicit-dirs=true"},
		{"--implicit-dirs=false"}}

	successCode := setup.ExecuteTestForFlags(flags, m)

	log.Printf("Test log: %s\n", setup.LogFile())
	os.Exit(successCode)
}
