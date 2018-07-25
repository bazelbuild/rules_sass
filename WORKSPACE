workspace(name = "io_bazel_rules_sass")
load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

# To use the JavaScript version of Sass, we need to first install nodejs
http_archive(
    name = "build_bazel_rules_nodejs",
    url = "https://github.com/bazelbuild/rules_nodejs/archive/0.10.0.zip",
    strip_prefix = "rules_nodejs-0.10.0",
    sha256 = "2f77623311da8b5009b1c7eade12de8e15fa3cd2adf9dfcc9f87cb2082b2211f",
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


#################################################
# Dependencies for bazel formatting and linting #
#################################################

BAZEL_BUILDTOOLS_VERSION = "c39a197f7d35aebb0e0b031d728fb918f73887d6"

# Bazel buildtools repo contains tools for BUILD file formatting ("buildifier") etc.
http_archive(
    name = "com_github_bazelbuild_buildtools",
    url = "https://github.com/bazelbuild/buildtools/archive/%s.zip" % BAZEL_BUILDTOOLS_VERSION,
    strip_prefix = "buildtools-%s" % BAZEL_BUILDTOOLS_VERSION,
    sha256 = "30c8e027d0ed7843651fbe2dbb6338171c963ba1184a7bce802ae4a30e223fd4",
)


IO_BAZEL_VERSION = "968f87900dce45a7af749a965b72dbac51b176b3"

# Fetching the Bazel source code allows us to compile/use the Skylark linter
http_archive(
    name = "io_bazel",
    url = "https://github.com/bazelbuild/bazel/archive/%s.zip" % IO_BAZEL_VERSION,
    strip_prefix = "bazel-%s" % IO_BAZEL_VERSION,
    sha256 = "e373d2ae24955c1254c495c9c421c009d88966565c35e4e8444c082cb1f0f48f",
)


# Parts of the build toolchain are written in Go, such as buildifier.
# Bazel doesn't support transitive WORKSPACE deps, so we must repeat them here.
http_archive(
    name = "io_bazel_rules_go",
    url = "https://github.com/bazelbuild/rules_go/releases/download/0.10.3/rules_go-0.10.3.tar.gz",
    sha256 = "feba3278c13cde8d67e341a837f69a029f698d7a27ddbb2a202be7a10b22142a",
)

load("@io_bazel_rules_go//go:def.bzl", "go_rules_dependencies", "go_register_toolchains")
go_rules_dependencies()
go_register_toolchains()
