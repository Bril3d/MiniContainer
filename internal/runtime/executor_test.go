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
}

func TestExecFull_CapturesStdoutAndExitCode(t *testing.T) {
	var out string
	var code int

	if runtime.GOOS == "windows" {
		out, code, _ = ExecFull("cmd", "/c", "echo", "test-output")
	} else {
		out, code, _ = ExecFull("echo", "test-output")
	}

	if code != 0 {
		t.Fatalf("expected exit code 0, got: %d", code)
	}
	if strings.TrimSpace(out) != "test-output" {
		t.Fatalf("expected 'test-output', got: '%s'", out)
	}
}

func TestExecFull_CapturesNonZeroExit(t *testing.T) {
	var code int

	if runtime.GOOS == "windows" {
		_, code, _ = ExecFull("cmd", "/c", "exit", "42")
	} else {
		_, code, _ = ExecFull("sh", "-c", "exit 42")
	}

	if code != 42 {
		t.Fatalf("expected exit code 42, got: %d", code)
	}
}
