package utils

import (
	"flag"
	"fmt"
	"golang.org/x/sys/unix"
	"io/fs"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"sync"
)

var (
	extra_dir       = flag.String("ge_findfile_additional_dir", "", "if set, findfile will search this directory as well")
	debug_find_file = flag.Bool("ge_debug_find_file", false, "debug fuzzy filename matches")
	find_file_cache = make(map[string]string)
	ffclock         sync.Mutex
	workingdir      string
)

func init() {
	var err error
	workingdir, err = os.Getwd()
	if err != nil {
		fmt.Printf("cannot get current working directory: %s\n", err)
	}
}

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

// search for a file in git repo and above. error if not found
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
			nname, err := filepath.Abs(nname)
			if err != nil {
				return "", err
			}
			return nname, nil
		}
		nname = "../" + nname
	}
	if *extra_dir != "" {
		nname := fmt.Sprintf("%s/%s", *extra_dir, name)
		nname, err := filepath.Abs(nname)
		if err != nil {
			return "", err
		}
		if FileExists(nname) {
			return nname, nil
		}
	}
	debugFind(name, "[not found]")
	return "", fmt.Errorf("File not found: %s", name)
}

// this will find a file or directory in the given working directory by traversing up the directory
// hierarchy but only upto the workingdir (not higher).
// the returned file is guaranteed to be within the workingdir.
// if it cannot be found, an error is returned.
func FindFileInWorkingDir(name string) (string, error) {
	name, err := FindFile(name)
	if err != nil {
		return "", err
	}
	name, err = filepath.Abs(name)
	if err != nil {
		return "", err
	}
	if !strings.HasPrefix(name, WorkingDir()) {
		return "", fmt.Errorf("file %s not in working dir %s", name, WorkingDir())
	}
	return name, nil
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

// return the "workingdir" (the one we were started in)
func WorkingDir() string {
	if workingdir == "" {
		var err error
		workingdir, err = os.Getwd()
		// this is critical, user wanted workingdir, but we do not
		// have one. how did this even happen?
		Bail("failed to get workingdir", err)
	}
	return workingdir
}

// get current home dir (from environment variable HOME, if that fails user.Current())
func HomeDir() (string, error) {
	he := os.Getenv("HOME")
	if he != "" && FileExists(he) {
		return he, nil
	}
	u, err := user.Current()
	if err != nil {
		return "", fmt.Errorf("no current user: %s", err)
	}
	if u == nil {
		return "", fmt.Errorf("No current user")
	}
	return u.HomeDir, nil
}

// removes dir, changes permissions if needs to
func RemoveAll(dir string) error {
	var err error

	err = os.RemoveAll(dir)
	if err == nil {
		return nil
	}

	// reset permissions and try again
	err = ChmodR(dir, 0777, 0777)
	if err != nil {
		return err
	}
	err = os.RemoveAll(dir)
	if err != nil {
		return err
	}

	return err
}

// recursively chmod a dir and files and subdirectories. applies 'dirmask' to dirs and 'filemask' to files
func ChmodR(dir string, dirmask, filemask fs.FileMode) error {
	err := os.Chmod(dir, dirmask)
	if err != nil {
		return err
	}
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}
	for _, f := range files {
		if f.IsDir() {
			err = ChmodR(dir+"/"+f.Name(), dirmask, filemask)
			if err != nil {
				return err
			}
			continue
		}
		err = os.Chmod(dir+"/"+f.Name(), filemask)
		if err != nil {
			return err
		}
	}
	return nil
}
