package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

func main() {

	// main is the entry point of the application
	StartPrompt(">> ", os.Stdin, os.Stdout)
}

// StartPrompt starts the interactive prompt with a given prompt string
// It reads commands from the input (in) and writes output to (out)

func StartPrompt(prompt string, in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	for {

		// Print the prompt
		fmt.Fprint(out, prompt)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		// Read the user input
		line := scanner.Text()
		if line == "exit" || line == "" {
			return
		}

		// Split the input into command and arguments
		parts := strings.Fields(line)
		if len(parts) == 0 {
			return
		}

		command := parts[0]
		args := parts[1:]

		//handler built-in commands
		switch command {
		case "ls":
			// Use 'dir' instead of 'ls' on Windows
			if runtime.GOOS == "windows" {
				command = "cmd"
				args = append([]string{"/C", "dir"}, args...)
			}
			executeCommand(command, args, in, out)
		case "mkdir":
			// Create a directory
			if len(args) != 1 {
				fmt.Fprintln(out, "Usage: mkdir <directory>")
			} else {
				err := os.Mkdir(args[0], 0755)
				if err != nil {
					fmt.Fprintln(out, "Error:", err.Error())
				}
			}
		case "touch":
			// Create a file
			if len(args) != 1 {
				fmt.Fprintln(out, "Usage: touch <file>")
			} else {
				file, err := os.Create(args[0])
				if err != nil {
					fmt.Fprintln(out, "Error:", err.Error())
				} else {
					file.Close()
				}
			}
		default:
				// Execute any other commands
			executeCommand(command, args, in, out)
		}
	}
}

// executeCommand runs the specified command with arguments, 
// redirecting input, output, and error streams
func executeCommand(command string, args []string, in io.Reader, out io.Writer) {
	runner := exec.Command(command, args...)
	runner.Stdin = in
	runner.Stdout = out
	runner.Stderr = out

	err := runner.Run()
	if err != nil {
		fmt.Fprintln(out, "Error:", err.Error())
	}
}
