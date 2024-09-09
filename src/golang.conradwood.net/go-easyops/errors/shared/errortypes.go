package shared

import (
	"fmt"
	"google.golang.org/grpc/status"
)

type StackError interface {
	Stack() ErrorStackTrace
}
type ErrorStackTrace interface {
	Positions() []*StackPos
}

type MyError struct {
	stack ErrorStackTrace
	err   error
}

type WrappedError struct {
	stack        ErrorStackTrace
	err          error
	extramessage string
}

func NewMyError(err error, st ErrorStackTrace) *MyError {
	return &MyError{err: err, stack: st}
}
func NewWrappedError(err error, st ErrorStackTrace) *WrappedError {
	return &WrappedError{err: err, stack: st}
}
func NewWrappedErrorWithString(err error, st ErrorStackTrace, s string) *WrappedError {
	return &WrappedError{err: err, stack: st, extramessage: s}
}
func (me *MyError) Stack() ErrorStackTrace {
	return me.stack
}
func (me *MyError) Error() string {
	return me.String()
}
func (me *MyError) String() string {
	return fmt.Sprintf("%s at %s", me.err.Error(), first_non_internal_pos(me.Stack().Positions()).String())
}

func (me *MyError) GRPCStatus() *status.Status {
	e, _ := status.FromError(me.err)
	return e
}
func (me *WrappedError) GRPCStatus() *status.Status {
	e, _ := status.FromError(me.err)
	return e
}

func (we *WrappedError) Stack() ErrorStackTrace {
	return we.stack
}
func (we *WrappedError) Error() string {
	return we.String()
}
func (we *WrappedError) String() string {
	if we == nil {
		return "[noerror]"
	}

	x := ""
	if we.extramessage != "" {
		x = we.extramessage + ": "
	}
	s := fmt.Sprintf("%s%v at %s", x, we.err, first_non_internal_pos(we.Stack().Positions()).String())
	wo, ok := we.err.(*WrappedError)
	if ok {
		s = wo.String()
	}

	me, ok := we.err.(*MyError)
	if ok {
		s = me.Error()
	}

	return s
}
