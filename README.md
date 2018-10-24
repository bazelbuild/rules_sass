[![Build status](https://badge.buildkite.com/accb37a80d88e0ffda97f55451d05eea2004ed8bbb80a27958.svg)](https://buildkite.com/bazel/rules-sass-postsubmit)

# Sass Rules for Bazel

## Rules
* [sass_binary]()
* [sass_library]()

## Overview
These build rules are used for building [Sass][sass] projects with Bazel.

[sass]: http://www.sass-lang.com

## Setup
To use the Sass rules, add the following to your
`WORKSPACE` file to add the external repositories for Sass, making sure to use the latest
published versions:

```py
# Fetch required transitive dependencies
load("@io_bazel_rules_sass//:package.bzl", "rules_sass_dependencies")
rules_sass_dependencies()

# Fetch required transitive dependencies
load("@io_bazel_rules_sass//:defs.bzl", "sass_repositories")
sass_repositories()
```

## Basic Example

Suppose you have the following directory structure for a simple Sass project:

```
[workspace]/
    WORKSPACE
    hello_world/
        BUILD
        main.scss
    shared/
        BUILD
        _fonts.scss
        _colors.scss
```

`shared/_fonts.scss`

```scss
$default-font-stack: Cambria, "Hoefler Text", serif;
$modern-font-stack: Constantia, "Lucida Bright", serif;
```

`shared/_colors.scss`

```scss
$example-blue: #0000ff;
$example-red: #ff0000;
```

`shared/BUILD`

```python
package(default_visibility = ["//visibility:public"])

load("@io_bazel_rules_sass//:defs.bzl", "sass_library")

sass_library(
    name = "colors",
    srcs = ["_colors.scss"],
)

sass_library(
    name = "fonts",
    srcs = ["_fonts.scss"],
)
```

`hello_world/main.scss`:

```scss
@import "shared/fonts";
@import "shared/colors";

html {
  body {
    font-family: $default-font-stack;
    h1 {
      font-family: $modern-font-stack;
      color: $example-red;
    }
  }
}
```

`hello_world/BUILD:`

```py
package(default_visibility = ["//visibility:public"])

load("@io_bazel_rules_sass//:defs.bzl", "sass_binary")

sass_binary(
    name = "hello_world",
    src = "main.scss",
    deps = [
         "//shared:colors",
         "//shared:fonts",
    ],
)
```

Build the binary:

```
$ bazel build //hello_world
INFO: Found 1 target...
Target //hello_world:hello_world up-to-date:
  bazel-bin/hello_world/hello_world.css
  bazel-bin/hello_world/hello_world.css.map
INFO: Elapsed time: 1.911s, Critical Path: 0.01s
```

## Build Rule Reference

<a name="reference-sass_binary"></a>
### sass_binary

```py
sass_binary(name, src, deps=[], output_style="compressed", include_paths=[], output_dir=".", output_name=<src_filename.css>)
```

`sass_binary` compiles a single CSS output from a single Sass entry-point file. The entry-point file
may have dependencies (`sass_library` rules, see below).


#### Implicit output targets
| Label            | Description                                                                  |
|------------------|------------------------------------------------------------------------------|
| **name**.css     | The generated CSS output                                                     |
| **name**.css.map | The [source map][] that can be used to debug the Sass source in-browser      |

[source map]: https://developers.google.com/web/tools/chrome-devtools/javascript/source-maps


| Attribute       | Description                                                                   |
|-----------------|-------------------------------------------------------------------------------|
| `name`          | Unique name for this rule (required)                                          |
| `src`           | Sass compilation entry-point (required).                                      |
| `deps`          | List of dependencies for the `src`. Each dependency is a `sass_library`       |
| `include_paths` | Additional directories to search when resolving imports                       |
| `output_dir`    | Output directory, relative to this package                                    |
| `output_name`   | Output file name, including .css extension. Defaults to `<src_name>.css`      |
| `output_style`  | [Output style][] for the generated CSS.                                       |
| `sourcemap`     | Whether to generate sourcemaps for the generated CSS. Defaults to True.       |

[Output style]: http://sass-lang.com/documentation/file.SASS_REFERENCE.html#output_style

### sass_library

```py
sass_library(name, src, deps=[])
```

Defines a collection of Sass files that can be depended on by a `sass_binary`. Does not generate
any outputs.

| Attribute | Description                                                                         |
|-----------|-------------------------------------------------------------------------------------|
| `name`    | Unique name for this rule (required)                                                |
| `srcs`    | Sass files included in this library. Each file should start with an underscore      |
| `deps`    | Dependencies for the `src`. Each dependency is a `sass_library`                     |
