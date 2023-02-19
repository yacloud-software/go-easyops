package main

import (
	"fmt"
	"golang.conradwood.net/go-easyops/linux"
	"golang.conradwood.net/go-easyops/utils"
	"strings"
	"time"
)

func main() {
	allps, err := linux.AllPids()
	utils.Bail("failed to get pids", err)
	fmt.Printf("Got %d pids\n", len(allps))
	for _, ps := range allps {
		fmt.Printf("Pid: %s, Parent: %d\n", ps, ps.ParentPid())
	}

	ps := linux.PidStatus(1)
	printTreeOf(ps)
	fmt.Printf("Pidstatus: %s\n", ps)
	lin := linux.New()
	fmt.Printf("My IP: %s\n", lin.MyIP())
	run([]string{"/bin/true"})
	check_with_duration(time.Duration(6)*time.Second, []string{"sleep", "3"})
	check_with_duration(time.Duration(6)*time.Second, []string{"sleep", "100"})

	//time.Sleep(time.Duration(6) * time.Second)
	//	TestExecuteContainer()
}
func run(com []string) {
	fmt.Printf("executing \"%s\"...", strings.Join(com, " "))
	lin := linux.New()
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
	err := linux.CopyDir("/tmp/x", "/tmp/y")
	utils.Bail("failed to copydir", err)
}

func printTreeOf(ps *linux.ProcessState) {
	printTree(ps, " ")
}
func printTree(ps *linux.ProcessState, prefix string) {
	fmt.Printf("%s%s\n", prefix, ps)
	children, err := ps.Children()
	utils.Bail("failed to get children", err)
	prefix = prefix + "   "
	for _, c := range children {
		//fmt.Printf("%s%s\n", prefix, c)
		printTree(c, prefix)
	}
}
