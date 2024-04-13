// Package files provides utility functions over document/template files.
package files

import (
	"path/filepath"
	"strings"
)

// HasKnownMarkdownExt returns true if file has a known Markdown extension.
func HasKnownMarkdownExt(file string) bool {
	ext := strings.ToLower(filepath.Ext(file))

	switch ext {
	case ".md", ".markdown":
		return true
	}

	return false
}

// HasKnownHTMLExt returns true if file has a known HTML extension.
func HasKnownHTMLExt(file string) bool {
	ext := strings.ToLower(filepath.Ext(file))

	switch ext {
	case ".html", ".htm", ".xhtml", ".xhtm", ".xht":
		return true
	}

	return false
}
