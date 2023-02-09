package shared

import (
	"context"
	"flag"
	"fmt"
)

var (
	debug = flag.Bool("ge_debug_context", false, "if true print context debug stuff")
)

func Debugf(ctx context.Context, format string, args ...interface{}) {
	if !*debug {
		return
	}
	s := fmt.Sprintf(format, args...)
	fmt.Printf("[go-easyops] ctx-debug %s\n", s)
}
