package main

import (
	l "golang.conradwood.net/go-easyops/linux"
	"golang.conradwood.net/go-easyops/utils"
)

func main() {
	TestExecuteContainer()
}

func TestExecuteContainer() {
	panic("no containers yet")
}

func copdir() {
	err := l.CopyDir("/tmp/x", "/tmp/y")
	utils.Bail("failed to copydir", err)
}
