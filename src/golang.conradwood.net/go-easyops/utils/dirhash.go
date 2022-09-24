package utils

import (
	"crypto/sha256"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
)

/*
return a hash across an entire directory tree, including the contents of all files
*/
func DirHash(ddir string) (string, error) {
	if ddir == "" {
		return "", fmt.Errorf("no config dir")
	}
	m, err := SHAAll(ddir)
	if err != nil {
		return "", err
	}
	sort.Slice(m, func(i, j int) bool {
		return m[i] < m[j]
	})
	s := ""
	for _, v := range m {
		s = s + v + ":"
	}
	res := fmt.Sprintf("%x", sha256.Sum256([]byte(s)))
	return res, nil
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
		x := sha256.Sum256(data)
		m = append(m, fmt.Sprintf("%x", x))
		return nil
	})
	if err != nil {
		return nil, err
	}
	return m, nil
}
