package main

import (
	"fmt"
	"os"
)

func main() {
	args := os.Args

	envInDir, err := ReadDir(args[1])
	if err != nil {
		fmt.Println(err)

		return
	}

	RunCmd(args[2:], envInDir)
}
