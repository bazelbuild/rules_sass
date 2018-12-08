package parser

import "fmt"

// Function is a function token.
type Function struct {
	// The text inside of the comment.
	Value string
}

func (f *Function) String() string {
	return fmt.Sprintf("<Function %q>", f.Value)
}
func (_ *Function) Type() string { return "Function" }

var _ Token = &Function{}
