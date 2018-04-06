// Copyright 2015 The Gogs Authors. All rights reserved.
// Copyright 2018 Unknwon. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package command

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"time"
)

var (
	// Debug enables verbose logging on everything.
	// This should be false in case Gogs starts in SSH mode.
	Debug  = true
	Prefix = "[command] "
)

func log(format string, args ...interface{}) {
	if !Debug {
		return
	}

	fmt.Print(Prefix)
	if len(args) == 0 {
		fmt.Println(format)
	} else {
		fmt.Printf(format+"\n", args...)
	}
}

func concatenateError(err error, stderr string) error {
	if len(stderr) == 0 {
		return err
	}
	return fmt.Errorf("%v - %s", err, stderr)
}

type ErrExecTimeout struct {
	Duration time.Duration
}

func IsErrExecTimeout(err error) bool {
	_, ok := err.(ErrExecTimeout)
	return ok
}

func (err ErrExecTimeout) Error() string {
	return fmt.Sprintf("execution is timeout [duration: %v]", err.Duration)
}

// Command represents a command with its subcommands or arguments.
type Command struct {
	name string
	args []string
	envs []string
}

func (c *Command) String() string {
	if len(c.args) == 0 {
		return c.name
	}
	return fmt.Sprintf("%s %s", c.name, strings.Join(c.args, " "))
}

// New creates and returns a new Command object based on given command and arguments.
func New(command string, args ...string) *Command {
	return &Command{
		name: command,
		args: args,
	}
}

// AddArguments adds new argument(s) to the command.
func (c *Command) AddArguments(args ...string) *Command {
	c.args = append(c.args, args...)
	return c
}

// AddEnvs adds new environment variables to the command.
func (c *Command) AddEnvs(envs ...string) *Command {
	c.envs = append(c.envs, envs...)
	return c
}

const DEFAULT_TIMEOUT = 60 * time.Second

// RunInDirTimeoutPipeline executes the command in given directory with given timeout,
// it pipes stdout and stderr to given io.Writer.
func (c *Command) RunInDirTimeoutPipeline(timeout time.Duration, dir string, stdout, stderr io.Writer) error {
	if timeout == -1 {
		timeout = DEFAULT_TIMEOUT
	}

	if len(dir) == 0 {
		log(c.String())
	} else {
		log("%s: %v", dir, c)
	}

	cmd := exec.Command(c.name, c.args...)
	if c.envs != nil {
		cmd.Env = append(os.Environ(), c.envs...)
	}
	cmd.Dir = dir
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	if err := cmd.Start(); err != nil {
		return err
	}

	done := make(chan error)
	go func() {
		done <- cmd.Wait()
	}()

	var err error
	select {
	case <-time.After(timeout):
		if cmd.Process != nil && cmd.ProcessState != nil && !cmd.ProcessState.Exited() {
			if err := cmd.Process.Kill(); err != nil {
				return fmt.Errorf("fail to kill process: %v", err)
			}
		}

		<-done
		return ErrExecTimeout{timeout}
	case err = <-done:
	}

	return err
}

// RunInDirTimeout executes the command in given directory with given timeout,
// and returns stdout in []byte and error (combined with stderr).
func (c *Command) RunInDirTimeout(timeout time.Duration, dir string) ([]byte, error) {
	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)
	if err := c.RunInDirTimeoutPipeline(timeout, dir, stdout, stderr); err != nil {
		return nil, concatenateError(err, stderr.String())
	}

	if stdout.Len() > 0 {
		log("stdout:\n%s", stdout.Bytes()[:1024])
	}
	return stdout.Bytes(), nil
}

// RunInDirPipeline executes the command in given directory,
// it pipes stdout and stderr to given io.Writer.
func (c *Command) RunInDirPipeline(dir string, stdout, stderr io.Writer) error {
	return c.RunInDirTimeoutPipeline(-1, dir, stdout, stderr)
}

// RunInDir executes the command in given directory
// and returns stdout in []byte and error (combined with stderr).
func (c *Command) RunInDirBytes(dir string) ([]byte, error) {
	return c.RunInDirTimeout(-1, dir)
}

// RunInDir executes the command in given directory
// and returns stdout in string and error (combined with stderr).
func (c *Command) RunInDir(dir string) (string, error) {
	stdout, err := c.RunInDirTimeout(-1, dir)
	if err != nil {
		return "", err
	}
	return string(stdout), nil
}

// RunTimeout executes the command in defualt working directory with given timeout,
// and returns stdout in string and error (combined with stderr).
func (c *Command) RunTimeout(timeout time.Duration) (string, error) {
	stdout, err := c.RunInDirTimeout(timeout, "")
	if err != nil {
		return "", err
	}
	return string(stdout), nil
}

// Run executes the command in defualt working directory
// and returns stdout in string and error (combined with stderr).
func (c *Command) Run() (string, error) {
	return c.RunTimeout(-1)
}
