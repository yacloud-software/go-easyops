package utils

import (
	"fmt"
	"runtime"
	"runtime/debug"
	"strings"
)

func NotImpl(format string, args ...interface{}) {
	s := fmt.Sprintf(format, args...)
	fmt.Printf("[goeasyops] Not Implemented: %s\n", s)
	debug.PrintStack()
}

// print a simpliefied stacktrace (filenames and linenumbers only)
func PrintStack(format string, args ...interface{}) {
	s := GetStack(format, args...)
	fmt.Print(s)

}

// get stack in a human readable format
func GetStack(format string, args ...interface{}) string {
	s := fmt.Sprintf("Stacktrace for: "+format+"\n", args...)
	pc := make([]uintptr, 128)
	num := runtime.Callers(0, pc)
	if num == 0 {
		return "[nostack]"
	}
	pc = pc[:num] // pass only valid pcs to runtime.CallersFrames
	frames := runtime.CallersFrames(pc)

	more := true
	var frame runtime.Frame
	ignore_functions := []string{
		"golang.conradwood.net/go-easyops/utils.PrintStack",
		"golang.conradwood.net/go-easyops/utils.GetStack",
	}
	for {
		// Check whether there are more frames to process after this one.
		if !more {
			break
		}
		frame, more = frames.Next()
		// Process this frame.
		//
		// To keep this example's output stable
		// even if there are changes in the testing package,
		// stop unwinding when we leave package runtime.
		if strings.Contains(frame.File, "runtime/") {
			continue
		}
		ign := false
		for _, ifu := range ignore_functions {
			if strings.Contains(frame.Function, ifu) {
				ign = true
				break
			}
		}
		if ign {
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
		s = s + fmt.Sprintf(" %s in %s:%d\n", name, fname, frame.Line)
	}
	s = s + "---end stacktrace\n"
	return s
}

// returns a single line with the calling function immedialy preceding the function which invoked this one
func CallingFunction() string {
	pc := make([]uintptr, 128)
	num := runtime.Callers(0, pc)
	if num == 0 {
		return "[no caller]"
	}
	pc = pc[:num] // pass only valid pcs to runtime.CallersFrames
	frames := runtime.CallersFrames(pc)
	res := "[unidentifiable caller]"
	more := true
	var frame runtime.Frame
	i := 0
	for {
		// Check whether there are more frames to process after this one.
		if !more {
			break
		}
		frame, more = frames.Next()
		// Process this frame.
		//
		// To keep this example's output stable
		// even if there are changes in the testing package,
		// stop unwinding when we leave package runtime.
		if strings.Contains(frame.File, "runtime/") {
			continue
		}
		if strings.Contains(frame.Function, "golang.conradwood.net/go-easyops/utils.CallingFunction") {
			continue
		}
		i++
		if i < 2 {
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
		res = fmt.Sprintf("%s in %s:%d", name, fname, frame.Line)
		break
		//		fmt.Printf("- more:%v | %s\n", more, frame.Function)

	}
	return res
}
