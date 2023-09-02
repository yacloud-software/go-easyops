package utils

import (
	"os"
	"testing"
)

func TestDeleteDirs(t *testing.T) {
	test_remove_all(t, "/tmp/testdir")
}
func test_remove_all(t *testing.T, dir string) {
	err := RecreateSafely(dir)
	if err != nil {
		t.Fatalf("failed to create dir %s: %s", dir, err)
	}
	subdir := dir + "/subdir"
	os.MkdirAll(subdir, 0777)
	fname := subdir + "/foobar"
	err = WriteFile(fname, []byte("foo"))
	if err != nil {
		t.Fatalf("failed to write to dir %s: %s", dir, err)
	}
	err = os.Chmod(fname, 000)
	if err != nil {
		t.Fatalf("failed to chmod %s: %s", fname, err)
	}
	err = os.Chmod(subdir, 000)
	if err != nil {
		t.Fatalf("failed to chmod %s: %s", subdir, err)
	}
	err = RemoveAll(dir)
	if err != nil {
		t.Fatalf("failed to removeall() dir %s: %s", dir, err)
	}
	if FileExists(dir) {
		t.Fatalf("failed to delete dir %s (no error)", dir)
	}
}
