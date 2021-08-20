package(default_visibility = ["//visibility:public"])

load("@io_bazel_skydoc//skylark:skylark.bzl", "skylark_doc")

skylark_doc(
    name = "docs",
    srcs = [
        "//sass:sass.bzl",
        "//sass:sass_repositories.bzl",
        "//sass:npm_sass_library.bzl",
    ],
    format = "html",
    site_root = "/rules_sass",
    strip_prefix = "sass/",
)
