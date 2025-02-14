package linux

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/user"
	"strconv"
	"syscall"
	"time"

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

   to move a process (e.g. bash) into a new cgroup, use:
   echo [PID] >/sys/fs/cgroup/LINUXCOM/foo/me/cgroup.procs
*/

type Command interface {
	SigInt() error  // -2
	SigKill() error // -9
	SetStdinWriter(r io.Writer)
	SetStdoutReader(r io.Reader)
	SetStderrReader(r io.Reader)
	IsRunning() bool
	SetDebug(bool)
	Start(ctx context.Context, com ...string) (ComInstance, error)
}
type ComInstance interface {
	Wait(ctx context.Context) error    // waits for main command to exit. might leave fork'ed children running
	WaitAll(ctx context.Context) error // waits for all children to exit as well
	Signal(signal syscall.Signal) error
	GetCommand() Command
}
type command struct {
	stdinwriter  io.Writer
	stdoutreader io.Reader
	stderrreader io.Reader
	instance     *cominstance
	debug        bool
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

func NewCommand() Command {
	return &command{}
}
func (c *command) SetDebug(b bool) {
	c.debug = b
}
func (c *cominstance) GetCommand() Command {
	return c.command
}

func (c *command) SetStdinWriter(r io.Writer) {
	c.stdinwriter = r
}
func (c *command) SetStdoutReader(r io.Reader) {
	c.stdoutreader = r
}
func (c *command) SetStderrReader(r io.Reader) {
	c.stderrreader = r
}
func (c *command) IsRunning() bool {
	return true
}

// failed to start, then error
func (c *command) Start(ctx context.Context, com ...string) (ComInstance, error) {
	//	n := newctr()
	//	cgroupdir_cmd := fmt.Sprintf("%s/com_%d", c.cgroupdir, n)
	cgroupdir_cmd, err := CreateStandardAdjacentCgroup()
	if err != nil {
		return nil, err
	}
	c.debugf("Created cgroup \"%s\"\n", cgroupdir_cmd)
	err = mkdir(cgroupdir_cmd + "/tasks")
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
	pids, err := get_pids_for_cgroup(ci.cgroupdir_cmd)
	if err == nil && len(pids) == 0 {
		remove_cgroup(ci.cgroupdir_cmd)
	}
	return err
}
func (ci *cominstance) WaitAll(ctx context.Context) error {
	com_err := ci.Wait(ctx)
	sig := syscall.SIGINT
	wait_started := time.Now()
	for {
		if time.Since(wait_started) > time.Duration(5)*time.Second {
			sig = syscall.SIGKILL
		}
		pids, err := get_pids_for_cgroup(ci.cgroupdir_cmd)
		if err != nil {
			fmt.Printf("Could not get pids for cgroup \"%s\": %s\n", ci.cgroupdir_cmd, err)
			return err
		}
		if len(pids) == 0 {
			break
		}
		for _, pid := range pids {
			ci.debugf("Sending signal %v to pid %d\n", sig, pid)
			err = syscall.Kill(int(pid), sig)
		}

		ci.debugf("Waiting for pid(s): %v\n", pids)
		waited := false
		proc, err := os.FindProcess(int(pids[0]))
		if err != nil {
			fmt.Printf("Failed to find proc: %s\n", err)
		} else {
			_, err := proc.Wait()
			if err != nil {
				ci.debugf("failed to wait for proc: %s\n", err)
			} else {
				waited = true
			}
		}
		if !waited {
			time.Sleep(time.Duration(1) * time.Second)
		}
	}
	if com_err != nil {
		return com_err
	}
	fmt.Printf("All processes exited, now removing cgroup dir (%s)\n", ci.cgroupdir_cmd)
	remove_cgroup(ci.cgroupdir_cmd)
	return nil
}

func (c *command) ExitCode() int {
	return 0
}
func (c *command) CombinedOutput() []byte {
	return nil
}
func (c *command) SigInt() error { // -2
	fmt.Printf("sending sigint\n")
	ci := c.instance
	if ci == nil {
		return errors.Errorf("no instance to send signal to")
	}
	return ci.Signal(syscall.SIGINT)

}
func (c *command) SigKill() error { // -9
	fmt.Printf("sending sigkill\n")
	ci := c.instance
	if ci == nil {
		return errors.Errorf("no instance to send signal to")
	}
	return ci.Signal(syscall.SIGKILL)
}

func (ci *cominstance) Signal(sig syscall.Signal) error {
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

func (ci *cominstance) debugf(format string, args ...any) {
	if !ci.command.debug {
		return
	}
	x := fmt.Sprintf(format, args...)
	prefix := fmt.Sprintf("[%s] ", ci.exe[0])
	fmt.Printf("%s%s", prefix, x)
}
func (c *command) debugf(format string, args ...any) {
	if !c.debug {
		return
	}
	x := fmt.Sprintf(format, args...)
	prefix := "[no instance] "
	if c.instance != nil && c.instance.exe != nil {
		prefix = fmt.Sprintf("[%s] ", c.instance.exe[0])
	}
	fmt.Printf("%s%s", prefix, x)
}
