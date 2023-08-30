package cmdline

import (
	"os"
)

var (
	envs []*env_var
)

func ENV(name, description string) *env_var {
	res := &env_var{name: name, description: description}
	envs = append(envs, res)
	return res
}

type env_var struct {
	name        string
	description string
}

func (e *env_var) Value() string {
	x := os.Getenv(e.name)
	return x
}
func (e *env_var) Name() string {
	return e.name
}
func render_env_help() string {
	longest := 0
	for _, e := range envs {
		if len(e.name) > longest {
			longest = len(e.name)
		}
	}
	res := ""
	for _, e := range envs {
		s := e.name
		for len(s) < longest {
			s = s + " "
		}
		s = s + ": " + e.description
		res = res + s + "\n"
	}
	return res
}
