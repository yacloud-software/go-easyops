package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/user"
	"strconv"
	"sync"
	"syscall"

	"golang.conradwood.net/go-easyops/errors"
	"golang.conradwood.net/go-easyops/utils"
	"golang.org/x/sys/unix"
)

/*
   CGROUP Permissions:

   clone3 returns -EACCESS (permission denied) unless:
   the user has write access to cgroup.procs in the nearest common ancestor director of calling process and cgroup of new process.

   EXAMPLES: (assuming only LINUXCOM and below is user-writeable)
   | CALLING_PROC                    | NEW_PROC                           | Result   |
   +---------------------------------+------------------------------------+----------+
   | /sys/fs/cgroup/LINUXCOM/me/     | /sys/fs/cgroup/LINUXCOM/com_1/     | EACCESS  |
   | /sys/fs/cgroup/LINUXCOM/foo/me/ | /sys/fs/cgroup/LINUXCOM/com_1/     | EACCESS  |
   | /sys/fs/cgroup/LINUXCOM/foo/me/ | /sys/fs/cgroup/LINUXCOM/foo/com_1/ | OK       |
*/

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
		cgroupdir: "/sys/fs/cgroup/LINUXCOM/ancestor/",
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
	c.instance = ci
	return ci, ci.start(ctx, com...)
}
func (ci *cominstance) start(ctx context.Context, com ...string) error {
	u, err := user.Current()
	if err != nil {
		return errors.Wrap(err)
	}

	uid, err := strconv.Atoi(u.Uid)
	if err != nil {
		return errors.Wrap(err)
	}

	gid, err := strconv.Atoi(u.Gid)
	if err != nil {
		return errors.Wrap(err)
	}

	// open cgroup filedescriptor
	cgroup_fd_path := ci.cgroupdir_cmd
	cgroup_fd, err := syscall.Open(cgroup_fd_path, unix.O_PATH, 0)
	if err != nil {
		return errors.Wrap(err)
	}
	fmt.Printf("CgroupFD for \"%s\": %d\n", cgroup_fd_path, cgroup_fd)
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

	ci.com.SysProcAttr = &syscall.SysProcAttr{
		Credential: &syscall.Credential{
			Uid:         uint32(uid),
			Gid:         uint32(gid),
			NoSetGroups: true,
		},
		UseCgroupFD: true,
		CgroupFD:    cgroup_fd,
	}
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
func (c *command) SigInt() error { // -2
	fmt.Printf("sending sigint\n")
	return c.sendsig(syscall.SIGINT)
}
func (c *command) SigKill() error { // -9
	fmt.Printf("sending sigkill\n")
	return c.sendsig(syscall.SIGKILL)
}

func (c *command) sendsig(sig syscall.Signal) error {
	ci := c.instance
	pids, err := get_pids_for_cgroup(ci.cgroupdir_cmd)
	if err != nil {
		fmt.Printf("Could not get pids for cgroup \"%s\": %s\n", ci.cgroupdir_cmd, err)
		return err
	}
	fmt.Printf("Cgroupdir \"%s\" has %d pids\n", ci.cgroupdir_cmd, len(pids))
	for _, pid := range pids {
		fmt.Printf("Sending signal %v to pid %d\n", sig, pid)
		err = syscall.Kill(int(pid), sig)
		if err != nil {
			return errors.Wrap(err)
		}
	}
	return nil
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
