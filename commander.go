package leaf

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"syscall"

	"github.com/kballard/go-shellquote"
)

var errEmptyCmd = fmt.Errorf("empty command")

// Commander has a set of commands that run in order
// and exit when the context is canceled.
type Commander struct {
	Commands []string

	OnStart func(*Command)
	OnError func(error)
	OnExit  func()

	ExitOnError bool

	done chan bool
}

// NewCommander creates a new commander.
func NewCommander(commander Commander) *Commander {
	return &Commander{
		Commands:    commander.Commands,
		OnStart:     commander.OnStart,
		OnError:     commander.OnError,
		OnExit:      commander.OnExit,
		ExitOnError: commander.ExitOnError,
		done:        make(chan bool, 1),
	}
}

// Done signals that the commander is done running the commands.
func (c *Commander) Done() <-chan bool {
	return c.done
}

// Run executes the commands in order. It stops when the
// context is canceled.
func (c *Commander) Run(ctx context.Context) {
	defer func() {
		// signal done when running commands is complete
		// or the function exits, eitherway.
		c.done <- true
		c.OnExit()
	}()

	for _, command := range c.Commands {
		cmd, err := NewCommand(command)
		if err != nil {
			c.OnError(err)
			return
		}

		select {
		case <-ctx.Done():
			return

		default:
			if cmd == nil {
				continue
			}

			if c.OnStart != nil {
				c.OnStart(cmd)
			}

			if err := cmd.Execute(ctx); err != nil {
				if c.OnError != nil {
					c.OnError(err)
				}
				if c.ExitOnError {
					return
				}
			}
		}
	}
}

// Command is an external command that can be executed.
type Command struct {
	Name string
	Args []string

	str string
}

// String returns the command in a human-readable format.
func (c *Command) String() string {
	return c.str
}

// NewCommand creates a new command from the string.
func NewCommand(cmd string) (*Command, error) {
	parsedCmd, err := shellquote.Split(cmd)
	if err != nil {
		return nil, err
	}

	if len(parsedCmd) == 0 {
		return nil, errEmptyCmd
	}

	name := parsedCmd[0]
	var args []string
	if len(parsedCmd) > 1 {
		args = parsedCmd[1:]
	}

	return &Command{
		Name: name,
		Args: args,
		str:  shellquote.Join(parsedCmd...),
	}, nil
}

// Execute runs the commands and exits elegantly when the
// context is canceled.
//
// This doesn't use the exec.CommandContext because we just
// don't want to kill the parent process but all the child
// processes too.
func (c *Command) Execute(ctx context.Context) error {
	stream := make(chan error)

	cmd := exec.Command(c.Name, c.Args...) // nolint:gosec
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
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
		err := syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
		if err != nil {
			return err
		}

		return nil

	case err := <-stream:
		return err
	}
}
