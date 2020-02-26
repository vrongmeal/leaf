package commander

import (
	"fmt"
	"os"
	"os/exec"
	"sync"
	"syscall"

	"github.com/kballard/go-shellquote"
	"github.com/sirupsen/logrus"
)

// Commander is a type with multiple commands and runs them in order.
type Commander struct {
	index int
	cmds  []string
	cmd   *exec.Cmd
	kill  chan bool
	wg    *sync.WaitGroup
}

// NewCommander returns a Commander with given commands.
func NewCommander(cmds []string) *Commander {
	return &Commander{
		cmds:  cmds,
		index: 0,
		kill:  make(chan bool),
		wg:    &sync.WaitGroup{},
	}
}

func newCmd(cmd string) (*exec.Cmd, error) {
	parsedCmd, err := shellquote.Split(cmd)
	if err != nil {
		return nil, err
	}

	if len(parsedCmd) == 0 {
		return nil, fmt.Errorf("command cannot be empty")
	}

	var c *exec.Cmd

	if len(parsedCmd) == 1 {
		c = exec.Command(parsedCmd[0]) // nolint:gosec
	} else {
		name := parsedCmd[0]
		args := parsedCmd[1:]
		c = exec.Command(name, args...) // nolint:gosec
	}

	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	c.Stdin = os.Stdin
	c.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	return c, nil
}

// Run executes the commands in order.
func (c *Commander) Run() error {
	c.wg.Add(1)
	defer c.wg.Done()

	var err error

	for _, command := range c.cmds {
		c.cmd, err = newCmd(command)
		if err != nil {
			c.reset()
			continue
		}

		logrus.Debugln("Running:", c.cmd.String())
		if err := c.cmd.Start(); err != nil {
			c.reset()
			return err
		}

		c.cmd.Wait() // nolint:errcheck,gosec
		select {
		case <-c.kill:
			goto killRun
		default:
			continue
		}
	}

killRun:
	c.reset()
	return nil
}

func (c *Commander) reset() {
	c.cmd = nil
}

// Kill stops the execution of current command and terminates the Run.
func (c *Commander) Kill() error {
	if c.cmd == nil {
		return nil
	}

	if err := syscall.Kill(-c.cmd.Process.Pid, syscall.SIGKILL); err != nil {
		return err
	}
	c.kill <- true
	c.wg.Wait()
	c.reset()
	return nil
}
