package utils

import (
	"fmt"
	"os"
	"sync"
)

type ByteStreamReceiver struct {
	sync.Mutex
	open_files map[string]*open_file
	last_file  *open_file
}

// the proto must be compatible with this interface
type StreamData interface {
	GetFilename() string
	GetData() []byte
}

func NewByteStreamReceiver() *ByteStreamReceiver {
	res := &ByteStreamReceiver{
		open_files: make(map[string]*open_file),
	}
	return res
}

// the result of srv.Recv()
func (bsr *ByteStreamReceiver) NewData(data StreamData) error {
	write_to := bsr.last_file
	if data.GetFilename() != "" {
		fmt.Printf("Receiving: \"%s\"\n", data.GetFilename())
		write_to = bsr.get_file_by_name(data.GetFilename())
		bsr.last_file = write_to
	}
	if write_to == nil {
		return fmt.Errorf("premature data received without filename")
	}
	b := data.GetData()
	err := write_to.Write(b)
	if err != nil {
		return err
	}
	return nil
}
func (bsr *ByteStreamReceiver) Close() error {
	bsr.Lock()
	defer bsr.Unlock()
	var err error
	for _, of := range bsr.open_files {
		xerr := of.Close()
		if xerr != nil {
			err = xerr
		}
	}
	return err
}

func (bsr *ByteStreamReceiver) get_file_by_name(name string) *open_file {
	bsr.Lock()
	defer bsr.Unlock()
	of, fd := bsr.open_files[name]
	if fd {
		return of
	}
	of = &open_file{filename: name}
	bsr.open_files[name] = of
	return of

}
func (bsr *ByteStreamReceiver) TotalBytesReceived() uint64 {
	bsr.Lock()
	defer bsr.Unlock()
	res := uint64(0)
	for _, of := range bsr.open_files {
		res = res + of.size
	}
	return res
}

type open_file struct {
	filename string
	size     uint64
	fd       *os.File
}

func (of *open_file) Write(buf []byte) error {
	if of.fd == nil {
		f, err := os.Create("/tmp/x/" + of.filename)
		if err != nil {
			return err
		}
		of.fd = f
	}
	of.size = of.size + uint64(len(buf))
	n, err := of.fd.Write(buf)
	if n != len(buf) {
		return fmt.Errorf("short write")
	}
	if err != nil {
		return err
	}
	return nil
}
func (of *open_file) Close() error {
	if of.fd != nil {
		err := of.fd.Close()
		of.fd = nil
		return err
	}
	return nil
}
