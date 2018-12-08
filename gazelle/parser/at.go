package parser

import (
	"fmt"
)

// At is an at keyword token.
type At struct {
	// The kind of at token.
	Ident IdentLike
}

func (a *At) String() string {
	return fmt.Sprintf("<At %q>", a.Ident.String())
}
func (a *At) Type() string { return "At" }

var _ Token = &At{}

func (s *Scanner) scanAt() *At {
	// Create a buffer and read the current character into it.
	return &At{Ident: s.scanIdent()}
}
