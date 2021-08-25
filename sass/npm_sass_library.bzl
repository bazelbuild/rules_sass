# Copyright 2021 The Bazel Authors. All rights reserved.
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

load("@build_bazel_rules_nodejs//:providers.bzl", "ExternalNpmPackageInfo")
load("@io_bazel_rules_sass//sass:sass.bzl", "SassInfo")

def _is_sass_file(filename):
    """Checks whether the specified filename resolves to a file supported by Sass."""
    if filename.endswith(".scss"):
        return True
    elif filename.endswith(".sass"):
        return True
    elif filename.endswith(".css"):
        return True
    return False

def _filter_sass_files(files):
    """Filters a list of files to only return files which are supported by Sass."""
    return [f for f in files if _is_sass_file(f.short_path)]


def _npm_sass_library_impl(ctx):
    """
    Rule that extracts Sass sources and its transitive dependencies from an npm
    package. The extracted source files are provided with the `SassInfo` provider
    so that they can be consumed directly as dependencies of other Sass libraries
    or Sass binaries.

    This rule is helpful when build targets rely on Sass files provided by an external
    npm package. In those cases, one wouldn't want to list out all individual source
    files of the npm package, but rather glob all needed Sass files from the npm package.
    """

    transitive_sources = []

    # Iterate through all specified dependencies and collect Sass files from build
    # targets that have the `ExternalNpmPackageInfo` provider set. The `yarn_install`
    # rule automatically sets these providers for individual targets in `@npm//<..>`.
    for dep in ctx.attr.deps:
        npm_info = dep[ExternalNpmPackageInfo]
        filered_files = _filter_sass_files(npm_info.sources.to_list())
        transitive_sources.append(depset(filered_files))

    # Convert the collected transitive Sass sources to a depset. This is necessary
    # for proper deduping of dependencies. Performance-wise it's not efficient that
    # we need to unwrap the depset for npm packages, but this is necessary as otherwise
    # many unused files would end up being action inputs. This ensures efficient runfile
    # trees and proper sandboxing for targets relying on providers of this rule.
    outputs = depset(transitive = transitive_sources)

    return [
        DefaultInfo(
            files = outputs,
            runfiles = ctx.runfiles(transitive_files = outputs),
        ),
        SassInfo(transitive_sources = outputs),
    ]

npm_sass_library = rule(
    implementation = _npm_sass_library_impl,
    attrs = {
        "deps": attr.label_list(
            allow_files = False,
            mandatory = True,
            providers = [ExternalNpmPackageInfo],
            doc = "List of npm package targets for which direct and transitive Sass files are collected."
        ),
    },
)
"""Rule that collects Sass files from npm package targets and exposes them for consumption
within `sass_binary` or `sass_library` targets."""
