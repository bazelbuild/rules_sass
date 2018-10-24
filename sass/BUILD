package(default_visibility = ["//visibility:public"])

load("@build_bazel_rules_nodejs//:defs.bzl", "nodejs_binary")

exports_files([
    "sass.bzl",
    "sass_repositories.bzl",
])

# Executable for the sass_binary rule
nodejs_binary(
    name = "sass",
    entry_point = "sass/sass.js",
    install_source_map_support = False,
    data = [
        "@build_bazel_rules_sass_deps//sass",
    ],
)
