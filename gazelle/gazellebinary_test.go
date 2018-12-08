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
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/bazelbuild/bazel-gazelle/testtools"
)

var (
	gazellePath = flag.String("gazelle", "", "path to gazelle binary")
	baseFiles   = []testtools.FileSpec{
		{Path: "WORKSPACE"},
		//{Path: "BUILD.bazel"},
		{Path: "hello_world/main.scss", Content: `
@import "shared/fonts";
@import "shared/colors";

html {
  body {
    font-family: $default-font-stack;
    h1 {
      font-family: $modern-font-stack;
      color: $example-red;
    }
  }
}
`},
		{Path: "shared/_fonts.scss", Content: `
$default-font-stack: Cambria, "Hoefler Text", serif;
$modern-font-stack: Constantia, "Lucida Bright", serif;
`},
		{Path: "shared/_colors.scss", Content: `
$example-blue: #0000ff;
$example-red: #ff0000;
`},
	}
)

func TestMain(m *testing.M) {
	_, ok := os.LookupEnv("TEST_TARGET")
	if !ok {
		// Skip all tests if we aren't run by Bazel
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

func TestGazelleBinary(t *testing.T) {
	files := append(baseFiles)
	dir, cleanup := testtools.CreateFiles(t, files)
	defer cleanup()

	cmd := exec.Command(*gazellePath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = dir
	if err := cmd.Run(); err != nil {
		t.Fatal(err)
	}

	testtools.CheckFiles(t, dir, []testtools.FileSpec{{
		Path: "hello_world/BUILD.bazel",
		Content: `load("@io_bazel_rules_sass//:defs.bzl", "sass_binary")

sass_binary(
    name = "hello_world",
    src = "main.scss",
    deps = [
        "//shared:colors",
        "//shared:fonts",
    ],
)`,
	}, {
		Path: "shared/BUILD.bazel",
		Content: `load("@io_bazel_rules_sass//:defs.bzl", "sass_library")

sass_library(
    name = "colors",
    srcs = ["_colors.scss"],
    visibility = ["//visibility:public"],
)

sass_library(
    name = "fonts",
    srcs = ["_fonts.scss"],
    visibility = ["//visibility:public"],
)
`,
	}})
}
