package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func toWindowsPath(path string) (string, error) {
	wslpath, err := exec.LookPath("wslpath")
	if err != nil {
		return "", errors.New("command 'wslpath' not found in your PATH. Please ensure WSL is installed correctly")
	}
	cmd := exec.Command(wslpath, "-w", path)
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("wslpath failed to convert '%s' to a Windows path: %w", path, err)
	}
	return strings.TrimSpace(string(output)), nil
}

func toWslPath(path string) (string, error) {
	wslpath, err := exec.LookPath("wslpath")
	if err != nil {
		return "", errors.New("command 'wslpath' not found in your PATH. Please ensure WSL is installed correctly")
	}
	cmd := exec.Command(wslpath, "-u", path)
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("wslpath failed to convert '%s' to a WSL path: %w", path, err)
	}
	return strings.TrimSpace(string(output)), nil
}

func getStartCmdPath() (string, error) {
	path, err := exec.LookPath("powershell.exe")
	if err == nil {
		return path, nil
	}

	systemRoot := os.Getenv("SystemRoot")
	if systemRoot == "" {
		systemRoot = "C:\\Windows"
	}

	defaultPath := systemRoot + "\\System32\\WindowsPowerShell\\v1.0\\powershell.exe"
	return toWslPath(defaultPath)
}

func run(args []string) error {
	fs := flag.NewFlagSet(args[0], flag.ExitOnError)
	quiet := fs.Bool("q", false, "Enable quiet mode, suppressing informational output.")
	dryRun := fs.Bool("x", false, "Perform a dry run, printing the command without executing it.")

	fs.Parse(args[1:])
	if fs.NArg() != 1 {
		fs.Usage()
		return fmt.Errorf("invalid arguments: exactly one URL or file path is required")
	}
	input := fs.Arg(0)

	var target string
	var err error

	if strings.HasPrefix(input, "http://") || strings.HasPrefix(input, "https://") {
		target = input
	} else {
		target, err = toWindowsPath(input)
		if err != nil {
			return fmt.Errorf("failed to convert WSL path '%s': %w", input, err)
		}
	}

	psExePath, err := getStartCmdPath()
	if err != nil {
		return fmt.Errorf("could not locate powershell.exe: %w", err)
	}

	sanitizedTarget := strings.ReplaceAll(target, "'", "''")
	command := fmt.Sprintf("Start-Process -FilePath '%s'", sanitizedTarget)
	cmd := exec.Command(psExePath, "-Command", command)

	if *dryRun {
		fmt.Println("Dry Run: Would execute command:", cmd.String())
		return nil
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		if len(output) > 0 {
			fmt.Fprintf(os.Stderr, "Command output:\n%s\n", string(output))
		}
		return fmt.Errorf("command execution failed for '%s': %w", target, err)
	}

	if !*quiet && len(output) > 0 {
		fmt.Print(string(output))
	}

	return nil
}

func main() {
	if err := run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
