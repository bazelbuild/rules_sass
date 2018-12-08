package(default_visibility = ["//visibility:public"])

load("//sass:sass.bzl", "sass_library")

# make a :colors target that any sass_binary rules can depend on.
sass_library(
    name = "colors",
    srcs = ["_colors.scss"],
    visibility = ["//visibility:public"],
)

# make a :fonts target that any sass_binary rules can depend on.
sass_library(
    name = "fonts",
    srcs = ["_fonts.scss"],
    visibility = ["//visibility:public"],
)
