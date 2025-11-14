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
	fmt.Print(t.ToPrettyString())
	fmt.Print(t.ToCSV())

	fmt.Print("Hiding Column #1\n")
	t.DisableColumn(1)
	fmt.Print(t.ToPrettyString())

	fmt.Print("Hiding Column #0\n")
	t.DisableColumn(0)
	fmt.Print(t.ToPrettyString())

	fmt.Print("Enabling all Columns again\n")
	t.EnableAllColumns()
	fmt.Print(t.ToPrettyString())

}
