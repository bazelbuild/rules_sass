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

load("//sass:sass.bzl", "sass_binary", "sass_library")

def _basename(p):
    """Gets the basename of a given path."""
    return p.rpartition("/")[-1]

def _strip_extension(p):
    """Strips th Sass extension of the given path."""
    return p[:-len(".scss")]

def multi_sass_binary(name, srcs, output_style = None, sourcemap = None, **kwargs):
    """`multi_sass_binary` compiles a list of Sass files and outputs the corresponding
      CSS files and optional sourcemaps."""
    sass_library(
        name = "%s_lib" % name,
        srcs = srcs,
    )

    targets = []

    # Iterate through all source files and build individual Sass binary
    # targets that will be later grouped in a filegroup. Partial files
    # starting with an underscore will be skipped.
    for idx, input in enumerate(srcs):
        if (_basename(input).startswith("_")):
            continue

        output_name = "%s.css" % _strip_extension(input)
        target_name = "%s_piece--%s" % (name, input)
        targets.append(target_name)

        sass_binary(
            name = target_name,
            src = input,
            deps = [":%s_lib" % name],
            include_paths = [native.package_name()],
            output_name = output_name,
            output_style = output_style,
            sourcemap = sourcemap,
            tags = ["manual"],
            visibility = [":__pkg__"],
        )

    native.filegroup(
        name = name,
        srcs = targets,
        **kwargs
    )
