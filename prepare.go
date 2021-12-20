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

	format := ""
	if engine_base == "podman" {
		format = " --format " + arg_format
	}

	for i, container := range containers {
		command := fmt.Sprintf("%s build -t %s -f %s%s %s",
			arg_engine, container.tag, container.file, format, container.dir)
		commands = append(commands, command)

		network := ""
		if i > 0 {
			network = " --network container:" + containers[0].name
		}
		environ := ""
		for _, env := range container.env {
			environ += fmt.Sprintf(" -e '%s'", env)
		}
		command = fmt.Sprintf("%s run -dt --name %s%s%s %s",
			arg_engine, container.name, network, environ, container.tag)
		commands = append(commands, command)
	}

	for _, testcase := range testcases {
		command := fmt.Sprintf("%s exec %s %s",
			arg_engine, getname(testcase.tag), testcase.command)
		commands = append(commands, command)
	}

	if !arg_nopush {
		for _, container := range containers {
			if container.push != "" {
				if engine_base == "podman" {
					command := fmt.Sprintf("%s push %s %s",
						arg_engine, container.tag, container.push)
					commands = append(commands, command)
				} else {
					command := fmt.Sprintf("%s tag %s %s",
						arg_engine, container.tag, container.push)
					commands = append(commands, command)

					command = fmt.Sprintf("%s push %s",
						arg_engine, container.push)
					commands = append(commands, command)
				}
			}
		}
	}

	if !arg_norm {
		for i := len(containers) - 1; i >= 0; i-- {
			cleanup := fmt.Sprintf("%s rm -f %s", arg_engine, containers[i].name)
			cleanups = append(cleanups, cleanup)
		}

		canrmi := false
		cleanup := arg_engine + " rmi"
		for _, container := range containers {
			if !container.keep {
				canrmi = true
				cleanup += " " + container.tag
			}
			if !arg_nopush && container.push != "" && engine_base != "podman" {
				cleanup += " " + container.push
			}
		}
		if canrmi {
			cleanups = append(cleanups, cleanup)
		}
	}
}

func showplan() {
	for _, command := range commands {
		fmt.Println(command)
	}
	for _, cleanup := range cleanups {
		fmt.Println(cleanup)
	}
}
