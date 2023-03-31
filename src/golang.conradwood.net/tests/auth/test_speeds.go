package main

import (
	"fmt"
	"golang.conradwood.net/go-easyops/utils"
	"time"
)

func test_speed(name string, f func()) {
	// test authremote.Context() speed
	started := time.Now()
	p := &utils.ProgressReporter{}
	fmt.Printf("%s() speed:\n", name)
	total := 0.0
	ct := 0
	for {
		s1 := time.Now()
		f()
		total = total + time.Since(s1).Seconds()
		ct++
		p.Add(1)
		p.Print()
		if time.Since(started) > time.Duration(5)*time.Second {
			break
		}
	}
	avg := total / float64(ct)
	fmt.Printf("Average speed of %s: %0.2fs\n", name, avg)
}
