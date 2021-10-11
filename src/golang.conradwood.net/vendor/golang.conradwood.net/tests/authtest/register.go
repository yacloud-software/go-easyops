package main

import (
	"fmt"
	"golang.conradwood.net/apis/create"
)

func Register() {
	p := *port + 1
	fmt.Printf("Creating service on port %d\n", p)
	create.NewEasyOpsTest(&testserver{}, p)
}
