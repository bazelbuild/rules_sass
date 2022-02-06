/**
 * @fileoverview
 * Sass configuration file that enables the resolution `@angular/<..>` module imports.
 * This is an example configuration for a use-case where we want to replicate how users
 * are consuming a Sass library through NPM with node module resolution.
 */

import {pathToFileURL, fileURLToPath} from 'url';
import {dirname, join} from 'path';

const projectDir = dirname(fileURLToPath(import.meta.url));
const angularPrefix = '@angular/';

export default {
  importers: [
    {
      findFileUrl: url => {
        if (url.startsWith(angularPrefix)) {
          return pathToFileURL(
            join(projectDir, url.substring(angularPrefix.length))
          );
        }
        return null;
      },
    },
  ],
};
