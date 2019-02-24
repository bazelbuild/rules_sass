/* Copyright 2019 The Bazel Authors. All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

   http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package gazelle

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"

	"github.com/bazelbuild/rules_sass/gazelle/parser"
)

// FileInfo contains metadata extracted from a .sass/.scss file.
type FileInfo struct {
	Path, Name string

	Imports []string
}

// sassFileInfo takes a dir and file name and parses the .sass file into
// the constituent components, extracting metadata like the set of
// imports that it has.
func sassFileInfo(dir, name string) FileInfo {
	info := FileInfo{
		Path: filepath.Join(dir, name),
		Name: name,
	}

	file, err := os.Open(filepath.Join(dir, name))
	if err != nil {
		log.Printf("%s: error reading sass file: %v", info.Path, err)
		return info
	}

	s := parser.New(file)

	// Abstract the scan so that it doesn't return whitespace.
	scan := func() parser.Token {
		for {
			t := s.Scan()
			if _, ok := t.(*parser.WhiteSpace); !ok {
				return t
			}
		}
	}

	// parseError is a convenience wrapper to print to stderr a parse error message
	// in a standardized format.
	parseError := func(msg string, args ...interface{}) {
		// Prepend the file.Name() so we tell the user which file had the parse error.
		args = append([]interface{}{file.Name()}, args...)
		fmt.Fprintf(os.Stderr, "%s: "+msg+"\n", args...)
	}

	for t := scan(); t.Type() != "EOF"; t = scan() {
		var imports []string
		switch v := t.(type) {
		case *parser.At:
			if i, ok := v.Ident.(*parser.Ident); ok && i.Value == "import" {
				// import_stmt ::= <At <Ident "import">> <String> (<Comma> <String>)+ <Semicolon>
				// Scan the next symbol. It should be either a <String>
				t = scan()
				if s, ok := t.(*parser.String); ok {
					imports = append(imports, s.Value)
				} else {
					parseError("expected a string but got a %v\n", t.Type())
					continue
				}

				// Scan the next symbol. It should be either a <Comma> or a <Semicolon>
				t = scan()

				// Loop through the next tokens in case the user provided a comma
				// separated list of files to include.
				for {
					if _, ok := t.(*parser.Comma); !ok {
						break
					}

					t = scan()
					if s, ok := t.(*parser.String); ok {
						imports = append(imports, s.Value)
					} else {
						break
					}

					t = scan()
				}

				if _, ok := t.(*parser.Semicolon); !ok {
					// If the last character isn't a semicolon then it isn't a valid
					// import statement. Don't add it to the list of imports.
					parseError("imports should end in semicolons\n")
					continue
				}
				if len(imports) > 0 {
					info.Imports = append(info.Imports, imports...)
				}
			}
		}
	}

	sort.Strings(info.Imports)

	return info
}

// unquoteSassString takes a string that has a complex quoting around it
// and returns a string without the complex quoting.
func unquoteSassString(q []byte) string {
	// Adjust quotes so that Unquote is happy. We need a double quoted string
	// without unescaped double quote characters inside.
	noQuotes := bytes.Split(q[1:len(q)-1], []byte{'"'})
	if len(noQuotes) != 1 {
		for i := 0; i < len(noQuotes)-1; i++ {
			if len(noQuotes[i]) == 0 || noQuotes[i][len(noQuotes[i])-1] != '\\' {
				noQuotes[i] = append(noQuotes[i], '\\')
			}
		}
		q = append([]byte{'"'}, bytes.Join(noQuotes, []byte{'"'})...)
		q = append(q, '"')
	}
	if q[0] == '\'' {
		q[0] = '"'
		q[len(q)-1] = '"'
	}

	s, err := strconv.Unquote(string(q))
	if err != nil {
		log.Panicf("unquoting string literal %s from sass: %v", q, err)
	}
	return s
}
