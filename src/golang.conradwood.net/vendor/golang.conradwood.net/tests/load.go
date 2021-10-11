package main

import (
	"fmt"
	"golang.conradwood.net/go-easyops/linux"
	"golang.conradwood.net/go-easyops/utils"
	"time"
)

func LoadClient() {
	la := linux.Loadavg{}
xloop:
	c, err := la.GetCPULoad(nil, nil)
	if err != nil {
		fmt.Printf("Fail: %s\n", err)
	} else {
		//	fmt.Printf("Sum=%d, IdleTime=%0.1f%%, Idle=%d,User=%d\n", c.Sum, c.IdleTime, c.Idle, c.User)
		fmt.Printf("[%20s] Load: CPUs %2d, Total: %0.2f, PerCPU: %0.2f, IdleTime: %0.2f%%\n", utils.TimeString(time.Now()), c.CPUCount, c.Avg1, c.PerCPU, c.IdleTime)
	}
	if *loop {
		time.Sleep(time.Duration(1000) * time.Millisecond)
		goto xloop
	}
	fmt.Printf("\n")
}
