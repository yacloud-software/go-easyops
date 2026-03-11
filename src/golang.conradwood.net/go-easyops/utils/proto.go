package utils

import (
	"encoding/base64"

	cm "golang.conradwood.net/go-easyops/common"

	"github.com/golang/protobuf/proto"
)

func MarshalBytes(req proto.Message) ([]byte, error) {
	data, err := proto.Marshal(req)
	return data, cm.Wrap(err)
}
func UnmarshalBytes(pdata []byte, req proto.Message) error {
	err := proto.Unmarshal(pdata, req)
	return cm.Wrap(err)
}

// take a proto and convert it into a base64 string
func Marshal(req proto.Message) (string, error) {
	data, err := proto.Marshal(req)
	if err != nil {
		return "", cm.Wrap(err)
	}
	b64 := base64.StdEncoding.EncodeToString(data)
	return b64, nil
}

// take a base64 string and convert it into the proto
func Unmarshal(b64string string, req proto.Message) error {
	pdata, err := base64.StdEncoding.DecodeString(b64string)
	if err != nil {
		return cm.Wrap(err)
	}
	err = proto.Unmarshal(pdata, req)
	if err != nil {
		return cm.Wrap(err)
	}

	return nil
}

// write a proto to disk
func WriteProto(filename string, req proto.Message) error {
	data, err := proto.Marshal(req)
	if err != nil {
		return cm.Wrap(err)
	}
	err = WriteFile(filename, data)
	if err != nil {
		return cm.Wrap(err)
	}
	return nil
}

// read a proto from disk
func ReadProto(filename string, req proto.Message) error {
	b, err := ReadFile(filename)
	if err != nil {
		return cm.Wrap(err)
	}
	err = UnmarshalBytes(b, req)
	if err != nil {
		return cm.Wrap(err)
	}
	return nil
}
