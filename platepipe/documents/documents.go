// Package documents defines procedures to load documents from files or
// io.Reader instances and optionally convert them from Markdown to HTML. These
// documents may have metadata headers.
//
// This package treats metadata blocks that fail to parse correctly as document
// content. For additional details on the formats and processing of the
// metadata headers, see the documentation for the metadata package.
package documents

import (
	"bytes"
	"io"
	"os"

	"cdop.pt/go/free/platepipe/documents/files"
	"cdop.pt/go/free/platepipe/documents/markdown"
	"cdop.pt/go/free/platepipe/metadata"
)

// FromFile loads a text or Markdown document from a file whose path is passed
// as the argument.
//
// If the file's extension indicates that the file contains Markdown, the
// content will be converted to HTML. Otherwise, no conversion is made.
func FromFile(file string) ([]byte, map[string]any, error) {
	if files.HasKnownMarkdownExt(file) {
		return FromMarkdownFile(file)
	}

	return FromTextFile(file)
}

// FromMarkdownFile loads content/metadata from the given file and returns
// the Markdown content converted to HTML.
func FromMarkdownFile(file string) ([]byte, map[string]any, error) {
	r, err := os.Open(file)
	if err != nil {
		return []byte{}, map[string]any{}, err
	}
	defer r.Close()

	return FromMarkdownStream(r)
}

// FromMarkdownStream loads content/metadata from the given io.Reader and
// returns the Markdown content converted to HTML.
func FromMarkdownStream(r io.Reader) ([]byte, map[string]any, error) {
	buf, data, err := FromTextStream(r)
	if err != nil {
		return []byte{}, map[string]any{}, err
	}

	var html bytes.Buffer
	err = markdown.ToHTML(buf, &html)
	if err != nil {
		return []byte{}, map[string]any{}, err
	}

	return html.Bytes(), data, err
}

// FromTextFile loads content/metadata from the given file. No content
// conversion is made.
func FromTextFile(file string) ([]byte, map[string]any, error) {
	r, err := os.Open(file)
	if err != nil {
		return []byte{}, map[string]any{}, err
	}
	defer r.Close()

	return FromTextStream(r)
}

// FromTextStream loads content/metadata from the given io.Reader. No content
// conversion is made.
func FromTextStream(r io.Reader) ([]byte, map[string]any, error) {
	buf, err := io.ReadAll(r)
	if err != nil {
		return []byte{}, map[string]any{}, err
	}

	present, splitPos := metadata.IsPresent(buf)
	if !present {
		return buf, map[string]any{}, nil
	}

	data, err := metadata.FromTomlBuffer(buf[:splitPos])
	if err != nil {
		// treat parse error as document content, see metadata package doc
		return buf, map[string]any{}, nil
	}

	return buf[splitPos:], data, nil
}
