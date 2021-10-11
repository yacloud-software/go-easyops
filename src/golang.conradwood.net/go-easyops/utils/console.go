package utils

import (
	"fmt"
)

func SetConsoleTitle(title string) {
	fmt.Printf("\033]0;%s\007", title)
}
