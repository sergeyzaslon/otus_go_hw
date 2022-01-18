package main

import (
	"bytes"
	"os"
	"path"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	env := make(Environment)
	for _, file := range files {
		fileInfo, err := file.Info()
		if err != nil {
			return nil, err
		}

		val := EnvValue{}

		if fileInfo.Size() == 0 {
			val.NeedRemove = true
			env[file.Name()] = val
			continue
		}

		content, err := os.ReadFile(path.Join(dir, file.Name()))
		if err != nil {
			return nil, err
		}

		content = bytes.Split(content, []byte("\n"))[0]
		content = bytes.TrimRight(content, " ")
		content = bytes.ReplaceAll(content, []byte{0x00}, []byte{'\n'})

		val.Value = string(content)

		env[file.Name()] = val
	}

	return env, nil
}
