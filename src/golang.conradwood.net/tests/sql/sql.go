package main

import (
	"flag"
	"golang.conradwood.net/go-easyops/sql"
	"golang.conradwood.net/go-easyops/utils"
)

func main() {
	flag.Parse()
	_, err := sql.Open()
	utils.Bail("failed to open", err)
}
