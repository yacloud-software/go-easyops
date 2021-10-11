package client

import (
	"context"
	os "golang.conradwood.net/apis/objectstore"
	"golang.conradwood.net/go-easyops/errors"
	"io"
	"time"
)

var (
	ostore os.ObjectStoreClient
)

func getostore() {
	if ostore != nil {
		return
	}
	ostore = os.NewObjectStoreClient(Connect("objectstore.ObjectStore"))
}

func PutWithID(ctx context.Context, key string, buf []byte) error {
	if key == "" {
		return errors.InvalidArgs(ctx, "missing key to store in objectstore", "missing key to store in objectstore")
	}
	return PutWithIDAndExpiry(ctx, key, buf, time.Time{})
}
func PutWithIDAndExpiry(ctx context.Context, key string, buf []byte, expiry time.Time) error {
	if key == "" {
		return errors.InvalidArgs(ctx, "missing key to store in objectstore", "missing key to store in objectstore")
	}
	getostore()
	stream, err := ostore.LPutWithID(ctx)
	if err != nil {
		return err
	}
	size := 8192
	repeat := true
	offset := 0
	for repeat {
		if offset+size > len(buf) {
			size = len(buf) - offset
			repeat = false
		}
		n := buf[offset : offset+size]
		offset = offset + size
		pwr := &os.PutWithIDRequest{ID: key, Content: n}
		if !expiry.IsZero() {
			pwr.Expiry = uint32(expiry.Unix())
		}
		err := stream.Send(pwr)
		if err != nil {
			return err
		}
	}
	_, err = stream.CloseAndRecv()
	if err != nil {
		return err
	}
	return err
}
func Get(ctx context.Context, key string) ([]byte, error) {
	if key == "" {
		return nil, errors.InvalidArgs(ctx, "missing key to retrieve from objectstore", "missing key to retrieve from objectstore")
	}
	getostore()
	gr := &os.GetRequest{ID: key}
	stream, err := ostore.LGet(ctx, gr)
	if err != nil {
		return nil, err
	}
	var buf []byte
	for {
		ct, err := stream.Recv()
		if err == nil {
			buf = append(buf, ct.Content...)
			continue
		}
		if err == io.EOF {
			break
		}
		return nil, err

	}
	return buf, nil
}
