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
	"fmt"
	"strconv"
	"strings"
)

type NumberKind int

const (
	Integer NumberKind = iota
	Number
)

var NumberKindName = map[NumberKind]string{
	Integer: "Integer",
	Number:  "Number",
}

// Percentage is a token.
type Percentage struct {
	Value float64
}

func (p *Percentage) String() string {
	return fmt.Sprintf("<Percentage %f>", p.Value)
}
func (_ *Percentage) Type() string { return "Percentage" }

var _ Token = &Percentage{}

// Dimension is a token.
type Dimension struct {
	Value float64
	Kind  NumberKind
	Unit  IdentLike
}

func (d *Dimension) String() string {
	return fmt.Sprintf("<Dimension (%s) %f>", NumberKindName[d.Kind], d.Value)
}
func (_ *Dimension) Type() string { return "Dimension" }

var _ Token = &Dimension{}

// Digit is an a token.
type Digit struct {
	Value float64
	Kind  NumberKind
}

func (d *Digit) String() string {
	return fmt.Sprintf("<Digit (%s) %f>", NumberKindName[d.Kind], d.Value)
}
func (_ *Digit) Type() string { return "Digit" }

var _ Token = &Digit{}

func (s *Scanner) scanNumber() (float64, NumberKind) {
	// This section describes how to consume a number from a stream of code points. It returns a numeric value, and a type which is either "Integer" or "number".

	// Note: This algorithm does not do the verification of the first few code points that are necessary to ensure a number can be obtained from the stream. Ensure that the stream starts with a number before calling this algorithm.

	// Execute the following steps in order:

	// Initially set type to "Integer". Let repr be the empty string.
	t := Integer
	repr := strings.Builder{}

	// If the next input code point is U+002B PLUS SIGN (+) or U+002D HYPHEN-MINUS (-), consume it and append it to repr.
	if ch := s.peek(); ch == '+' || ch == '-' {
		repr.WriteRune(s.read())
	}

	// While the next input code point is a digit, consume it and append it to repr.
	for isDigit(s.peek()) {
		repr.WriteRune(s.read())
	}

	if s.peek() == '.' && isDigit(s.peekOffset(1)) {
		// If the next 2 input code points are U+002E FULL STOP (.) followed by a digit, then:
		// Consume them.
		// Append them to repr.
		// Set type to "Number".
		repr.WriteRune(s.read())
		repr.WriteRune(s.read())
		t = Number
	}

	// While the next input code point is a digit, consume it and append it to repr.
	for ch := s.peek(); isDigit(ch); ch = s.peek() {
		repr.WriteRune(s.read())
	}

	// If the next 2 or 3 input code points are U+0045 LATIN CAPITAL LETTER E (E) or U+0065 LATIN SMALL LETTER E (e), optionally followed by U+002D HYPHEN-MINUS (-) or U+002B PLUS SIGN (+), followed by a digit, then:
	if buf := s.peekN(2); (buf == "e+" || buf == "e-" || buf == "E+" || buf == "E-") && isDigit(s.peekOffset(2)) {
		// Consume them.
		// Append them to repr.
		repr.WriteRune(s.read())
		repr.WriteRune(s.read())

		// Set type to "Number".

		// Write 2 chars out (e and the +/-) and then let the next loop handle
		// writing out the remaining digits.

		t = Number
	}

	// While the next input code point is a digit, consume it and append it to repr.
	for ch := s.peek(); isDigit(ch); ch = s.peek() {
		repr.WriteRune(s.read())
	}
	// Convert repr to a Number, and set the value to the returned value.
	v, _ := strconv.ParseFloat(repr.String(), 64)

	// Return value and type.
	return v, t
}

func (s *Scanner) scanDigit() Token {
	// This section describes how to consume a numeric token from a stream of code points. It returns either a <number-token>, <percentage-token>, or <dimension-token>.

	// Consume a number and let number be the result.
	number, t := s.scanNumber()

	// If the next 3 input code points would start an identifier, then:
	if s.checkIfThreeCodePointsWouldStartAnIdentifier() {
		// Create a <dimension-token> with the same value and type flag as Number, and a unit set initially to the empty string.
		// Consume a name. Set the <dimension-token>’s unit to the returned value.
		// Return the <dimension-token>.
		return &Dimension{Value: number, Kind: t, Unit: s.scanIdent()}
	} else if s.peek() == '%' {
		// Otherwise, if the next input code point is U+0025 PERCENTAGE SIGN (%), consume it.
		s.read()
		//Create a <percentage-token> with the same value as number, and return it.
		return &Percentage{Value: number}
	}

	// Otherwise, create a <number-token> with the same value and type flag as number, and return it.
	return &Digit{Value: number, Kind: t}
}
