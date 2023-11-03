/*
 * Copyright (c) 2023. yimincai(Neil) <bravc29229@gmail.com>.
 */

package command

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"

	"github.com/fatih/color"
)

// By default, shell will be using bash
func commandOut(command string, shell string) (string, string, error) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	if shell == "" {
		shell = "bash"
	}

	cmd := exec.Command(shell, "-c", command)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	return stdout.String(), stderr.String(), err
}

// Run executes a command and prints it to stdout
func Run(command string, showCommand bool) {
	if showCommand {
		color.HiBlue("$ %v", command)
	}
	out, errOut, err := commandOut(command, "")
	if err != nil || errOut != "" {
		if errOut != "" {
			color.Yellow("warning: %v", errOut)
		}
		if err != nil {
			color.Red("error: %v", err)
			os.Exit(1)
		}
	}
	if out != "" {
		fmt.Println(out)
	}
}

// Return executes a command and returns the output
func Return(command string, showCommand bool) string {
	if showCommand {
		color.HiBlue("$ %v", command)
	}
	out, errOut, err := commandOut(command, "")
	if err != nil || errOut != "" {
		if errOut != "" {
			color.Yellow("warning: %v", errOut)
		}
		if err != nil {
			color.Red("error: %v", err)
			os.Exit(1)
		}
	}
	if out != "" {
		fmt.Println(out)
	}

	return out
}
