/**
 * @license
 * Copyright Google LLC. All Rights Reserved.
 * Use of this source code is governed by the Apache 2.0 license that can be found in the LICENSE.txt file.
 */

import fs from 'fs';
import url from 'url';
import {importNpmModuleFromSassDeps} from './sass_deps_helper.mjs';

// Import these external NPM modules from our Sass Bazel-managed deps.
const sass = await importNpmModuleFromSassDeps('sass/sass.default.dart.js');
const yargs = await importNpmModuleFromSassDeps('yargs/index.mjs');

/**
 * Performs a Sass build. Expects raw command line arguments as
 * constructed by the Sass Bazel rule.
 */
export async function invokeBuild(args) {
  await yargs(args)
    // Ensures that Yargs does not print anything to stdout which is used
    // for communication in persistent workers. This would cause errors.
    .help(false)
    .showHelpOnFail(false)
    .fail(false)
    .strict()
    .demandCommand()
    .command(
      '* <inputExecpath> <outputExecpath>',
      'Performs a Sass build',
      args =>
        args
          .positional('inputExecpath', {
            type: 'string',
            description: 'Execpath for the input Sass file.',
            demandOption: true,
          })
          .positional('outputExecpath', {
            type: 'string',
            demandOption: true,
            description: 'Execpath where the output CSS should be written to.',
          })
          .option('configFile', {
            type: 'string',
            description:
              'Optional execpath pointing to a configuration JavaScript file.',
          })
          .option('sourceMap', {
            type: 'boolean',
            description: 'Whether a sourcemap file should be generated.',
            default: true,
          })
          .option('embedSources', {
            type: 'boolean',
            description:
              'Whether source file contents should be embedded directly.',
          })
          .option('loadPath', {
            type: 'array',
            description:
              'Load paths (as execpaths) to be added to the compilation.',
          })
          .option('style', {
            type: 'string',
            description:
              'Output Sass style to emit. i.e. compressed or expanded.',
          }),
      args => performSyncSassBuild(args)
    )
    .parseAsync();
}

/**
 * Performs a synchronous Sass build for the given options. The output
 * files are written asynchronously to the file system.
 */
async function performSyncSassBuild(args) {
  const {
    sourceMap,
    embedSources,
    style,
    loadPath,
    inputExecpath,
    outputExecpath,
    configFile,
  } = args;

  const userConfigOptions =
    configFile !== undefined
      ? (await import(url.pathToFileURL(configFile))).default
      : {};

  const result = sass.compile(inputExecpath, {
    ...userConfigOptions,
    style: style,
    loadPaths: loadPath,
    sourceMap: sourceMap,
    sourceMapIncludeSources: embedSources,
  });

  const writeTasks = [fs.promises.writeFile(outputExecpath, result.css)];

  if (sourceMap === true) {
    // The source map will be written next to the CSS output file.
    const mapOutPath = `${outputExecpath}.map`;
    writeTasks.push(
      fs.promises.writeFile(mapOutPath, JSON.stringify(result.sourceMap))
    );
  }

  // Wait for both CSS and source-map (if enabled) to be written.
  await Promise.all(writeTasks);
}
