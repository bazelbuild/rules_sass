load("@bazel_tools//tools/build_rules:test_rules.bzl", "file_test", "rule_test")

def custom_importer_test():
    rule_test(
        name = "file_generation_expectations",
        generates = ["theme.css", "theme.css.map"],
        rule = "//sass/test/custom_importer:theme",
    )

    file_test(
        name = "test_fixture_red_color_expectation",
        file = "//sass/test/custom_importer:theme.css",
        regexp = "color:red",
        matches = 1,
    )
