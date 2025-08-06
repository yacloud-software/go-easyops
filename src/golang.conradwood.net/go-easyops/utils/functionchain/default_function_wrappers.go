package functionchain

import "context"

// shortcut: add a pair of enable/disable functions
func (fc *FunctionChain) AddFuncs(enable, disable func(context.Context) error) *function_ref {
	frw := &function_ref_wrapper{
		enable_function:  enable,
		disable_function: disable,
	}
	return fc.Add(frw)
}

type function_ref_wrapper struct {
	enable_function  func(ctx context.Context) error
	disable_function func(ctx context.Context) error
}

func (frw *function_ref_wrapper) SetTo(ctx context.Context, b bool) error {
	if b {
		return frw.enable_function(ctx)
	}
	return frw.disable_function(ctx)
}

/*
 ***********************************************************************************************
 */

type function_ref_setter_wrapper struct {
	setter func(context.Context, bool) error
}

// shortcut: add a pair of enable/disable functions
func (fc *FunctionChain) AddSetter(setter func(context.Context, bool) error) *function_ref {
	frw := &function_ref_setter_wrapper{
		setter: setter,
	}
	return fc.Add(frw)
}
func (frw *function_ref_setter_wrapper) SetTo(ctx context.Context, b bool) error {
	return frw.setter(ctx, b)
}
