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

"Install Sass toolchain dependencies"

load("@build_bazel_rules_nodejs//:index.bzl", "yarn_install", "check_rules_nodejs_version")

def sass_repositories():
    """Set up environment for Sass compiler.
    """

    # 0.31.1: entry_point attribute of rules_nodejs is now a label
    # 0.32.0: @npm//node_modules/foobar:foobar.js labels changed to @npm//:node_modules/foobar/foobar.js with fix for bazelbuild/rules_nodejs#802.
    check_rules_nodejs_version("0.32.0")

    yarn_install(
        name = "build_bazel_rules_sass_deps",
        package_json = "@io_bazel_rules_sass//sass:package.json",
        yarn_lock = "@io_bazel_rules_sass//sass:yarn.lock",
        # Do not symlink node_modules as when used in downstream repos we should not create
        # node_modules folders in the @io_bazel_rules_sass external repository. This is
        # not supported by managed_directories.
        symlink_node_modules = False,
    )
