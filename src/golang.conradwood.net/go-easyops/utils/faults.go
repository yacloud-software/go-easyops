package utils

import (
	"context"
	"fmt"

	"golang.yacloud.eu/apis/faultindicator"
	"google.golang.org/grpc"
)

var (
	// public as implementation detail
	Client_connector ClientConnector
)

type ClientConnector interface {
	Connect(string) *grpc.ClientConn
}

// This logs a fault to the faultindicator service.
func LogFault(ctx context.Context, name, desc string) {
	cc := Client_connector
	if cc == nil {
		fmt.Printf("[go-easyops] *********** FAULT LOG: %s -> %s\n", name, desc)
		return
	}
	con := cc.Connect("faultindicator.FaultIndicator")
	fic := faultindicator.NewFaultIndicatorClient(con)
	_, err := fic.LogFault(ctx, &faultindicator.LogFaultRequest{
		Name:        name,
		Description: desc,
	})
	if err != nil {
		fmt.Printf("[go-easyops] --- failed to log fault: %s\n", ErrorString(err))
	}
	con.Close()

}
