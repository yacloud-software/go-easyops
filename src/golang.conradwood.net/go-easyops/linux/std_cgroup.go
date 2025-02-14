package linux

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"golang.conradwood.net/go-easyops/errors"
	"golang.conradwood.net/go-easyops/utils"
)

var (
	cgroup_check_lock sync.Mutex
)

func MyCgroup() (string, error) {
	b, err := utils.ReadFile("/proc/self/cgroup")
	if err != nil {
		return "", errors.Wrap(err)
	}
	s := string(b)
	idx := strings.Index(s, "/")
	if idx == -1 {
		return "", errors.Errorf("odd /proc/self/cgroup line: \"%s\"", s)
	}
	res := s[idx:]
	return res, nil

}

// if caller is in cgroup "/LINUXCOM/ancestor/me", new cgrop will be "/LINUXCOM/ancestor/com_1"
func CreateStandardAdjacentCgroup() (string, error) {
	myc, err := MyCgroup()
	if err != nil {
		return "", err
	}
	ancestor := filepath.Dir(myc)
	ctr := 0
	cgroup_check_lock.Lock()
	defer cgroup_check_lock.Unlock()
	name := ""
	for {
		ctr++
		name = fmt.Sprintf("/sys/fs/cgroup/%s/com_%d", ancestor, ctr)
		if !utils.FileExists(name) {
			break
		}
	}
	if name == "" {
		return "", errors.Errorf("Unable to determine adjacent cgroup")
	}
	err = CreateStandardCgroup(name)
	if err != nil {
		return "", err
	}
	return name, nil
}

func CreateStandardCgroup(dir string) error {
	err := mkdir(dir)
	if err != nil {
		return errors.Errorf("failed to create parent cgroup (%s): %s", dir, err)
	}
	taskdir := dir + "/tasks"
	err = mkdir(taskdir)
	if err != nil {
		return errors.Errorf("failed to create cgroup tasks (%s): %s", taskdir, err)
	}
	return nil
}

func get_pids_for_cgroup(cgroupdir string) ([]uint64, error) {
	fname := cgroupdir + "/cgroup.procs"
	if !utils.FileExists(fname) {
		return nil, nil
	}
	b, err := utils.ReadFile(fname)
	if err != nil {
		return nil, errors.Wrap(err)
	}
	var res []uint64
	for _, line := range strings.Split(string(b), "\n") {
		if line == "" {
			continue
		}
		pid, err := strconv.ParseUint(line, 10, 64)
		if err != nil {
			return nil, errors.Wrap(err)
		}
		res = append(res, pid)
	}
	return res, nil
}

// remove cgroup (and child cgroups)
func remove_cgroup(cgroupdir string) error {
	rmdir(cgroupdir + "/tasks")
	rmdir(cgroupdir)
	return nil
}
func rmdir(dir string) error {
	err := os.Remove(dir)
	if err != nil {
		fmt.Printf("failed to remove dir: %s\n", err)
	}
	return err
}
