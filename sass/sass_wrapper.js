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

const args = process.argv.slice(2);
if (runAsWorker(args)) {
  debug('Starting Sass compiler persistent worker...');
  runWorkerLoop(args =>
    sass.cli_pkg_main_0_(args)
        .then(() => true)
        .catch(() => false)
  );
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
