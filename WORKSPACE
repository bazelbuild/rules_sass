workspace(name = "io_bazel_rules_sass")

load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

# Dependency for running Skylint.
http_archive(
    name = "io_bazel",
    sha256 = "978f7e0440dd82182563877e2e0b7c013b26b3368888b57837e9a0ae206fd396",
    strip_prefix = "bazel-0.18.0",
    url = "https://github.com/bazelbuild/bazel/archive/0.18.0.zip",
)

# Required for the Buildtool repository.
http_archive(
    name = "io_bazel_rules_go",
    sha256 = "7be7dc01f1e0afdba6c8eb2b43d2fa01c743be1b9273ab1eaf6c233df078d705",
    url = "https://github.com/bazelbuild/rules_go/releases/download/0.16.5/rules_go-0.16.5.tar.gz",
)

# Bazel buildtools repo contains tools for BUILD file formatting ("buildifier") etc.
http_archive(
    name = "com_github_bazelbuild_buildtools",
    sha256 = "a82d4b353942b10c1535528b02bff261d020827c9c57e112569eddcb1c93d7f6",
    strip_prefix = "buildtools-0.17.2",
    url = "https://github.com/bazelbuild/buildtools/archive/0.17.2.zip",
)

# Needed in order to generate documentation
http_archive(
    name = "io_bazel_skydoc",
    sha256 = "75fd965a71ca1f0d0406d0d0fb0964d24090146a853f58b432761a1a6c6b47b9",
    strip_prefix = "skydoc-82fdbfe797c6591d8732df0c0389a2b1c3e50992",
    url = "https://github.com/bazelbuild/skydoc/archive/82fdbfe797c6591d8732df0c0389a2b1c3e50992.zip",  # 2018-12-12
)

http_archive(
    name = "build_bazel_rules_nodejs",
    sha256 = "c077680a307eb88f3e62b0b662c2e9c6315319385bc8c637a861ffdbed8ca247",
    urls = ["https://github.com/bazelbuild/rules_nodejs/releases/download/5.1.0/rules_nodejs-5.1.0.tar.gz"],
)

load("@build_bazel_rules_nodejs//:repositories.bzl", "build_bazel_rules_nodejs_dependencies")

build_bazel_rules_nodejs_dependencies()

load("@rules_nodejs//nodejs:repositories.bzl", "nodejs_register_toolchains")
load("@rules_nodejs//nodejs:yarn_repositories.bzl", "yarn_repositories")

nodejs_register_toolchains(
    name = "nodejs",
    node_version = "16.13.2",
)

yarn_repositories(
    name = "yarn",
    yarn_version = "1.22.17",
)

load("//:defs.bzl", "sass_repositories")

sass_repositories()

#############################################
# Required dependencies for docs generation
#############################################

load("@io_bazel_rules_go//go:def.bzl", "go_register_toolchains", "go_rules_dependencies")

go_rules_dependencies()

go_register_toolchains()

load("@io_bazel_skydoc//skylark:skylark.bzl", "skydoc_repositories")

skydoc_repositories()
