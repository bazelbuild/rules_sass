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

load(
    "//sass:internal.bzl",
    _SassInfo = "SassInfo",
    _sass_library = "sass_library",
    _sass_binary_attrs = "sass_binary_attrs",
    _sass_binary_impl = "sass_binary_impl",
    _sass_binary_outputs = "sass_binary_outputs",
    _strip_extension = "strip_extension",
    _ALLOWED_SRC_FILE_EXTENSIONS = "ALLOWED_SRC_FILE_EXTENSIONS",
    _COMPILER_ATTR_DOC = "COMPILER_ATTR_DOC",
)

SassInfo = _SassInfo
sass_library = _sass_library

sass_binary = rule(
    implementation = _sass_binary_impl,
    attrs = _sass_binary_attrs,
    outputs = _sass_binary_outputs,
)

def _multi_sass_binary_impl(ctx):
  """multi_sass_binary accepts a list of sources and compile all in one pass.

  Args:
    ctx: The Bazel build context

  Returns:
    The multi_sass_binary rule.
  """

  inputs = ctx.files.srcs
  outputs = []
  # Every non-partial Sass file will produce one CSS output file and,
  # optionally, one sourcemap file.
  for f in inputs:
    # Sass partial files (prefixed with an underscore) do not produce any
    # outputs.
    if f.basename.startswith("_"):
      continue
    name = _strip_extension(f.basename)
    outputs.append(ctx.actions.declare_file(
      name + ".css",
      sibling = f,
    ))
    if ctx.attr.sourcemap:
      outputs.append(ctx.actions.declare_file(
        name + ".css.map",
        sibling = f,
      ))

  # Use the package directory as the compilation root given to the Sass compiler
  root_dir = ctx.label.package

  # Declare arguments passed through to the Sass compiler.
  # Start with flags and then expected program arguments.
  args = ctx.actions.args()
  args.add("--style", ctx.attr.output_style)
  args.add("--load-path", root_dir)

  if not ctx.attr.sourcemap:
    args.add("--no-source-map")

  args.add(root_dir + ":" + ctx.bin_dir.path + '/' + root_dir)
  args.use_param_file("@%s", use_always = True)
  args.set_param_file_format("multiline")

  if inputs:
    ctx.actions.run(
        inputs = inputs,
        outputs = outputs,
        executable = ctx.executable.compiler,
        arguments = [args],
        mnemonic = "SassCompiler",
        progress_message = "Compiling Sass",
    )

  return [DefaultInfo(files = depset(outputs))]

multi_sass_binary = rule(
  implementation = _multi_sass_binary_impl,
  attrs = {
    "srcs": attr.label_list(
      doc = "A list of Sass files and associated assets to compile",
      allow_files = _ALLOWED_SRC_FILE_EXTENSIONS,
      allow_empty = True,
      mandatory = True,
    ),
    "sourcemap": attr.bool(
      doc = "Whether sourcemaps should be emitted",
      default = True,
    ),
    "output_style": attr.string(
      doc = "How to style the compiled CSS",
      default = "compressed",
      values = [
        "expanded",
        "compressed",
      ],
    ),
    "compiler": attr.label(
      doc = _COMPILER_ATTR_DOC,
      default = Label("//sass"),
      executable = True,
      cfg = "host",
    ),
  }
)
