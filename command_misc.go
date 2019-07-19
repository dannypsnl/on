package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func runPipeline(cmdPipe ...*exec.Cmd) error {
	if len(cmdPipe) <= 0 {
		return fmt.Errorf("no commands for pipeline")
	}
	command := cmdPipe[0]
	if len(cmdPipe) > 1 {
		for _, curCmd := range cmdPipe[1:] {
			curCmd.Stdin, _ = command.StdoutPipe()
			err := command.Start()
			if err != nil {
				return err
			}
			command = curCmd
		}
	}
	command = cmdPipe[len(cmdPipe)-1]
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	return command.Run()
}

func newCommand(commandText string) *exec.Cmd {
	cmd := make([]string, 0)
	for _, elem := range strings.Split(commandText, " ") {
		if elem != "" {
			cmd = append(cmd, elem)
		}
	}
	return exec.Command(cmd[0], cmd[1:]...)
}
