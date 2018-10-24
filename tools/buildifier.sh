#!/usr/bin/env bash

bazel build --noshow_progress @com_github_bazelbuild_buildtools//buildifier

find . -name "BUILD" -or -name "BUILD.bazel" -or -iname "*.bzl" \
 | xargs $(bazel info bazel-bin)/external/com_github_bazelbuild_buildtools/buildifier/*/buildifier
