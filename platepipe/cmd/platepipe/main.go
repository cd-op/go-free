package main

import (
	"flag"
	"io"
	"os"
	"path"
	"strings"

	"cdop.pt/go/free/platepipe/documents/markdown"
	"cdop.pt/go/free/platepipe/variables"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
)

var progname = path.Base(os.Args[0])

func init() {
	flag.Usage = usage

	opts := getEnabledGoldmarkExtensions()
	gm := goldmark.New(opts...)
	markdown.ToHTML = func(buf []byte, w io.Writer) error {
		return gm.Convert(buf, w)
	}
}

func getEnabledGoldmarkExtensions() []goldmark.Option {
	value, defined := os.LookupEnv("PLATEPIPE_GOLDMARK_EXTS")
	if !defined {
		return []goldmark.Option{}
	}

	valparts := strings.Split(value, ",")
	exts := []goldmark.Extender{}

	for _, v := range valparts {
		switch strings.ToLower(v) {
		case "table":
			exts = append(exts, extension.Table)
		case "strikethrough":
			exts = append(exts, extension.Strikethrough)
		case "linkify":
			exts = append(exts, extension.Linkify)
		case "tasklist":
			exts = append(exts, extension.TaskList)
		case "gfm":
			exts = append(exts, extension.GFM)
		case "definitionlist":
			exts = append(exts, extension.DefinitionList)
		case "footnote":
			exts = append(exts, extension.Footnote)
		case "typographer":
			exts = append(exts, extension.Typographer)
		case "cjk":
			exts = append(exts, extension.CJK)
		}
	}

	return []goldmark.Option{goldmark.WithExtensions(exts...)}
}

func main() {
	opts := &options{}
	args := opts.Parse()
	argc := len(args)

	if opts.help {
		help()
	}

	if argc < 1 {
		usageError("no document specified")
	}

	if argc < 2 {
		usageError("no templates specified")
	}

	documentBuffer, documentMetadata, htmlSafe := loadDocumentAndMetadata(
		args[0],
		opts.docFmt,
	)

	templateChain, metadataChain := loadTemplateAndMetadataChains(
		args[1:],
		opts.tplFmt,
	)

	dataOverrides := loadVariables(opts.vOverrides)
	dataDefaults := loadVariables(opts.vDefaults)

	data := variables.Coalesce(
		programMetadata(args),
		dataOverrides,
		documentMetadata,
		variables.Coalesce(metadataChain...),
		dataDefaults,
	)

	runTemplatePipeline(string(documentBuffer), htmlSafe, templateChain, data)
}

type options struct {
	help       bool
	vDefaults  string
	vOverrides string
	docFmt     string
	tplFmt     string
}

func (opts *options) Parse() []string {
	flag.BoolVar(&opts.help, "h", false, "show this help")
	flag.StringVar(&opts.vDefaults, "vd", "", "variable defaults, metadata variables from this file will be used if not defined anywhere in the rendering pipeline")
	flag.StringVar(&opts.vOverrides, "vo", "", "variable overrides, metadata variables from this file will supersede variables from the rendering pipeline")
	flag.StringVar(&opts.docFmt, "df", "", `document format, "txt" or "md", default: autodetect (txt for stdin)`)
	flag.StringVar(&opts.tplFmt, "tf", "", `template format, "txt" or "html", default: autodetect (txt for stdin)`)

	flag.Parse()

	return flag.Args()
}

func help() {
	eprintln(`Usage: %s [OPTION]... [DOCUMENT] [TEMPLATE]...

When [DOCUMENT] is -, read document from standard input and assume`+
		` plaintext format with TOML metadata header

OPTIONS:`,
		progname)

	flag.PrintDefaults()

	eprintln(`
EXAMPLES:
  %[1]s doc.txt template.txt
    	render doc.txt through template.txt, use variables from headers in that order

  %[1]s doc.txt template1.txt template2.txt ... templateN.txt
    	render doc.txt through all templates, in order, earlier header variables supercede later definitions of the same variable

  %[1]s doc.md template.html
    	render doc.md through template.html after converting Markdown to HTML

  %[1]s -df txt doc.md template.txt
    	do not convert doc.md to HTML before rendering

  %[1]s -tf txt doc.txt template.html
    	treat template.html as a plaintext template

  %[1]s -tf txt doc.txt template.html
    	treat template.html as a plaintext template

  %[1]s -tf txt doc.txt template.html
    	treat template.html as a plaintext template`,
		progname,
	)

	os.Exit(0)
}

func usage() {
	eprintln(`Usage: %s [OPTION]... [DOCUMENT] [TEMPLATE]...`, progname)
	eprintln("Try '%s -h' for more information.", progname)
}
