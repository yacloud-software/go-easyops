package main

import (
	"fmt"
	l "golang.conradwood.net/go-easyops/linux"
	"golang.conradwood.net/go-easyops/utils"
	"strings"
	"time"
)

func main() {
	lin := l.New()
	fmt.Printf("My IP: %s\n", lin.MyIP())
	run([]string{"/bin/true"})
	run([]string{"sleep", "100"})
	time.Sleep(time.Duration(6) * time.Second)
	//	TestExecuteContainer()
}
func run(com []string) {
	fmt.Printf("executing \"%s\"...", strings.Join(com, " "))
	lin := l.New()
	out, err := lin.SafelyExecute(com, nil)
	if err != nil {
		fmt.Printf("Output:\n%s\n", out)
		utils.Bail("failed to execute", err)
	}
	fmt.Printf("Done\n")
}
func TestExecuteContainer() {
	panic("no containers yet")
}

func copdir() {
	err := l.CopyDir("/tmp/x", "/tmp/y")
	utils.Bail("failed to copydir", err)
}
