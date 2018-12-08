package parser

import (
	"bytes"
	"fmt"
	"strings"
)

type CommentKind int

const (
	// BlockComment is a comment that starts with /*.
	BlockComment CommentKind = iota
	// LineComment is a comment that starts with //.
	LineComment
)

// Comment is a token that represents either BlockComments or LineComments.
type Comment struct {
	// Kind represents the textual form of this comment. Either BlockComment or
	// LineComment.
	Kind CommentKind
	// The text inside of the comment.
	Value string
}

func (c *Comment) String() string {
	return fmt.Sprintf("<comment %v %q>", CommentKindName[c.Kind], c.Value)
}
func (c *Comment) Type() string { return "Comment" }

var _ Token = &Comment{}

// CommentKindName maps from CommentKind to a human digestable name.
var CommentKindName = map[CommentKind]string{
	BlockComment: "BlockComment",
	LineComment:  "LineComment",
}

// scanLineComment consumes the current rune and the remainder of the line.
func (s *Scanner) scanLineComment() *Comment {
	// Create a buffer and read the current character into it.
	var buf bytes.Buffer
	buf.WriteRune(s.read())

	// Read every subsequent non-newline character into the buffer.
	for {
		switch ch := s.read(); ch {
		case eof:
			goto done
		case '\r':
			// If the next char is a \n consume it, otherwise unread it.
			if s.read() != '\n' {
				s.unread()
			}
			goto done
		case '\n', '\f':
			goto done
		default:
			buf.WriteRune(ch)
		}
	}

done:
	return &Comment{
		Kind:  LineComment,
		Value: strings.Trim(buf.String(), " "),
	}
}

// scanBlockComment consumes the current rune and the remainder of the block comment.
func (s *Scanner) scanBlockComment() *Comment {
	// Create a buffer and read the current character into it.
	var buf bytes.Buffer
	buf.WriteRune(s.read())

	// Read every subsequent whitespace character into the buffer.
	// Non-whitespace characters and EOF will cause the loop to exit.
	for {
		if ch := s.read(); ch == eof {
			break
		} else if ch == '*' {
			if ch := s.read(); ch == '/' {
				break
			} else {
				buf.WriteRune('*')
				buf.WriteRune(ch)
			}
		} else {
			buf.WriteRune(ch)
		}
	}

	return &Comment{
		Kind:  BlockComment,
		Value: strings.Trim(buf.String(), " "),
	}
}
