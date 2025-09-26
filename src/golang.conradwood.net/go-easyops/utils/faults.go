package utils

import (
	"context"
	"fmt"
	"sync"
	"time"

	"golang.yacloud.eu/apis/faultindicator"
	"google.golang.org/grpc"
)

var (
	fault_lock sync.Mutex
	// public as implementation detail
	Client_connector ClientConnector
	faultcc          faultindicator.FaultIndicatorClient
)

type ClientConnector interface {
	Connect(string) *grpc.ClientConn
}

// get the faultindicator client.
func GetFaultIndicatorClient() faultindicator.FaultIndicatorClient {
	fault_lock.Lock()
	defer fault_lock.Unlock()
	if faultcc != nil {
		return faultcc
	}
	cc := Client_connector
	if cc == nil {
		fmt.Printf("No fault indicator client!\n")
		time.Sleep(time.Duration(3) * time.Second)
		return nil
	}
	con := cc.Connect("faultindicator.FaultIndicator")
	faultcc = faultindicator.NewFaultIndicatorClient(con)
	return faultcc
}

// This logs a fault to the faultindicator service.
func LogFault(ctx context.Context, name, desc string) {

	_, err := GetFaultIndicatorClient().LogFault(ctx, &faultindicator.LogFaultRequest{
		Name:        name,
		Description: desc,
	})
	if err != nil {
		fmt.Printf("[go-easyops] --- failed to log fault: %s\n", ErrorString(err))
	}
}
