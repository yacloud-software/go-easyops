package utils

import (
	"crypto/sha256"
	"fmt"
	"hash"
	"io/ioutil"
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
	dh := &DirHasher{root: ddir, hasher: sha256.New()}
	err := DirWalk(ddir, dh.AddEntry)
	if err != nil {
		return "", err
	}
	sum := fmt.Sprintf("%x", dh.hasher.Sum(nil))
	return sum, nil
}
func (dh *DirHasher) AddEntry(root string, relname string) error {
	_, err := dh.hasher.Write([]byte(relname))
	if err != nil {
		return err
	}
	b, err := ioutil.ReadFile(root + "/" + relname)
	if err != nil {
		return err
	}
	_, err = dh.hasher.Write(b)
	if err != nil {
		return err
	}

	return nil
}
