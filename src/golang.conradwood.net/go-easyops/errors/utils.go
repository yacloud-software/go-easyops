package errors

import (
	"fmt"
	"runtime"
	//	"runtime/debug"
	"golang.conradwood.net/go-easyops/errors/shared"
	"strings"
)

type stacktrace struct {
	frames    *runtime.Frames
	positions []*shared.StackPos
}

// returns a single line with the calling function immedialy preceding the function which invoked this one (copied from utils/not_impl.go)
func callingFunction() (string, *stacktrace) {
	pc := make([]uintptr, 128)
	num := runtime.Callers(0, pc)
	if num == 0 {
		return "[no caller]", nil
	}
	pc = pc[:num] // pass only valid pcs to runtime.CallersFrames
	frames := runtime.CallersFrames(pc)
	st := &stacktrace{frames: frames}
	res := "[unidentifiable caller]"
	more := true
	var frame runtime.Frame
	stop := false
	for {
		// Check whether there are more frames to process after this one.
		if !more {
			break
		}
		frame, more = frames.Next()
		pos := &shared.StackPos{Function: frame.Function, Line: frame.Line, Filename: frame.File}
		st.positions = append(st.positions, pos)
		//fmt.Printf("FUNCTION: %s:%d\n", frame.Function, frame.Line)
		// Process this frame.
		//
		// To keep this example's output stable
		// even if there are changes in the testing package,
		// stop unwinding when we leave package runtime.
		if strings.Contains(frame.File, "runtime/") {
			continue
		}
		if strings.Contains(frame.Function, "golang.conradwood.net/go-easyops/errors") {
			continue
		}

		name := frame.Function
		n := strings.LastIndex(name, ".")
		if n != -1 {
			name = name[n+1:] + "()"
		}
		fname := frame.File
		n = strings.LastIndex(fname, "/src/")
		if n != -1 {
			fname = fname[n+5:]
		}
		if !stop {
			res = fmt.Sprintf("%s in %s:%d", name, fname, frame.Line)
		}
		stop = true
		//		fmt.Printf("- more:%v | %s\n", more, frame.Function)

	}
	return res, st
}

func (st *stacktrace) Positions() []*shared.StackPos {
	return st.positions
}
