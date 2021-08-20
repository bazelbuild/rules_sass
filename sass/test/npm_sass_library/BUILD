load("//:defs.bzl", "npm_sass_library", "sass_binary")
load(":test_assertions.bzl", "npm_sass_library_test")

npm_sass_library(
    name = "test_pkg_sass_deps",
    deps = ["//sass/test/npm_sass_library/test_external_pkg"],
)

sass_binary(
    name = "test_fixture",
    src = "test_fixture.scss",
    deps = [":test_pkg_sass_deps"],
)

npm_sass_library_test()
