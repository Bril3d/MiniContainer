package runtime

import (
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/fatih/color"
)

var debugMode bool

// SetDebug enables or disables global debug output for the runtime.
func SetDebug(enabled bool) {
	debugMode = enabled
}

// Exec runs a command and returns its combined output.
// If debug mode is on, it prints the command and execution time.
func Exec(cmd string, args ...string) (string, error) {
	if debugMode {
		color.Cyan("DEBUG: %s %s", cmd, strings.Join(args, " "))
	}

	start := time.Now()
	out, err := exec.Command(cmd, args...).CombinedOutput()
	duration := time.Since(start)

	if debugMode {
		color.Cyan("DEBUG: Finished in %v", duration)
	}

	return strings.TrimSpace(string(out)), err
}

// ExecFull runs a command and returns its combined output and exit code.
// If debug mode is on, it prints the command and execution time.
func ExecFull(cmd string, args ...string) (string, int, error) {
	if debugMode {
		color.Cyan("DEBUG: %s %s", cmd, strings.Join(args, " "))
	}

	start := time.Now()
	command := exec.Command(cmd, args...)
	out, err := command.CombinedOutput()
	duration := time.Since(start)

	if debugMode {
		color.Cyan("DEBUG: Finished in %v", duration)
	}

	exitCode := 0
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			exitCode = exitError.ExitCode()
		} else {
			exitCode = 1
		}
	}

	return strings.TrimSpace(string(out)), exitCode, err
}

// ExecInteractive runs a command attached to the user's terminal (stdin/stdout/stderr).
// Used for interactive containers (e.g. `mini run -i alpine sh`).
func ExecInteractive(command string, args ...string) error {
	if debugMode {
		color.Cyan("DEBUG (Interactive): %s %s", command, strings.Join(args, " "))
	}
	cmd := exec.Command(command, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// ExecStream runs a command and streams stdout/stderr to os.Stdout/os.Stderr in real time.
// Used for operations like `podman pull` where you want live progress output.
func ExecStream(command string, args ...string) error {
	if debugMode {
		color.Cyan("DEBUG (Stream): %s %s", command, strings.Join(args, " "))
	}
	cmd := exec.Command(command, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
