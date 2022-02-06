/**
 * @license
 * Copyright Google LLC. All Rights Reserved.
 * Use of this source code is governed by the Apache 2.0 license that can be found in the LICENSE.txt file.
 *
 * A Sass compiler wrapper that supports bazel persistent worker protocol.
 *
 * Bazel can spawn a persistent worker process that handles multiple invocations.
 * It can also be invoked with an argument file to run once and exit.
 */
'use strict';

const {debug, runAsWorker, runWorkerLoop} = require('@bazel/worker');
const fs = require('fs');
/**
 * Entry-point for the Sass build action.
 *
 * This is a small wrapper function dealing with the Bazel persistent worker
 * command line unwrapping. The actual build is performed in `invokeBuild`.
 */
async function main(args) {
  const {invokeBuild} = await import('./sass_build.mjs');

  if (runAsWorker(args)) {
    debug('Starting Sass compiler persistent worker...');
    await runWorkerLoop(args =>
      // The worker loop expects a graceful promise completion with a
      // boolean indicating success or failure.
      invokeBuild(args)
        .then(() => true)
        .catch(error => {
          // Note: Error should be printed to `stderr` to not break the Bazel
          // persistent worker communication over `stdout`.
          console.error(error);

          return false;
        })
    );
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

    await invokeBuild(configContent.split('\n'));
  }
}

main(process.argv.slice(2)).catch(e => {
  console.error(e);
  process.exitCode = 1;
});
