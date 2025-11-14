package main

import (
	"fmt"
	"os"

	"golang.conradwood.net/go-easyops/utils"
)

func TestMulti() {
	test_multi(500, "foo\n", "prefix", "prefixfoo\n")
	test_multi(500, "foo\nbar\n", "prefix", "prefixfoo\nprefixbar\n")
	test_multi(500, "foo\nbar\n\n", "prefix", "prefixfoo\nprefixbar\n")
	test_multi(500, "\nfoo\nbar\n\n", "prefix", "prefix\nprefixfoo\nprefixbar\n")
	test_multi(3, "\nf1of2of3o\nb1rb2r\n\n", "prefix", "prefix\nprefixf1o\nprefixf2o\nprefixf3o\nprefixb1r\nprefixb2r\n")
}

func test_multi(maxlen int, input, prefix, expected string) {
	x := utils.MultiLinePrefix(input, prefix, maxlen)
	if x != expected {
		fmt.Printf("For input:\n<%s>\nexpected:\n<%s>\nbut got:\n<%s>\n", input, expected, x)
		os.Exit(10)
	}
}
