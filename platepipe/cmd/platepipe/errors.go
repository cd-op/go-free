package main

import (
	"fmt"
	"os"
)

func usageError(msg string) {
	eprintln("%s: %s", progname, msg)
	usage()
	os.Exit(1)
}

func fail(msg string) {
	eprintln("%s: %s", progname, msg)
	os.Exit(2)
}

func eprintln(format string, args ...any) {
	fmt.Fprintf(os.Stderr, format, args...)
	fmt.Fprintln(os.Stderr)
}
