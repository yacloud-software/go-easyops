package utils

import (
	"flag"
	"fmt"
)

var (
	debug_packetizer = flag.Bool("debug_packetizer", false, "debug mode")
)

func check_valid(start, escape, stop byte) error {
	if start == escape {
		return fmt.Errorf("start byte (0x%02X) must not be identical to escape byte", start)
	}
	if stop == escape {
		return fmt.Errorf("stop byte (0x%02X) must not be identical to escape byte", stop)
	}
	return nil
}
