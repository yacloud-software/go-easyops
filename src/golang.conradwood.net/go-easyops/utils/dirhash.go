package utils

import (
	"crypto/sha256"
	"fmt"
	"hash"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type DirHasher struct {
	root   string
	hasher hash.Hash
}

/*
return a hash across an entire directory tree, including the contents of all files
*/
func DirHash(ddir string) (string, error) {
	if ddir == "" {
		return "", fmt.Errorf("no config dir")
	}
	ddir = strings.TrimSuffix(ddir, "/")
	dh := &DirHasher{root: ddir, hasher: sha256.New()}
	err := dh.WalkDir("/")
	if err != nil {
		return "", err
	}
	sum := fmt.Sprintf("%x", dh.hasher.Sum(nil))
	return sum, nil
}

// path is relative to "root"
func (dh *DirHasher) WalkDir(path string) error {
	path = strings.TrimPrefix(path, "/")
	fpath := fmt.Sprintf("%s/%s", dh.root, path)
	entries, err := ioutil.ReadDir(fpath)
	if err != nil {
		return err
	}
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Name() < entries[j].Name()
	})
	// hash all filenames
	for _, e := range entries {
		relname := path + "/" + e.Name()
		fmt.Printf("file: \"%s\"\n", relname)
		_, err = dh.hasher.Write([]byte(relname))
		if err != nil {
			return err
		}
	}

	// do files first
	for _, e := range entries {
		if !e.Mode().IsRegular() {
			continue
		}
		fname := fpath + "/" + e.Name()
		b, err := ioutil.ReadFile(fname)
		if err != nil {
			return err
		}
		_, err = dh.hasher.Write(b)
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
		err := dh.WalkDir(path + "/" + fname)
		if err != nil {
			return err
		}
	}
	return nil
}

// MD5All reads all the files in the file tree rooted at root and returns a map
// from file path to the MD5 sum of the file's contents.  If the directory walk
// fails or any read operation fails, MD5All returns an error.
func SHAAll(root string) ([]string, error) {
	var m []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.Mode().IsRegular() {
			return nil
		}
		data, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}
		data = append(data, ([]byte(path))...) // append the filename to data, so that change of filename also changes dirhash
		x := sha256.Sum256(data)
		m = append(m, fmt.Sprintf("%x", x))
		return nil
	})
	if err != nil {
		return nil, err
	}
	return m, nil
}
