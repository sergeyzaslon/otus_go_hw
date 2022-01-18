package main

import (
	"fmt"
	"os"
)

const exitError = 1

func main() {
	args := os.Args
	if len(args) < 3 {
		usage()
	}

	env, err := ReadDir(args[1])
	if err != nil {
		fmt.Printf("ERR: failed to read %s: %s", args[1], err.Error())
		os.Exit(exitError)
	}

	code := RunCmd(args[2:], env)

	os.Exit(code)
}

func usage() {
	fmt.Println("Usage: go-envdir /path/to/env/dir /path/to/command command arguments")
	os.Exit(0)
}
