package main

import (
	"bytes"
	"fmt"
	"golang.conradwood.net/go-easyops/cmdline"
	"golang.conradwood.net/go-easyops/errors"
	"golang.conradwood.net/go-easyops/utils"
	"io"
	"os"
	"sort"
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
	builder_start bool
	builder_error bool
	stdout_writer io.Writer
	stdout_buf    *bytes.Buffer
}

func NewTest(format string, args ...interface{}) *test {
	t := &test{
		id:            newid(),
		prefix:        fmt.Sprintf(format, args...),
		dc_start:      cmdline.Datacenter(),
		builder_start: cmdline.ContextWithBuilder(),
		stdout_buf:    &bytes.Buffer{},
	}
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
	t.builder_error = cmdline.ContextWithBuilder()

	t.err = err
	fmt.Printf("%s Failed (%s)\n", t.Prefix(), err)
}
func (t *test) getstdout() string {
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
		return tests[i].prefix < tests[j].prefix
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
			fmt.Println(t.getstdout())
		}
		ta := utils.Table{}
		//		ta.SetMaxLen(5, 30)
		ta.AddHeaders("name", "dc (start)", "dc (error)", "builder (start)", "builder (error)", "error", "long")
		for _, t := range failed_tests {
			s := utils.ErrorString(t.err) + "\n" + t.getstdout()
			ge := errors.UnmarshalError(t.err)
			se := ge.MultilineError()
			ta.AddString(t.prefix)
			ta.AddBool(t.dc_start)
			ta.AddBool(t.dc_error)
			ta.AddBool(t.builder_start)
			ta.AddBool(t.builder_error)
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
}

func (t *test) pipe_reader(r *os.File) {
	io.Copy(t.stdout_writer, r)

}
