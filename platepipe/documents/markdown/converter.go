// Package markdown defines the name, signature and default implementation of
// the Markdown conversion procedure.
package markdown

import (
	"io"

	"github.com/yuin/goldmark"
)

// ToHTML is the default Markdown conversion procedure.
//
// It can be swapped with any other conversion procedure with the same
// signature. This allows the use of other Markdown conversion libraries.
//
// Changing the conversion procedure should be done as early as possible in the
// main program.
var ToHTML func([]byte, io.Writer) error

func init() {
	gm := goldmark.New()

	ToHTML = func(buf []byte, w io.Writer) error {
		return gm.Convert(buf, w)
	}
}
