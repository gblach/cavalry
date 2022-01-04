package main

import (
	"bytes"
	"fmt"
	"io"
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
var buffer bytes.Buffer
var output io.Writer

func run(args []string) int {
	fmt.Fprintln(output, white+"=>", arg_engine, strings.Join(args, " "), reset)

	cmd := exec.Command(arg_engine, args...)
	cmd.Stdout = output
	cmd.Stderr = output
	err := cmd.Run()

	if err != nil {
		fmt.Fprint(output, red+"-> ", err, reset, "\n\n")
		return err.(*exec.ExitError).ExitCode()
	}

	fmt.Fprintln(output, "")
	return 0
}

func sendmail(exitcode int) {
	sendmail_cmd := os.Getenv("SENDMAIL_CMD")
	if sendmail_cmd == "" {
		sendmail_cmd = "/usr/sbin/sendmail"
	}

	subject := "Subject: " + strings.Join(os.Args, " ")
	if exitcode == 0 {
		subject += " (success)"
	} else {
		subject += fmt.Sprintf(" (fail: %d)", exitcode)
	}

	cmd := exec.Command(sendmail_cmd, arg_mailto)
	cmd.Stdin = io.MultiReader(
		strings.NewReader(subject+"\n\n"),
		bytes.NewReader(buffer.Bytes()))
	cmd.Run()
}

func execute() {
	var exitcode = 0

	output = io.MultiWriter(os.Stdout, &buffer)

	for _, command := range commands {
		exitcode = run(command)
		if exitcode != 0 {
			break
		}
	}

	for _, cleanup := range cleanups {
		run(cleanup)
	}

	if arg_mailto != "" && (arg_mailalw || exitcode != 0) {
		sendmail(exitcode)
	}

	os.Exit(exitcode)
}
