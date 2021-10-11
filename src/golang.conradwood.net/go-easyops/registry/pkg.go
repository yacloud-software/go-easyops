package registry

import (
	pb "golang.conradwood.net/apis/registry"
	"golang.conradwood.net/go-easyops/utils"
	"google.golang.org/grpc"
	"strings"
	"time"
)

const (
	CONST_CALL_TIMEOUT = 4
)

// get a registry client from a specific ip
func Client(ip string) (pb.RegistryClient, error) {
	if !strings.Contains(ip, ":") {
		ip = ip + ":5000"
	}
	conn, err := grpc.Dial(ip, grpc.WithInsecure(), grpc.WithTimeout(time.Duration(CONST_CALL_TIMEOUT)*time.Second))
	if err != nil {
		return nil, err
	}
	r := pb.NewRegistryClient(conn)
	return r, nil
}

// get a registry client from a specific ip, panic if it cannot
func ClientOrPanic(ip string) pb.RegistryClient {
	conn, err := grpc.Dial(ip, grpc.WithInsecure(), grpc.WithTimeout(time.Duration(CONST_CALL_TIMEOUT)*time.Second))
	utils.Bail("Failed to get registry client at ip "+ip, err)
	r := pb.NewRegistryClient(conn)
	return r
}
