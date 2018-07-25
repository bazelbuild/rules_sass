#!/usr/bin/env bash

# rules_go (used to compile `buildifier`) still uses the deprecated APIs for
# args.add, so we explicitly set --noincompatible_disallow_old_style_args_add
# to avoid errors, since this project's bazelrc has this set the other way.
bazel build --noincompatible_disallow_old_style_args_add --noshow_progress @com_github_bazelbuild_buildtools//buildifier

find . -name "BUILD" -or -name "BUILD.bazel" -or -iname "*.bzl" | xargs $(bazel info bazel-bin)/external/com_github_bazelbuild_buildtools/buildifier/*/buildifier
