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
"use strict";

const {debug, runAsWorker, runWorkerLoop} = require('@bazel/worker');
const sass = require('sass');
const fs = require('fs');
const minimist = require('minimist');

const args = process.argv.slice(2);
if (runAsWorker(args)) {
  debug('Starting Sass compiler persistent worker...');
  runWorkerLoop(args => {
    const argv = minimist(args);
    const positionalArgs = argv['_'];
    const output = positionalArgs[1];
    const input = positionalArgs[0];

    const sourceMap =
      typeof argv['source-map'] === 'boolean' ? argv['source-map'] : true;
    const embedSources =
      typeof argv['embed-sources'] === 'boolean'
        ? argv['embed-sources']
        : false;

    try {
      const result = sass.renderSync({
        file: input,
        outFile: output,
        includePaths: argv['load-path'],
        outputStyle: argv['style'],
        sourceMap,
        sourceMapContents: embedSources
      });

      fs.writeFileSync(output, result.css);
      if (sourceMap) {
        fs.writeFileSync(output + '.map', result.map);
      }
      return true;
    } catch (e) {
      console.error(e.message);
      return false;
    }
  });
  // Note: intentionally don't process.exit() here, because runWorkerLoop
  // is waiting for async callbacks from node.
} else {
  debug('Running a single build...');
  if (args.length === 0) throw new Error('Not enough arguments');
  if (args.length !== 1) {
    throw new Error('Expected one argument: path to flagfile');
  }

  // Bazel worker protocol expects the only arg to be @<path_to_flagfile>.
  // When we are running a single build, we remove the @ prefix and read the list 
  // of actual arguments line by line.
  const configFile = args[0].replace(/^@+/, '');
  const configContent = fs.readFileSync(configFile, 'utf8').trim();
  sass.cli_pkg_main_0_(configContent.split('\n'));
}

process.exitCode = 0;
