package main

import (
	"fmt"
	"os"
	"os/exec"
)

var execCommand = exec.Command

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	for envName, envVal := range env {
		delEnvByExist(envName)

		if envVal.NeedRemove {
			continue
		}

		_ = os.Setenv(envName, envVal.Value)
	}

	oneCmd := execCommand(cmd[0], cmd[1:]...)
	oneCmd.Stdout = os.Stdout

	if err := oneCmd.Run(); err != nil {
		fmt.Println(err)

		return 1
	}

	return 0
}

func delEnvByExist(envName string) {
	if _, ok := os.LookupEnv(envName); !ok {
		return
	}

	_ = os.Unsetenv(envName)
}
