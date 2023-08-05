package main

import (
	"flag"
	"fmt"
	"golang.conradwood.net/go-easyops/logger"
	"golang.conradwood.net/go-easyops/utils"
	"io"
)

func main() {
	flag.Parse()
	lo, err := logger.NewAsyncLogQueue("test", 50, 1, "test", "test", "foodeplid")
	var w io.Writer
	w = lo // test if asignment works
	w.Write([]byte("foo\n"))
	utils.Bail("failed to create logger", err)
	for i := 0; i < 10; i++ {
		lo.Log("testing", "Line %d logged", i)
	}
	for i := 0; i < 10; i++ {
		lo.LogCommandStdout(fmt.Sprintf("stdout - Line %d logged", i), "testing")
	}
	lo.Close(0)
	fmt.Printf("Done\n")
}
