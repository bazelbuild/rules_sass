workspace(name = "io_bazel_rules_sass")

load("//:package.bzl", "rules_sass_dependencies", "rules_sass_dev_dependencies")

rules_sass_dependencies();
rules_sass_dev_dependencies();

load("@build_bazel_rules_nodejs//:defs.bzl", "node_repositories")
node_repositories()

load("//:defs.bzl", "sass_repositories")
sass_repositories()

#############################################
# Required dependencies for docs generation
#############################################

load("@io_bazel_rules_go//go:def.bzl", "go_rules_dependencies", "go_register_toolchains")
go_rules_dependencies()
go_register_toolchains()

load("@io_bazel_skydoc//skylark:skylark.bzl", "skydoc_repositories")
skydoc_repositories()
