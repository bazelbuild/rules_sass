#!/usr/bin/env bash

bazel build --noshow_progress @io_bazel//src/tools/skylark/java/com/google/devtools/skylark/skylint:Skylint
find . -iname "*.bzl" | xargs $(bazel info bazel-bin)/external/io_bazel/src/tools/skylark/java/com/google/devtools/skylark/skylint/Skylint
