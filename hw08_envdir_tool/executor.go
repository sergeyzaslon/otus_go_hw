package main

import (
	"errors"
	"os"
	"os/exec"
	"strings"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	if len(cmd) == 0 {
		return 0
	}

	command := cmd[0]
	args := cmd[1:]

	// read current environment variables
	envVars := make(map[string]string)
	for _, v := range os.Environ() {
		parts := strings.Split(v, "=")
		envVars[parts[0]] = strings.Join(parts[1:], "=")
	}

	// modify environment variables based on "env" rules
	for varName := range env {
		if env[varName].NeedRemove {
			delete(envVars, varName)
			continue
		}

		envVars[varName] = env[varName].Value
	}

	// Convert map to []strings
	envVarStrings := make([]string, 0)
	for varName := range envVars {
		envVarStrings = append(envVarStrings, varName+"="+envVars[varName])
	}

	c := exec.Command(command, args...)
	c.Env = envVarStrings
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	c.Stdin = os.Stdin

	err := c.Run()
	if err != nil {
		var e *exec.ExitError
		if errors.As(err, &e) {
			return e.ExitCode()
		}
		return 1
	}

	return 0
}
