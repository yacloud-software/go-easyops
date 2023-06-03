package main

import (
	"fmt"
	"golang.conradwood.net/go-easyops/utils"
	"time"
)

func main() {
	fmt.Printf("h\n")
	testTimer([]uint32{30, 15, 10, 5})
}
func testTimer(secs []uint32) {
	pt := utils.NewPeriodicTimer(secs, timerCallback)
	pt.Start()
	pt.Wait()
}
func timerCallback(pt *utils.PeriodicTimer, secsLapsed uint32) error {
	sc := time.Since(pt.LastStarted()).Seconds()
	desc := fmt.Sprintf("%v", pt.Secs())
	fmt.Printf("timercallback (%s) after %0.1fs (%d)\n", desc, sc, secsLapsed)
	if (sc >= 9) && (sc < 17) {
		return fmt.Errorf("foo")
	}
	return nil
}
