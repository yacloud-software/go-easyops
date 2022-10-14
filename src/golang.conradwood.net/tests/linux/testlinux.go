package main

import (
	"fmt"
	l "golang.conradwood.net/go-easyops/linux"
	"golang.conradwood.net/go-easyops/utils"
)

func main() {
	lin := l.New()
	fmt.Printf("My IP: %s\n", lin.MyIP())
	//	TestExecuteContainer()
}

func TestExecuteContainer() {
	panic("no containers yet")
}

func copdir() {
	err := l.CopyDir("/tmp/x", "/tmp/y")
	utils.Bail("failed to copydir", err)
}
