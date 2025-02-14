package linux

import (
	"fmt"
	"io"
)

type LinePrefixPipe struct {
	reader io.Reader
}

func NewLinePrefixPipe(r io.Reader) *LinePrefixPipe {
	return &LinePrefixPipe{reader: r}
}

func (l *LinePrefixPipe) Read([]byte) (int, error) {
	buf := make([]byte, 8192)
	var out_err error
	for {
		n, read_err := l.reader.Read(buf)
		if n != 0 {
			xerr := l.newBytes(buf[:n])
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
	return 0, nil
}

func (l *LinePrefixPipe) newBytes([]byte) error {
	return nil
}
