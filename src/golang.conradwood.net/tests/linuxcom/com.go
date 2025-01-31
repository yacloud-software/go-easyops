package main

import (
	"fmt"
	"golang.conradwood.net/go-easyops/errors"
	"golang.conradwood.net/go-easyops/utils"
	"io"
	"os"
	"sync"
)

var (
	ctr     = 0
	ctrlock sync.Mutex
)

type Command interface {
}
type command struct {
	exe           []string
	cgroupdir     string
	cgroupdir_cmd string
	stdinwriter   io.Writer
	stdoutreader  io.Reader
	stderrreader  io.Reader
}

func NewCommand() *command {
	return &command{
		cgroupdir: "/sys/fs/cgroup/LINUXCOM",
	}
}

// e.g. /sys/fs/cgroup/LINUXCOM
func (c *command) SetCGroupDir(dir string) {
	c.cgroupdir = dir
}

func (c *command) SetExecutable(com ...string) {
	c.exe = com
}

func (c *command) StdinWriter(r io.Writer) {
	c.stdinwriter = r
}
func (c *command) StdoutReader(r io.Reader) {
	c.stdoutreader = r
}
func (c *command) StderrReader(r io.Reader) {
	c.stderrreader = r
}
func (c *command) IsRunning() bool {
	return true
}

// failed to start, then error
func (c *command) Start() error {
	n := newctr()
	c.cgroupdir_cmd = fmt.Sprintf("%s/com_%d", c.cgroupdir, n)
	err := mkdir(c.cgroupdir_cmd + "/tasks")
	if err != nil {
		return err
	}

	return nil
}

// if unable to wait, it returns error
func (c *command) Wait() error {
	return nil
}
func (c *command) ExitCode() int {
	return 0
}
func (c *command) CombinedOutput() []byte {
	return nil
}
func (c *command) SigInt() { // -2
}
func (c *command) SigKill() { // -9
}

func newctr() int {
	ctrlock.Lock()
	ctr++
	res := ctr
	ctrlock.Unlock()
	return res
}

func mkdir(dir string) error {
	if utils.FileExists(dir) {
		return nil
	}
	err := os.MkdirAll(dir, 0777)
	if err != nil {
		return errors.Wrap(err)
	}
	if !utils.FileExists(dir) {
		return errors.Errorf("failed to create \"%s\"", dir)
	}
	return nil
}
