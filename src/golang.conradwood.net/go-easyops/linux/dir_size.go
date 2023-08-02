package linux

import (
	"golang.conradwood.net/go-easyops/utils"
	"os"
)

func DirSize(dir string) (uint64, error) {
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
