package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"syscall"
)

type Command struct {
	Name   string
	Args   []string
	Stdin  *os.File
	Stdout *os.File
	Stderr *os.File
}

func (cmd Command) String() string {
	return fmt.Sprintf("%s %s", cmd.Name, strings.Join(cmd.Args, " "))
}

// Run runs command with `Exec` on platforms except Windows
// which only supports `Spawn`
func (cmd *Command) Run() error {
	if isWindows() {
		return cmd.Spawn()
	} else {
		return cmd.Exec()
	}
}

func isWindows() bool {
	return runtime.GOOS == "windows" || detectWSL()
}

// https://github.com/Microsoft/WSL/issues/423#issuecomment-221627364
func detectWSL() bool {
	var detectedWSLContents string
	b := make([]byte, 1024)
	f, err := os.Open("/proc/version")
	if err == nil {
		f.Read(b)
		f.Close()
		detectedWSLContents = string(b)
	}
	return strings.Contains(detectedWSLContents, "Microsoft")
}

// Spawn runs command with spawn(3)
func (cmd *Command) Spawn() error {
	c := exec.Command(cmd.Name, cmd.Args...)
	c.Stdin = cmd.Stdin
	c.Stdout = cmd.Stdout
	c.Stderr = cmd.Stderr
	return c.Run()
}

// Exec runs command with exec(3)
// Note that Windows doesn't support exec(3): http://golang.org/src/pkg/syscall/exec_windows.go#L339
func (cmd *Command) Exec() error {
	binary, err := exec.LookPath(cmd.Name)
	if err != nil {
		return &exec.Error{
			Name: cmd.Name,
			Err:  fmt.Errorf("command not found"),
		}
	}
	args := []string{binary}
	args = append(args, cmd.Args...)
	return syscall.Exec(binary, args, os.Environ())
}
