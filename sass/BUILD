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
    data = [
        "sass_bazel_main.js",
        "sass_build.mjs",
        "sass_deps_helper.mjs",
        "@build_bazel_rules_sass_deps//@bazel/runfiles",
        "@build_bazel_rules_sass_deps//@bazel/worker",
        "@build_bazel_rules_sass_deps//sass",
        "@build_bazel_rules_sass_deps//yargs",
    ],
    entry_point = "sass_bazel_main.js",
    # Disables the Bazel node modules linker. The node module linker is unreliable for the
    # persistent worker executable, as it would rely on the `node_modules/` folder in the
    # execroot that can be shared in non-sandbox environments or for persistent workers.
    # https://docs.bazel.build/versions/main/command-line-reference.html#flag--worker_sandboxing.
    templated_args = ["---nobazel_run_linker"],
)
