load("@build_bazel_rules_nodejs//:index.bzl", "js_library")

js_library(
    name = "test_external_pkg",
    # Special attribute that instructs `js_library` to expose the `ExternalNpmPackageInfo` provider
    # for this library. This emulates an external NPM package installed through `yarn_install`.
    external_npm_package = True,
    srcs = ["_index.scss", "_other_file.scss"],
    visibility = ["//sass/test/npm_sass_library:__pkg__"],
)

