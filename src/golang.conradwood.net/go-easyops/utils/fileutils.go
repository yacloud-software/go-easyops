package utils

import (
	"flag"
	"fmt"
	"golang.org/x/sys/unix"
	"io/ioutil"
	"os"
	"os/user"
	"sync"
)

var (
	debug_find_file = flag.Bool("ge_debug_find_file", false, "debug fuzzy filename matches")
	find_file_cache = make(map[string]string)
	ffclock         sync.Mutex
)

// given an arbitrary string, will remove all unsafe characters. result may be safely used as a filename
func MakeSafeFilename(name string) string {
	// allowed chars
	foo := ".0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ_=-"
	res := ""
	isdup := false
	for _, x := range name {
		isgood := false
		for _, allowed := range foo {
			if allowed == x {
				isgood = true
				break
			}
		}
		if isgood {
			isdup = false
		} else {
			if isdup {
				continue
			}
			isdup = true
			x = '_'
		}
		res = res + fmt.Sprintf("%c", x)
	}
	return res
}

// like ioutil - but with open permissions to share
func WriteFile(filename string, content []byte) error {
	unix.Umask(000)
	err := ioutil.WriteFile(filename, content, 0666)
	return err
}

// like ioutil - but with open permissions to share
func OpenWriteFile(filename string) (*os.File, error) {
	unix.Umask(000)
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	return file, err
}

// I can never remember how to do this, so here's a helper:
// return true if file exists, false otherwise
func FileExists(filename string) bool {
	_, err := os.Stat(filename)
	if err != nil {
		return false
	}
	return true
}
func debugFind(name, nname string) {
	if !*debug_find_file {
		return
	}
	fmt.Printf("File: %s resolved to %s\n", name, nname)
	return
}
func FindFile(name string) (string, error) {
	ffclock.Lock()
	nname, foo := find_file_cache[name]
	ffclock.Unlock()
	if foo {
		return nname, nil
	}
	nname = name
	for i := 0; i < 20; i++ {
		if FileExists(nname) {
			debugFind(name, nname)
			ffclock.Lock()
			find_file_cache[name] = nname
			ffclock.Unlock()
			return nname, nil
		}
		nname = "../" + nname
	}
	debugFind(name, "[not found]")
	return "", fmt.Errorf("File not found: %s", name)
}

// read file (uses some magic to find it too)
func ReadFile(filename string) ([]byte, error) {
	fn, err := FindFile(filename)
	if err != nil {
		return nil, err
	}
	b, err := ioutil.ReadFile(fn)
	return b, err
}

// get current home dir
func HomeDir() (string, error) {
	u, err := user.Current()
	if err != nil {
		return "", fmt.Errorf("no current user: %s", err)
	}
	if u == nil {
		return "", fmt.Errorf("No current user")
	}
	return u.HomeDir, nil
}
