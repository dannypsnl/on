package main

import (
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_command_cut(t *testing.T) {
	result := cutCommand("kubectl get po | grep xxx")

	expected := []*exec.Cmd{
		newCommand("kubectl get po"),
		newCommand("grep xxx"),
	}

	assert.Equal(t, expected, result)
}
