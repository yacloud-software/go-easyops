/*
Package linux provides methods to execute commands on linux
*/
package linux

import (
	"context"
	"flag"
	"fmt"
	"golang.conradwood.net/go-easyops/auth"
	"golang.conradwood.net/go-easyops/ctx"
	"io"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"
)

const (
	add_serialised_context = false
)

var (
	cmdLock    sync.Mutex
	curCmd     string
	LogExe     = flag.Bool("ge_debug_exe", false, "debug execution of third party binaries")
	maxRuntime = flag.Duration("ge_default_max_runtime_exe", time.Duration(5)*time.Second, "m̀ax_runtime for external binaries")
)

type linux struct {
	Runtime          time.Duration
	AllowConcurrency bool
	ctx              context.Context
	context_set      bool // if user-supplied context
	envs             []string
	lastcmd          []string
	runforever       bool
}

type Linux interface {
	SafelyExecute(cmd []string, stdin io.Reader) (string, error)
	SafelyExecuteWithDir(cmd []string, dir string, stdin io.Reader) (string, error)
	MyIP() string
	SetMaxRuntime(time.Duration)
	SetRunForever() // incompatible with setmaxruntime
	SetAllowConcurrency(bool)
	SetEnvironment([]string)
}

func NewWithContext(ctx context.Context) Linux {
	l := New()
	ln := l.(*linux)
	ln.context_set = true
	ln.ctx = ctx
	return l
}
func New() Linux {
	res := &linux{
		Runtime:          *maxRuntime,
		AllowConcurrency: false,
	}
	res.recalc_context_from_timeout()
	return res
}

func (l *linux) recalc_context_from_timeout() {
	if l.runforever {
		l.ctx = context.Background()
		return
	}
	cb := ctx.NewContextBuilder()
	cb.WithTimeout(l.Runtime)
	l.ctx = cb.ContextWithAutoCancel()
}

// execute a command...
// print stdout/err (so it ends up in the logs)
// also we add a timeout - if program hangs we return an error
// rather than 'hanging' forever
// and we use a low-level lock to avoid calling binaries at the same time
func (l *linux) SafelyExecute(cmd []string, stdin io.Reader) (string, error) {
	return l.SafelyExecuteWithDir(cmd, "", stdin)
}

/*
execute a command within a working directory
*/
func (l *linux) SafelyExecuteWithDir(cmd []string, dir string, stdin io.Reader) (string, error) {
	// avoid possible segfaults (afterall it's called 'safely...')
	if len(cmd) == 0 {
		return "", fmt.Errorf("no command specified for execute.")
	}
	l.lastcmd = cmd
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
		fmt.Printf("[go-easyops] preparing to execute below command:\n%s\n", l.ComWithParas())
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
	output, err := l.syncExecute(c, l.Runtime, !l.runforever)
	if *LogExe {
		printOutput(l.ComName(), output)
	}
	curCmd = ""
	if err != nil {
		return output, err
	}
	return output, nil
}

// execute with timeout.
// sends SIGKILL to process on timeout and returns error
func (l *linux) syncExecute(c *exec.Cmd, timeout time.Duration, hastimeout bool) (string, error) {
	running := false
	killed := false
	if hastimeout {
		timer1 := time.NewTimer(timeout)
		go func() {
			<-timer1.C
			if running {
				if c.Process == nil {
					fmt.Printf("[go-easyops] no process to kill after %0.2fs\n", timeout.Seconds())
					return
				}
				if !running {
					return
				}
				c.Process.Kill()
				killed = true
				if *LogExe {
					fmt.Printf("[go-easyops] process killed after %0.2fs\n", timeout.Seconds())
				}
			}
		}()
	}
	// racecondition - timer might expire between
	// setting flag and starting process.
	// (if timer is really short)
	running = true
	if *LogExe {
		fmt.Printf("[go-easyops] executing command %s (timeout=%0.2fs)\n", l.ComName(), timeout.Seconds())
	}
	b, err := c.CombinedOutput()
	if *LogExe {
		fmt.Printf("[go-easyops] process terminated\n")
	}
	running = false
	if killed {
		err = fmt.Errorf("Process killed after %0.2f seconds", timeout.Seconds())
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
func (l *linux) SetRunForever() {
	l.runforever = true
	if !l.context_set {
		l.recalc_context_from_timeout()
	}
}
func (l *linux) SetMaxRuntime(d time.Duration) {
	l.runforever = false
	l.Runtime = d
	if !l.context_set {
		l.recalc_context_from_timeout()
	}
}
func (l *linux) SetAllowConcurrency(b bool) {
	l.AllowConcurrency = b
}

// add context to environment
func (l *linux) env(c *exec.Cmd) error {
	if l.context_set {
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
	}
	for _, e := range l.envs {
		c.Env = append(c.Env, e)
	}
	return nil
}

func (l *linux) ComWithParas() string {
	if len(l.lastcmd) == 0 {
		return "<no command executed>"
	}
	return strings.Join(l.lastcmd, " ")
}
func (l *linux) ComName() string {
	if len(l.lastcmd) == 0 || l.lastcmd[0] == "" {
		return "<no command executed>"
	}
	s := l.lastcmd[0]
	if strings.Contains(s, "sudo") && len(l.lastcmd) > 1 {
		return "sudo " + l.lastcmd[1]
	}
	return l.lastcmd[0]
}
