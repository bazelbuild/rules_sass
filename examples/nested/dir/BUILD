package(default_visibility = ["//visibility:public"])

load("//sass:sass.bzl", "sass_binary")

# Import our shared colors and fonts so we can generate a CSS file.
sass_binary(
    name = "dir",
    src = "main.scss",
    deps = [
        "//examples/shared:colors",
        "//examples/shared:fonts",
    ],
)
