# Copyright 2018 The Bazel Authors. All rights reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#    http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
"Compile Sass files to CSS"

_FILETYPES = [".sass", ".scss", ".svg", ".png", ".gif"]

# Documentation for switching which compiler is used
_COMPILER_ATTR_DOC = """Choose which Sass compiler binary to use.

By default, we use the JavaScript-transpiled version of the
dart-sass library, based on https://github.com/sass/dart-sass.
This is the canonical compiler under active development by the Sass team.
This compiler is convenient for frontend developers since it's released
as JavaScript and can run natively in NodeJS without being locally built.
However, it is the slowest option. If you have a substantial Sass
codebase, consider the two other options for this attribute:

1. `compiler = "@sassc"` uses the libsass compiler written in C++.
   As of 2018, this is the most commonly used compiler.
   It requires your Bazel setup has a working C++ compilation toolchain.

   NOTE: future releases of rules_sass may remove this option, if it
   becomes obsolete or if the maintenance burden is too high.

2. Dart Sass runs the new compiler implementation natively in the Dart
   VM. This option requires a change to the Dart Bazel rules which is
   not yet available as of May 2018.
"""

SassInfo = provider(
    doc = "Collects files from sass_library for use in downstream sass_binary",
    fields = {
        "transitive_sources": "Sass sources for this target and its dependencies",
    })

def _collect_transitive_sources(srcs, deps):
  "Sass compilation requires all transitive .sass source files"
  return depset(
      srcs,
      transitive=[dep[SassInfo].transitive_sources for dep in deps],
      # Provide .sass sources from dependencies first
      order="postorder")

def _sass_library_impl(ctx):
  """sass_library collects all transitive sources for given srcs and deps.
  It doesn't execute any actions."""
  transitive_sources = _collect_transitive_sources(
      ctx.files.srcs, ctx.attr.deps)
  return [SassInfo(transitive_sources=transitive_sources)]

def _run_sass(ctx, input, css_output, map_output):
  """run_sass performs an action to compile a single Sass file into CSS."""
  # The Sass CLI expects inputs like
  # sass <flags> <input_filename> <output_filename>
  args = ctx.actions.args()

  # Flags (see https://github.com/sass/node-sass#options)
  # Note, the command line is compatible with Dart sass.
  args.add(["--style", ctx.attr.output_style], join_with="=")
  # FIXME: sassc requires --sourcemap to produce the .map file, otherwise this rule fails
  # However dart sass only accepts --source-map (Extra hyphen)
  # https://github.com/sass/dart-sass/commit/234aa12e081c7cc873549fee1c5f37a12564d1b1#r28862590
  # so using compiler="@sassc" is currently broken.
  #args.add(["--sourcemap"])

  # Sources for compilation may exist in the source tree, in bazel-bin, or bazel-genfiles.
  for prefix in [".", ctx.var['BINDIR'], ctx.var['GENDIR']]:
    args.add("--load-path=%s/" % prefix)
    for include_path in ctx.attr.include_paths:
      args.add("--load-path=%s/%s" % (prefix, include_path))

  # Last arguments are input and output paths
  # Note that the sourcemap is implicitly written to a path the same as the
  # css with the added .map extension.
  args.add([input.path, css_output.path])

  ctx.actions.run(
      mnemonic = "SassCompiler",
      executable = ctx.executable.compiler,
      inputs = [ctx.executable.compiler] +
          list(_collect_transitive_sources([input], ctx.attr.deps)),
      arguments = [args],
      outputs = [css_output, map_output],
  )

def _sass_binary_impl(ctx):
  _run_sass(ctx, ctx.file.src, ctx.outputs.css_file, ctx.outputs.map_file)

  # Make sure the output CSS is available in runfiles if used as a data dep.
  return DefaultInfo(runfiles = ctx.runfiles(files = [
      ctx.outputs.css_file,
      ctx.outputs.map_file,
  ]))

def _sass_binary_outputs(src, output_name, output_dir):
  """Get map of sass_binary outputs, which includes the generated css file
  and its sourcemap.
  Note that the arguments to this function are named after attributes on the rule."""
  output_name = output_name or "%{src}.css"
  css_file = "/".join([p for p in [output_dir, output_name] if p])
  return {
      "css_file": css_file,
      "map_file": "%s.map" % css_file,
  }

sass_deps_attr = attr.label_list(
    doc = "sass_library targets to include in the compilation",
    providers = [SassInfo],
    allow_files = False,
)

sass_library = rule(
    implementation = _sass_library_impl,
    attrs = {
        "srcs": attr.label_list(
            doc = "Sass source files",
            allow_files = _FILETYPES,
            non_empty = True,
            mandatory = True,
        ),
        "deps": sass_deps_attr,
    },
)
"""Defines a group of Sass include files.
"""

_sass_binary_attrs = {
    "src": attr.label(
        doc = "Sass entrypoint file",
        allow_files = _FILETYPES,
        mandatory = True,
        single_file = True,
    ),
    "include_paths": attr.string_list(
        doc = "Additional directories to search when resolving imports"),
    "output_dir": attr.string(
        doc = "Output directory, relative to this package.",
        default = ""),
    "output_name": attr.string(
        doc = """Name of the output file, including the .css extension.

By default, this is based on the `src` attribute: if `styles.scss` is
the `src` then the output file is `styles.css.`.
You can override this to be any other name.
Note that some tooling may assume that the output name is derived from
the input name, so use this attribute with caution.""",
        default = ""),
    "output_style": attr.string(
        doc = "How to style the compiled CSS",
        default = "compressed",
        values = [
            "expanded",
            "compressed",
        ],
    ),
    "deps": sass_deps_attr,
    "compiler": attr.label(
        doc = _COMPILER_ATTR_DOC,
        default = Label("//sass"),
        executable = True,
        cfg = "host",
    ),
}

sass_binary = rule(
    implementation = _sass_binary_impl,
    attrs = _sass_binary_attrs,
    outputs = _sass_binary_outputs,
)
