package utils

import (
	"fmt"
	//	"io/fs"
	"io/ioutil"
	"sort"
	"strings"
)

type dirwalker struct {
	root string
	/*
		this function is called with "root" being whatever DirWalk has been invoked with and relative_filename is the filename
		within that directory. The full path thus can be constructed by root+"/"+relative_filename
	*/
	fn func(root string, relative_filename string) error
}

// walk a directory tree and call function for each file (but not each dir)
func DirWalk(dir string, fn func(root string, relative_filename string) error) error {
	dw := &dirwalker{root: dir, fn: fn}
	return dw.Walk("")
}
func (dw *dirwalker) Walk(relative_path string) error {
	path := strings.TrimPrefix(relative_path, "/")
	fpath := fmt.Sprintf("%s/%s", dw.root, path)
	entries, err := ioutil.ReadDir(fpath)
	if err != nil {
		return err
	}
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Name() < entries[j].Name()
	})

	// do files first
	for _, e := range entries {
		m := e.Mode()
		if !m.IsRegular() {
			continue
		}
		s := path + "/" + e.Name()
		if path == "" {
			s = e.Name()
		}
		err := dw.fn(dw.root, s)
		if err != nil {
			return err
		}
	}
	// do dirs now
	for _, e := range entries {
		if !e.IsDir() {
			//			fmt.Printf("not a dir: %s\n", e.Name())
			continue
		}
		fname := e.Name()
		err := dw.Walk(path + "/" + fname)
		if err != nil {
			return err
		}
	}
	return nil
}
