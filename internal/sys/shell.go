package sys

import (
	"os"
	"os/exec"
	"runtime"
)

// ShellCommand returns an exec.Cmd that evaluates the command string using the host's default shell.
func ShellCommand(commandStr string) *exec.Cmd {
	if runtime.GOOS == "windows" {
		return exec.Command("cmd", "/C", commandStr)
	}
	return exec.Command("sh", "-c", commandStr)
}

// DefaultShell returns the default interactive shell for the host OS.
func DefaultShell() string {
	shell := os.Getenv("SHELL")
	if shell == "" {
		if runtime.GOOS == "windows" {
			return "cmd.exe"
		}
		return "sh"
	}
	return shell
}
