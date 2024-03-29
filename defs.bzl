# Copyright 2018 The Bazel Authors. All rights reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#    http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

""" Public API is re-exported here."""

load("//sass:sass_repositories.bzl", _sass_repositories = "sass_repositories")
load(
  "//sass:sass.bzl",
  _SassInfo = "SassInfo",
  _sass_binary = "sass_binary",
  _sass_library = "sass_library",
  _multi_sass_binary = "multi_sass_binary",
)
load("//sass:npm_sass_library.bzl", _npm_sass_library = "npm_sass_library")

sass_repositories = _sass_repositories

sass_library = _sass_library
sass_binary = _sass_binary
multi_sass_binary = _multi_sass_binary
npm_sass_library = _npm_sass_library

# Expose the SassInfo provider so that people can make their own custom rules
# that expose sass library outputs.
SassInfo = _SassInfo
