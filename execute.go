package main

import (
	"fmt"
	"os"
	"os/exec"
)

const (
	white = "\033[1;37m"
	red   = "\033[1;31m"
	reset = "\033[0m"
)

var commands []string
var cleanups []string

func system(command string) int {
	cmd := exec.Command("/bin/sh", "-c", command)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		fmt.Printf("%s-> %s%s\n", red, err, reset)
		return err.(*exec.ExitError).ExitCode()
	}
	return 0
}

func execute() {
	var exitcode = 0

	for i, command := range commands {
		if i > 0 {
			fmt.Println("")
		}
		fmt.Printf("%s=> %s%s\n", white, command, reset)
		exitcode = system(command)
		if exitcode != 0 {
			break
		}
	}

	for _, cleanup := range cleanups {
		fmt.Printf("\n%s=> %s%s\n", white, cleanup, reset)
		system(cleanup)
	}

	os.Exit(exitcode)
}
