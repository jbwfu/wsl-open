package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
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

func isWSL() bool {
	return os.Getenv("WSL_INTEROP") != ""
}

func setupUsage(fs *flag.FlagSet) {
	progName := filepath.Base(fs.Name())

	fs.Usage = func() {
		output := fs.Output()

		fmt.Fprintf(output, "%s - A utility to open files, directories, and URLs from WSL in Windows.\n\n", progName)
		fmt.Fprintf(output, "USAGE:\n")
		fmt.Fprintf(output, "    %s [OPTIONS] <URL_or_FILE_PATH>\n\n", progName)
		fmt.Fprintf(output, "ARGUMENTS:\n")
		fmt.Fprintf(output, "    <URL_or_FILE_PATH>\n")
		fmt.Fprintf(output, "        The target to open. This can be a WSL path to a file or directory\n")
		fmt.Fprintf(output, "        (e.g., './document.txt', '.') or a full URL (e.g., 'https://google.com').\n\n")
		fmt.Fprintf(output, "OPTIONS:\n")
		fs.PrintDefaults()
		fmt.Fprintf(output, "\nEXAMPLES:\n")
		fmt.Fprintf(output, "    # Open a file in its default Windows application:\n")
		fmt.Fprintf(output, "    %s notes.txt\n\n", progName)
		fmt.Fprintf(output, "    # Open the current directory in Windows File Explorer:\n")
		fmt.Fprintf(output, "    %s .\n\n", progName)
		fmt.Fprintf(output, "    # Open a URL in the default Windows browser:\n")
		fmt.Fprintf(output, "    %s https://github.com\n", progName)
	}
}

func run(args []string) error {
	fs := flag.NewFlagSet(args[0], flag.ExitOnError)
	quiet := fs.Bool("q", false, "Enable quiet mode, suppressing informational output.")
	dryRun := fs.Bool("x", false, "Perform a dry run, printing the command without executing it.")
	setupUsage(fs)

	fs.Parse(args[1:])

	if !isWSL() {
		return errors.New("this tool requires a WSL environment with Windows interoperability enabled")
	}

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
