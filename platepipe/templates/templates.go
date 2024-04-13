// Package templates provides procedures and wrapper types for working with
// Go's standard library template packages in a more straightforward way,
// as well as enabling templates to include metadata, which the standard
// library packages do not support.
//
// For details on the formats and processing of the metadata headers, see the
// documentation for the metadata package.
package templates

import (
	"crypto/md5"
	"fmt"
	"io"

	htemplate "html/template"
	ttemplate "text/template"

	"cdop.pt/go/free/platepipe/documents"
	"cdop.pt/go/free/platepipe/documents/files"
)

// FromFile loads a text or HTML template from a file whose path is passed as
// the argument.
//
// If the file's extension indicates that the file contains HTML, the
// template will be parsed by html/template. Otherwise, it will be parsed by
// text/template.
func FromFile(file string) (*Template, map[string]any, error) {
	if files.HasKnownHTMLExt(file) {
		return HTMLTemplateFromFile(file)
	}

	return TextTemplateFromFile(file)
}

// HTMLTemplateFromFile loads an html/template and its metadata (if any) from
// the given file.
func HTMLTemplateFromFile(file string) (*Template, map[string]any, error) {
	buf, data, err := documents.FromTextFile(file)
	if err != nil {
		return nil, nil, err
	}

	t, err := newHTMLTemplate(buf)
	if err != nil {
		return nil, nil, err
	}

	return t, data, nil
}

// HTMLTemplateFromStream loads an html/template and its metadata (if any) from
// the given io.Reader.
func HTMLTemplateFromStream(r io.Reader) (*Template, map[string]any, error) {
	buf, data, err := documents.FromTextStream(r)
	if err != nil {
		return nil, nil, err
	}

	t, err := newHTMLTemplate(buf)
	if err != nil {
		return nil, nil, err
	}

	return t, data, nil
}

// TextTemplateFromFile loads an html/template and its metadata (if any) from
// the given file.
func TextTemplateFromFile(file string) (*Template, map[string]any, error) {
	buf, data, err := documents.FromTextFile(file)
	if err != nil {
		return nil, nil, err
	}

	t, err := newTextTemplate(buf)
	if err != nil {
		return nil, nil, err
	}

	return t, data, nil
}

// TextTemplateFromStream loads an html/template and its metadata (if any) from
// the given io.Reader.
func TextTemplateFromStream(r io.Reader) (*Template, map[string]any, error) {
	buf, data, err := documents.FromTextStream(r)
	if err != nil {
		return nil, nil, err
	}

	t, err := newTextTemplate(buf)
	if err != nil {
		return nil, nil, err
	}

	return t, data, nil
}

func newHTMLTemplate(buf []byte) (*Template, error) {
	t, err := htemplate.New(hash(buf)).Parse(string(buf))
	if err != nil {
		return nil, err
	}

	return &Template{t}, nil
}

func newTextTemplate(buf []byte) (*Template, error) {
	t, err := ttemplate.New(hash(buf)).Parse(string(buf))
	if err != nil {
		return nil, err
	}

	return &Template{t}, nil
}

func hash(buf []byte) string {
	return fmt.Sprintf("%x", (md5.New().Sum(buf))[0:4])
}

// Template is a wrapper type around the template types provided by both
// text/template and html/template.
type Template struct {
	stdTemplate interface {
		Execute(io.Writer, any) error
	}
}

// Apply renders the template, with the provided data, to an io.Writer.
func (t *Template) Apply(w io.Writer, data map[string]any) error {
	return t.stdTemplate.Execute(w, data)
}
