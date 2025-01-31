package main

import (
	"fmt"
	"io"
)

type comDefaultReader struct {
	pipe            io.ReadCloser
	notused_printer *LinePrefixPipe
}

func newDefaultReader(pipe io.ReadCloser) *comDefaultReader {
	c := &comDefaultReader{pipe: pipe, notused_printer: NewLinePrefixPipe(pipe)}
	go c.read_loop()
	return nil
}
func (r *comDefaultReader) Read([]byte) (int, error) {
	return 0, nil
}
func (r *comDefaultReader) newBytes(buf []byte) error {
	s := string(buf)
	fmt.Print(s)
	return nil
}
func (r *comDefaultReader) read_loop() {
	buf := make([]byte, 8192)
	var out_err error
	for {
		n, read_err := r.pipe.Read(buf)
		if n != 0 {
			xerr := r.newBytes(buf[:n])
			if xerr != nil {
				out_err = xerr
				break
			}
		}
		if read_err != nil {
			if read_err != io.EOF {
				out_err = read_err
			}
			break
		}
	}
	if out_err != nil {
		fmt.Printf("Error reading: %s\n", out_err)
	}

}
