package(default_visibility = ["//visibility:public"])

load("//:defs.bzl", "sass_binary")

# Import our shared colors and fonts so we can generate a CSS file.
sass_binary(
    name = "hello_world",
    src = "main.scss",
    deps = [
        "//examples/shared",
    ],
)

sass_binary(
    name = "hello_world_no_sourcemap",
    src = "main.scss",
    output_name = "main-no-sourcemap.css",
    sourcemap = False,
    deps = [
        "//examples/shared:colors",
        "//examples/shared:fonts",
    ],
)
