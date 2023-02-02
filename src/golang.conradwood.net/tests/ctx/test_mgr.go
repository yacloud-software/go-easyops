package main

import (
	"fmt"
	"golang.conradwood.net/go-easyops/cmdline"
)

var (
	tests []*test
)

type test struct {
	err    error
	prefix string
	dc     bool
}

func NewTest(format string, args ...interface{}) *test {
	t := &test{
		prefix: fmt.Sprintf(format, args...),
		dc:     cmdline.Datacenter(),
	}
	fmt.Printf("%s -------- STARTING\n", t.Prefix())
	tests = append(tests, t)
	return t
}

func (t *test) Prefix() string {
	v := fmt.Sprintf("%v", cmdline.ContextWithBuilder())
	d := fmt.Sprintf("%v", t.dc)
	return fmt.Sprintf("[dc=%5s %s (builder=%5s)]", d, t.prefix, v)
}

func (t *test) Printf(format string, args ...interface{}) {
	fmt.Printf(t.Prefix()+" "+format, args...)
}
func (t *test) Error(err error) {
	if err == nil {
		return
	}
	t.err = err
	fmt.Printf("%s Failed (%s)\n", t.Prefix(), err)
}
func (t *test) Done() {
	if t.err != nil {
		fmt.Printf("%s -------- FAILURE\n", t.Prefix())
		return
	}
	fmt.Printf("%s -------- SUCCESS\n", t.Prefix())
}

func PrintResult() {
	failed := 0
	succeeded := 0
	for _, t := range tests {
		if t.err != nil {
			failed++
		} else {
			succeeded++
		}
	}
	fmt.Printf("Overall Result: %d tests suceeded, %d tests failed\n", succeeded, failed)
	if failed > 0 {
		fmt.Printf("TESTS FAILED\n")
	}
}
