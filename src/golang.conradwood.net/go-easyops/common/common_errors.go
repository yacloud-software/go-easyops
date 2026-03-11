package common

import (
	"flag"
	"fmt"
)

var (
	wrapper                    func(error) error
	new_error                  func(format string, args ...any) error
	err2str                    func(error) string
	err2strWithSingleLineStack func(error) string
	err2strWithStackTrace      func(error) string
	with_stack                 = flag.Bool("ge_err_with_stack", false, "if true print stacktrace on errors")
)

func Wrap(err error) error {
	if wrapper == nil {
		return err
	}
	return wrapper(err)
}
func Errorf(format string, args ...any) error {
	if new_error == nil {
		return fmt.Errorf(format, args...)
	}
	return new_error(format, args...)
}
func Err2Str(err error) string {
	if err2str == nil {
		return fmt.Sprintf("%s", err)
	}
	if *with_stack {
		return err2strWithStackTrace(err)
	}
	return err2str(err)
}
func RegisterErr(f1, f2, f3 func(error) string, f4 func(error) error, f5 func(format string, args ...any) error) {
	err2str = f1
	err2strWithSingleLineStack = f2
	err2strWithStackTrace = f3
	wrapper = f4
	new_error = f5
}
