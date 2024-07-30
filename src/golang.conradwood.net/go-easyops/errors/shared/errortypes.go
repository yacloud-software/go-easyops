package shared

import (
	"fmt"
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
	stack ErrorStackTrace
	err   error
}

func NewMyError(err error, st ErrorStackTrace) *MyError {
	return &MyError{err: err, stack: st}
}
func NewWrappedError(err error, st ErrorStackTrace) *WrappedError {
	return &WrappedError{err: err, stack: st}
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

	s := fmt.Sprintf("%v at %s", we.err, first_non_internal_pos(we.Stack().Positions()).String())
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
