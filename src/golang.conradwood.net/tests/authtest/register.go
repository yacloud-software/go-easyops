package main

import (
	"fmt"
)

func Register() {
	p := *port + 1
	fmt.Printf("Creating service on port %d\n", p)
}
