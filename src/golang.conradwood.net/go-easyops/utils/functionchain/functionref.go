package functionchain

import (
	"context"
	"sync"
)

type function_ref struct {
	sync.Mutex
	functions            []Function // one or more!
	last_state_requested bool
	last_run_err         error
}

func (fr *function_ref) Add(f Function) *function_ref {
	fr.Lock()
	defer fr.Unlock()
	fr.functions = append(fr.functions, f)
	return fr
}
func (fr *function_ref) SetTo(ctx context.Context, b bool) error {
	if fr.last_state_requested == b && fr.last_run_err == nil {
		return nil
	}
	fr.last_state_requested = b
	fr.last_run_err = nil

	for _, f := range fr.functions {
		err := f.SetTo(ctx, b)
		if err != nil {
			fr.last_run_err = err
		}
	}
	return nil
}
