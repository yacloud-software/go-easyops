package main

import (
	"golang.conradwood.net/go-easyops/errors"
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
