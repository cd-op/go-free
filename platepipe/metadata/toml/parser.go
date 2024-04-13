// Package toml defines the name, signature and default implementation of the
// TOML parser procedure.
package toml

import "github.com/BurntSushi/toml"

// Parse is the default parser procedure.
//
// It can be swapped with any other parser procedure with the same signature.
// This allows the use of other TOML parser libraries.
//
// Changing the parser procedure should be done as early as possible in the
// main program.
var Parse func([]byte, any) error

func init() {
	Parse = toml.Unmarshal
}
