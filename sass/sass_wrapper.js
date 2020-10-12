/**
 * @license
 * Copyright Google Inc. All Rights Reserved.
 * Use of this source code is governed by the Apache 2.0 license that can be found in the LICENSE.txt file.
 *
 * A Sass compiler wrapper that supports bazel persistent worker protocol.
 *
 * Bazel can spawn a persistent worker process that handles multiple invocations.
 * It can also be invoked with an argument file to run once and exit.
 */
'use strict';

const worker = require('@bazel/worker');
const sass = require('sass');
const yargs = require('yargs/yargs');
const fs = require('fs');
const path = require('path');

if (require.main === module) {
  // Bazel will pass a special argument to the program when it's running us as a worker
  if (worker.runAsWorker(process.argv)) {
    worker.debug('Running as a Bazel worker');

    worker.runWorkerLoop(main);
  } else {
    // Running standalone so stdout is available as usual
    console.debug('Running as a standalone process');

    // The first argument to the program is prefixed with '@'
    // because Bazel does that for param files. Strip it first.
    const paramFile = process.argv[2].replace(/^@/, '');
    const args = fs.readFileSync(paramFile, 'utf-8').trim().split('\n');

    // Bazel is just running the program as a single action, don't act like a worker
    if (!main(args)) {
      process.exitCode = 1;
    }
  }
}

/**
 * Main function that passes the arguments from the worker or standalone to the
 * dart sass compiler
 * @param {string[]} argv The parsed command line args
 * @returns {boolean} Returns true if the compilation was successful.
 */
function main(argv) {
  // IMPORTANT don't log with console.out - stdout is reserved for the worker protocol.
  // This is true for any code running in the program, even if it comes from a third-party library.
  const { files, style, loadPaths, noSourceMap, embedSources } = yargs(
    argv
  ).options({
    files: { array: true, default: [] },
    style: { string: true, default: 'compressed' },
    noSourceMap: { boolean: true, default: false },
    'embed-sources': { boolean: true },
    'load-paths': { array: true, default: [] },
  }).argv;

  for (let i = 0, max = files.length; i < max; i++) {
    const [input, outFile] = files[i].split(':');
    compileDartSass({
      style,
      outFile,
      input,
      embedSources,
      sourceMap: !noSourceMap,
      loadPaths,
    });
  }

  return true;
}

/**
 * Function that uses the dart sass nodeJS API to compile sass to css
 * https://sass-lang.com/documentation/js-api
 * @param {object} config Configuration that should be passed to the render function
 * @param {string} config.input The sass file that should be compiled
 * @param {string} config.outFile The css file that should be written
 * @param {boolean} config.embedSources If the source Maps should be embedded
 * @param {boolean} config.sourceMap If source maps should be written (default `true`)
 * @param {string[]} config.loadPaths This array of strings option provides load paths for Sass to look for imports.
 * @param {'compressed' | 'expanded'} config.style Output style of the resulting css
 * @returns {void}
 */
function compileDartSass(config) {
  // IMPORTANT don't log with console.out - stdout is reserved for the worker protocol.
  // This is true for any code running in the program, even if it comes from a third-party library.

  // use renderSync() as it is almost twice as fast as render() according to the
  // official documentation.
  const result = sass.renderSync({
    style: config.style,
    // This option defines one or more additional handlers for loading files
    // when a @use rule or an @import rule is encountered. This should handle the
    // common WebPack import style from node_modules starting with ~
    importer: function (url) {
      if (url.startsWith('~')) {
        const resolvedFile = resolveScssFile(
          path.resolve(url.replace('~', '../../external/npm/node_modules/'))
        );

        if (resolvedFile) {
          return {
            file: resolvedFile,
          };
        }
      }
      // null, which indicates that it doesn't recognize the URL and another
      // importer should be tried instead.
      return null;
    },
    sourceMap: config.sourceMap,
    sourceMapContents: config.embedSources,
    sourceMapEmbed: config.embedSources,
    includePaths: [...config.loadPaths, '../../external/npm/node_modules/'],
    file: config.input,
    outFile: config.outFile,
  });

  fs.writeFileSync(config.outFile, result.css);

  if (config.sourceMap) {
    fs.writeFileSync(`${config.outFile}.map`, result.map);
  }
}

/**
 * Function to resolve a path for a supported import style file ending.
 * https://sass-lang.com/documentation/at-rules/import#finding-the-file
 * @param {string} importPath The import path where it should try to resolve the
 * style file that can be imported or referenced via use
 * @returns {string|undefined} Returns the resolved path if the file exists or undefined
 * if it cannot resolve the path.
 */
function resolveScssFile(importPath) {
  const fileName = path.basename(importPath);
  const variants = [
    `${importPath}.scss`,
    `${importPath}.sass`,
    `${importPath}.css`,
    `${importPath}/${fileName}.scss`,
    `${importPath}/${fileName}.sass`,
    `${importPath}/_${fileName}.scss`,
    `${importPath}/_${fileName}.sass`,
    `${importPath}/_index.scss`,
    `${importPath}/_index.sass`,
  ];

  for (let i = 0, max = variants.length; i < max; i++) {
    // return the variant if the file exist
    if (fs.existsSync(variants[i]) && fs.lstatSync(variants[i]).isFile()) {
      return variants[i];
    }
  }
}
