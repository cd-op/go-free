package templates_test

import (
	"bytes"
	"fmt"
	"os"
	"testing"
	"testing/iotest"

	"cdop.pt/go/free/platepipe/templates"
	. "cdop.pt/go/open/assertive"
)

func TestHTMLTemplateFiles(t *testing.T) {
	t.Run("non existent html template", func(t *testing.T) {
		_, _, err := templates.FromFile("non-existent-file.html")

		Need(t, err != nil)
		Want(t, err.Error() ==
			"open non-existent-file.html: no such file or directory")
	})

	t.Run("bad html template", func(t *testing.T) {
		f := mkTestFile(t, "bad-html-template*.html", "<p>malformed {{")
		defer os.Remove(f)

		tpl, data, err := templates.FromFile(f)

		Need(t, err != nil)
		errmsg := err.Error()
		Want(t, errmsg[:10] == "template: ")
		Want(t, errmsg[len(errmsg)-15:] == "unclosed action")
		Want(t, len(data) == 0)
		Want(t, tpl == nil)
	})

	t.Run("good html template", func(t *testing.T) {
		f := mkTestFile(t, "good-html-template-*.html",
			"key = 'value'\n\n<p>text {{ .key }}</p>")
		defer os.Remove(f)

		tpl, data, err := templates.FromFile(f)

		Need(t, err == nil)
		Want(t, len(data) == 1)
		Want(t, tpl != nil)
	})
}

func TestHTMLTemplateStreams(t *testing.T) {
	t.Run("html stream error", func(t *testing.T) {
		msg := "error reading html template"
		r := iotest.ErrReader(fmt.Errorf(msg))

		_, _, err := templates.HTMLTemplateFromStream(r)

		Need(t, err != nil)
		Want(t, err.Error() == msg)
	})

	t.Run("html stream template error", func(t *testing.T) {
		r := bytes.NewBuffer([]byte("<p>malformed {{"))

		tpl, data, err := templates.HTMLTemplateFromStream(r)

		Need(t, err != nil)
		errmsg := err.Error()
		Want(t, errmsg[:10] == "template: ")
		Want(t, errmsg[len(errmsg)-15:] == "unclosed action")
		Want(t, len(data) == 0)
		Want(t, tpl == nil)
	})

	t.Run("html stream success", func(t *testing.T) {
		r := bytes.NewBuffer([]byte("key = 'value'\n\n<p>text {{ .key }}</p>"))

		tpl, data, err := templates.HTMLTemplateFromStream(r)

		Need(t, err == nil)
		Want(t, len(data) == 1)
		Want(t, tpl != nil)
	})
}

func TestTextTemplateFiles(t *testing.T) {
	t.Run("non existent text template", func(t *testing.T) {
		_, _, err := templates.FromFile("non-existent-file.txt")

		Need(t, err != nil)
		Want(t, err.Error() ==
			"open non-existent-file.txt: no such file or directory")
	})

	t.Run("bad text template", func(t *testing.T) {
		f := mkTestFile(t, "bad-text-template*.txt", "malformed {{")
		defer os.Remove(f)

		tpl, data, err := templates.FromFile(f)

		Need(t, err != nil)
		errmsg := err.Error()
		Want(t, errmsg[:10] == "template: ")
		Want(t, errmsg[len(errmsg)-15:] == "unclosed action")
		Want(t, len(data) == 0)
		Want(t, tpl == nil)
	})

	t.Run("good text template", func(t *testing.T) {
		f := mkTestFile(t, "good-text-template-*.txt", "key = 'value'\n\ntext {{ .key }}")
		defer os.Remove(f)

		tpl, data, err := templates.FromFile(f)

		Need(t, err == nil)
		Want(t, len(data) == 1)
		Want(t, tpl != nil)
	})
}

func TestTextTemplateStreams(t *testing.T) {
	t.Run("text stream error", func(t *testing.T) {
		msg := "error reading text template"

		r := iotest.ErrReader(fmt.Errorf(msg))

		_, _, err := templates.TextTemplateFromStream(r)

		Need(t, err != nil)
		Want(t, err.Error() == msg)
	})

	t.Run("text stream template error", func(t *testing.T) {
		r := bytes.NewBuffer([]byte("malformed {{"))

		tpl, data, err := templates.TextTemplateFromStream(r)

		Need(t, err != nil)
		errmsg := err.Error()
		Want(t, errmsg[:10] == "template: ")
		Want(t, errmsg[len(errmsg)-15:] == "unclosed action")
		Want(t, len(data) == 0)
		Want(t, tpl == nil)
	})

	t.Run("text stream success", func(t *testing.T) {
		r := bytes.NewBuffer([]byte("key = 'value'\n\ntext {{ .key }}"))

		tpl, data, err := templates.TextTemplateFromStream(r)

		Need(t, err == nil)
		Want(t, len(data) == 1)
		Want(t, tpl != nil)
	})
}

func TestTemplateApplication(t *testing.T) {
	t.Run("apply text template", func(t *testing.T) {
		f := mkTestFile(t, "template-*.txt",
			"tkey0 = 'v0'\ntkey1 = 1\n\nhere are all the keys:\n\n"+
				"template key 0 is {{.tkey0}}\n"+
				"template key 1 is {{.tkey1}}\n"+
				"var key 0 is {{.vkey0}}\n"+
				"var key 1 is {{.vkey1}}\n")
		defer os.Remove(f)

		tpl, data, err := templates.FromFile(f)

		Need(t, err == nil)
		Need(t, tpl != nil)

		data["vkey0"] = 2.43
		data["vkey1"] = []int{1, 2, 3}

		buf := bytes.NewBuffer([]byte{})
		err = tpl.Apply(buf, data)

		Need(t, err == nil)
		Want(t, buf.String() == "here are all the keys:\n\n"+
			"template key 0 is v0\n"+
			"template key 1 is 1\n"+
			"var key 0 is 2.43\n"+
			"var key 1 is [1 2 3]\n")
	})

	t.Run("apply html template", func(t *testing.T) {
		f := mkTestFile(t, "template-*.html",
			"tkey0 = 'v0'\ntkey1 = 1\n\n<p>here are all the keys:<ul>"+
				"<li>template key 0 is {{.tkey0}}</li>"+
				"<li>template key 1 is {{.tkey1}}</li>"+
				"<li>var key 0 is {{.vkey0}}</li>"+
				"<li>var key 1 is {{.vkey1}}</li>"+
				"</ul></p>")
		defer os.Remove(f)

		tpl, data, err := templates.FromFile(f)

		Need(t, err == nil)
		Need(t, tpl != nil)

		data["vkey0"] = 2.43
		data["vkey1"] = []int{1, 2, 3}

		buf := bytes.NewBuffer([]byte{})
		err = tpl.Apply(buf, data)

		Need(t, err == nil)
		Want(t, buf.String() == "<p>here are all the keys:"+
			"<ul><li>template key 0 is v0</li><li>template key 1 is 1</li>"+
			"<li>var key 0 is 2.43</li><li>var key 1 is [1 2 3]</li></ul></p>")
	})
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
