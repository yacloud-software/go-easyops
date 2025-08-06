package main

import (
	"context"
	"testing"
	"time"

	"golang.conradwood.net/go-easyops/errors"
	"golang.conradwood.net/go-easyops/utils/functionchain"
)

func TestFuncChain(t *testing.T) {
	fc := functionchain.NewFunctionChain()
	xt := &tester{
		t:                t,
		disable_fail_ctr: 5,
	}
	fc.Add(xt)
	xt.enable_fail_max = 5
	testfc(t, xt, true, true)
	testfc(t, xt, false, false)
}

type tester struct {
	t                 *testing.T
	enable_fail_max   int // how often to fail
	enable_fail_ctr   int // how often actually failed
	enable_requested  int
	enable_failed     int
	disable_fail_max  int // how often to fail
	disable_fail_ctr  int // how often actually failed
	disable_requested int
	disable_failed    int
}

func (ts *tester) SetTo(ctx context.Context, b bool) error {
	if b {
		ts.enable_requested++
		if ts.enable_fail_ctr < ts.enable_fail_max {
			ts.enable_fail_ctr++
			return errors.Errorf("enable fail")
		}
	} else {
		ts.disable_requested++
		if ts.disable_fail_ctr < ts.disable_fail_max {
			ts.disable_fail_ctr++
			return errors.Errorf("disable fail")
		}
	}
	return nil
}

func testfc(t *testing.T, xt *tester, enable, enable_result bool) {
	fc := functionchain.NewFunctionChain()
	fc.Add(xt)
	started := time.Now()
	err := fc.SetTo(context.Background(), enable)
	diff := time.Since(started)

	if enable_result && err != nil {
		// failed
		t.Errorf("Failed enable=%v after %0.1fs", enable, diff.Seconds())
		return

	}

	if !enable_result && err == nil {
		// failed
		t.Errorf("Expected to fail enable=%v, but it did not after %0.1fs", enable, diff.Seconds())
		return

	}

}
