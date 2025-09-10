package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"
)

// fakeExecCommand creates a new exec.Cmd that executes the current test binary
// again, enabling the TestHelperProcess to act as a mock for the real command.
func fakeExecCommand(command string, args ...string) *exec.Cmd {
	cs := []string{"-test.run=TestHelperProcess", "--", command}
	cs = append(cs, args...)
	cmd := exec.Command(os.Args[0], cs...)
	// This environment variable is used to signal to the re-executed test
	// binary that it should behave as the helper process.
	cmd.Env = []string{"GO_WANT_HELPER_PROCESS=1"}
	return cmd
}

// TestHelperProcess isn't a real test. It's a helper process that acts as a mock
// for external commands like 'wslpath'. It's executed by fakeExecCommand.
// It checks the command it's supposed to impersonate and prints the expected
// output to stdout or stderr, then exits.
func TestHelperProcess(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}
	defer os.Exit(0)

	// The actual command and its arguments are passed after "--".
	args := os.Args
	for len(args) > 0 {
		if args[0] == "--" {
			args = args[1:]
			break
		}
		args = args[1:]
	}
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "No command to mock")
		os.Exit(1)
	}

	command, args := args[0], args[1:]

	if strings.HasSuffix(command, "wslpath") {
		if len(args) < 2 {
			fmt.Fprintln(os.Stderr, "wslpath mock needs arguments")
			os.Exit(1)
		}
		pathArg := args[1]
		if pathArg == "/mnt/c/Users" {
			fmt.Fprint(os.Stdout, `C:\Users`)
			return
		}
		if pathArg == "C:\\Windows\\System32\\WindowsPowerShell\\v1.0\\powershell.exe" {
			fmt.Fprint(os.Stdout, `/mnt/c/Windows/System32/WindowsPowerShell/v1.0/powershell.exe`)
			return
		}
		if pathArg == "/path/to/fail" {
			fmt.Fprintln(os.Stderr, "wslpath error: simulated failure")
			os.Exit(1)
		}
	}

	fmt.Fprintf(os.Stderr, "Unknown command to mock or wrong arguments: %s %v\n", command, args)
	os.Exit(1)
}

func TestIsWSL(t *testing.T) {
	tests := []struct {
		name     string
		envValue string
		setEnv   bool
		want     bool
	}{
		{"WSL_INTEROP is set", "/run/WSL/1_interop", true, true},
		{"WSL_INTEROP is empty", "", true, false},
		{"WSL_INTEROP is not set", "", false, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setEnv {
				t.Setenv("WSL_INTEROP", tt.envValue)
			} else {
				os.Unsetenv("WSL_INTEROP")
			}

			if got := isWSL(); got != tt.want {
				t.Errorf("isWSL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToWindowsPath(t *testing.T) {
	app := &App{
		LookPath: func(name string) (string, error) {
			return "/fake/path/to/" + name, nil
		},
		Command: fakeExecCommand,
	}

	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{"Valid WSL path", "/mnt/c/Users", `C:\Users`, false},
		{"Path that causes wslpath to fail", "/path/to/fail", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := app.toWindowsPath(tt.input)

			if (err != nil) != tt.wantErr {
				t.Fatalf("toWindowsPath() error = %v, wantErr %v", err, tt.wantErr)
			}
			if got != tt.want {
				t.Errorf("toWindowsPath() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestGetStartCmdPath(t *testing.T) {
	app := &App{
		LookPath: func(name string) (string, error) {
			if name == "powershell.exe" {
				return "", errors.New("not found")
			}
			return "/fake/path/to/wslpath", nil
		},
		Command: fakeExecCommand,
	}

	t.Run("Fallback finds powershell via SystemRoot", func(t *testing.T) {
		t.Setenv("SystemRoot", "C:\\Windows")

		got, err := app.getStartCmdPath()
		if err != nil {
			t.Fatalf("getStartCmdPath() returned an unexpected error: %v", err)
		}

		want := "/mnt/c/Windows/System32/WindowsPowerShell/v1.0/powershell.exe"
		if got != want {
			t.Errorf("getStartCmdPath() = %q, want %q", got, want)
		}
	})
}
