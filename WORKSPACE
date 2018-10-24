workspace(name = "io_bazel_rules_sass")

load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

# To use the JavaScript version of Sass, we need to first install nodejs
http_archive(
    name = "build_bazel_rules_nodejs",
    sha256 = "2f77623311da8b5009b1c7eade12de8e15fa3cd2adf9dfcc9f87cb2082b2211f",
    strip_prefix = "rules_nodejs-0.10.0",
    url = "https://github.com/bazelbuild/rules_nodejs/archive/0.10.0.zip",
)

load("@build_bazel_rules_nodejs//:defs.bzl", "node_repositories")

node_repositories(package_json = [])

load("//sass:sass_repositories.bzl", "sass_repositories")

sass_repositories()

#################################################
# Dependencies for generating documentation     #
#################################################

http_archive(
    name = "bazel_skylib",
    strip_prefix = "bazel-skylib-0.3.1",
    url = "https://github.com/bazelbuild/bazel-skylib/archive/0.3.1.zip",
)

http_archive(
    name = "io_bazel_skydoc",
    strip_prefix = "skydoc-b91215b629c67527a98466e61c11c77530fe85c1",
    url = "https://github.com/bazelbuild/skydoc/archive/b91215b629c67527a98466e61c11c77530fe85c1.zip",
)

load("@io_bazel_skydoc//skylark:skylark.bzl", "skydoc_repositories")

skydoc_repositories()

#################################################
# Dependencies for bazel formatting and linting #
#################################################

BAZEL_BUILDTOOLS_VERSION = "c39a197f7d35aebb0e0b031d728fb918f73887d6"

# Bazel buildtools repo contains tools for BUILD file formatting ("buildifier") etc.
http_archive(
    name = "com_github_bazelbuild_buildtools",
    sha256 = "30c8e027d0ed7843651fbe2dbb6338171c963ba1184a7bce802ae4a30e223fd4",
    strip_prefix = "buildtools-%s" % BAZEL_BUILDTOOLS_VERSION,
    url = "https://github.com/bazelbuild/buildtools/archive/%s.zip" % BAZEL_BUILDTOOLS_VERSION,
)

IO_BAZEL_VERSION = "968f87900dce45a7af749a965b72dbac51b176b3"

# Fetching the Bazel source code allows us to compile/use the Skylark linter
http_archive(
    name = "io_bazel",
    sha256 = "e373d2ae24955c1254c495c9c421c009d88966565c35e4e8444c082cb1f0f48f",
    strip_prefix = "bazel-%s" % IO_BAZEL_VERSION,
    url = "https://github.com/bazelbuild/bazel/archive/%s.zip" % IO_BAZEL_VERSION,
)

# Parts of the build toolchain are written in Go, such as buildifier.
# Bazel doesn't support transitive WORKSPACE deps, so we must repeat them here.
http_archive(
    name = "io_bazel_rules_go",
    sha256 = "ee5fe78fe417c685ecb77a0a725dc9f6040ae5beb44a0ba4ddb55453aad23a8a",
    url = "https://github.com/bazelbuild/rules_go/releases/download/0.16.0/rules_go-0.16.0.tar.gz",
)

load("@io_bazel_rules_go//go:def.bzl", "go_rules_dependencies", "go_register_toolchains")

go_rules_dependencies()

go_register_toolchains()
