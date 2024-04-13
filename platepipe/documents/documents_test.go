package documents_test

import (
	"fmt"
	"io"
	"os"
	"testing"
	"testing/iotest"

	"cdop.pt/go/free/platepipe/documents"
	"cdop.pt/go/free/platepipe/documents/markdown"
	. "cdop.pt/go/open/assertive"
)

func TestMarkdown(t *testing.T) {
	t.Run("no file", func(t *testing.T) {
		_, _, err := documents.FromFile("non-existent-file.md")

		Need(t, err != nil)
		Want(t, err.Error() ==
			"open non-existent-file.md: no such file or directory")
	})

	t.Run("markdown to html error", func(t *testing.T) {
		f := mkTestFile(t, "md-processing-error-*.md", "[link](broken")
		defer os.Remove(f)

		oldConv := markdown.ToHTML
		markdown.ToHTML = func([]byte, io.Writer) error {
			return fmt.Errorf("markdown: error processing buffer")
		}
		text, data, err := documents.FromFile(f)
		markdown.ToHTML = oldConv

		Need(t, err != nil)
		Want(t, err.Error() == "markdown: error processing buffer")
		Want(t, len(data) == 0)
		Want(t, len(text) == 0)
	})

	t.Run("success", func(t *testing.T) {
		f := mkTestFile(t, "success-*.md",
			"key = 'value'\n\n# header\n\nparagraph")
		defer os.Remove(f)

		text, data, err := documents.FromFile(f)

		Need(t, err == nil)
		Want(t, fmt.Sprint(data) == fmt.Sprint(map[string]any{"key": "value"}))
		Want(t, string(text) == "<h1>header</h1>\n<p>paragraph</p>\n")
	})
}

func TestText(t *testing.T) {
	t.Run("no file", func(t *testing.T) {
		_, _, err := documents.FromFile("non-existent-file.txt")

		Need(t, err != nil)
		Want(t, err.Error() ==
			"open non-existent-file.txt: no such file or directory")
	})

	t.Run("no metadata", func(t *testing.T) {
		f := mkTestFile(t, "no-metadata-*.txt", "no metadata in this file")
		defer os.Remove(f)

		text, data, err := documents.FromFile(f)

		Need(t, err == nil)
		Want(t, len(data) == 0)
		Want(t, string(text) == "no metadata in this file")
	})

	t.Run("metadata error", func(t *testing.T) {
		f := mkTestFile(t, "metadata-error-*.txt", "key: value\n\ntext")
		defer os.Remove(f)

		text, data, err := documents.FromFile(f)

		Need(t, err == nil)
		Want(t, len(data) == 0)
		Want(t, string(text) == "key: value\n\ntext")
	})

	t.Run("success", func(t *testing.T) {
		f := mkTestFile(t, "success-*.txt", "key = 'value'\n\ntext")
		defer os.Remove(f)

		text, data, err := documents.FromFile(f)

		Need(t, err == nil)
		Want(t, fmt.Sprint(data) == fmt.Sprint(map[string]any{"key": "value"}))
		Want(t, string(text) == "text")
	})
}

func TestReadErrors(t *testing.T) {
	msg := "error reading document"
	r := iotest.ErrReader(fmt.Errorf(msg))

	_, _, err := documents.FromMarkdownStream(r)

	Need(t, err != nil)
	Want(t, err.Error() == msg)
}

func mkTestFile(t *testing.T, pattern, content string) string {
	f, err := os.CreateTemp("", pattern)
	Need(t, err == nil)

	name := f.Name()

	_, err = f.WriteString(content)
	Need(t, err == nil)

	err = f.Sync()
	Need(t, err == nil)

	err = f.Close()
	Need(t, err == nil)

	return name
}
