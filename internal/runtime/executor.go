package runtime

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// ExecResult holds the result of a command execution.
type ExecResult struct {
	Stdout   string
	Stderr   string
	ExitCode int
}

// Exec runs a command and captures stdout, stderr, and exit code.
// Returns the trimmed stdout on success, or an error with stderr details.
func Exec(command string, args ...string) (string, error) {
	cmd := exec.Command(command, args...)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	outStr := strings.TrimSpace(stdout.String())
	errStr := strings.TrimSpace(stderr.String())

	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return outStr, fmt.Errorf("command failed (exit %d): %s", exitErr.ExitCode(), errStr)
		}
		return outStr, fmt.Errorf("command execution error: %w", err)
	}

	return outStr, nil
}

// ExecFull runs a command and returns the full ExecResult.
func ExecFull(command string, args ...string) ExecResult {
	cmd := exec.Command(command, args...)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	result := ExecResult{
		Stdout: strings.TrimSpace(stdout.String()),
		Stderr: strings.TrimSpace(stderr.String()),
	}

	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			result.ExitCode = exitErr.ExitCode()
		} else {
			result.ExitCode = -1
		}
	}

	return result
}

// ExecInteractive runs a command attached to the user's terminal (stdin/stdout/stderr).
// Used for interactive containers (e.g. `mini run -i alpine sh`).
func ExecInteractive(command string, args ...string) error {
	cmd := exec.Command(command, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// ExecStream runs a command and streams stdout/stderr to os.Stdout/os.Stderr in real time.
// Used for operations like `podman pull` where you want live progress output.
func ExecStream(command string, args ...string) error {
	cmd := exec.Command(command, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
