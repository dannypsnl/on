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
	if len(cmdPipe) == 1 {
		onlyCmd := cmdPipe[0]
		onlyCmd.Stderr = os.Stderr
		onlyCmd.Stdout = os.Stdout
		onlyCmd.Stdin = os.Stdin
		return onlyCmd.Run()
	}
	prevCmd := cmdPipe[0]
	for _, curCmd := range cmdPipe[1:] {
		curCmd.Stdin, _ = prevCmd.StdoutPipe()
		err := prevCmd.Start()
		if err != nil {
			return err
		}
		prevCmd = curCmd
	}
	lastCmd := cmdPipe[len(cmdPipe)-1]
	lastCmd.Stdout = os.Stdout
	lastCmd.Stderr = os.Stderr
	return lastCmd.Run()
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
