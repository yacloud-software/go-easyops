package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"sync"

	"golang.conradwood.net/go-easyops/errors"
	"golang.conradwood.net/go-easyops/utils"
)

var (
	ctr     = 0
	ctrlock sync.Mutex
)

type Command interface {
}
type command struct {
	cgroupdir    string
	stdinwriter  io.Writer
	stdoutreader io.Reader
	stderrreader io.Reader
	instance     *cominstance
}
type cominstance struct {
	exe             []string
	command         *command
	cgroupdir_cmd   string
	com             *exec.Cmd
	stdout_pipe     io.ReadCloser
	stderr_pipe     io.ReadCloser
	defStdoutReader *comDefaultReader
	defStderrReader *comDefaultReader
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
func (c *command) Start(ctx context.Context, com ...string) (*cominstance, error) {
	n := newctr()
	cgroupdir_cmd := fmt.Sprintf("%s/com_%d", c.cgroupdir, n)
	err := mkdir(cgroupdir_cmd + "/tasks")
	if err != nil {
		return nil, err
	}
	ci := &cominstance{command: c, cgroupdir_cmd: cgroupdir_cmd}
	return ci, ci.start(ctx, com...)
}
func (ci *cominstance) start(ctx context.Context, com ...string) error {
	var err error
	ci.com = exec.CommandContext(ctx, com[0], com[1:]...)
	ci.stdout_pipe, err = ci.com.StdoutPipe()
	if err != nil {
		return err
	}
	ci.defStdoutReader = newDefaultReader(ci.stdout_pipe)
	ci.stderr_pipe, err = ci.com.StderrPipe()
	if err != nil {
		return err
	}
	ci.defStderrReader = newDefaultReader(ci.stdout_pipe)

	err = ci.com.Start()
	if err != nil {
		return errors.Wrap(err)
	}
	return nil
}

func (ci *cominstance) Wait(ctx context.Context) error {
	if ci.com == nil {
		return nil
	}
	err := ci.com.Wait()
	return err
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
