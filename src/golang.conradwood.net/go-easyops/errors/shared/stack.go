package shared

import (
	"fmt"
	"strings"
)

var (
	STACK_FILTER = []string{
		"golang.conradwood.net/go-easyops/errors",
		"runtime",
	}
)

type StackPos struct {
	Filename string
	Function string
	Line     int
}

func first_non_internal_pos(sp []*StackPos) *StackPos {
	if len(sp) == 0 {
		return nil
	}
	for _, sp := range sp {
		if sp.IsFiltered() {
			continue
		}
		return sp
	}
	return sp[0]
}

func (sp *StackPos) String() string {
	if sp == nil {
		return "[no stack]"
	}
	return fmt.Sprintf("%s:%d", sp.Function, sp.Line)
}

func (sp *StackPos) IsFiltered() bool {
	for _, s := range STACK_FILTER {
		if strings.Contains(sp.Function, s) {
			return true
		}
	}
	return false
}
