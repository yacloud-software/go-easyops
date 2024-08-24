package utils

import (
	"fmt"
)

type Send_data_type func(b []byte) error
type Send_new_file_type func(key, filename string) error

// sends a bunch of bytes down a grpc stream
type ByteStreamSender struct {
	new_file     Send_new_file_type
	send_data    Send_data_type
	file_counter int
}

/*
create a new bytestream sender.
(key and filename are opaque to this sender)
f1 - a function that sends a message on the stream to indicate start of a new file and key
f2 - a function that sends a bunch of data on the stream

once stream sender is created, call SendBytes with filename and content.
it will break the file into small pieces and call f2() with small arrays suitable for sending in a packet
*/
func NewByteStreamSender(f1 Send_new_file_type, f2 Send_data_type) *ByteStreamSender {
	res := &ByteStreamSender{
		new_file:  f1,
		send_data: f2,
	}
	return res
}

// how many files were sent?
func (bss *ByteStreamSender) FileCount() int {
	return bss.file_counter
}

func (bss *ByteStreamSender) SendBytes(key, filename string, b []byte) error {
	err := bss.new_file(key, filename)
	if err != nil {
		return err
	}
	bss.file_counter++
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
