package main

import (
	"fmt"
	ge "golang.conradwood.net/apis/getestservice"
	"golang.conradwood.net/go-easyops/authremote"
	"golang.conradwood.net/go-easyops/cmdline"
	"sync"
	"time"
)

var (
	sleep_wg sync.WaitGroup
)

func sleepTests() {
	cmdline.SetContextBuilderVersion(CONTEXT_VERSION)
	sleep_wg.Add(1)
	go sleepTest1(time.Duration(11) * time.Second)
	time.Sleep(time.Duration(1) * time.Second) // give time to start test with global parameters

	sleep_wg.Add(1)
	go sleepTest1(time.Duration(20) * time.Second)
	time.Sleep(time.Duration(1) * time.Second) // give time to start test with global parameters

	cmdline.SetContextBuilderVersion(0)
	sleep_wg.Add(1)
	go sleepTest1(time.Duration(11) * time.Second)
	time.Sleep(time.Duration(1) * time.Second) // give time to start test with global parameters

	sleep_wg.Add(1)
	go sleepTest1(time.Duration(20) * time.Second)
	time.Sleep(time.Duration(1) * time.Second) // give time to start test with global parameters

	cmdline.SetContextBuilderVersion(0)
	sleep_wg.Add(1)
	go sleepTest2(time.Duration(20) * time.Second)
	time.Sleep(time.Duration(1) * time.Second) // give time to start test with global parameters

	sleep_wg.Add(1)
	go sleepTest2(time.Duration(20) * time.Second)
	time.Sleep(time.Duration(1) * time.Second) // give time to start test with global parameters

	cmdline.SetContextBuilderVersion(CONTEXT_VERSION)
	sleep_wg.Add(1)
	go sleepTest2(time.Duration(20) * time.Second)
	time.Sleep(time.Duration(1) * time.Second) // give time to start test with global parameters

	sleep_wg.Add(1)
	go sleepTest2(time.Duration(20) * time.Second)
	time.Sleep(time.Duration(1) * time.Second) // give time to start test with global parameters

	cmdline.SetContextBuilderVersion(CONTEXT_VERSION)
	sleep_wg.Add(1)
	go sleepTest3()
	time.Sleep(time.Duration(1) * time.Second) // give time to start test with global parameters

	cmdline.SetContextBuilderVersion(0)
	sleep_wg.Add(1)
	go sleepTest3()
	time.Sleep(time.Duration(1) * time.Second) // give time to start test with global parameters

	sleep_wg.Wait()
}
func waitForSleepTests() {
}

func sleepTest1(dur time.Duration) {
	t := NewTest(fmt.Sprintf("sleep test 1 for %0.2fs", dur.Seconds()))
	ctx := authremote.ContextWithTimeout(dur)
	sl_seecs := dur.Seconds() - 2.0
	_, err := ge.GetCtxTestClient().Sleep(ctx, &ge.SleepRequest{Seconds: sl_seecs})
	t.Error(err)
	t.Done()
	sleep_wg.Done()
}
func sleepTest2(dur time.Duration) {
	t := NewTest(fmt.Sprintf("sleep test 2 for %0.2fs", dur.Seconds()))
	sl_seecs := dur.Seconds() - 2.0
	ctx := authremote.ContextWithTimeout(dur)
	ctx = authremote.DerivedContextWithRouting(ctx, make(map[string]string), true)
	_, err := ge.GetCtxTestClient().Sleep(ctx, &ge.SleepRequest{Seconds: sl_seecs})
	t.Error(err)
	t.Done()
	sleep_wg.Done()
}

// check if it _actually_ times out
func sleepTest3() {
	sl_seecs := 20.0
	t := NewTest(fmt.Sprintf("sleep test 3 for %0.2fs", sl_seecs))
	ctxdur := time.Duration(5) * time.Second
	ctx := authremote.ContextWithTimeout(ctxdur)
	ctx = authremote.DerivedContextWithRouting(ctx, make(map[string]string), true)
	started := time.Now()
	_, err := ge.GetCtxTestClient().Sleep(ctx, &ge.SleepRequest{Seconds: sl_seecs})
	if err == nil {
		t.Error(fmt.Errorf("Context with %0.2fs and sleep of %0.2fs did not time out", ctxdur.Seconds(), sl_seecs))
	}
	dur := time.Since(started)
	deviation := perc_diff(dur.Seconds(), ctxdur.Seconds())
	if deviation > 10 {
		t.Error(fmt.Errorf("Context with %0.2fs, but timeout occured after %0.2fs (deviation %f%%)", ctxdur.Seconds(), dur.Seconds(), deviation))
	}
	t.Done()
	sleep_wg.Done()
}

func perc_diff(a, b float64) float64 {
	diff := a / b
	if a < b {
		diff = b / a
	}
	res := diff * 100
	if res > 100 {
		return res - 100
	}
	return 100 - res
}
