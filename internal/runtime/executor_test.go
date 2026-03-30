package runtime

import (
	"runtime"
	"strings"
	"testing"
)

func TestExec_SimpleCommand(t *testing.T) {
	var out string
	var err error

	if runtime.GOOS == "windows" {
		out, err = Exec("cmd", "/c", "echo", "hello")
	} else {
		out, err = Exec("echo", "hello")
	}

	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if out != "hello" {
		t.Fatalf("expected 'hello', got: '%s'", out)
	}
}

func TestExec_CommandNotFound(t *testing.T) {
	_, err := Exec("nonexistent_command_12345")
	if err == nil {
		t.Fatal("expected error for nonexistent command, got nil")
	}
}

func TestExec_CommandFailure(t *testing.T) {
	var _, err error

	if runtime.GOOS == "windows" {
		_, err = Exec("cmd", "/c", "exit", "1")
	} else {
		_, err = Exec("sh", "-c", "exit 1")
	}

	if err == nil {
		t.Fatal("expected error for failing command, got nil")
	}
	if !strings.Contains(err.Error(), "exit") {
		t.Fatalf("expected error to mention exit, got: %v", err)
	}
}

func TestExecFull_CapturesStdoutAndExitCode(t *testing.T) {
	var result ExecResult

	if runtime.GOOS == "windows" {
		result = ExecFull("cmd", "/c", "echo", "test-output")
	} else {
		result = ExecFull("echo", "test-output")
	}

	if result.ExitCode != 0 {
		t.Fatalf("expected exit code 0, got: %d", result.ExitCode)
	}
	if result.Stdout != "test-output" {
		t.Fatalf("expected 'test-output', got: '%s'", result.Stdout)
	}
}

func TestExecFull_CapturesNonZeroExit(t *testing.T) {
	var result ExecResult

	if runtime.GOOS == "windows" {
		result = ExecFull("cmd", "/c", "exit", "42")
	} else {
		result = ExecFull("sh", "-c", "exit 42")
	}

	if result.ExitCode != 42 {
		t.Fatalf("expected exit code 42, got: %d", result.ExitCode)
	}
}
