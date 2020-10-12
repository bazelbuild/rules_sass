load("@build_bazel_rules_nodejs//:index.bzl", "nodejs_binary")

package(default_visibility = ["//visibility:public"])

exports_files([
    "sass.bzl",
    "sass_repositories.bzl",
])

# Executable for the sass_binary rule
nodejs_binary(
    name = "sass",
    data = [
        ":sass_wrapper.js",
        "@build_bazel_rules_sass_deps//@bazel/worker",
        "@build_bazel_rules_sass_deps//sass",
        "@build_bazel_rules_sass_deps//yargs",
    ],
    entry_point = "sass_wrapper",
)
