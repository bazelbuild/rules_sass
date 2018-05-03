workspace(name = "build_bazel_rules_sass")

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
