package parser

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

var (
	updateGoldens = flag.Bool("update_goldens", false, "Set to true to update the goldens that are checked in. Note that this only works when running under `go test`")
)

const (
	testFileSuffix = ".in.scss"
	goldenSuffix   = ".out.txt"
)

func TestTokenizing(t *testing.T) {
	tests := []struct {
		in  string
		out []Token
	}{
		{
			in:  "",
			out: []Token{&EOF{}},
		},
		{
			in: "ident",
			out: []Token{&Ident{
				Value: "ident",
			}, &EOF{}},
		},
		{
			in: "ident",
			out: []Token{&Ident{
				Value: "ident",
			}, &EOF{}},
		},
		{
			in: `"string"`,
			out: []Token{&String{
				Value: "string",
			}, &EOF{}},
		},
		{
			in: `'string'`,
			out: []Token{&String{
				Value: "string",
			}, &EOF{}},
		},
		{
			// This is a total head fake testcase but I think it is correct.
			// If the string is not terminated by the starting delimiter, but
			// instead is terminated by eof it is a valid string!
			in: `'string`,
			out: []Token{&String{
				Value: "string",
			}, &EOF{}},
		},
		{
			in:  "'badstring\n",
			out: []Token{&BadString{}, &EOF{}},
		},
		{
			in:  "'badstring\r\n'goodstring'",
			out: []Token{&BadString{}, &String{Value: "goodstring"}, &EOF{}},
		},
		{
			in:  "'string\\",
			out: []Token{&String{Value: "string\uFFFD"}, &EOF{}},
		},
		{
			in:  "'string\f",
			out: []Token{&BadString{}, &EOF{}},
		},
		{
			in: "@at",
			out: []Token{
				&At{
					Ident: &Ident{
						Value: "at",
					},
				},
				&EOF{}},
		},
		{
			in: "@at moo",
			out: []Token{
				&At{
					Ident: &Ident{
						Value: "at",
					},
				},
				&WhiteSpace{Value: " "},
				&Ident{Value: "moo"},
				&EOF{}},
		},
		{
			in: "#abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ01234567890_-",
			out: []Token{
				&Hash{
					Value: "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ01234567890_-",
				},
				&EOF{}},
		},
		{
			in: `'asdf\123456'`,
			out: []Token{&String{
				Value: `asdf123456`,
			}, &EOF{}},
		},
		{
			in: "'-\\123456A'",
			out: []Token{
				&String{
					Value: "-123456A",
				},
				&EOF{}},
		},
		{
			in: "\\7890AB",
			out: []Token{&Ident{
				Value: "7890AB",
			}, &EOF{}},
		},
		{
			in: "\\ABCDEF",
			out: []Token{&Ident{
				Value: "ABCDEF",
			}, &EOF{}},
		},
		{
			in: "\\ABC moo",
			out: []Token{
				&Ident{
					Value: "ABC",
				},
				&Ident{
					Value: "moo",
				},
				&EOF{}},
		},
		{
			in: "// \\r\\n newline \r\nident",
			out: []Token{
				&Comment{
					Kind:  LineComment,
					Value: "\\r\\n newline",
				},
				&Ident{Value: "ident"}, &EOF{}},
		},
		{
			in: "// \\f newline \fident",
			out: []Token{
				&Comment{
					Kind:  LineComment,
					Value: "\\f newline",
				},
				&Ident{Value: "ident"}, &EOF{}},
		},
		{
			in: "// \\r newline \rident",
			out: []Token{
				&Comment{
					Kind:  LineComment,
					Value: "\\r newline",
				},
				&Ident{Value: "ident"}, &EOF{}},
		},
		{
			in: "/* block comment */",
			out: []Token{&Comment{
				Kind:  BlockComment,
				Value: "block comment",
			}, &EOF{}},
		},
		{
			in: "/* * /* /* nested block comment */",
			out: []Token{&Comment{
				Kind:  BlockComment,
				Value: "* /* /* nested block comment",
			}, &EOF{}},
		},
		{
			in: "/* comment */ ident",
			out: []Token{
				&Comment{
					Kind:  BlockComment,
					Value: "comment",
				},
				&WhiteSpace{
					Value: " ",
				},
				&Ident{
					Value: "ident",
				}, &EOF{}},
		},
		{
			in: "// comment\nident",
			out: []Token{
				&Comment{
					Kind:  LineComment,
					Value: "comment",
				},
				&Ident{
					Value: "ident",
				}, &EOF{}},
		},
		{
			in: "{",
			out: []Token{
				&LeftCurlyBracket{},
				&EOF{}},
		},
		{
			in: "}",
			out: []Token{
				&RightCurlyBracket{},
				&EOF{}},
		},
		{
			in: "[",
			out: []Token{
				&LeftSquareBracket{},
				&EOF{}},
		},
		{
			in: "]",
			out: []Token{
				&RightSquareBracket{},
				&EOF{}},
		},
		{
			in: ";",
			out: []Token{
				&Semicolon{},
				&EOF{}},
		},
		{
			in: ":",
			out: []Token{
				&Colon{},
				&EOF{}},
		},
		{
			in: ",",
			out: []Token{
				&Comma{},
				&EOF{}},
		},
		{
			in: "(",
			out: []Token{
				&LeftParenthesis{},
				&EOF{}},
		},
		{
			in: "))",
			out: []Token{
				&RightParenthesis{},
				&RightParenthesis{},
				&EOF{}},
		},
		{
			in: "{",
			out: []Token{
				&LeftCurlyBracket{},
				&EOF{}},
		},
		{
			in: "url",
			out: []Token{
				&Ident{Value: "url"},
				&EOF{}},
		},
		{
			in: "url( '",
			out: []Token{
				&Function{Value: "url"},
				&String{},
				&EOF{}},
		},
		{
			in: "url( \"",
			out: []Token{
				&Function{Value: "url"},
				&String{},
				&EOF{}},
		},
		{
			in: "url('",
			out: []Token{
				&Function{Value: "url"},
				&String{},
				&EOF{}},
		},
		{
			in: "function(",
			out: []Token{
				&Function{Value: "function"},
				&EOF{}},
		},
		{
			in: "url(asdf)",
			out: []Token{
				&URL{Value: "asdf"},
				&EOF{}},
		},
		{
			in: "url(asdf.com)",
			out: []Token{
				&URL{Value: "asdf.com"},
				&EOF{}},
		},
		{
			in: "url(asdf.com/moo.css)",
			out: []Token{
				&URL{Value: "asdf.com/moo.css"},
				&EOF{}},
		},
		{
			in: "url(\"asdf\")",
			out: []Token{
				&Function{Value: "url"},
				&String{Value: "asdf"},
				&RightParenthesis{},
				&EOF{}},
		},
		{
			in: "url(asdf)",
			out: []Token{
				&URL{Value: "asdf"},
				&EOF{}},
		},
		{
			in: "url('asdf')",
			out: []Token{
				&Function{Value: "url"},
				&String{Value: "asdf"},
				&RightParenthesis{},
				&EOF{}},
		},
		{
			in: "moo -->moo",
			out: []Token{
				&Ident{Value: "moo"},
				&WhiteSpace{Value: " "},
				&CDC{},
				&Ident{Value: "moo"},
				&EOF{}},
		},
		{
			in: "-",
			out: []Token{
				&Delim{Value: "-"},
				&EOF{}},
		},
		{
			in: "-123",
			out: []Token{
				&Digit{Value: -123},
				&EOF{}},
		},
		{
			in: "123%",
			out: []Token{
				&Percentage{Value: 123},
				&EOF{}},
		},
		{
			in: "123.123%",
			out: []Token{
				&Percentage{Value: 123.123},
				&EOF{}},
		},
		{
			in: ".",
			out: []Token{
				&Delim{Value: "."},
				&EOF{}},
		},
		{
			in: ".25",
			out: []Token{
				&Digit{Value: 0.25, Kind: Number},
				&EOF{}},
		},
		{
			in: "-123",
			out: []Token{
				&Digit{Value: -123},
				&EOF{}},
		},
		{
			in: "123",
			out: []Token{
				&Digit{Value: 123},
				&EOF{}},
		},
		/* Test ice chest. Move stubborn things in here while testing
		/* */
	}

	for _, test := range tests {
		s := New(strings.NewReader(test.in))
		got, err := s.ScanAll()
		if err != nil {
			t.Errorf("Error parsing input: %v", err)
		}

		if diff := cmp.Diff(test.out, got); diff != "" {
			t.Errorf("ScanAll(%q) = Diff (-want, +got):\n%s\nWant: %v\nGot:  %v", test.in, diff, test.out, got)
		}
	}
}

func TestScanEscape(t *testing.T) {
	tests := []struct {
		in  string
		out string
	}{
		{"", "\uFFFD"},
		{"AAAAAA", "AAAAAA"},
		{"FFGGGG", "FF"},
		{"1234567", "123456"},
	}

	for _, test := range tests {
		s := New(strings.NewReader(test.in))
		got := s.scanEscape()
		if got != test.out {
			t.Errorf("scanEscape(\"%s\") = \"%s\" want \"%s\"", test.in, got, test.out)
		}
	}
}

func TestPeek(t *testing.T) {
	tests := []struct {
		in  string
		out rune
	}{
		{"世界", '世'},
		{"FFGGGG", 'F'},
		{"1234567", '1'},
	}

	for _, test := range tests {
		s := New(strings.NewReader(test.in))
		got := s.peek()
		if got != test.out {
			t.Errorf("peek(\"%s\") = \"%c\" want \"%c\"", test.in, got, test.out)
		}
	}
}

func TestPeekOffset(t *testing.T) {
	tests := []struct {
		in  string
		n   int
		out rune
	}{
		{"世界", 0, '世'},
		{"世界", 1, '界'},
		{"FFGGGG", 3, 'G'},
		{"1234567", 2, '3'},
	}

	for _, test := range tests {
		s := New(strings.NewReader(test.in))
		got := s.peekOffset(test.n)
		if got != test.out {
			t.Errorf("over %q peekOffset(%d) = \"%c\" want \"%c\"", test.in, test.n, got, test.out)
		}

		// Now validate that the buffer wasn't ruined.
		buf := strings.Builder{}
		for {
			r := s.read()
			if r == eof {
				break
			}
			buf.WriteRune(r)
		}

		if ch := s.read(); ch != eof {
			t.Errorf("In testcase %q. Expected eof, got: %v", test.in, ch)
		}

		if got != test.out {
			t.Errorf("over %q peekOffset(%d) data was consumed.\nGot:  %q\nWant: %q", test.in, test.n, got, test.out)
		}
	}
}

func TestPeekN(t *testing.T) {
	tests := []struct {
		in  string
		n   int
		out string
	}{
		{"世界", 1, "世"},
		{"世界", 2, "世界"},
		{"FFGGGG", 2, "FF"},
		{"1234567", 3, "123"},
	}

	for _, test := range tests {
		s := New(strings.NewReader(test.in))
		got := s.peekN(test.n)
		if got != test.out {
			t.Errorf("over %q peekN(%d) = %+q want %+q", test.in, test.n, got, test.out)
		}

		// Now validate that the buffer wasn't ruined.
		buf := strings.Builder{}
		// Extract a rune ones for every rune in input.
		for _ = range test.in {
			r := s.read()
			if r == eof {
				break
			}
			buf.WriteRune(r)
		}

		if got != test.out {
			t.Errorf("over %q peekN(%d), data was consumed. Expected no buffer diff.\nGot:  %q\nWant: %q", test.in, test.n, got, test.out)
		}

	}
}

func TestCheckIfTwoCodePointsAreValidEscape(t *testing.T) {
	tests := []struct {
		in  string
		out bool
	}{
		{"世界", false},
		{"FFGGGG", false},
		{"1234567", false},
		{"\\n", true},
		{"\\A", true},
	}
	for _, test := range tests {
		s := New(strings.NewReader(test.in))
		if s.checkIfTwoCodePointsAreValidEscape() != test.out {
			t.Errorf("checkIfTwoCodePointsAreValidEscape(%q) == %v want %v", test.in, !test.out, test.out)
		}
	}
}

func TestCheckIfThreeCodePointsWouldStartAnIdentifier(t *testing.T) {
	tests := []struct {
		in  string
		out bool
	}{
		{"世界", true},
		{"FFGGGG", true},
		{"--1234567", false},
	}
	for _, test := range tests {
		s := New(strings.NewReader(test.in))
		if s.checkIfThreeCodePointsWouldStartAnIdentifier() != test.out {
			t.Errorf("checkIfThreeCodePointsWouldStartAnIdentifier(%q) == %v want %v", test.in, !test.out, test.out)
		}
	}
}

func TestCheckIfThreeCodePointsWouldStartANumber(t *testing.T) {
	tests := []struct {
		in  string
		out bool
	}{
		{"世界", false},
		{"世界", false},
		{".n", false},
		{"1.2", true},
		{"0.234567", true},
	}
	for _, test := range tests {
		s := New(strings.NewReader(test.in))
		if s.checkIfThreeCodePointsWouldStartANumber() != test.out {
			t.Errorf("checkIfThreeCodePointsWouldStartANumber(%q) == %v want %v", test.in, !test.out, test.out)
		}
	}
}

func TestScanNumber(t *testing.T) {
	tests := []struct {
		in   string
		out  float64
		kind NumberKind
	}{
		{"123", 123, Integer},
		{"123.4", 123.4, Number},
	}
	for _, test := range tests {
		s := New(strings.NewReader(test.in))
		got, gotKind := s.scanNumber()
		if got != test.out {
			t.Errorf("scanNumber(%q).value == %v want %v", test.in, got, test.out)
		}
		if gotKind != test.kind {
			t.Errorf("scanNumber(%q).kind == %v want %v", test.in, gotKind, test.kind)
		}

		if ch := s.read(); ch != eof {
			t.Errorf("Expected eof after calling scanNumber(%q). Got %c", test.in, ch)
		}
	}
}

func TestLargeInputs(t *testing.T) {
	found := false
	if err := filepath.Walk("testdata", func(path string, info os.FileInfo, err error) error {
		// Don't parse the root, or directories.
		if info == nil || info.IsDir() {
			return nil
		}

		// Don't parse files that don't end in .in.css
		if !strings.HasSuffix(path, testFileSuffix) {
			return nil
		}

		found = true

		t.Run(fmt.Sprintf("%s", path), func(t *testing.T) {
			walkTestData(t, path)
		})

		return nil
	}); err != nil {
		t.Errorf("Unable to walk testdata directory: %v", err)
	}

	if !found {
		t.Errorf("Didn't match any files")
	}

}

func walkTestData(t *testing.T, path string) {
	in, err := os.Open(path)
	if err != nil {
		t.Errorf("Unable to open input file: %v", err)
	}
	s := New(in)
	tokens, err := s.ScanAll()
	if err != nil {
		t.Errorf("Unable to scan tokens: %v", err)
	}
	buf := strings.Builder{}
	for _, token := range tokens {
		buf.WriteString(fmt.Sprintf("%s\n", token.String()))
		t.Logf("Token: %v", token)
	}

	wantPath := strings.TrimSuffix(path, testFileSuffix) + goldenSuffix
	wantBytes, err := ioutil.ReadFile(wantPath)
	if err != nil {
		t.Errorf("Unable to read output file: %v", err)
	}
	want := string(wantBytes)
	got := buf.String()
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("Tokenized output differed from expected (-want, +got):\n%s", diff)

		if *updateGoldens {
			if err := ioutil.WriteFile(wantPath, []byte(got), 0644); err != nil {
				t.Errorf("Unable to write output file: %v", err)
			}
		}
	}
}
