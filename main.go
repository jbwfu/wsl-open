package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func toWindowsPath(path string) (string, error) {
	cmd := exec.Command("wslpath", "-w", path)
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to execute 'wslpath' command. Is it in your PATH? Error: %w", err)
	}
	return strings.TrimSpace(string(output)), nil
}

func toWslPath(path string) (string, error) {
	cmd := exec.Command("wslpath", "-u", path)
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to execute 'wslpath' command. Is it in your PATH? Error: %w", err)
	}
	return strings.TrimSpace(string(output)), nil
}

func getStartCmdPath(quiet bool) (string, error) {
	path, err := exec.LookPath("powershell.exe")
	if err == nil {
		return path, nil
	}

	if !quiet {
		fmt.Fprintln(os.Stderr, "Info: 'powershell.exe' not found in PATH. Falling back to default known location.")
	}

	defaultPath := "C:\\Windows\\System32\\WindowsPowerShell\\v1.0\\powershell.exe"
	return toWslPath(defaultPath)
}

func main() {
	var quiet bool
	flag.BoolVar(&quiet, "q", false, "Enable quiet mode, suppressing informational output.")
	flag.Parse()

	if flag.NArg() != 1 {
		fmt.Fprintf(os.Stderr, "Usage: %s [-q] <URL_or_FILE_PATH>\n", os.Args[0])
		os.Exit(1)
	}
	input := flag.Arg(0)

	var target string
	var err error

	if strings.HasPrefix(input, "http://") || strings.HasPrefix(input, "https://") {
		target = input
	} else {
		target, err = toWindowsPath(input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: Failed to convert WSL path: %v\n", err)
			os.Exit(1)
		}
	}

	psExePath, err := getStartCmdPath(quiet)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Could not locate powershell.exe: %v\n", err)
		os.Exit(1)
	}

	cmd := exec.Command(psExePath, "start", "\""+target+"\"")
	cmd.Dir, err = toWslPath("C:\\")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Unable to locate a non-UNC path as the execution path")
	}

	output, err := cmd.CombinedOutput()

	if !quiet && len(output) > 0 {
		fmt.Print(string(output))
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Command execution failed for '%s': %v\n", target, err)
		os.Exit(1)
	}
}
