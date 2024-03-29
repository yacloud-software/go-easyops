package main

import (
	"bytes"
	"fmt"
	"golang.conradwood.net/go-easyops/cmdline"
	"golang.conradwood.net/go-easyops/errors"
	"golang.conradwood.net/go-easyops/utils"
	"html/template"
	"io"
	"os"
	"sort"
	"strings"
	"sync"
)

var (
	old_stdout *os.File
	testctr    = 0
	tests      []*test
	newidlock  sync.Mutex
)

type test struct {
	err           error
	id            int
	prefix        string
	dc_start      bool
	dc_error      bool
	builder_start int
	builder_error int
	stdout_writer io.Writer
	stdout_buf    *bytes.Buffer
}

func NewTest(format string, args ...interface{}) *test {
	t := &test{
		id:            newid(),
		prefix:        fmt.Sprintf(format, args...),
		dc_start:      cmdline.Datacenter(),
		builder_start: cmdline.GetContextBuilderVersion(),
		stdout_buf:    &bytes.Buffer{},
	}
	t.builder_error = t.builder_start
	t.dc_error = t.dc_start
	if old_stdout == nil {
		old_stdout = os.Stdout
	}
	r, w, err := os.Pipe()
	utils.Bail("failed to open pipe for stdout", err)
	wrprefix := fmt.Sprintf("TEST %s ", t.Prefix())
	t.stdout_writer = NewTee(t.stdout_buf, old_stdout, wrprefix)
	os.Stdout = w
	go t.pipe_reader(r)
	fmt.Printf("%s -------- STARTING\n", t.Prefix())
	tests = append(tests, t)
	return t
}
func newid() int {
	newidlock.Lock()
	testctr++
	newid := testctr
	newidlock.Unlock()
	return newid
}
func (t *test) Prefix() string {
	v := fmt.Sprintf("%v", t.builder_start)
	d := fmt.Sprintf("%v", t.dc_start)
	return fmt.Sprintf("[#%02d dc=%5s %s (builder=%5s)]", t.id, d, t.prefix, v)
}

func (t *test) Printf(format string, args ...interface{}) {
	fmt.Printf(t.Prefix()+" "+format, args...)
}
func (t *test) Error(err error) {
	if err == nil {
		return
	}
	if t.err != nil {
		return
	}
	t.dc_error = cmdline.Datacenter()
	t.builder_error = cmdline.GetContextBuilderVersion()

	t.err = err
	fmt.Printf("%s Failed (%s)\n", t.Prefix(), err)
}
func (t *test) Getstdout() string {
	return t.stdout_buf.String()
}
func (t *test) Done() {
	if t.err != nil {
		fmt.Printf("%s -------- FAILURE\n", t.Prefix())
		return
	}
	fmt.Printf("%s -------- SUCCESS\n", t.Prefix())
}

func PrintResult() {
	os.Stdout = old_stdout
	failed := 0
	succeeded := 0
	sort.Slice(tests, func(i, j int) bool {
		if tests[i].prefix != tests[j].prefix {
			return tests[i].prefix < tests[j].prefix
		}
		if tests[i].builder_start != tests[j].builder_start {
			return tests[i].builder_start < tests[j].builder_start
		}
		if tests[i].builder_error != tests[j].builder_error {
			return tests[i].builder_error < tests[j].builder_error
		}
		return tests[i].id < tests[j].id
	})
	var failed_tests []*test
	for _, t := range tests {
		if t.err != nil {
			failed++
			failed_tests = append(failed_tests, t)
		} else {
			succeeded++
		}
	}
	if failed > 0 {
		fmt.Printf("List of failed tests:\n")
		for _, t := range failed_tests {
			fmt.Println(t.Getstdout())
		}
		ta := utils.Table{}
		//		ta.SetMaxLen(5, 30)
		ta.AddHeaders("name", "dc (start)", "dc (error)", "builder (start)", "builder (error)", "error", "long")
		for _, t := range failed_tests {
			s := utils.ErrorString(t.err) + "\n" + t.Getstdout()
			ge := errors.UnmarshalError(t.err)
			se := ge.MultilineError()
			ta.AddString(t.prefix)
			ta.AddBool(t.dc_start)
			ta.AddBool(t.dc_error)
			ta.AddInt(t.builder_start)
			ta.AddInt(t.builder_error)
			ta.AddString(se)
			ta.AddString(s)
			ta.NewRow()
		}
		fmt.Println(ta.ToPrettyString())
	}
	fmt.Printf("Overall Result: %d tests suceeded, %d tests failed\n", succeeded, failed)

	if failed > 0 {
		fmt.Printf("TESTS FAILED\n")
	}

	b, err := render_tests_to_html(tests)
	if err != nil {
		fmt.Printf("failed to render to html: %s", err)
	} else {
		fname := "/tmp/tests.html"
		err = utils.WriteFile(fname, b)
		utils.Bail("failed to write file", err)
		fmt.Printf("HTML written to %s\n", fname)
	}
}

func (t *test) pipe_reader(r *os.File) {
	io.Copy(t.stdout_writer, r)

}
func (t *test) BuilderStart() int {
	return t.builder_start
}
func (t *test) BuilderError() int {
	return t.builder_error
}
func (t *test) GetError() error {
	return t.err
}
func (t *test) ID() int {
	return t.id
}
func (t *test) Name() string {
	return t.prefix
}
func (t *test) HtmlErrorDetails() template.HTML {
	lines := strings.Split(t.Getstdout(), "\n")
	res := ""
	for _, l := range lines {
		res = res + l + "<br/>"
	}
	return template.HTML(res)
}
func (t *test) DCStart() bool {
	return t.dc_start
}
func (t *test) DCError() bool {
	return t.dc_error
}
