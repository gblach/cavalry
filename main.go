package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path"
)

const version = "0.1.0"

var arg_chdir string
var arg_engine string
var arg_format string
var arg_nopush bool
var arg_norm bool
var arg_plan bool
var arg_version bool
var arg_help bool

func main() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}()

	def_engine, _ := exec.LookPath("docker")
	if def_engine == "" {
		def_engine, _ = exec.LookPath("podman")
	}
	if def_engine != "" {
		def_engine = path.Base(def_engine)
	}

	def_format := "docker"
	if def_engine == "podman" {
		def_format = "oci"
	}

	flag.StringVar(&arg_chdir, "c", ".", "Change working directory.")
	flag.StringVar(&arg_engine, "e", def_engine, "Choose the engine: podman or docker.")
	flag.StringVar(&arg_format, "f", def_format, "Choose the image format: oci or docker.")
	flag.BoolVar(&arg_nopush, "np", false, "Do not push images.")
	flag.BoolVar(&arg_norm, "nr", false, "Do not remove containers and images.")
	flag.BoolVar(&arg_plan, "p", false, "Show plan instead of executing them.")
	flag.BoolVar(&arg_version, "v", false, "Show version and exit.")
	flag.BoolVar(&arg_help, "h", false, "Show this message.")
	flag.Parse()

	if arg_help {
		arg0 := path.Base(os.Args[0])
		fmt.Fprint(os.Stderr, arg0, "[-c dir] [-e <podman|docker>] [-f <oci|docker>]")
		fmt.Fprintln(os.Stderr, " [-np] [-nr] [-p] [Cavalryfile]")
		fmt.Fprintln(os.Stderr, arg0, "-v")
		fmt.Fprintln(os.Stderr, arg0, "-h")
		flag.PrintDefaults()
		os.Exit(0)
	}

	if arg_version {
		fmt.Println("cavalry", version)
		os.Exit(0)
	}

	_, err := exec.LookPath(arg_engine)
	if err != nil {
		panic(err)
	}

	if arg_chdir != "." {
		err := os.Chdir(arg_chdir)
		if err != nil {
			panic(err)
		}
	}

	cavalryfile := "Cavalryfile"
	if flag.Arg(0) != "" {
		cavalryfile = flag.Arg(0)
	}

	loadfile(cavalryfile)
	makeplan()

	if arg_plan {
		showplan()
	} else {
		execute()
	}
}
