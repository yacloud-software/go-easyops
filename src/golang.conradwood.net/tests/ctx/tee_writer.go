package main

import (
	"io"
)

type tee struct {
	out1    io.Writer
	out2    io.Writer
	prefix2 string
}

func NewTee(out1, out2 io.Writer, prefix2 string) *tee {
	t := &tee{out1: out1, out2: out2, prefix2: prefix2}
	return t
}
func (t *tee) Write(buf []byte) (int, error) {
	tbuf := t.addPrefix(buf, t.prefix2)
	t.out2.Write(tbuf)
	n, err := t.out1.Write(buf)
	return n, err
}

// add a prefix to each line
func (t *tee) addPrefix(buf []byte, prefix string) []byte {
	if prefix == "" {
		return buf
	}
	pbuf := []byte(prefix)
	res := append(pbuf, buf...)
	return res
}
