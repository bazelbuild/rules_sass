package(default_visibility = ["//visibility:public"])

load("@build_bazel_rules_nodejs//:index.bzl", "nodejs_binary")

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
        "@build_bazel_rules_sass_deps//minimist",
        "@build_bazel_rules_sass_deps//sass",
    ],
    entry_point = "sass_wrapper.js",
)
