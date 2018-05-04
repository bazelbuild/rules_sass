package(default_visibility = ["//visibility:public"])

load("@build_bazel_rules_nodejs//:defs.bzl", "nodejs_binary")

# Executable for the sass_binary rule
nodejs_binary(
    name = "sass",
    node_modules = "@build_bazel_rules_sass_compiletime_deps//:node_modules",
    entry_point = "sass/sass.js",
)

load("@io_bazel_skydoc//skylark:skylark.bzl", "skylark_doc")

skylark_doc(
    name = "docs",
    srcs = [
        ":sass.bzl",
        ":sass_repositories.bzl",
    ],
    format = "html",
    site_root = "/rules_sass",
    #strip_prefix = "/",
)
