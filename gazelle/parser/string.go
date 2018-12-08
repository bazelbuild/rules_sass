package parser

import (
	"bytes"
	"fmt"
)

type StringLike interface {
	Token
}

type BadString struct {
}

func (i *BadString) String() string {
	return fmt.Sprintf("<BadString>")
}
func (_ *BadString) Type() string { return "BadString" }

// String is an at keyword token.
type String struct {
	Value string
}

func (i *String) String() string {
	return fmt.Sprintf("<String %q>", i.Value)
}
func (_ *String) Type() string { return "String" }

var _ Token = &String{}

func (s *Scanner) scanString(delimiter rune) StringLike {
	// This section describes how to consume a string token from a stream of code points. It returns either a <string-token> or <bad-string-token>.

	// This algorithm may be called with an ending code point, which denotes the code point that ends the string. If an ending code point is not specified, the current input code point is used.

	// Initially create a <string-token> with its value set to the empty string.
	var buf bytes.Buffer

	// Repeatedly consume the next input code point from the stream:
	for {
		if ch := s.read(); ch == eof {
			// EOF
			// This is a parse error. Return the <string-token>.
			return &String{Value: buf.String()}
		} else if ch == delimiter {
			// ending code point
			// Return the <string-token>.
			return &String{Value: buf.String()}
		} else if ch == '\n' {
			// This is a parse error. Reconsume the current input code point, create a <bad-string-token>, and return it.
			return &BadString{}
		} else if ch == '\\' {
			// U+005C REVERSE SOLIDUS (\)
			// If the next input code point is EOF,
			ch := s.read()
			if ch == eof {
				// do nothing.
			} else if ch == '\n' || ch == '\f' {
				// Otherwise, if the next input code point is a newline, consume it.
			} else if ch == '\r' {
				if s.read() != '\n' {
					s.unread()
				}
			} else {
				s.unread()
			}

			// Otherwise, (the stream starts with a valid escape) consume an escaped code point and append the returned code point to the <string-token>’s value.
			buf.WriteString(s.scanEscape())
		} else {
			// anything else
			// Append the current input code point to the <string-token>’s value.
			buf.WriteRune(ch)
		}
	}
}
