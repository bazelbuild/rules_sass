package parser

import (
	"bytes"
	"fmt"
)

type WhiteSpace struct {
	Value string
}

func (w *WhiteSpace) String() string { return fmt.Sprintf("<WhiteSpace %q>", w.Value) }
func (_ *WhiteSpace) Type() string   { return "WhiteSpace" }

var _ Token = &WhiteSpace{}

// EOF is a token that represents the end of the file.
type EOF struct{}

func (e *EOF) String() string { return fmt.Sprintf("<eof>") }
func (_ *EOF) Type() string   { return "EOF" }

var _ Token = &EOF{}

func isWhitespace(ch rune) bool {
	// A newline, U+0009 CHARACTER TABULATION, or U+0020 SPACE.
	return ch == '\n' || ch == '\t' || ch == ' '
}

// scanWhitespace consumes the current rune and all contiguous whitespace.
func (s *Scanner) scanWhitespace() Token {
	// Create a buffer and read the current character into it.
	var buf bytes.Buffer
	buf.WriteRune(s.read())

	// Read every subsequent whitespace character into the buffer.
	// Non-whitespace characters and EOF will cause the loop to exit.
	for {
		if ch := s.read(); ch == eof {
			break
		} else if !isWhitespace(ch) {
			s.unread()
			break
		} else {
			buf.WriteRune(ch)
		}
	}

	return &WhiteSpace{Value: buf.String()}
}

type Delim struct {
	Value string
}

func (d *Delim) String() string { return fmt.Sprintf("<delim %q>", d.Value) }
func (_ *Delim) Type() string   { return "Delim" }

var _ Token = &Delim{}
