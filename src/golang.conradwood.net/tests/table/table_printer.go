package main

import (
	"fmt"
	"golang.conradwood.net/go-easyops/utils"
)

func main() {
	t := &utils.Table{}
	t.AddHeaders("col0", "column1", "last column")
	t.AddStrings("foo", "bar", "foobar")
	t.NewRow()
	t.AddStrings("foo", "bar", "foobar")
	fmt.Printf(t.ToPrettyString())
	fmt.Printf(t.ToCSV())

	fmt.Printf("Hiding Column #1\n")
	t.DisableColumn(1)
	fmt.Printf(t.ToPrettyString())

	fmt.Printf("Hiding Column #0\n")
	t.DisableColumn(0)
	fmt.Printf(t.ToPrettyString())

	fmt.Printf("Enabling all Columns again\n")
	t.EnableAllColumns()
	fmt.Printf(t.ToPrettyString())

}
