package main

import (
	"strconv"
	"strings"

	"golang.conradwood.net/go-easyops/errors"
	"golang.conradwood.net/go-easyops/utils"
)

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
