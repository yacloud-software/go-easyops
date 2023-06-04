package main

import (
	"flag"
	"fmt"
	"golang.conradwood.net/go-easyops/utils"
	"time"
)

func main() {
	flag.Parse()
	fmt.Printf("h\n")
	testTimer([]uint32{30, 15, 10, 5})
}
func testTimer(secs []uint32) {
	var sd []time.Duration
	for _, s := range secs {
		sd = append(sd, time.Duration(s)*time.Second)
	}
	pt := utils.NewPeriodicTimer(sd, timerCallback)
	pt.Start()
	pt.Wait()
}
func timerCallback(pt *utils.PeriodicTimer, secsLapsed time.Duration) error {
	sc := time.Since(pt.LastStarted()).Seconds()
	desc := fmt.Sprintf("%v", pt.Secs())
	fmt.Printf("timercallback (%s) after %0.1fs (%0.1fs)\n", desc, sc, secsLapsed.Seconds())
	if (sc >= 9) && (sc < 17) {
		return fmt.Errorf("foo")
	}
	return nil
}
