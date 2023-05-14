package main

import (
	"fmt"
	"golang.conradwood.net/go-easyops/linux"
	"strings"
	"time"
)

func check_with_duration(dur time.Duration, com []string) {
	coms := strings.Join(com, " ")
	fmt.Printf("executing \"%s\", max runtime %0.2fs\n", coms, dur.Seconds())
	started := time.Now()
	l := linux.New()
	l.SetMaxRuntime(dur)
	out, err := l.SafelyExecute(com, nil)
	if err != nil {
		fmt.Printf("\"%s\" failed (%s)\n", coms, err)
		fmt.Println(out)
	}
	exectime := time.Since(started)
	fmt.Printf("Execution time: %0.2fs\n", exectime.Seconds())
}
