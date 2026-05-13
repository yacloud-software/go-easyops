package utils

import (
	"flag"
	"fmt"
	"strings"
)

func DebugFlag(name string) *debugFlag {
	res := &debugFlag{name: name}
	helptext := fmt.Sprintf("Enable debug mode for \"%s\"", name)
	flag.Var(res, "debug_"+name, helptext)
	return res
}

type debugFlag struct {
	name  string
	value bool
	isset bool
}

func (t *debugFlag) Set(s string) error {
	x := strings.ToLower(s)
	if x == "true" || x == "on" {
		t.value = true
	} else if x == "false" || x == "off" {
		t.value = false
	} else {
		return fmt.Errorf("value \"%s\" not valid for a boolean value", s)
	}
	t.isset = true
	return nil
}
func (t *debugFlag) String() string {
	if t == nil {
		return "debugflag description"
	}
	if !t.isset {
		return "undefined"
	}
	return fmt.Sprintf("%v", t.value)
}
func (t *debugFlag) Value() string {
	return fmt.Sprintf("%v", t.value)
}
func (t *debugFlag) IsBoolFlag() bool {
	return true
}

func (t *debugFlag) IsSet() bool {
	return t.isset
}
func (t *debugFlag) BoolValue() bool {
	return t.value

}

func (t *debugFlag) Printf(format string, args ...any) {
	prefix := fmt.Sprintf("[%s] ", t.name)
	txt := fmt.Sprintf(format, args...)
	fmt.Print(prefix + txt)
}
func (t *debugFlag) Debugf(format string, args ...any) {
	if !t.value {
		return
	}
	t.Printf(format, args...)
}
