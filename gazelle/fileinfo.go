/* Copyright 2018 The Bazel Authors. All rights reserved.

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
	"io/ioutil"
	"log"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

// FileInfo contains metadata extracted from a .sass/.scss file.
type FileInfo struct {
	Path, Name string

	Imports []string
}

var sassRe = buildProtoRegexp()

// sassFileInfo takes a dir and file name and parses the .sass file into
// the constituent components, extracting metadata like the set of
// imports that it has.
func sassFileInfo(dir, name string) FileInfo {
	info := FileInfo{
		Path: filepath.Join(dir, name),
		Name: name,
	}
	content, err := ioutil.ReadFile(info.Path)
	if err != nil {
		log.Printf("%s: error reading sass file: %v", info.Path, err)
		return info
	}

	for _, match := range sassRe.FindAllSubmatch(content, -1) {
		switch {
		case match[importSubexpIndex] != nil:
			imp := match[importSubexpIndex]
			info.Imports = append(info.Imports, unquoteSASSString(imp))

		default:
			// Comment matched. Nothing to extract.
		}
	}
	sort.Strings(info.Imports)

	return info
}

// unquoteSASSString takes a string that has a complex quoting around it
// and returns a string without the complex quoting.
func unquoteSASSString(q []byte) string {
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

const (
	importSubexpIndex = 1
)

// Based on https://pygments-doc.readthedocs.io/en/latest/_modules/pygments/lexers/web.html#SassLexer
func buildProtoRegexp() *regexp.Regexp {
	hexEscape := `\\[xX][0-9a-fA-f]{2}`
	octEscape := `\\[0-7]{3}`
	charEscape := `\\[abfnrtv'"\\]`
	charValue := strings.Join([]string{hexEscape, octEscape, charEscape, "[^\x00\\'\\\"\\\\]"}, "|")
	strLit := `'(?:` + charValue + `|")*'|"(?:` + charValue + `|')*"`
	importStmt := `\bimport\s*(?P<import>` + strLit + `)\s*;`
	sassReSrc := strings.Join([]string{importStmt}, "|")
	return regexp.MustCompile(sassReSrc)
}
