package main

import (
	"fmt"
	"golang.conradwood.net/go-easyops/utils"
	"os"
	"time"
)

func TestTime() {
	test_one_time("2023-05-01", 1682899200)
	test_one_time("2023-01-15", 1673740800)
	test_one_time("14/05/2023", 1684022400)
	test_one_time("2023-05-01 16:03:34", 1682957014)
	test_one_time("2023-01-15 16:03:34", 1673798614)
	test_one_time("1684234688", 1684234688)
	test_one_time_loc("2023-05-15T14:35:00+00:00", "Europe/Paris", 1684161300)
	test_one_time_loc("2023-05-15T14:35:00+01:00", "Europe/Paris", 1684157700)
	test_one_time_loc("2023-05-15T14:35:00+00:00", "Europe/London", 1684161300)
	test_one_time_loc("2023-05-15T14:35:00", "Europe/London", 1684157700)
	test_one_time_loc("2023-05-15T14:35:00", "Europe/Paris", 1684154100)

}

func test_one_time_loc(ts, loc string, exp uint32) {
	t, err := utils.ParseTimestampWithLocation(ts, loc)
	if err != nil {
		fmt.Printf("Parsing \"%s\" failed: %s\n", ts, err)
		os.Exit(10)
	}
	comp_time(t, ts, exp)
}
func test_one_time(ts string, exp uint32) {
	t, err := utils.ParseTimestamp(ts)
	if err != nil {
		fmt.Printf("Parsing \"%s\" failed: %s\n", ts, err)
		os.Exit(10)
	}
	comp_time(t, ts, exp)
}

func comp_time(t uint32, ts string, exp uint32) {
	if t != exp {
		fmt.Printf("Parsing \"%s\" was expected to result in %d, but got %d\n", ts, exp, t)
		os.Exit(10)
	}
	nt := time.Unix(int64(t), 0)
	nts := fmt.Sprintf("DD/MM/YYYY %d/%d/%d", nt.Day(), nt.Month(), nt.Year())
	xts := nt.Format(time.RFC850)
	fmt.Printf("Determined that %v == %s == %s == %s\n", ts, utils.TimestampString(t), nts, xts)
}
