package main

import (
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

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <URL_or_FILE_PATH>\n", os.Args[0])
		os.Exit(1)
	}
	input := os.Args[1]

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

	cmdExePath, err := toWslPath("C:\\Windows\\System32\\cmd.exe")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Could not locate cmd.exe via wslpath: %v\n", err)
		os.Exit(1)
	}

	cmd := exec.Command(cmdExePath, "/C", "start", target)
	cmd.Dir, err = toWslPath("C:\\")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Unable to locate a non-UNC path as the execution path")
	}

	output, err := cmd.CombinedOutput()

	if len(output) > 0 {
		fmt.Print(string(output))
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Command execution failed for '%s': %v\n", target, err)
		os.Exit(1)
	}
}
