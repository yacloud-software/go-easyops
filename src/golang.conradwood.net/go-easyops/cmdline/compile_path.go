package cmdline

import (
	"runtime"
	"strings"
)

// returns the sourceode path from which main() was compiled into this binary
func SourceCodePath() string {
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
		if frame.Function != "main.main" {
			continue
		}
		//		fmt.Printf("Frame function: %s in %s\n", frame.Function, frame.File)
		if strings.Contains(frame.Function, "go-easyops") {
			continue
		}
		if strings.Contains(frame.File, "/opt/yacloud") { //internal go stuff
			continue
		}
		res = frame.File
		break
		//		fmt.Printf("- more:%v | %s\n", more, frame.Function)

	}
	return res
}
