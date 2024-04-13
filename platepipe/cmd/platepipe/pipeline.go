package main

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"os"
	"time"

	"cdop.pt/go/free/platepipe/documents"
	"cdop.pt/go/free/platepipe/documents/files"
	"cdop.pt/go/free/platepipe/metadata"
	"cdop.pt/go/free/platepipe/templates"
)

func loadDocumentAndMetadata(filePath, format string) (
	[]byte, map[string]any, bool,
) {
	var buf []byte
	var data map[string]any
	var htmlSafe bool
	var err error

	switch format {
	case "md":
		htmlSafe = true
		if filePath == "-" {
			buf, data, err = documents.FromMarkdownStream(os.Stdin)
		} else {
			buf, data, err = documents.FromMarkdownFile(filePath)
		}
	case "html":
		htmlSafe = true
		if filePath == "-" {
			buf, data, err = documents.FromTextStream(os.Stdin)
		} else {
			buf, data, err = documents.FromTextFile(filePath)
		}
	case "txt":
		htmlSafe = false
		if filePath == "-" {
			buf, data, err = documents.FromTextStream(os.Stdin)
		} else {
			buf, data, err = documents.FromTextFile(filePath)
		}
	case "": // autodetect file type, assume plain text for stdin
		if filePath == "-" {
			htmlSafe = false
			buf, data, err = documents.FromTextStream(os.Stdin)
		} else {
			htmlSafe = files.HasKnownHTMLExt(filePath) ||
				files.HasKnownMarkdownExt(filePath)
			buf, data, err = documents.FromFile(filePath)
		}
	default:
		usageError("unknown document format")
	}

	if err != nil {
		fail("error reading document: " + err.Error())
	}

	return buf, data, htmlSafe
}

func loadTemplateAndMetadataChains(
	filePaths []string, format string,
) (
	[]*templates.Template, []map[string]any,
) {
	var loader func(string) (*templates.Template, map[string]any, error)

	switch format {
	case "html":
		loader = templates.HTMLTemplateFromFile
	case "txt":
		loader = templates.TextTemplateFromFile
	case "": // autodetect
		loader = templates.FromFile
	default:
		usageError("unknown template format")
	}

	tchain := []*templates.Template{}
	dchain := []map[string]any{}

	for _, p := range filePaths {
		t, data, err := loader(p)
		if err != nil {
			fail("error loading template: " + err.Error())
		}

		tchain = append(tchain, t)
		dchain = append(dchain, data)
	}

	return tchain, dchain
}

func loadVariables(filePath string) map[string]any {
	if filePath == "" {
		return map[string]any{}
	}

	r, err := os.Open(filePath)
	if err != nil {
		fail("error loading variables: " + err.Error())
	}
	defer r.Close()

	buf, err := io.ReadAll(r)
	if err != nil {
		fail("error loading variables: " + err.Error())
	}

	ret, err := metadata.FromTomlBuffer(buf)
	if err != nil {
		fail("error loading variables: " + err.Error())
	}

	return ret
}

func programMetadata(args []string) map[string]any {
	dir, _ := os.Getwd()

	return map[string]any{
		"platepipe": map[string]any{
			"time":      time.Now(),
			"document":  args[0],
			"templates": args[1:],
			"directory": dir,
		},
	}
}

func runTemplatePipeline(
	doc string, safe bool, ts []*templates.Template, data map[string]any,
) {
	data["content"] = markSafeAsNeeded(doc, safe)

	buf := new(bytes.Buffer)
	for _, t := range ts {
		buf.Reset()

		err := t.Apply(buf, data)
		if err != nil {
			fail("error applying template: " + err.Error())
		}

		data["content"] = markSafeAsNeeded(buf.String(), safe)
	}

	fmt.Fprint(os.Stdout, data["content"])
}

func markSafeAsNeeded(s string, safe bool) any {
	if safe {
		return template.HTML(s)
	}

	return s
}
