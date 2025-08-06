package functionchain

import (
	"context"
	"sync"
	"time"

	"golang.conradwood.net/go-easyops/errors"
)

/*
a "function chain" is a chain of idempotent functions that enable or disable something. The chain will attempt to keep ALL functions enabled or disabled. In other words, it attempts to keep all functions either disabled or enabled.
*/
type FunctionChain struct {
	sync.Mutex
	functions              []*function_ref
	signal_stop_to_enable  bool
	signal_stop_to_disable bool
}

type Function interface {
	SetTo(ctx context.Context, activate bool) error
}

func NewFunctionChain() *FunctionChain {
	return &FunctionChain{}
}

func (fc *FunctionChain) Add(f Function) *function_ref {
	fc.Lock()
	defer fc.Unlock()
	fr := &function_ref{functions: []Function{f}}
	fc.functions = append(fc.functions, fr)
	return fr
}

func (fc *FunctionChain) SetTo(ctx context.Context, b bool) error {
	fc.Lock()
	defer fc.Unlock()
	if len(fc.functions) == 0 {
		return nil
	}
	stop_signal := &fc.signal_stop_to_disable
	if b {
		stop_signal = &fc.signal_stop_to_enable
	}
	*stop_signal = false
	start_time := time.Now()
	end_time := start_time.Add(time.Duration(30) * time.Second) // max runtime
	repeat := true
	for repeat {
		if *stop_signal {
			break
		}
		if time.Now().After(end_time) {
			return errors.Errorf("timeout enabling functionchain (%0.1fs)", time.Since(start_time).Seconds())
		}
		repeat = false
		for _, frw := range fc.functions {
			if *stop_signal {
				break
			}
			frw.Lock()
			if time.Now().After(end_time) {
				frw.Unlock()
				return errors.Errorf("timeout enabling functionchain (%0.1fs)", time.Since(start_time).Seconds())
			}
			err := frw.SetTo(ctx, b)
			frw.Unlock()
			if err != nil {
				repeat = true
			}
		}
		if repeat {
			time.Sleep(time.Duration(100) * time.Millisecond) // wait a little before retry
		}

	}
	return nil
}
func (fc *FunctionChain) Enable(ctx context.Context) error {
	fc.signal_stop_to_disable = true
	return fc.SetTo(ctx, true)

}
func (fc *FunctionChain) Disable(ctx context.Context) error {
	fc.signal_stop_to_enable = true
	return fc.SetTo(ctx, false)
}
