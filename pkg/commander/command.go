package commander

import (
	"context"
	"errors"
	"os"
	"os/exec"
	"syscall"
)

// Command is an external command that can be executed.
type Command struct {
	Name string
	Args []string

	Stdin, Stdout, Stderr *os.File
}

// String returns the command in a human readable format.
func (c Command) String() string {
	str := c.Name
	for i, arg := range c.Args {
		str = str + arg
		if i != len(c.Args)-1 {
			str = str + " "
		}
	}

	return str
}

// Execute runs the commands and exits elegantly when the context is canceled.
func (c Command) Execute(ctx context.Context) error {
	if c.Name == "" {
		return errors.New("command name cannot be empty")
	}

	stream := make(chan error, 1)

	cmd := exec.Command(c.Name, c.Args...) // nolint:gosec
	cmd.Stdout = c.Stdout
	cmd.Stderr = c.Stderr
	cmd.Stdin = c.Stdin
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	if err := cmd.Start(); err != nil {
		return err
	}

	go func(ex *exec.Cmd, err chan<- error) {
		err <- ex.Wait()
	}(cmd, stream)

	select {
	case <-ctx.Done():
		// Elegantly close the parent along-with the children.
		if err := syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL); err != nil {
			return err
		}

		return nil

	case err := <-stream:
		return err
	}
}
