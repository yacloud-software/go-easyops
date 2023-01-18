package linux

import (
	"golang.conradwood.net/go-easyops/utils"
)

// read the proc filename (filename must be relative to /proc, e.g. "uptime" or "net/route"
func readProc(filename string) ([]byte, error) {
	b, err := utils.ReadFile("/proc/" + filename)
	return b, err
}
