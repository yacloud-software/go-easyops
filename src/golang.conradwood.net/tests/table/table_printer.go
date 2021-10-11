package main

import (
	"fmt"
	"golang.conradwood.net/go-easyops/utils"
)

func main() {
	t := &utils.Table{}
	t.AddHeaders("col1", "column2", "last column")
	t.AddStrings("foo", "bar", "foobar")
	t.NewRow()
	t.AddStrings("foo", "bar", "foobar")
	fmt.Printf(t.ToPrettyString())

}
