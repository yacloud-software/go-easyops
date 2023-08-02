package linux

import (
	"errors"
	"golang.conradwood.net/go-easyops/utils"
	"os"
)

// if dir is the name of a directory, it will recursively calculate the size. if it is a file, it will stat it and return the filesize
func DirSize(dir string) (uint64, error) {
	fi, err := os.Stat(dir)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return 0, nil
		}
		return 0, err
	}
	if !fi.IsDir() {
		return uint64(fi.Size()), nil
	}
	res := uint64(0)
	utils.DirWalk(dir, func(root, relfile string) error {
		fname := root + "/" + relfile
		fi, err := os.Stat(fname)
		if err != nil {
			return err
		}
		res = res + uint64(fi.Size())
		return nil
	})
	return res, nil
}
