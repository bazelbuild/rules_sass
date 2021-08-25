load("@bazel_tools//tools/build_rules:test_rules.bzl", "file_test", "rule_test")

def npm_sass_library_test():
    rule_test(
        name = "file_generation_expectations",
        generates = ["test_fixture.css", "test_fixture.css.map"],
        rule = "//sass/test/npm_sass_library:test_fixture",
    )

    file_test(
        name = "test_fixture_red_color_expectation",
        file = "//sass/test/npm_sass_library:test_fixture.css",
        regexp = "color:red",
        matches = 1,
    )

    file_test(
        name = "test_fixture_blue_button_expectation",
        file = "//sass/test/npm_sass_library:test_fixture.css",
        regexp = "background:blue",
        matches = 1,
    )

    file_test(
        name = "test_fixture_blue_input_expectation",
        file = "//sass/test/npm_sass_library:test_fixture.css",
        regexp = "color:blue",
        matches = 1,
    )
