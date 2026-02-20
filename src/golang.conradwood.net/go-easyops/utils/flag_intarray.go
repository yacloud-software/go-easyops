package utils

import (
	"flag"
	"strconv"
	"strings"
)

// a flag to allow for a list of uint64 and a convenient method to check if an integer is part of the flag
func IntArrayFlag(name string, help string) *intarrayFlag {
	res := &intarrayFlag{name: name, help: help}
	flag.Var(res, name, help)
	return res

}

type intarrayFlag struct {
	name string
	help string
	ints []uint64
}

func (t *intarrayFlag) Set(s string) error {
	var i []uint64
	for _, sx := range strings.Split(s, ",") {
		sx = strings.Trim(sx, " ")
		n, err := strconv.ParseUint(sx, 10, 64)
		if err != nil {
			return err
		}
		i = append(i, n)
	}
	t.ints = i
	return nil
}
func (t *intarrayFlag) String() string {
	return t.name
}
func (t *intarrayFlag) Contains(val uint64) bool {
	for _, i := range t.ints {
		if i == val {
			return true
		}
	}
	return false
}
