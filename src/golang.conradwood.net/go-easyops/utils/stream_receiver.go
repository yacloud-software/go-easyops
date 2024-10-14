package utils

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"sync"
)

type ByteStreamReceiver struct {
	sync.Mutex
	open_files      map[string]*open_file
	last_file       *open_file
	path            string
	custom_function func(filename string, content []byte) error
}

// the proto must be compatible with this interface
type StreamData interface {
	GetFilename() string
	GetData() []byte
}

// report each file to the function. This usually happens when Close() is called, prior to that the
// receiver has no means of telling if the file is completed. Perhaps a flaw in the protocol?
func NewByteStreamReceiverWithFunction(newfile func(filename string, content []byte) error) *ByteStreamReceiver {
	res := NewByteStreamReceiver("")
	res.custom_function = newfile
	return res
}

func NewByteStreamReceiver(path string) *ByteStreamReceiver {
	p, err := filepath.Abs(path)
	if err != nil {
		fmt.Printf("[go-easyops] byte-stream receiver failed filepath.Abs(%s): %s", path, err)
		return nil
	}
	for strings.HasSuffix(p, "/") {
		p = p[:len(p)-1]
	}
	res := &ByteStreamReceiver{
		path:       p,
		open_files: make(map[string]*open_file),
	}
	return res
}

// the result of srv.Recv()
func (bsr *ByteStreamReceiver) NewData(data StreamData) error {
	if data == nil || reflect.ValueOf(data).IsNil() {
		return nil
	}
	write_to := bsr.last_file
	if data.GetFilename() != "" {
		//		fmt.Printf("Receiving: \"%s\"\n", data.GetFilename())
		write_to = bsr.get_file_by_name(data.GetFilename())
		bsr.last_file = write_to
		err := write_to.Write(bsr.path, make([]byte, 0)) //create file
		if err != nil {
			return err
		}
	}
	if len(data.GetData()) == 0 {
		return nil
	}
	if write_to == nil {
		return fmt.Errorf("premature data received without filename")
	}
	b := data.GetData()
	err := write_to.Write(bsr.path, b)
	if err != nil {
		return err
	}
	return nil
}

// how many files were retrieved?
func (bsr *ByteStreamReceiver) FileCount() int {
	bsr.Lock()
	defer bsr.Unlock()
	return len(bsr.open_files)
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
	of = &open_file{bsr: bsr, filename: name, content: &bytes.Buffer{}}
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
	bsr      *ByteStreamReceiver
	filename string
	size     uint64
	fd       *os.File
	content  *bytes.Buffer
}

func (of *open_file) Write(path string, buf []byte) error {
	if of.bsr.custom_function != nil {
		_, err := of.content.Write(buf)
		if err != nil {
			return err
		}
	}
	if of.bsr.path != "" {
		if of.fd == nil {
			if strings.Contains(of.filename, "..") {
				return fmt.Errorf("Error: filename contains '..'")
			}
			os.MkdirAll(filepath.Dir(path+"/"+of.filename), 0777)
			f, err := os.Create(path + "/" + of.filename)
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
	}
	return nil
}
func (of *open_file) Close() error {
	cf := of.bsr.custom_function
	if cf != nil {
		err := cf(of.filename, of.content.Bytes())
		if err != nil {
			return err
		}
	}
	if of.fd != nil {
		err := of.fd.Close()
		of.fd = nil
		return err
	}
	return nil
}
