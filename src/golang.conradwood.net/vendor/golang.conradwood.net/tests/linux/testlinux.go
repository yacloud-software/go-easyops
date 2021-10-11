package main

import (
	l "golang.conradwood.net/go-easyops/linux"
	"golang.conradwood.net/go-easyops/utils"
)

func main() {
	err := l.CopyDir("/tmp/x", "/tmp/y")
	utils.Bail("failed to copydir", err)
}
