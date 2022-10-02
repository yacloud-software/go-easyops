package utils

import (
	"fmt"
	"io/ioutil"
	"sort"
	"strings"
)

type dirwalker struct {
	root string
	fn   func(root string, relative_filename string) error
}

// walk a directory tree and call function for each file
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
		if !e.Mode().IsRegular() {
			continue
		}
		err := dw.fn(dw.root, path+"/"+e.Name())
		if err != nil {
			return err
		}
	}
	// do dirs now
	for _, e := range entries {
		if !e.IsDir() {
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
