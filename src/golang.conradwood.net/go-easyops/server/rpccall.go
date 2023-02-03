package server

import (
	"fmt"
	"time"
)

type rpccall struct {
	ServiceName string
	MethodName  string
	Started     time.Time
}

func (r *rpccall) FullMethod() string {
	return fmt.Sprintf("%s/%s", r.ServiceName, r.MethodName)
}
