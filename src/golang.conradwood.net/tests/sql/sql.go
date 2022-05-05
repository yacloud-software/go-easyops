package main

import (
	"golang.conradwood.net/go-easyops/sql"
	"golang.conradwood.net/go-easyops/utils"
)

func main() {
	_, err := sql.Open()
	utils.Bail("failed to open", err)
}
