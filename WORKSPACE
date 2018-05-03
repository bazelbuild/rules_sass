workspace(name = "io_bazel_rules_sass")

# To use the JavaScript version of Sass, we need to first install nodejs
http_archive(
    name = "build_bazel_rules_nodejs",
    url = "https://github.com/bazelbuild/rules_nodejs/archive/0.8.0.zip",
    strip_prefix = "rules_nodejs-0.8.0",
    sha256 = "4e40dd49ae7668d245c3107645f2a138660fcfd975b9310b91eda13f0c973953",
)
load("@build_bazel_rules_nodejs//:defs.bzl", "node_repositories")
node_repositories(package_json = [])

load("//sass:sass_repositories.bzl", "sass_repositories")
sass_repositories()

#############################################
# Dependencies for generating documentation #
#############################################

http_archive(
    name = "bazel_skylib",
    url = "https://github.com/bazelbuild/bazel-skylib/archive/0.3.1.zip",
    strip_prefix = "bazel-skylib-0.3.1",
)

http_archive(
    name = "io_bazel_skydoc",
    url = "https://github.com/bazelbuild/skydoc/archive/0ef7695c9d70084946a3e99b89ad5a99ede79580.zip",
    strip_prefix = "skydoc-0ef7695c9d70084946a3e99b89ad5a99ede79580",
)
load("@io_bazel_skydoc//skylark:skylark.bzl", "skydoc_repositories")
skydoc_repositories()
