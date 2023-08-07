package main

import (
	"flag"
	"fmt"
	"golang.conradwood.net/go-easyops/linux"
	"golang.conradwood.net/go-easyops/utils"
	"os"
	"strings"
	"time"
)

var (
	ds = flag.String("dirsize", "", "if set, calc dirsize")
)

func main() {
	flag.Parse()
	if *ds != "" {
		size, err := linux.DirSize(*ds)
		utils.Bail("failed to get dirsize", err)
		fmt.Printf("Dirsize of \"%s\": %s\n", *ds, utils.PrettyNumber(size))
		os.Exit(0)
	}
	allps, err := linux.AllPids()
	utils.Bail("failed to get pids", err)
	fmt.Printf("Got %d pids\n", len(allps))
	for _, ps := range allps {
		fmt.Printf("Pid: %s, Parent: %d, Cgroup: %s\n", ps, ps.ParentPid(), ps.Cgroup())
	}

	ps := linux.PidStatus(1)
	printTreeOf(ps)
	fmt.Printf("Pidstatus: %s\n", ps)
	lin := linux.New()
	fmt.Printf("My IP: %s\n", lin.MyIP())
	run([]string{"/bin/true"})
	check_with_duration(time.Duration(6)*time.Second, []string{"sleep", "3"})
	check_with_duration(time.Duration(6)*time.Second, []string{"sleep", "100"})

	run([]string{"sleep", "300"})
	//time.Sleep(time.Duration(6) * time.Second)
	//	TestExecuteContainer()
}
func run(com []string) {
	fmt.Printf("executing \"%s\"...", strings.Join(com, " "))
	lin := linux.New()
	started := time.Now()
	out, err := lin.SafelyExecute(com, nil)
	if err != nil {
		fmt.Printf("Output:\n%s\n", out)
		utils.Bail("failed to execute", err)
	}
	fmt.Printf("Done (%0.1fs)\n", time.Since(started).Seconds())
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
	fmt.Printf("%s%s (Cgroup: \"%s\")\n", prefix, ps, ps.Cgroup())
	children, err := ps.Children()
	utils.Bail("failed to get children", err)
	prefix = prefix + "   "
	for _, c := range children {
		//fmt.Printf("%s%s\n", prefix, c)
		printTree(c, prefix)
	}
}
