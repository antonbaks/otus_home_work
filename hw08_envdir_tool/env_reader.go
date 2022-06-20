package main

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"strings"
)

type Environment map[string]EnvValue

var ErrIncorrectFileName = errors.New("incorrect file name")

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	dirEntries, errReadDir := os.ReadDir(dir)

	if errReadDir != nil {
		return nil, errReadDir
	}

	envInDir := make(map[string]EnvValue)

	for _, dirEntry := range dirEntries {
		needRemove := false

		if dirEntry.IsDir() {
			continue
		}

		dirEntityInfo, errInfo := dirEntry.Info()

		if errInfo != nil {
			continue
		}

		if strings.Contains(dirEntityInfo.Name(), "=") {
			return envInDir, ErrIncorrectFileName
		}

		if dirEntityInfo.Size() <= 0 {
			needRemove = true
		}

		fileString, errRead := os.ReadFile(filepath.Join(dir, dirEntityInfo.Name()))

		if errRead != nil {
			continue
		}

		fileString = []byte(strings.Split(string(fileString), "\n")[0])

		fileString = []byte(strings.TrimRight(string(fileString), "\t "))

		fileString = bytes.ReplaceAll(
			fileString,
			[]byte{0},
			[]byte("\n"))

		envInDir[dirEntry.Name()] = EnvValue{
			Value:      string(fileString),
			NeedRemove: needRemove,
		}
	}

	return envInDir, nil
}
