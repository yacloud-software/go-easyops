package linux

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"syscall"
)

func CopyDir(srcDir, dest string) error {
	err := os.MkdirAll(dest, 0777)
	if err != nil {
		return err
	}
	c := copy{permIgnore: true}
	err = c.copyDirInternal(srcDir, dest)
	return err
}

type copy struct {
	permIgnore bool
}

func (c *copy) copyDirInternal(srcDir, dest string) error {
	entries, err := ioutil.ReadDir(srcDir)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		sourcePath := filepath.Join(srcDir, entry.Name())
		destPath := filepath.Join(dest, entry.Name())

		fileInfo, err := os.Stat(sourcePath)
		if err != nil {
			return fmt.Errorf("stat(%s) failed: %s", sourcePath, err)
		}

		stat, ok := fileInfo.Sys().(*syscall.Stat_t)
		if !ok {
			return fmt.Errorf("failed to get raw syscall.Stat_t data for '%s'", sourcePath)
		}

		switch fileInfo.Mode() & os.ModeType {
		case os.ModeDir:
			if err := CreateIfNotExists(destPath, 0755); err != nil {
				return err
			}
			if err := c.copyDirInternal(sourcePath, destPath); err != nil {
				return err
			}
		case os.ModeSymlink:
			if err := CopySymLink(sourcePath, destPath); err != nil {
				return err
			}
		default:
			if err := Copy(sourcePath, destPath); err != nil {
				return err
			}
		}

		// requires root:
		err = os.Lchown(destPath, int(stat.Uid), int(stat.Gid))
		if !c.permIgnore && err != nil {
			return fmt.Errorf("lchown(file=%s,uid=%d,gid=%d) failed: %s", destPath, stat.Uid, stat.Gid, err)
		}
		isSymlink := entry.Mode()&os.ModeSymlink != 0
		if !isSymlink {
			err := os.Chmod(destPath, entry.Mode())
			if !c.permIgnore && err != nil {
				return fmt.Errorf("chmod(%s) failed: %s", destPath, err)
			}
		}
	}
	return nil
}

func Copy(srcFile, dstFile string) error {
	out, err := os.Create(dstFile)
	if err != nil {
		return fmt.Errorf("failed to create(%s): %s", dstFile, err)
	}

	defer out.Close()

	in, err := os.Open(srcFile)
	if err != nil {
		return fmt.Errorf("failed to open(%s): %s", srcFile, err)
	}
	defer in.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}

	return nil
}

func Exists(filePath string) bool {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return false
	}

	return true
}

func CreateIfNotExists(dir string, perm os.FileMode) error {
	if Exists(dir) {
		return nil
	}

	if err := os.MkdirAll(dir, perm); err != nil {
		return fmt.Errorf("failed to create directory: '%s', error: '%s'", dir, err.Error())
	}

	return nil
}

func CopySymLink(source, dest string) error {
	link, err := os.Readlink(source)
	if err != nil {
		return err
	}
	return os.Symlink(link, dest)
}
