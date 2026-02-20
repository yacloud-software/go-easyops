package utils

import (
	"flag"
	"fmt"
	"sort"
	"strings"
)

// very weird string flag thing
func FlagStrings(name, help string, values map[string]string) *flag_string_t {
	t := &flag_string_t{name: name, help: help, values: values}
	var vals []string
	for k, _ := range values {
		vals = append(vals, k)
	}
	sort.Slice(vals, func(i, j int) bool {
		return vals[i] < vals[j]
	})
	s := " [" + strings.Join(vals, "|") + "]"
	flag.Var(t, name, help+s)
	return t
}

type flag_string_t struct {
	name   string
	help   string
	values map[string]string
	value  string
}

func (f *flag_string_t) Set(s string) error {
	for k, _ := range f.values {
		if k == s {
			f.value = s
			return nil
		}
	}
	return fmt.Errorf("not a valid flag value \"%s\" for -%s.", s, f.name)
}
func (f *flag_string_t) String() string {
	return f.name
}
