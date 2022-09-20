package cmdline

import (
	"flag"
	"os"
)

type impara struct {
	f    func()
	name string
	desc string
}

func ImmediatePara(name string, desc string, f func()) *impara {
	t := &impara{name: name, desc: desc, f: f}
	flag.Var(t, name, desc)
	return t
}

func (i *impara) Set(s string) error {
	i.f()
	os.Exit(0)
	return nil
}
func (i *impara) String() string {
	return i.name + " " + i.desc
}
