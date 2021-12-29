package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

const (
	white = "\033[1;37m"
	red   = "\033[1;31m"
	reset = "\033[0m"
)

var commands [][]string
var cleanups [][]string

func run(args []string) int {
	fmt.Println(white+"=>", arg_engine, strings.Join(args, " "), reset)

	cmd := exec.Command(arg_engine, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()

	if err != nil {
		fmt.Println(red+"->", err, reset, "\n")
		return err.(*exec.ExitError).ExitCode()
	}

	fmt.Println("")
	return 0
}

func execute() {
	var exitcode = 0

	for _, command := range commands {
		exitcode = run(command)
		if exitcode != 0 {
			break
		}
	}

	for _, cleanup := range cleanups {
		run(cleanup)
	}

	os.Exit(exitcode)
}
