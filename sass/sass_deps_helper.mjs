/**
 * @license
 * Copyright Google LLC. All Rights Reserved.
 * Use of this source code is governed by the Apache 2.0 license that can be found in the LICENSE.txt file.
 */

import * as path from 'path';
import * as url from 'url';
import {runfiles} from '@bazel/runfiles';

// Path to the Bazel-managed node modules for the Sass rules. See `sass_repositories.bzl`.
const sassDepsWorkspacePath = runfiles.resolve('build_bazel_rules_sass_deps');
const sassNodeModules = path.join(sassDepsWorkspacePath, 'node_modules/');

/**
 * Helper function that can be used to import NPM modules from the Sass
 * dependency external workspace.
 *
 * Without this helper, module requests would either be reliant on the node modules
 * linker from the Bazel NodeJS rules, or on the legacy patched module resolution
 * that does not work with ECMAScript modules anyway. The node module linker is unreliable
 * for the persistent worker in general, as it would rely on the `node_modules/` folder in the
 * execroot that can be shared in non-sandbox environments or for worker instances.
 * https://docs.bazel.build/versions/main/command-line-reference.html#flag--worker_sandboxing.
 *
 * For additional context: In some situations, the use of ESM is necessary to workaround
 * issues with the legacy patched module resolution (that occurs in CJS). The patched resolution
 * will always prioritize `.mjs` files over `.js` leading to `ERR_REQUIRE_ESM` errors.
 */
export async function importNpmModuleFromSassDeps(modulePath) {
  const targetPath = path.join(sassNodeModules, modulePath);
  return (await import(url.pathToFileURL(targetPath))).default;
}
