package(default_visibility = ["//visibility:public"])

load("@build_bazel_rules_nodejs//:index.bzl", "nodejs_binary")

exports_files([
    "sass.bzl",
    "sass_repositories.bzl",
    "npm_sass_library.bzl",
])

# Executable for the sass_binary rule
nodejs_binary(
    name = "sass",
    entry_point = "sass_wrapper.js",
    data = [
        ":sass_wrapper.js",
        "@build_bazel_rules_sass_deps//sass",
        "@build_bazel_rules_sass_deps//@bazel/worker",
    ],
    # Opt-in to the patched require() function like google3 has
    # rules_nodejs 3.0 will flip the default for this flag which breaks rules_sass users
    templated_args = [
        "--nobazel_node_patches",
        "--bazel_patch_module_resolver",
    ],
)
