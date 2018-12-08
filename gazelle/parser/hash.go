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
