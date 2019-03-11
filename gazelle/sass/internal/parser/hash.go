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
package parser

import (
	"bytes"
	"fmt"
)

// Hash is a hash keyword token.
type Hash struct {
	// The value of the hash token.
	Value string
}

func (h *Hash) String() string {
	return fmt.Sprintf("<Hash %q>", h.Value)
}
func (_ *Hash) Type() string { return "Hash" }

var _ Token = &Hash{}

// scanHashToken consumes the current rune and the remainder of the line.
func (s *Scanner) scanHashToken() *Hash {
	// Create a buffer and read the current character into it.
	var buf bytes.Buffer
	buf.WriteRune(s.read())

	for {
		if ch := s.read(); ch == '\\' {
			buf.WriteString(s.scanEscape())
		} else if isDigit(ch) || isLetter(ch) || ch == '_' || ch == '-' {
			buf.WriteRune(ch)
		} else {
			break
		}
	}

	return &Hash{
		Value: buf.String(),
	}
}
