package utils

import (
	"fmt"
)

type send_data_type func(b []byte) error
type send_new_file_type func(key, filename string) error

// sends a bunch of bytes down a grpc stream
type ByteStreamSender struct {
	new_file  send_new_file_type
	send_data send_data_type
}

func NewByteStreamSender(f1 send_new_file_type, f2 send_data_type) *ByteStreamSender {
	res := &ByteStreamSender{
		new_file:  f1,
		send_data: f2,
	}
	return res
}
func (bss *ByteStreamSender) SendBytes(key, filename string, b []byte) error {
	err := bss.new_file(key, filename)
	if err != nil {
		return err
	}
	size := 8192
	offset := 0
	for {
		if size+offset > len(b) {
			size = len(b) - offset
		}
		if size == 0 {
			break
		}
		//bss.debugf("Sending %s [%d - %d]\n", filename, offset, size)
		err := bss.send_data(b[offset : offset+size])
		if err != nil {
			return err
		}
		offset = offset + size
	}
	return nil
}

func (bss *ByteStreamSender) debugf(format string, args ...interface{}) {
	fmt.Printf("[bss] "+format, args...)
}
