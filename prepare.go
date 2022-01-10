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
	"bufio"
	"crypto/rand"
	"fmt"
	"os"
	"path"
	"strings"
)

type container_t struct {
	name string
	tag  string
	dir  string
	file string
	env  []string
	keep bool
	push string
}

type testcase_t struct {
	tag     string
	command string
}

var containers []container_t
var testcases []testcase_t
var container_names = map[string]string{}

func genname(tag string) string {
	s := strings.Replace(tag, ":", ".", -1)
	x, _ := rand.Prime(rand.Reader, 80)
	container_names[tag] = fmt.Sprintf("%s.%x", s, x)
	return container_names[tag]
}

func getname(tag string) string {
	return container_names[tag]
}

func loadfile(cavalryfile string) {
	file, err := os.Open(cavalryfile)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	var container container_t
	var testcase testcase_t

	i := 0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		i++
		line := strings.Join(strings.Fields(scanner.Text()), " ")

		if strings.HasPrefix(line, "CONTAINER") {
			if container.tag != "" {
				containers = append(containers, container)
			}
			container = container_t{dir: ".", file: "Dockerfile", keep: false}
		}

		n, _ := fmt.Sscanf(line, "CONTAINER %s", &container.tag)
		if n > 0 {
			container.name = genname(container.tag)
			continue
		}

		n, _ = fmt.Sscanf(line, "DIR %s", &container.dir)
		if n > 0 {
			continue
		}

		n, _ = fmt.Sscanf(line, "FILE %s", &container.file)
		if n > 0 {
			continue
		}

		n, _ = fmt.Sscanf(line, "PUSH %s", &container.push)
		if n > 0 {
			continue
		}

		if line == "KEEP" {
			container.keep = true
			continue
		}

		if strings.HasPrefix(line, "ENV ") {
			container.env = append(container.env, line[4:])
			continue
		}

		n, _ = fmt.Sscanf(line, "EXEC %s", &testcase.tag)
		if n > 0 {
			if getname(testcase.tag) == "" {
				panic(fmt.Sprintf("No such container: %s (line %d)", testcase.tag, i))
			}

			start := len(testcase.tag) + 6
			testcase.command = line[start:]
			testcases = append(testcases, testcase)
			continue
		}

		if len(line) > 0 && line[0] != '#' {
			panic(fmt.Sprintf("Wrong input (line %d): %s", i, line))
		}
	}

	containers = append(containers, container)
}

func makeplan() {
	engine_base := path.Base(arg_engine)

	format := []string{}
	if engine_base == "podman" && arg_format != "" {
		format = []string{"--format", arg_format}
	}

	network := []string{"--network", "container:" + containers[0].name}

	for i, container := range containers {
		command := []string{"build", "-t", container.tag, "-f", container.file}
		command = append(command, format...)
		command = append(command, container.dir)
		commands = append(commands, command)

		environ := []string{}
		for _, env := range container.env {
			environ = append(environ, "-e", env)
		}
		command = []string{"run", "-dt", "--name", container.name}
		if i > 0 {
			command = append(command, network...)
		}
		command = append(command, environ...)
		command = append(command, container.tag)
		commands = append(commands, command)
	}

	for _, testcase := range testcases {
		command := []string{"exec", getname(testcase.tag)}
		command = append(command, strings.Fields(testcase.command)...)
		commands = append(commands, command)
	}

	if !arg_nopush {
		for _, container := range containers {
			if container.push != "" {
				if engine_base == "podman" {
					command := []string{"push", container.tag, container.push}
					commands = append(commands, command)
				} else {
					command := []string{"tag", container.tag, container.push}
					commands = append(commands, command)

					command = []string{"push", container.push}
					commands = append(commands, command)
				}
			}
		}
	}

	if !arg_norm {
		for i := len(containers) - 1; i >= 0; i-- {
			cleanup := []string{"rm", "-f", containers[i].name}
			cleanups = append(cleanups, cleanup)
		}

		canrmi := false
		cleanup := []string{"rmi"}
		for _, container := range containers {
			if !container.keep {
				canrmi = true
				cleanup = append(cleanup, container.tag)
			}
			if !arg_nopush && container.push != "" && engine_base != "podman" {
				cleanup = append(cleanup, container.push)
			}
		}
		if canrmi {
			cleanups = append(cleanups, cleanup)
		}
	}
}

func showplan() {
	for _, command := range commands {
		fmt.Println(arg_engine, strings.Join(command, " "))
	}
	for _, cleanup := range cleanups {
		fmt.Println(arg_engine, strings.Join(cleanup, " "))
	}
}
