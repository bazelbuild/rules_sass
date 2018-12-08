package parser

import (
	"fmt"
	"strings"
)

type IdentLike interface {
	Token
}

type BadURL struct {
	Value string
}

func (u *BadURL) String() string {
	return fmt.Sprintf("<BadURL %q>", u.Value)
}
func (_ *BadURL) Type() string { return "BadURL" }

var _ Token = &BadURL{}

type URL struct {
	Value string
}

func (u *URL) String() string {
	return fmt.Sprintf("<URL %q>", u.Value)
}
func (_ *URL) Type() string { return "URL" }

var _ Token = &URL{}

// Ident is an ident token.
type Ident struct {
	// The text inside of the comment.
	Value string
}

func (i *Ident) String() string {
	return fmt.Sprintf("<Ident %q>", i.Value)
}
func (_ *Ident) Type() string { return "Ident" }

var _ Token = &Ident{}

// scanIdent consumes the current rune and all contiguous ident runes.
func (s *Scanner) scanIdent() IdentLike {
	// This section describes how to consume an ident-like token from a stream of code points. It returns an <ident-token>, <function-token>, <url-token>, or <bad-url-token>.

	// Consume a name, and let string be the result.
	str := s.consumeAName()

	// If string’s value is an ASCII case-insensitive match for "url"
	if strings.ToLower(str) == "url" {
		// and the next input code point is U+0028 LEFT PARENTHESIS ((), consume
		// it.
		if ch := s.read(); ch != '(' {
			s.unread()
		} else {
			// While the next two input code points are whitespace, consume the
			// next input code point.
			for i := 0; i < 2; i++ {
				if ch := s.peek(); isWhitespace(ch) {
					s.read()
				}
			}

			// If the next one or two input code points are
			// U+0022 QUOTATION MARK ("), U+0027 APOSTROPHE ('), or whitespace followed
			// by U+0022 QUOTATION MARK (") or U+0027 APOSTROPHE ('), then create a
			// <function-token> with its value set to string and return it.
			if ch := s.peek(); ch == '"' || ch == '\'' {
				return &Function{Value: str}
			} else if isWhitespace(ch) {
				s.read()
				if ch := s.peek(); ch == '"' || ch == '\'' {
					return &Function{Value: str}
				}
			}

			// Otherwise, consume a url token, and return it.
			return s.scanURL()
		}
	} else {
		// Otherwise, if the next input code point is U+0028 LEFT PARENTHESIS ((),
		// consume it. Create a <function-token> with its value set to string and
		// return it.
		if ch := s.read(); ch != '(' {
			s.unread()
		} else {
			return &Function{Value: str}
		}
	}

	// Otherwise, create an <ident-token> with its value set to string and return it.
	return &Ident{Value: str}
}

func (s *Scanner) scanURL() Token {
	// Note: This algorithm assumes that the initial "url(" has already been
	// consumed. This algorithm also assumes that it’s being called to consume an
	// "unquoted" value, like url(foo). A quoted value, like url("foo"), is
	// parsed as a <function-token>. Consume an ident-like token automatically
	// handles this distinction; this algorithm shouldn’t be called directly
	// otherwise.

	// Initially create a <url-token> with its value set to the empty string.
	buf := strings.Builder{}

	// Consume as much whitespace as possible.
	for ch := s.peek(); isWhitespace(ch); ch = s.peek() {
		s.read()
	}

	/* test code to print out all the runes that haven't been processed DO NOT SUBMIt
	for {
		if ch := s.read(); ch == eof {
			return &URL{Value: buf.String()}
		} else {
			buf.WriteRune(ch)
		}
	}
	/**/

	// Repeatedly consume the next input code point from the stream:
	for {
		if ch := s.read(); ch == ')' {
			// U+0029 RIGHT PARENTHESIS ())
			// Return the <url-token>.
			return &URL{Value: buf.String()}
		} else if ch == eof {
			// EOF
			// This is a parse error. Return the <url-token>.
			return &URL{Value: buf.String()}
		} else if isWhitespace(ch) {
			// whitespace
			// Consume as much whitespace as possible.
			for ch := s.peek(); isWhitespace(ch); ch = s.peek() {
				s.read()
			}

			// If the next input code point is U+0029 RIGHT PARENTHESIS ()) or EOF, consume it and return the <url-token>
			if ch := s.read(); ch == ')' {
				return &URL{Value: buf.String()}
			} else if ch == eof {
				// (if EOF was encountered, this is a parse error);
				// We choose to return a BadURL here.
				return &BadURL{Value: "firstone"}
				//return &BadURL{Value: s.scanBadURL()}
			} else {
				// otherwise, consume the remnants of a bad url, create a <bad-url-token>, and return it.
				return &BadURL{Value: "otherone"}
				//return &BadURL{Value: s.scanBadURL()}
			}
		} else if ch == '"' || ch == '\'' || ch == '(' || isNonASCIICodePoint(ch) {
			// U+0022 QUOTATION MARK (")
			// U+0027 APOSTROPHE (')
			// U+0028 LEFT PARENTHESIS (()
			// non-printable code point
			// This is a parse error. Consume the remnants of a bad url, create a <bad-url-token>, and return it.
			return &BadURL{Value: "thisone"}
			//return &BadURL{Value: s.scanBadURL()}
		} else if ch == '\\' {
			// U+005C REVERSE SOLIDUS (\)
			// If the stream starts with a valid escape, consume an escaped code point and append the returned code point to the <url-token>’s value.
			buf.WriteString(s.scanEscape())
		} else {
			buf.WriteRune(ch)
		}
	}

	// anything else
	// Append the current input code point to the <url-token>’s value.
	return &URL{Value: ""}
}

func (s *Scanner) scanBadURL() string {
	// This section describes how to consume the remnants of a bad url from a
	// stream of code points, "cleaning up" after the tokenizer realizes that
	// it’s in the middle of a <bad-url-token> rather than a <url-token>. It
	// returns nothing; its sole use is to consume enough of the input stream to
	// reach a recovery point where normal tokenizing can resume.

	buf := strings.Builder{}

	// Repeatedly consume the next input code point from the stream:
	for {
		if ch := s.read(); ch == '(' || ch == eof {
			// U+0029 RIGHT PARENTHESIS ())
			// EOF
			// Return.
			return buf.String()
		} else if ch == '\\' {
			// the input stream starts with a valid escape
			// Consume an escaped code point. This allows an escaped right
			// parenthesis ("\)") to be encountered without ending the
			// <bad-url-token>. This is otherwise identical to the "anything else"
			// clause.
			return s.scanEscape()
		} else {
			// anything else
			// Do nothing.
			buf.WriteRune(ch)
		}
	}
}
