package utils

import (
	"fmt"
	"io/ioutil"
	"os"
)

var ()

// if dir does not exist: create it and create a file ".filelayouter"
// if dir DOES exist, check for existence of a file  ".filelayouter", if so, recreate
// otherwise error
func RecreateSafely(dirname string) error {
	fname := dirname + "/.goeasyops-dir"

	// check if directory exists and is empty
	if FileExists(dirname) {
		files, err := ioutil.ReadDir(dirname)
		if err == nil {
			if len(files) == 0 {
				err = os.Chmod(dirname, 0777)
				if err != nil {
					return err
				}
				err = WriteFile(fname, make([]byte, 0))
				return err
			}
		}
	}

	// check if directory exists and has goeasyops marker
	var err error
	if FileExists(dirname) {
		if FileExists(fname) {
			err = RemoveAll(dirname)
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
