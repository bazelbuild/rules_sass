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

load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

"""Fetches transitive dependencies required for using the Sass rules"""

def _include_if_not_defined(repo_rule, name, **kwargs):
    if name not in native.existing_rules():
        repo_rule(name = name, **kwargs)

def rules_sass_dependencies():
    # Since we use the Dart version of Sass, we need to be able to run NodeJS binaries.
    _include_if_not_defined(
        http_archive,
        name = "build_bazel_rules_nodejs",
        url = "https://github.com/bazelbuild/rules_nodejs/archive/0.15.3.zip",
        strip_prefix = "rules_nodejs-0.15.3",
        sha256 = "05afbbc13b0b7d5056e412d66c98853978bd46a94bc8e7b71c7fba4349b77eef",
    )

    # Dependencies from the NodeJS rules. We don't want to use the "package.bzl" dependency macro
    # of the NodeJS rules here because we just want to fetch dependencies and not load from them.
    # Loading the transitive dependencies here would mean that developers have no possibility
    # to overwrite individual transitive dependencies after `rule_sass_dependencies` ran.
    _include_if_not_defined(
        http_archive,
        name = "bazel_skylib",
        url = "https://github.com/bazelbuild/bazel-skylib/archive/0.5.0.zip",
        strip_prefix = "bazel-skylib-0.5.0",
        sha256 = "ca4e3b8e4da9266c3a9101c8f4704fe2e20eb5625b2a6a7d2d7d45e3dd4efffd",
    )

"""Fetches dependencies which are required **only** for development"""

def rules_sass_dev_dependencies():
    # Dependency for running Skylint.
    _include_if_not_defined(
        http_archive,
        name = "io_bazel",
        sha256 = "978f7e0440dd82182563877e2e0b7c013b26b3368888b57837e9a0ae206fd396",
        strip_prefix = "bazel-0.18.0",
        url = "https://github.com/bazelbuild/bazel/archive/0.18.0.zip",
    )

    # Required for the Buildtool repository.
    _include_if_not_defined(
        http_archive,
        name = "io_bazel_rules_go",
        urls = ["https://github.com/bazelbuild/rules_go/archive/cbc1e32fba771845305f15e341fa26595d4a136d.zip"],
        strip_prefix = "rules_go-cbc1e32fba771845305f15e341fa26595d4a136d",
        sha256 = "d02b1d8d11fb67fb1e451645256e58a1542170eedd6e2ba160c8540c96f659da",
    )

    # Bazel buildtools repo contains tools for BUILD file formatting ("buildifier") etc.
    _include_if_not_defined(
        http_archive,
        name = "com_github_bazelbuild_buildtools",
        sha256 = "a82d4b353942b10c1535528b02bff261d020827c9c57e112569eddcb1c93d7f6",
        strip_prefix = "buildtools-0.17.2",
        url = "https://github.com/bazelbuild/buildtools/archive/0.17.2.zip",
    )

    # Needed in order to generate documentation
    _include_if_not_defined(
        http_archive,
        name = "io_bazel_skydoc",
        url = "https://github.com/bazelbuild/skydoc/archive/b91215b629c67527a98466e61c11c77530fe85c1.zip",
        strip_prefix = "skydoc-b91215b629c67527a98466e61c11c77530fe85c1",
        sha256 = "3e16df2e6785a241248e6b92c657eb2c40b7a527aa48a1804b9a29806c8d5490",
    )
