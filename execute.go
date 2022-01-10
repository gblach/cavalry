/* Copyright 2022 Grzegorz Blach

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/
package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"regexp"
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

	re := regexp.MustCompile("\033\\[[0-9;]+m")
	buffer_bytes := re.ReplaceAll(buffer.Bytes(), []byte{})

	cmd := exec.Command(sendmail_cmd, arg_mailto)
	cmd.Stdin = io.MultiReader(
		strings.NewReader(subject+"\n\n"),
		bytes.NewReader(buffer_bytes))
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
