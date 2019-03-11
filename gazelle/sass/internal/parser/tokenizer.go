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
// Package parser provides utilities useful for parsing a sass file.
package parser

// I used
// https://blog.gopheracademy.com/advent-2014/parsers-lexers/
// heavily for reference when creating this file.

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strings"
	"unicode/utf8"
)

var eof = rune(0)

type Token interface {
	String() string
	Type() string
}

func isHexDigit(ch rune) bool {
	// A digit, or a code point between U+0041 LATIN CAPITAL LETTER A (A) and U+0046 LATIN CAPITAL LETTER F (F) inclusive, or a code point between U+0061 LATIN SMALL LETTER A (a) and U+0066 LATIN SMALL LETTER F (f) inclusive.
	return isDigit(ch) || (ch >= 'a' && ch <= 'f') || (ch >= 'A' && ch <= 'F')
}
func isDigit(ch rune) bool {
	// A code point between U+0030 DIGIT ZERO (0) and U+0039 DIGIT NINE (9) inclusive.
	return (ch >= '0' && ch <= '9')
}
func isLetter(ch rune) bool {
	// An uppercase letter or a lowercase letter.
	return isLowercaseLetter(ch) || isUppercaseLetter(ch)
}
func isLowercaseLetter(ch rune) bool {
	// A code point between U+0061 LATIN SMALL LETTER A (a) and U+007A LATIN SMALL LETTER Z (z) inclusive.
	return ch >= 'a' && ch <= 'z'
}
func isUppercaseLetter(ch rune) bool {
	// A code point between U+0041 LATIN CAPITAL LETTER A (A) and U+005A LATIN CAPITAL LETTER Z (Z) inclusive.
	return ch >= 'A' && ch <= 'Z'
}
func isNonASCIICodePoint(ch rune) bool {
	// A code point with a value equal to or greater than U+0080 <control>.
	return ch >= '\u0080'
}
func isNameStartCodePoint(ch rune) bool {
	// A letter, a non-ASCII code point, or U+005F LOW LINE (_).
	return isLetter(ch) || isNonASCIICodePoint(ch) || ch == '_'
}
func isNameCodePoint(ch rune) bool {
	// A name-start code point, a digit, or U+002D HYPHEN-MINUS (-).
	return isNameStartCodePoint(ch) || isDigit(ch) || ch == '-'
}

// New creates a scanner that reads the provided io.Reader.
func New(r io.Reader) *Scanner {
	return &Scanner{r: bufio.NewReader(r)}
}

// Scanner is responsible for reading an input file.
type Scanner struct {
	r *bufio.Reader
}

func (s *Scanner) peek() rune {
	if r := []rune(s.peekNOffset(1, 0)); len(r) == 1 {
		return r[0]
	}
	return eof
}

// peek gives you the offsetth rune.
func (s *Scanner) peekOffset(offset int) rune {
	if r := []rune(s.peekNOffset(1, offset)); len(r) == 1 {
		return r[0]
	}
	return eof
}

// peekN gives you the next n runes.
func (s *Scanner) peekN(n int) string {
	return s.peekNOffset(n, 0)
}

// peekN gives you the next n runes offet by offset. Note that as n and offset increase this becomes increasingly less efficient.
func (s *Scanner) peekNOffset(n, offset int) string {
	// https://blog.golang.org/strings
	// Some people think Go strings are always UTF-8, but they are not: only
	// string literals are UTF-8. As we showed in the previous section,
	// string values can contain arbitrary bytes; as we showed in this one,
	// string literals always contain UTF-8 text as long as they have no
	// byte-level escapes.

	if offset != 0 {
		// if the offset is nonzero, read that many characters and then skip
		// that many bytes from the buffer.
		offset = len([]byte(s.peekNOffset(offset, 0)))
	}

	var ret []rune
	// Start reading from offset + 1 so that we don't include extra letters.
	for i := offset + 1; len(ret) < n; i++ {
		bytes, err := s.r.Peek(i)
		if err != nil {
			return ""
		}
		r, _ := utf8.DecodeLastRune(bytes)
		if r == '\uFFFD' {
			// The rune is invalid and we need more bytes.
			continue
		}

		ret = append(ret, r)
	}

	// We have now encoded the string as UTF8.
	return string(ret)
}

func (s *Scanner) checkIfTwoCodePointsAreValidEscape() bool {
	// This section describes how to check if two code points are a valid escape. The algorithm described here can be called explicitly with two code points, or can be called with the input stream itself. In the latter case, the two code points in question are the current input code point and the next input code point, in that order.

	// Note: This algorithm will not consume any additional code point.

	// If the first code point is not U+005C REVERSE SOLIDUS (\), return false.
	if s.peek() != '\\' {
		return false
	} else if s.peekN(2) == "\\\n" {
		// Otherwise, if the second code point is a newline, return false.
		return false
	}

	// Otherwise, return true.
	return true
}

func (s *Scanner) checkIfThreeCodePointsWouldStartAnIdentifier() bool {
	// This section describes how to check if three code points would start an identifier. The algorithm described here can be called explicitly with three code points, or can be called with the input stream itself. In the latter case, the three code points in question are the current input code point and the next two input code points, in that order.

	// Note: This algorithm will not consume any additional code points.

	// Look at the first code point.
	if ch := s.peek(); ch == '-' {
		// U+002D HYPHEN-MINUS
		// If the second code point is a name-start code point or a U+002D HYPHEN-MINUS,
		if ch := s.peekOffset(2); isNameStartCodePoint(ch) || ch == '-' {
			return true
		}
		// Read a character so that checkIfTwoCodePointsAreValidEscape is operating on the 2nd character.
		s.read()
		// Unread the read that just happened so we are non-destructive.
		defer s.unread()
		// or the second and third code points are a valid escape, return true. Otherwise, return false.
		return s.checkIfTwoCodePointsAreValidEscape()
	} else if isNameCodePoint(ch) {
		// name-start code point
		// Return true.
		return true
	} else if ch == '\\' {
		// U+005C REVERSE SOLIDUS (\)
		// If the first and second code points are a valid escape, return true. Otherwise, return false.
		return s.checkIfTwoCodePointsAreValidEscape()
	} else {
		// anything else
		// Return false.
		return false
	}
}

func (s *Scanner) checkIfThreeCodePointsWouldStartANumber() bool {
	//This section describes how to check if three code points would start a number. The algorithm described here can be called explicitly with three code points, or can be called with the input stream itself. In the latter case, the three code points in question are the current input code point and the next two input code points, in that order.

	//Note: This algorithm will not consume any additional code points.

	//Look at the first code point:

	if ch := s.peek(); ch == '+' || ch == '-' {
		//U+002B PLUS SIGN (+)
		//U+002D HYPHEN-MINUS (-)
		if isDigit(s.peekOffset(1)) {
			//If the second code point is a digit, return true.
			return true
		} else if s.peekOffset(1) == '.' && isDigit(s.peekOffset(2)) {
			//Otherwise, if the second code point is a U+002E FULL STOP (.) and the third code point is a digit, return true.
			return true
		}

		//Otherwise, return false.
		return false
	} else if ch == '.' {
		//U+002E FULL STOP (.)
		if s.peekOffset(1) == '.' {
			//If the second code point is a digit, return true.
			return true
		}
		// Otherwise, return false.
		return false
	} else if isDigit(ch) {
		//digit
		//Return true.
		return true
	}
	//anything else
	//Return false.
	return false
}

// read reads the next rune from the bufferred reader.
// Returns the rune(0) if an error occurs (or io.EOF is returned).
func (s *Scanner) read() rune {
	ch, _, err := s.r.ReadRune()
	if err != nil {
		return eof
	}

	switch ch {
	// Replace any U+000D CARRIAGE RETURN (CR) code points,
	case '\u000D':
		next, _, err := s.r.ReadRune()
		if err != nil {
			return eof
		}
		// pairs of U+000D CARRIAGE RETURN (CR) followed by U+000A LINE FEED (LF), by a single U+000A LINE FEED (LF) code point.
		if next != '\u000A' {
			s.r.UnreadRune()
		}
		return '\u000A'
	// U+000C FORM FEED (FF) code points,
	case '\u000C':
		return '\u000A'

	// Replace any U+0000 NULL or surrogate code points with U+FFFD REPLACEMENT CHARACTER (�).
	case '\u0000':
		return '\uFFFD'
	}

	return ch
}

// unread places the previously read rune back on the reader.
func (s *Scanner) unread() error { return s.r.UnreadRune() }

func (s *Scanner) consumeAName() string {
	// Note: This algorithm does not do the verification of the first few code points that are necessary to ensure the returned code points would constitute an <ident-token>. If that is the intended use, ensure that the stream starts with an identifier before calling this algorithm.

	// Let result initially be an empty string.
	res := &strings.Builder{}

	// Repeatedly consume the next input code point from the stream:

	for {
		if ch := s.read(); isNameCodePoint(ch) {
			// name code point
			// Append the code point to result.
			res.WriteRune(ch)
		} else if ch == '\\' {
			// the stream starts with a valid escape
			// Consume an escaped code point. Append the returned code point to result.
			return s.scanEscape()
		} else {
			// anything else
			// Reconsume the current input code point. Return result.
			s.unread()
			return res.String()
		}
	}
}

// scanEscape scans the buffer and consumes the entire escape string to completion.
func (s *Scanner) scanEscape() string {
	var buf bytes.Buffer
	for i := 0; i < 6; i++ {
		ch := s.read()
		if isHexDigit(ch) {
			// Consume as many hex digits as possible, but no more than 5. Note
			// that this means 1-6 hex digits have been consumed in total. If the
			// next input code point is whitespace, consume it as well. Interpret
			// the hex digits as a hexadecimal number. If this number is zero, or
			// is for a surrogate, or is greater than the maximum allowed code
			// point, return U+FFFD REPLACEMENT CHARACTER (�). Otherwise, return
			// the code point with that value.
			buf.WriteRune(ch)
		} else if ch == eof {
			// EOF
			// This is a parse error. Return U+FFFD REPLACEMENT CHARACTER (�).
			return "\uFFFD"
		} else {
			// anything else
			// Return the current input code point.
			return buf.String()
		}
	}
	return buf.String()
}

// Scan returns the next scanned token.
// For the scanning rules see: https://drafts.csswg.org/css-syntax-3/#consume-token
func (s *Scanner) Scan() Token {
	// Read the next rune.
	ch := s.read()

	// Consume comments.
	if ch == '/' {
		// Possibly a comment
		switch s.read() {
		case '/':
			return s.scanLineComment()
		case '*':
			return s.scanBlockComment()
		default:
			// Not a comment, it's probably a division operator.
			s.unread()
		}
	}

	// Consume the next input code point.
	// https://drafts.csswg.org/css-syntax-3/#next-input-code-point
	// I don't think there is any parser action here.

	// Consume as much whitespace as possible. Return a <whitespace-token>.
	if isWhitespace(ch) {
		s.unread()
		return s.scanWhitespace()
	}

	switch ch {
	case '"':
		// U+0022 QUOTATION MARK (")
		// Consume a string token and return it.
		return s.scanString('"')
	case '#':
		// U+0023 NUMBER SIGN (#)
		// If the next input code point is a name code point or the next two input code points are a valid escape, then:
		// Create a <hash-token>.
		// If the next 3 input code points would start an identifier, set the <hash-token>’s type flag to "id".
		// Consume a name, and set the <hash-token>’s value to the returned string.
		// Return the <hash-token>.
		return s.scanHashToken()
		// Otherwise, return a <delim-token> with its value set to the current input code point.
	case '\'':
		// U+0027 APOSTROPHE (')
		// Consume a string token and return it.
		return s.scanString('\'')

	case '(':
		// U+0028 LEFT PARENTHESIS (()
		// Return a <(-token>.
		return &LeftParenthesis{}
	case ')':
		// U+0029 RIGHT PARENTHESIS ())
		// Return a <)-token>.
		return &RightParenthesis{}
		// U+002B PLUS SIGN (+)
		// If the input stream starts with a number, reconsume the current input code point, consume a numeric token and return it.
		// Otherwise, return a <delim-token> with its value set to the current input code point.

	case ',':
		// U+002C COMMA (,)
		// Return a <comma-token>.
		return &Comma{}
	case '-':
		// U+002D HYPHEN-MINUS (-)
		if ch := s.peek(); isDigit(ch) {
			// scanDigit expects the - to be in the buffer, put it back in.
			s.unread()
			// If the input stream starts with a number, reconsume the current input code point, consume a numeric token, and return it.
			return s.scanDigit()
		} else if s.peekN(2) == "\u002d\u003E" {
			// Otherwise, if the next 2 input code points are U+002D HYPHEN-MINUS U+003E GREATER-THAN SIGN (->), consume them and return a <CDC-token>.
			s.read()
			s.read()
			return &CDC{}
		}
		// Otherwise, if the input stream starts with an identifier, reconsume the current input code point, consume an ident-like token, and return it.

		// Otherwise, return a <delim-token> with its value set to the current input code point.
		return &Delim{Value: "-"}

	case '.':
		// U+002E FULL STOP (.)
		if isDigit(s.peek()) {
			// If the input stream starts with a number, reconsume the current input code point, consume a numeric token, and return it.
			s.unread()
			return s.scanDigit()
		}
		// Otherwise, return a <delim-token> with its value set to the current input code point.
		return &Delim{Value: "."}

	case ':':
		// U+003A COLON (:)
		// Return a <colon-token>.
		return &Colon{}
	case ';':
		// U+003B SEMICOLON (;)
		// Return a <semicolon-token>.
		return &Semicolon{}
		// U+003C LESS-THAN SIGN (<)
		// If the next 3 input code points are U+0021 EXCLAMATION MARK U+002D HYPHEN-MINUS U+002D HYPHEN-MINUS (!--), consume them and return a <CDO-token>.
		// Otherwise, return a <delim-token> with its value set to the current input code point.

	case '@':
		// U+0040 COMMERCIAL AT (@)
		// If the next 3 input code points would start an identifier, consume a name, create an <at-keyword-token> with its value set to the returned value, and return it.
		// Otherwise, return a <delim-token> with its value set to the current input code point.
		return s.scanAt()

	case '[':
		// U+005B LEFT SQUARE BRACKET ([)
		// Return a <[-token>.
		return &LeftSquareBracket{}

	case '\\':
		// U+005C REVERSE SOLIDUS (\)
		// If the input stream starts with a valid escape, reconsume the current input code point, consume an ident-like token, and return it.
		// Otherwise, this is a parse error. Return a <delim-token> with its value set to the current input code point.
		/*
			if isLetter(ch) || ch == '\\' {
				s.unread()
				ident := s.scanIdent()
				if ch := s.read(); ch == '(' {
					return &Function{
						Ident: ident,
					}
				}
				// If the next value isn't a '(' then restore it for future parsing.
				s.unread()
				return ident
			}
		*/

		return &Ident{Value: s.scanEscape()}

	case ']':
		// U+005D RIGHT SQUARE BRACKET (])
		// Return a <]-token>.
		return &RightSquareBracket{}
	case '{':
		// U+007B LEFT CURLY BRACKET ({)
		// Return a <{-token>.
		return &LeftCurlyBracket{}
	case '}':
		// U+007D RIGHT CURLY BRACKET (})
		// Return a <}-token>.
		return &RightCurlyBracket{}
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		// digit
		// Reconsume the current input code point, consume a numeric token, and return it.
		s.unread()
		return s.scanDigit()
	case eof:
		// EOF
		// Return an <EOF-token>.
		return &EOF{}
	default:
		if isNameStartCodePoint(ch) {
			// name-start code point
			// Reconsume the current input code point, consume an ident-like token, and return it.
			// Restore the last character
			s.unread()
			return s.scanIdent()
		}
		// anything else
		// Return a <delim-token> with its value set to the current input code point.
		return &Delim{Value: string(ch)}
	}

	panic("This line should never be reached. All codepaths should travel " +
		"through the above switch statement")
}

// ScanAll returns all remaining tokens in the file.
func (s *Scanner) ScanAll() ([]Token, error) {
	var tokens []Token
	for {
		token := s.Scan()
		tokens = append(tokens, token)
		if _, ok := token.(*EOF); ok {
			return tokens, nil
		}

		if len(tokens) > 200 {
			// DO NOT SUBMIT
			// When you end up in tight loops this helps give debug info.
			return tokens, fmt.Errorf("Token count overflow")
		}
	}
}
