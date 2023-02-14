package main

import (
	"bytes"
	//	"fmt"
	"golang.conradwood.net/go-easyops/utils"
	"html/template"
	"os"
)

const (
	h_template = `<html>
<head>
<style>
.mytable {
  border: 1px;    
  outline: 1px solid #CCCCCC;
}
tr.desc {
    background-color: #F0F0F0;
    border-style: solid;
    border-colour: #000000;
    border-bottom: 1pt solid black;
    margin-top: 5px;
}
tr.details {
    background-color: #F0F0F0;
    border-style: solid;
    border-colour: #000000;
    border-bottom: 1pt solid black;
    margin-top: 5px;
   display:solid;
}
tr.invisible {
   display:none;
}
</style>

<script>
function togglevis(id) {
let tgt="details_"+id;
console.log("toggling ",tgt)
if (document.getElementById("chk_"+id).checked) {
document.getElementById(tgt).className="details";
} else {
document.getElementById(tgt).className="invisible";
}

}
</script>
</head>
<body>
<table class="mytable">
<tr class="desc">
<td>ID</td>
<td>Name</td>
<td>DC</td>
<td>Builder</td>
<td>Error</td>
</tr>
{{ range .Tests }}
<tr class="desc">
<td><input type="checkbox" id="chk_{{.ID}}" onclick="togglevis({{.ID}})">{{ .ID }}</td>
<td> {{ .Name }} </td>
<td> {{ .DCStart }} / {{ .DCError }} </td>
<td> {{ .BuilderStart }} / {{ .BuilderError }} </td>
<td> {{ .GetError }} </td>
</tr>
<tr id="details_{{.ID}}" class="invisible">
<td colspan="5"> {{ .HtmlErrorDetails }}</td>
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
func testrenderer_rendertest() {
	ts := []*test{
		NewTest("foo1"),
	}
	t := NewTest("foo2")
	t.id = 10
	t.stdout_buf = &bytes.Buffer{}
	t.stdout_buf.Write([]byte(`there is some stuff to see here
but not too much
`))
	ts = append(ts, t)
	for _, t := range ts {
		t.Done()
	}
	PrintResult()
	os.Exit(0)
}
