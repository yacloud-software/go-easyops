package utils

import (
	"flag"
	"fmt"
	"strings"
)

var (
	tribool tribool_t
)

func Tribool(name, help string) *tribool_t {
	t := &tribool_t{name: name, help: help}
	flag.Var(t, name, help)
	return t
}

type tribool_t struct {
	name  string
	help  string
	isset bool
	value bool
}

func (t *tribool_t) Set(s string) error {
	x := strings.ToLower(s)
	if x == "true" {
		t.value = true
	} else if x == "false" {
		t.value = false
	} else {
		return fmt.Errorf("value \"%s\" not valid for a boolean value", s)
	}
	t.isset = true
	return nil
}
func (t *tribool_t) String() string {
	if t == nil {
		return "tribool description"
	}
	if !t.isset {
		return "undefined"
	}
	return fmt.Sprintf("%v", t.value)
}
func (t *tribool_t) Value() string {
	return fmt.Sprintf("%v", t.value)
}
func (t *tribool_t) IsBoolFlag() bool {
	return true
}

func (t *tribool_t) IsSet() bool {
	return t.isset
}
func (t *tribool_t) BoolValue() bool {
	return t.value

}
