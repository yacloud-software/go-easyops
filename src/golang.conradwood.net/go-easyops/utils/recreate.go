package utils

import (
	"fmt"
	"os"
)

var ()

// if dir does not exist: create it and create a file ".filelayouter"
// if dir DOES exist, check for existence of a file  ".filelayouter", if so, recreate
// otherwise error
func RecreateSafely(dirname string) error {
	var err error
	fname := dirname + "/.goeasyops-dir"
	if FileExists(dirname) {
		if FileExists(fname) {
			err = os.RemoveAll(dirname)
			if err != nil {
				return err
			}
		} else {
			return fmt.Errorf("not recreating dir \"%s\", it exists already, but .goeasyops-dir does not.", dirname)
		}
	}
	if FileExists(dirname) {
		return fmt.Errorf("Attempt to delete dir \"%s\" failed", dirname)
	}
	err = os.MkdirAll(dirname, 0777)
	if err != nil {
		return err
	}
	err = WriteFile(fname, make([]byte, 0))
	return err
}
