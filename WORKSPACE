workspace(name = "io_bazel_rules_sass")

load("//:package.bzl", "rules_sass_dependencies", "rules_sass_dev_dependencies")

rules_sass_dependencies()

rules_sass_dev_dependencies()

load("//:defs.bzl", "sass_repositories")

sass_repositories()

#############################################
# Required dependencies for docs generation
#############################################

load("@io_bazel_rules_go//go:deps.bzl", "go_rules_dependencies", "go_register_toolchains")

go_rules_dependencies()

go_register_toolchains()

load("@bazel_gazelle//:deps.bzl", "gazelle_dependencies")

gazelle_dependencies()

load("@bazel_gazelle//:deps.bzl", "go_repository")

go_repository(
    name = "com_github_google_go_cmp",
    commit = "2248b49eaa8e1c8c0963ee77b40841adbc19d4ca",
    importpath = "github.com/google/go-cmp",
)

load("@io_bazel_skydoc//skylark:skylark.bzl", "skydoc_repositories")

skydoc_repositories()
