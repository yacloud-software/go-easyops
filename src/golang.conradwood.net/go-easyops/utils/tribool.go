package utils

import (
	"flag"
	"fmt"
	"strings"
)

var (
	tribool tribool_t
)

func TriboolNoFlag() *tribool_t {
	return &tribool_t{}
}
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
func (t *tribool_t) SetBool(b bool) {
	t.value = b
	t.isset = true
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

// if not set it is always false, otherwise compares b with bool
func (t *tribool_t) IsEqualTo(b bool) bool {
	if !t.IsSet() {
		return false
	}
	return b == t.BoolValue()
}
