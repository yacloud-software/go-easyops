package main

import (
	"bytes"
	"golang.conradwood.net/go-easyops/utils"
	"html/template"
)

const (
	h_template = `<html><body>
<table>
{{ range .Tests }}
<tr>
<td> {{ .ID }} </td>
<td> {{ .Name }} </td>
<td> {{ .BuilderStart }} </td>
<td> {{ .BuilderError }} </td>
<td> {{ .GetError }} </td>
</tr>
{{ end }}
</table>
</body></html>
`
)

var (
	html_template *template.Template
)

func init() {
	tpl := template.New("foo")
	p, err := tpl.Parse(h_template)
	utils.Bail("failed to parse", err)
	html_template = p
}

type htmlrender struct {
	Tests []*test
}

func render_tests_to_html(tests []*test) ([]byte, error) {
	hr := &htmlrender{Tests: tests}
	buf := &bytes.Buffer{}
	err := html_template.Execute(buf, hr)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
