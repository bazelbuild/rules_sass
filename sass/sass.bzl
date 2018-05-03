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

SassInfo = provider("transitive_sources")

def collect_transitive_sources(srcs, deps):
    return depset(
        srcs,
        transitive=[dep[SassInfo].transitive_sources for dep in deps],
        order="postorder")

def _sass_library_impl(ctx):
    transitive_sources = collect_transitive_sources(
        ctx.files.srcs, ctx.attr.deps)
    return [SassInfo(transitive_sources=transitive_sources)]

def _sass_binary_impl(ctx):
    args = ctx.actions.args()
    args.add(["--style", ctx.attr.output_style], join_with="=")
    for prefix in [".", ctx.var['BINDIR'], ctx.var['GENDIR']]:
      args.add("--load-path=%s/" % prefix)
      for include_path in ctx.attr.include_paths:
        args.add("--load-path=%s/%s" % (prefix, include_path))
    args.add([ctx.file.src.path, ctx.outputs.css_file.path])

    ctx.actions.run(
        mnemonic = "SassCompiler",
        executable = ctx.executable._binary,
        inputs = [ctx.executable._binary] +
            list(collect_transitive_sources([ctx.file.src], ctx.attr.deps)),
        arguments = [args],
        outputs = [ctx.outputs.css_file, ctx.outputs.map_file],
    )

    # Make sure the output CSS is available in runfiles if used as a data dep.
    return DefaultInfo(
        runfiles = ctx.runfiles(
            files = [
                ctx.outputs.css_file,
                ctx.outputs.map_file,
            ]))

def _sass_binary_outputs(output_name, output_dir):
  output_name = output_name or "%{name}.css"
  css_file = ("%s/%s" % (output_dir, output_name) if output_dir
              else output_name)
  outputs = {
      "css_file": css_file,
      "map_file": "%s.map" % css_file,
  }

  return outputs

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
        doc = "Sass source files",
        allow_files = _FILETYPES,
        mandatory = True,
        single_file = True,
    ),
    "include_paths": attr.string_list(
        doc = "Additional directories to search when resolving imports"),
    "output_dir": attr.string(
        doc = "Output directory, relative to xxx",
        default = ""),
    "output_name": attr.string(
        doc = "Name of the output file, if not specified it is xxx",
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
    "_binary": attr.label(
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
