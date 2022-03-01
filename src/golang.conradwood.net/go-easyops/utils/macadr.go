package utils

import (
	"fmt"
)

func Mac2Str(mac uint64) string {
	res := ""
	deli := ""
	for i := 0; i < 6; i++ {
		x := mac & 0xFF
		s := fmt.Sprintf("%s%02X", deli, x)
		res = res + s
		deli = ":"
		mac = mac >> 8
	}
	return res
}
