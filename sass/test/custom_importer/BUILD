load("//:defs.bzl", "sass_binary")
load(":test_assertions.bzl", "custom_importer_test")

sass_binary(
    name = "theme",
    testonly = True,
    src = "theme.scss",
    config_file = "sass_config.mjs",
    deps = ["//sass/test/custom_importer/cdk"],
)

custom_importer_test()
