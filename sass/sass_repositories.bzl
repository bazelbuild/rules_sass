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

load("@build_bazel_rules_nodejs//:defs.bzl", "yarn_install")

LIBSASS_BUILD_FILE = """
package(default_visibility = ["@sassc//:__pkg__"])
filegroup(
    name = "srcs",
    srcs = glob([
         "src/**/*.h*",
         "src/**/*.c*",
    ]),
)
# Includes directive may seem unnecessary here, but its needed for the weird
# interplay between libsass/sassc projects. This is intentional.
cc_library(
    name = "headers",
    includes = ["include"],
    hdrs = glob(["include/**/*.h"]),
)
"""

SASSC_BUILD_FILE = """
package(default_visibility = ["//visibility:public"])
cc_binary(
    name = "sassc",
    srcs = [
        "@libsass//:srcs",
        "sassc.c",
        "sassc_version.h",
    ] + select({
        "@bazel_tools//src/conditions:windows": glob([
            "win/**/*.c",
            "win/**/*.h",
        ]),
        "//conditions:default": [],
    }),
    includes = select({
        "@bazel_tools//src/conditions:windows": ["win/posix"],
        "//conditions:default": [],
    }),
    linkopts = select({
        "@bazel_tools//src/conditions:windows": ["-DEFAULTLIB:shell32.lib"],
        "//conditions:default": [
            "-ldl",
            "-lm",
        ],
    }),
    deps = ["@libsass//:headers"],
)
"""

def sass_repositories():
  """Provide different toolchains depending which compiler the user chooses.

  Note that Bazel will only download the tools needed for this compilation, so
  typically we won't execute all the rules below.
  """

  yarn_install(
      name = "build_bazel_rules_sass_compiletime_deps",
      package_json = "//sass:package.json",
      yarn_lock = "//sass:yarn.lock",
  )

  native.new_http_archive(
      name = "libsass",
      url = "http://bazel-mirror.storage.googleapis.com/github.com/sass/libsass/archive/3.3.0-beta1.tar.gz",
      sha256 = "6a4da39cc0b585f7a6ee660dc522916f0f417c890c3c0ac7ebbf6a85a16d220f",
      build_file_content = LIBSASS_BUILD_FILE,
      strip_prefix = "libsass-3.3.0-beta1",
  )

  native.new_http_archive(
      name = "sassc",
      url = "http://bazel-mirror.storage.googleapis.com/github.com/sass/sassc/archive/3.3.0-beta1.tar.gz",
      sha256 = "87494218eea2441a7a24b40f227330877dbba75c5fa9014ac6188711baed53f6",
      build_file_content = SASSC_BUILD_FILE,
      strip_prefix = "sassc-3.3.0-beta1",
  )
