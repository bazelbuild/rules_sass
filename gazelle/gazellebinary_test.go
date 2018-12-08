/* Copyright 2018 The Bazel Authors. All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

   http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package gazelle

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bazelbuild/bazel-gazelle/testtools"
)

var (
	gazellePath = flag.String("gazelle", "", "path to gazelle binary")
)

func runGazelle(dir string) error {
	cmd := exec.Command(*gazellePath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = dir

	err := cmd.Run()

	// Flush stdout/stderr since this makes log output show up more often.
	os.Stdout.Sync()
	os.Stderr.Sync()
	return err
}

func TestMain(m *testing.M) {
	if os.Getenv("TEST_TARGET") != "//gazelle:go_default_test" {
		fmt.Printf("This test only works in Bazel. Check to see that you're invoking it through bazel and that the target is correct. To invoke this test run `bazel test //gazelle:go_default_test`")
		return
	}

	flag.Parse()
	if abs, err := filepath.Abs(*gazellePath); err != nil {
		log.Fatalf("unable to find absolute path for gazelle: %v\n", err)
		os.Exit(1)
	} else {
		*gazellePath = abs
	}
	os.Exit(m.Run())
}

func TestTestdata(t *testing.T) {
	testDataDir := filepath.Join("gazelle", "testdata")
	testDataFiles, err := ioutil.ReadDir(testDataDir)
	if err != nil {
		t.Fatalf("Error enumerating test modes: %s", err)
	}

	var testSuites []string
	for _, candidate := range testDataFiles {
		// Test suites are dirs, not files.
		if candidate.IsDir() {
			testSuites = append(testSuites, candidate.Name())
		}
	}

	if len(testSuites) == 0 {
		t.Fatalf("There should be one or more test suites defined in `testdata`. Please see the `README.md` for more information. Found candidates %v", testDataFiles)
	}

	for _, testSuite := range testSuites {
		t.Run(testSuite, func(t *testing.T) {
			var files []testtools.FileSpec
			var want []testtools.FileSpec
			// Walk testSuite.Name() to find all the input files and add them to the files or the expectations
			testSuitePath := filepath.Join(testDataDir, testSuite)
			err := filepath.Walk(testSuitePath, func(path string, info os.FileInfo, err error) error {
				// There is no need to process directories. Skip them.
				if info.IsDir() {
					return nil
				}

				content, err := ioutil.ReadFile(path)
				if err != nil {
					return fmt.Errorf("Unable to read file %q. Err: %v", path, err)
				}

				// content is a []byte not a string so it has to be typecast and we can't define the filespec at the beginning.
				fileSpec := testtools.FileSpec{Path: strings.TrimPrefix(path, testSuitePath), Content: string(content)}

				if strings.HasSuffix(path, "/BUILD.bazel.in") {
					fileSpec.Path = strings.TrimSuffix(strings.TrimPrefix(path, testSuitePath), "BUILD.bazel.in") + "BUILD.bazel"
					files = append(files, fileSpec)
				} else if strings.HasSuffix(path, "/BUILD.bazel.out") {
					fileSpec.Path = strings.TrimSuffix(strings.TrimPrefix(path, testSuitePath), "BUILD.bazel.out") + "BUILD.bazel"
					want = append(want, fileSpec)
				} else {
					files = append(files, fileSpec)
				}

				return nil
			})
			if err != nil {
				t.Errorf("Error walking %q %v", testSuitePath, err)
			}

			if len(files) == 0 {
				t.Fatalf("Test suites should have nonzero input files")
			}
			if len(want) == 0 {
				t.Fatalf("Test suites should have nonzero output files")
			}

			dir, cleanup := testtools.CreateFiles(t, files)
			defer cleanup()

			t.Logf("Dir: %q", dir)
			t.Logf("Files: %v", files)
			t.Logf("Want: %v", want)

			if err := runGazelle(dir); err != nil {
				t.Fatal(err)
			}

			testtools.CheckFiles(t, dir, want)
		})
	}
}
