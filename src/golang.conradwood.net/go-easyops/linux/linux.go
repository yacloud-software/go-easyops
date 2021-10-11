package linux

import (
	"context"
	"flag"
	"fmt"
	"golang.conradwood.net/go-easyops/auth"
	"io"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"
)

var (
	cmdLock    sync.Mutex
	curCmd     string
	LogExe     = flag.Bool("ge_debug_exe", false, "debug execution of third party binaries")
	maxRuntime = flag.Int("ge_default_max_runtime_exe", 5, "m̀ax_runtime in seconds for external binaries")
)

type linux struct {
	Runtime          int
	AllowConcurrency bool
	ctx              context.Context
	envs             []string
}

type Linux interface {
	SafelyExecute(cmd []string, stdin io.Reader) (string, error)
	SafelyExecuteWithDir(cmd []string, dir string, stdin io.Reader) (string, error)
	MyIP() string
	SetRuntime(int)
	SetAllowConcurrency(bool)
	SetEnvironment([]string)
}

func NewWithContext(ctx context.Context) Linux {
	l := New()
	ln := l.(*linux)
	ln.ctx = ctx
	return l
}
func New() Linux {
	res := &linux{
		Runtime:          *maxRuntime,
		AllowConcurrency: false,
		ctx:              context.TODO(),
	}
	return res
}

// execute a command...
// print stdout/err (so it ends up in the logs)
// also we add a timeout - if program hangs we return an error
// rather than 'hanging' forever
// and we use a low-level lock to avoid calling binaries at the same time
func (l *linux) SafelyExecute(cmd []string, stdin io.Reader) (string, error) {
	return l.SafelyExecuteWithDir(cmd, "", stdin)
}
func (l *linux) SafelyExecuteWithDir(cmd []string, dir string, stdin io.Reader) (string, error) {
	// avoid possible segfaults (afterall it's called 'safely...')
	if len(cmd) == 0 {
		return "", fmt.Errorf("no command specified for execute.")
	}
	if !l.AllowConcurrency {
		if curCmd != "" {
			if *LogExe {
				fmt.Printf("Waiting for %s to complete...\n", curCmd)
			}
		}
		cmdLock.Lock()
		defer cmdLock.Unlock()
	}
	curCmd = cmd[0]
	if curCmd == "sudo" {
		if len(curCmd) < 2 {
			return "", fmt.Errorf("sudo without parameters not allowed")
		}
		curCmd = cmd[1]
	}
	// execute
	if *LogExe {
		fmt.Printf("Executing %s\n", curCmd)
	}
	c := exec.CommandContext(l.ctx, cmd[0], cmd[1:]...)
	if dir != "" {
		c.Dir = dir
	}
	if stdin != nil {
		c.Stdin = stdin
	}
	// set environment
	c.Env = os.Environ()
	l.env(c)
	output, err := l.syncExecute(c, l.Runtime)
	if *LogExe {
		printOutput(curCmd, output)
	}
	curCmd = ""
	if err != nil {
		return output, err
	}
	return output, nil
}

// execute with timeout.
// sends SIGKILL to process on timeout and returns error
func (l *linux) syncExecute(c *exec.Cmd, timeout int) (string, error) {
	running := false
	killed := false
	timer1 := time.NewTimer(time.Second * time.Duration(timeout))
	go func() {
		<-timer1.C
		if running {
			c.Process.Kill()
			killed = true
		}
	}()
	// racecondition - timer might expire between
	// setting flag and starting process.
	// (if timer is really short)
	running = true
	b, err := c.CombinedOutput()
	running = false
	if killed {
		err = fmt.Errorf("Process killed after %d seconds", timeout)
	}
	return string(b), err
}

func printOutput(cmd string, output string) {
	fmt.Printf("====BEGIN OUTPUT OF %s====\n", cmd)
	fmt.Printf("%s\n", output)
	fmt.Printf("====END OUTPUT OF %s====\n", cmd)
}
func (l *linux) SetEnvironment(sx []string) {
	l.envs = sx
}
func (l *linux) SetRuntime(r int) {
	l.Runtime = r
}
func (l *linux) SetAllowConcurrency(b bool) {
	l.AllowConcurrency = b
}

// add context to environment
func (l *linux) env(c *exec.Cmd) error {
	nc, err := auth.SerialiseContextToString(l.ctx)
	if err != nil {
		return err
	}
	ncs := fmt.Sprintf("GE_CTX=%s", nc)

	for i, e := range c.Env {
		if strings.HasPrefix(e, "GE_CTX=") {
			c.Env[i] = ncs
			return nil
		}
	}
	c.Env = append(c.Env, ncs)
	for _, e := range l.envs {
		c.Env = append(c.Env, e)
	}
	return nil
}
