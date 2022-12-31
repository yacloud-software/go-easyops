package utils

import (
	"fmt"
	"runtime/debug"
)

func NotImpl(format string, args ...interface{}) {
	s := fmt.Sprintf(format, args...)
	fmt.Printf("[goeasyops] Not Implemented: %s\n", s)
	debug.PrintStack()
}
